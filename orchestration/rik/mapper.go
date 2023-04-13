package rik

import (
	"github.com/polyxia-org/morty-gateway/pkg/rik"
	"github.com/polyxia-org/morty-gateway/types"
)

// mapRegisteredWorkloadToFn is a helper function that maps a RIK Workload to a Morty function
func mapRegisteredWorkloadToFn(wk *rik.RegisteredWorkload) *types.Function {
	return &types.Function{
		Id:       wk.WorkloadId,
		Name:     wk.Name,
		ImageURL: wk.Workload.Spec.Function.Executor.Rootfs,
	}
}

// mapFnToWorkload is a helper function that maps a Morty function to a RIK Workload
func mapFnToWorkload(fn *types.Function) *rik.Workload {
	return &rik.Workload{
		ApiVersion: "v0",
		Kind:       rikFunctionKind,
		Name:       fn.Name,
		Spec: rik.Spec{
			Function: rik.Fn{
				Executor: rik.Executor{
					Rootfs: fn.ImageURL,
				},
				Exposure: rik.Exposure{
					Type: "NodePort",
				},
			},
		},
	}
}
