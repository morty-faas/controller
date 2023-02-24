package rik

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/polyxia-org/morty-gateway/config"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
)

type ControllerClient struct {
	// c is the HTTP client used to communicate with the RIK controller
	c *http.Client
	// baseUrl is the base URL of the RIK controller
	baseUrl *url.URL
	// l is the logger used to log messages
	l *logrus.Entry
}

func NewControllerClient(l *logrus.Entry, config config.Config) (*ControllerClient, error) {
	client := &http.Client{}
	return &ControllerClient{
		c:       client,
		baseUrl: config.RIKController,
		l:       l.WithField("component", "rik-controller-client"),
	}, nil
}

func (client ControllerClient) formatUrl(path string) (string, error) {
	return url.JoinPath(client.baseUrl.String(), path)
}

func (client ControllerClient) mustFormatUrl(path string) string {
	formattedUrl, err := client.formatUrl(path)
	if err != nil {
		panic(err)
	}
	return formattedUrl
}

func (client ControllerClient) GetExistingFunctions() (map[string]string, error) {
	instances := map[string]string{}
	client.l.Debug("Get existing functions")

	path := client.mustFormatUrl("/api/v0/workloads.list")
	res, err := client.c.Get(path)

	var workloads []struct {
		WorkloadId string   `json:"id"`
		Name       string   `json:"name"`
		Workload   workload `json:"workload"`
	}

	if err != nil {
		client.l.WithError(err).Error("Could not get existing functions")
		return instances, err
	}

	if res != nil {
		err = json.NewDecoder(res.Body).Decode(&workloads)
		if err != nil {
			client.l.WithError(err).Error("Could not decode response body")
			return instances, err
		}
	}

	// Because with golang we cannot do a map on data...
	for _, entry := range workloads {
		instances[entry.Name] = entry.WorkloadId
	}

	client.l.Infof("Found %d existing functions", len(instances))

	return instances, nil
}

// CreateFunction takes a FunctionRequest request and returns the response given
// by the RIK controller
func (client ControllerClient) CreateFunction(functionBody FunctionRequest) (string, error) {
	l := client.l.WithFields(logrus.Fields{
		"functionName": functionBody.Name,
	})
	l.Debug("Create function")
	fullUrl := client.mustFormatUrl("/api/v0/workloads.create")

	workloadBody := newWorkloadRequest(functionBody.Name, functionBody.Rootfs)
	requestBody, err := json.Marshal(workloadBody)

	res, err := client.c.Post(fullUrl, "application/json", bytes.NewReader(requestBody))
	if err != nil {
		l.WithError(err).Error("Failed to create function, could not send request to RIK controller")
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		l.WithError(err).Error("Failed to create function, RIK controller returned non-OK status code")
		return "", fmt.Errorf("RIK controller returned non-OK status code: %d", res.StatusCode)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		l.WithError(err).Error("Could not read response body")
		return "", err
	}

	workloadResponse := new(createWorkloadResponse)
	err = json.Unmarshal(b, &workloadResponse)

	l.WithField("name", functionBody.Name).Debug("Function created successfully")
	return workloadResponse.WorkloadId, nil
}

func (client ControllerClient) FetchInstances(workloadId string) ([]Instance, error) {
	l := client.l.WithFields(logrus.Fields{
		"workloadId": workloadId,
	})
	l.Debug("Fetch instances")
	path := fmt.Sprintf("/api/v0/workloads.instances/%s", workloadId)
	fullUrl, err := client.formatUrl(path)
	if err != nil {
		l.WithError(err).Error("Could not format URL to fetch instances")
		return []Instance{}, err
	}

	res, err := client.c.Get(fullUrl)
	if err != nil {
		l.WithError(err).Error("Failed to fetch instances, could not send request to RIK controller")
		return []Instance{}, err
	}

	if res.StatusCode == http.StatusNotFound {
		l.WithError(err).Error("Failed to fetch instances, RIK controller returned 404")
		return []Instance{}, fmt.Errorf("Workload does not exist")
	}

	if res.StatusCode == http.StatusNoContent {
		l.WithError(err).Warn("Failed to fetch instances, RIK controller returned 204")
		return []Instance{}, nil
	}

	if res.StatusCode != http.StatusOK {
		l.WithError(err).Error("Failed to fetch instances, RIK controller returned non-OK status code")
		return []Instance{}, fmt.Errorf("RIK controller returned non-OK status code: %d", res.StatusCode)
	}

	instances := new(fetchInstancesResponse)
	err = json.NewDecoder(res.Body).Decode(instances)
	if err != nil {
		l.WithError(err).Error("Failed to fetch instances, could not decode response body")
		return []Instance{}, err
	}

	l.WithField("response", res).Debug("Instances fetched successfully")
	return instances.Instances, nil
}

func (client ControllerClient) CreateWorkloadInstance(workloadId string) error {
	instanceName := generateInstanceName()
	l := client.l.WithField("workloadId", workloadId).WithField("instanceName", instanceName)
	l.Debug("Create workload instance")
	path, err := client.formatUrl("/api/v0/instances.create")
	if err != nil {
		l.WithError(err).Error("Could not format URL to create instance")
		return err
	}

	req := newInstanceRequest(workloadId, instanceName)
	reqBody, err := json.Marshal(req)
	if err != nil {
		l.WithError(err).Warn("Failed to create instance, could not marshal request body")
		return err
	}

	res, err := client.c.Post(path, "application/json", bytes.NewReader(reqBody))

	if err != nil {
		l.WithError(err).Error("Failed to create instance, could not send request to RIK controller")
		return err
	}

	if res.StatusCode != http.StatusCreated {
		l.WithError(err).Error("Failed to create instance, RIK controller returned non-OK status code")
		return fmt.Errorf("RIK controller returned non-OK status code: %d", res.StatusCode)
	}

	l.WithField("response", res).Debug("Instance created successfully")
	return nil
}
