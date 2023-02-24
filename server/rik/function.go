package rik

type FunctionRequest struct {
	Name   string `json:"name"`
	Rootfs string `json:"rootfs"`
}

// FunctionResponse is the request body for a FunctionRequest invocation
// This is typically served by a runtime of RIK
type FunctionResponse struct {
	// Considered as string as we don't do anything to this field, we'll just return it
	Payload string `json:"payload"`
	// ProcessMetadata contains metadata about the function execution
	ProcessMetadata functionProcessMetadata `json:"process_metadata"`
}

type functionProcessMetadata struct {
	ExecutionTimeMs int `json:"execution_time_ms"`
	// Array of strings containing the logs of the function execution
	// it might be equal to nil if the runtime doesn't support logs
	Logs []string `json:"logs"`
}
