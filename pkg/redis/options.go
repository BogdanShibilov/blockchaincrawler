package redis

import "github.com/redis/go-redis/v9"

type Option func(*redis.Options)

func Address(url string) Option {
	return func(o *redis.Options) {
		o.Addr = url
	}
}

func Password(pass string) Option {
	return func(o *redis.Options) {
		o.Password = pass
	}
}
