package rik

import (
	"fmt"
	"net/url"
	"time"

	"github.com/goombaio/namegenerator"
)

type workloadListResponse struct {
	entries []struct {
		WorkloadId string   `json:"id"`
		Name       string   `json:"name"`
		Workload   workload `json:"workload"`
	}
}

type createWorkloadResponse struct {
	WorkloadId string `json:"id"`
}

type workload struct {
	ApiVersion string       `json:"apiVersion"`
	Kind       string       `json:"kind"`
	Name       string       `json:"name"`
	Spec       workloadSpec `json:"spec"`
}

type workloadExecutor struct {
	Rootfs string `json:"rootfs"`
}

type workloadSpecFunction struct {
	Executor workloadExecutor `json:"execution"`
}

type workloadSpec struct {
	Containers []interface{}        `json:"containers"`
	Function   workloadSpecFunction `json:"function"`
}

func newWorkloadRequest(name string, rootfs string) workload {
	return workload{
		ApiVersion: "v0",
		Kind:       "Function",
		Name:       name,
		Spec: workloadSpec{
			Containers: []interface{}{},
			Function: workloadSpecFunction{
				Executor: workloadExecutor{
					Rootfs: rootfs,
				},
			},
		},
	}
}

type createInstanceRequest struct {
	WorkloadId   string `json:"workload_id"`
	InstanceName string `json:"name"`
}

func newInstanceRequest(workloadId string, instanceName string) createInstanceRequest {
	return createInstanceRequest{
		WorkloadId:   workloadId,
		InstanceName: instanceName,
	}
}

type Instance struct {
	ID         string `json:"id"`
	Kind       string `json:"kind"`
	Namespace  string `json:"namespace"`
	Spec       Spec   `json:"spec"`
	Status     string `json:"status"`
	WorkloadID string `json:"workload_id"`
}
type FunctionExecution struct {
	Rootfs string `json:"rootfs"`
}
type FunctionExposure struct {
	Port       int    `json:"port"`
	TargetPort int    `json:"targetPort"`
	Type       string `json:"type"`
}
type Function struct {
	Execution FunctionExecution `json:"execution"`
	Exposure  FunctionExposure  `json:"exposure"`
}
type Spec struct {
	Containers []any    `json:"containers"`
	Function   Function `json:"function"`
}

func (i *Instance) GetRuntimeUrl() *url.URL {
	urlStr := fmt.Sprintf("http://%s:%d", "127.0.0.1", i.Spec.Function.Exposure.Port)
	fullUrl, _ := url.Parse(urlStr)
	return fullUrl
}

type fetchInstancesResponse struct {
	Instances []Instance `json:"instances"`
}

func generateInstanceName() string {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)

	return nameGenerator.Generate()
}
