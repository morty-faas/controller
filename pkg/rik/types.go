package rik

type (
	getWorkloadsResponse []RegisteredWorkload

	createWorkloadResponse struct {
		WorkloadId string `json:"id"`
	}

	CreateInstanceRequest struct {
		WorkloadId   string `json:"workload_id"`
		InstanceName string `json:"name"`
	}

	fetchInstancesResponse struct {
		Instances []Instance `json:"instances"`
	}
)
