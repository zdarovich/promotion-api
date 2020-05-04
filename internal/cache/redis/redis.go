package redis

import (
	"time"

	"github.com/zdarovich/promotion-api/internal/config"

	"github.com/go-redis/redis"
)

type (
	// Redis struct
	Redis struct {
		Configuration *config.Configuration
		Client        *redis.Client
	}
	// IRedis interface
	IRedis interface {
		Exists(key string) (interface{}, error)
		Get(key string) (interface{}, error)
		Set(key string, value interface{}) error
		SetX(key string, value interface{}, ttl time.Duration) error
	}
)

// New returns configured redis struct
func New(configuration *config.Configuration) IRedis {

	client := redis.NewClient(&redis.Options{
		Addr:         configuration.Redis.Server,
		DB:           configuration.Redis.DB,
		Password:     configuration.Redis.Password,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  time.Duration(configuration.Redis.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(configuration.Redis.WriteTimeout) * time.Second,
	})

	return &Redis{
		Configuration: configuration,
		Client:        client,
	}
}

// Exists returns true or false depending if the key exists
func (r *Redis) Exists(key string) (interface{}, error) {

	result := r.Client.Exists(key)
	return result.Val(), result.Err()
}

// Get returns the value of the key from the redis cache
func (r *Redis) Get(key string) (interface{}, error) {

	result := r.Client.Get(key)
	return result.Val(), result.Err()
}

// Set sets the new value to database
func (r *Redis) Set(key string, value interface{}) error {

	res := r.Client.Set(key, value, 0)
	return res.Err()
}

// SetX sets the new value to database with a timeout
func (r *Redis) SetX(key string, value interface{}, ttl time.Duration) error {

	res := r.Client.Set(key, value, ttl)
	return res.Err()
}
