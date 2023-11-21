package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

const (
	_defaultAddr     = "localhost:6379"
	_defaultPassword = ""
	_defaultDb       = 0
)

type Redis struct {
	*redis.Client
}

func New(opts ...Option) (*Redis, error) {
	o := &redis.Options{
		Addr:     _defaultAddr,
		Password: _defaultPassword,
		DB:       _defaultDb,
	}

	for _, opt := range opts {
		opt(o)
	}

	client := redis.NewClient(o)
	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return &Redis{client}, nil
}
