package rik

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"

	httpclient "github.com/polyxia-org/morty-gateway/pkg/client"
	"github.com/polyxia-org/morty-gateway/pkg/serdejson"
	"github.com/sirupsen/logrus"
)

type Client struct {
	c *httpclient.Client
}

var (
	ErrWorkloadNotFound = errors.New("workload not found")
)

// NewClient initiate a new client to communicate with a RIK cluster.
func NewClient(baseURL string) *Client {
	return &Client{httpclient.NewClient(baseURL)}
}

// GetWorkloads fetch all the workloads from the RIK cluster
func (c *Client) GetWorkloads(ctx context.Context) (*getWorkloadsResponse, error) {
	res, err := c.c.Get(ctx, "api/v0/workloads.list", nil)
	if err != nil {
		return nil, err
	}

	return serdejson.Deserialize[getWorkloadsResponse](res.Body)
}

// CreateWorkload create the given workload into the RIK Cluster
func (c *Client) CreateWorkload(ctx context.Context, workload *Workload) (*createWorkloadResponse, error) {
	by, err := serdejson.Serialize(workload)
	if err != nil {
		return nil, err
	}

	res, err := c.c.Post(ctx, "api/v0/workloads.create", bytes.NewReader(by), nil)
	if err != nil {
		return nil, err
	}

	return serdejson.Deserialize[createWorkloadResponse](res.Body)
}

// GetWorkloadInstances retrieve all the instances for the given workload
func (c *Client) GetWorkloadInstances(ctx context.Context, id string) ([]Instance, error) {
	res, err := c.c.Get(ctx, fmt.Sprintf("api/v0/workloads.instances/%s", id), nil)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusNotFound {
		logrus.Errorf("RIK workload not found: '%s'", id)
		return []Instance{}, ErrWorkloadNotFound
	}

	if res.StatusCode == http.StatusNoContent {
		logrus.Warnf("No instances for workload '%s'", id)
		return []Instance{}, nil
	}

	data, err := serdejson.Deserialize[fetchInstancesResponse](res.Body)
	if err != nil {
		return nil, err
	}

	return data.Instances, nil
}

// CreateWorkloadInstance schedule an instance for the given workload in the RIK cluster
func (c *Client) CreateWorkloadInstance(ctx context.Context, input *CreateInstanceRequest) error {
	by, err := serdejson.Serialize(input)
	if err != nil {
		return err
	}

	res, err := c.c.Post(ctx, "api/v0/instances.create", bytes.NewReader(by), nil)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("RIK returned non-OK status code: %v", res.StatusCode)
	}

	return nil
}
