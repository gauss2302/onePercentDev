package config

import "fmt"

type RedisConfig struct {
	URL      string `mapstructure:"url" description:"Redis connection URL"`
	Password string `mapstructure:"password" description:"Redis password"`
	DB       int    `mapstructure:"db" description:"Redis database number"`
}

func (r *RedisConfig) Validate() error {
	if r.URL == "" {
		return fmt.Errorf("redis.url is required")
	}

	if r.DB < 0 {
		return fmt.Errorf("redis.db must be a positive number")
	}
	return nil
}

func (r *RedisConfig) GetOptions() map[string]interface{} {
	options := map[string]interface{}{
		"addr":     r.URL,
		"password": r.Password,
		"db":       r.DB,
	}
	return options
}
