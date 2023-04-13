package rik

type RegisteredWorkload struct {
	WorkloadId string   `json:"id"`
	Name       string   `json:"name"`
	Workload   Workload `json:"value"`
}

type Workload struct {
	ApiVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Name       string `json:"name"`
	Spec       Spec   `json:"spec"`
}

type Spec struct {
	Function Fn `json:"function"`
}

type Fn struct {
	Executor Executor `json:"execution"`
	Exposure Exposure `json:"exposure"`
}

type Executor struct {
	Rootfs string `json:"rootfs"`
}

type Exposure struct {
	Port       int    `json:"port"`
	TargetPort int    `json:"targetPort"`
	Type       string `json:"type"`
}

type Instance struct {
	Id         string `json:"id"`
	Kind       string `json:"kind"`
	Namespace  string `json:"namespace"`
	Spec       Spec   `json:"spec"`
	Status     string `json:"status"`
	WorkloadID string `json:"workload_id"`
}
