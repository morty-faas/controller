package memory

import (
	"context"

	"github.com/polyxia-org/morty-gateway/state"
	log "github.com/sirupsen/logrus"
)

// adapter is an implementation of the state.State interface
type adapter struct {
	store map[string]string
}

var _ state.State = (*adapter)(nil)

// NewState initializes a new state adapter for Memory engine.
func NewState() state.State {
	log.Info("State engine 'memory' successfully initialized")
	return &adapter{
		store: make(map[string]string),
	}
}

func (a *adapter) Get(ctx context.Context, key string) (string, error) {
	log.Tracef("state/memory: retrieving value for key '%s'", key)
	v, exists := a.store[key]
	if !exists {
		return "", state.ErrKeyNotFound
	}
	return v, nil
}

func (a *adapter) Set(ctx context.Context, key, value string) error {
	log.Tracef("state/memory: setting value '%s' for key '%s'", key, value)
	a.store[key] = value
	return nil
}

func (a *adapter) SetMultiple(ctx context.Context, tuples map[string]string) []error {
	for k, v := range tuples {
		a.Set(ctx, k, v)
	}
	return nil
}
