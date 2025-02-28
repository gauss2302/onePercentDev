package auth

import (
	"fmt"
	"time"
)

type JWTConfig struct {
	AccessSecret  string        `mapstructure:"access_secret" description:"Secret key for signing access tokens"`
	RefreshSecret string        `mapstructure:"refresh_secret" description:"Secret key for signing refresh tokens"`
	AccessExpire  time.Duration `mapstructure:"access_expire" description:"Expiration duration for access tokens"`
	RefreshExpire time.Duration `mapstructure:"refresh_expire" description:"Expiration duration for refresh tokens"`
}

func (j *JWTConfig) Validate() error {
	if j.AccessSecret == "" {
		return fmt.Errorf("JWT access secret is required")
	}
	if j.RefreshSecret == "" {
		return fmt.Errorf("JWT refresh secret is required")
	}
	if j.AccessExpire <= 0 {
		return fmt.Errorf("JWT access token expiration must be greater than 0")
	}
	if j.RefreshExpire <= 0 {
		return fmt.Errorf("JWT refresh token expiration must be greater than 0")
	}
	return nil
}
