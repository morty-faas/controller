package rik

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"time"

	"github.com/polyxia-org/morty-gateway/orchestration"
	"github.com/polyxia-org/morty-gateway/pkg/rik"
	"github.com/polyxia-org/morty-gateway/types"
	log "github.com/sirupsen/logrus"
)

type adapter struct {
	cfg    *Config
	client *rik.Client
}

type Config struct {
	// Cluster is the address of the RIK controller
	Cluster string `yaml:"cluster"`
}

const rikFunctionKind = "Function"

var _ orchestration.Orchestrator = (*adapter)(nil)

// NewOrchestrator initializes the RIK orchestrator adapter.
func NewOrchestrator(cfg *Config) (orchestration.Orchestrator, error) {
	log.Info("Orchestrator engine 'rik' successfully initialized")
	return &adapter{cfg, rik.NewClient(cfg.Cluster)}, nil
}

func (a *adapter) GetFunctions(ctx context.Context) ([]*types.Function, error) {
	workloads, err := a.client.GetWorkloads(ctx)
	if err != nil {
		return nil, err
	}

	var functions []*types.Function
	for _, workload := range *workloads {
		// Filter on function elements only
		if workload.Workload.Kind == rikFunctionKind {
			functions = append(functions, &types.Function{
				Id:       workload.WorkloadId,
				Name:     workload.Name,
				ImageURL: workload.Workload.Spec.Function.Executor.Rootfs,
			})
		}
	}

	return functions, nil
}

func (a *adapter) CreateFunction(ctx context.Context, fn *types.Function) (*types.Function, error) {
	res, err := a.client.CreateWorkload(ctx, mapFnToWorkload(fn))
	if err != nil {
		return nil, err
	}

	fn.Id = res.WorkloadId
	return fn, nil
}

func (a *adapter) GetFunctionInstance(ctx context.Context, fn *types.Function) (*types.FnInstance, error) {
	instances, err := a.client.GetWorkloadInstances(ctx, fn.Id)
	if err != nil {
		return nil, err
	}

	if len(instances) == 0 {
		log.Debugf("Deploying new instance for function: %+v", fn)

		in := &rik.CreateInstanceRequest{
			WorkloadId:   fn.Id,
			InstanceName: fn.Name,
		}

		if err := a.client.CreateWorkloadInstance(ctx, in); err != nil {
			err := fmt.Errorf("Failed to create instance: %v", err)
			log.Error(err)
			return nil, err
		}

		time.Sleep(500 * time.Millisecond)

		instances, err = a.client.GetWorkloadInstances(ctx, fn.Id)
		if err != nil {
			return nil, err
		}
	}

	log.Debugf("%d instance(s)", len(instances))

	rikIn := instances[rand.Intn(len(instances))]

	url, _ := url.Parse(a.cfg.Cluster)
	url, _ = url.Parse(fmt.Sprintf("%s://%s:%d", url.Scheme, url.Hostname(), rikIn.Spec.Function.Exposure.Port))

	instance := &types.FnInstance{
		Function: fn,
		Endpoint: url,
	}

	return instance, nil
}
