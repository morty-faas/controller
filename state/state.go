package state

import (
	"context"
	"errors"

	"github.com/polyxia-org/morty-gateway/types"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

// State is a generic interface for our controller state
type State interface {
	// Get retrieve the value associated to the given key.
	// If the key doesn't exists, an error ErrKeyNotFound will be returned
	Get(ctx context.Context, key string) (*types.Function, error)
	// Set a tuple of key/value in the state
	Set(ctx context.Context, fn *types.Function) error
	// SetMultiple set multiple keys in one call
	SetMultiple(ctx context.Context, functions []*types.Function) []error
}
