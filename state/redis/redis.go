package redis

import (
	"context"

	"github.com/polyxia-org/morty-gateway/state"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

// adapter is an implementation of the state.State interface
type adapter struct {
	client *redis.Client
}

// Config hold the configuration about the Redis state adapter
type Config struct {
	Addr string `yaml:"addr"`
}

var _ state.State = (*adapter)(nil)

// NewState initializes a new state adapter for Redis based on the given configuration.
// An error could be returned if any errors happens during the adapter initialization.
func NewState(cfg *Config) (state.State, error) {
	log.Debugf("Bootstrapping Redis state adapter with options: %#v", cfg)
	client := redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
		DB:   0,
	})

	// Enable Keyspace events as we will need them to handle function instances expiration
	if _, err := client.ConfigSet(context.Background(), "notify-keyspace-events", "KEA").Result(); err != nil {
		log.Errorf("Failed to enable Redis Keyspace Events: %v", err)
		return nil, err
	}

	log.Info("State engine 'redis' successfully initialized")
	return &adapter{client}, nil
}

func (a *adapter) Get(ctx context.Context, key string) (string, error) {
	r := a.client.Get(ctx, key)
	log.Tracef("state/redis: %s", r.String())
	return r.Result()
}

func (a *adapter) Set(ctx context.Context, key, value string) error {
	// 0 as the key doesn't expire at the moment
	r := a.client.Set(ctx, key, value, 0)
	log.Tracef("state/redis: %s", r.String())
	_, err := r.Result()
	return err
}

func (a *adapter) SetMultiple(ctx context.Context, tuples map[string]string) []error {
	errors := []error{}
	for k, v := range tuples {
		if err := a.Set(ctx, k, v); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}
