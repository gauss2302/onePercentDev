package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"onepercentdev_server/pkg/auth"
	"os"
	"strings"
)

type Config struct {
	DB     PostgresConfig `mapstructure:"postgres"`
	Server ServerConfig   `mapstructure:"server"`
	JWT    auth.JWTConfig `mapstructure:"jwt"` // Используем JWTConfig из pkg/auth
	Redis  RedisConfig    `mapstructure:"redis"`
}

func LoadConfig(configPath string) (*Config, error) {
	if os.Getenv("ENVIRONMENT") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}

	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if configPath != "" {
		viper.SetConfigFile(configPath)
		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %v", err)
		}
	}

	bindEnvs()

	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", "8090")

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %v", err)
	}

	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %v", err)
	}

	log.Println("Configuration loaded successfully")
	return &config, nil
}

func bindEnvs() {
	envBindings := map[string]string{
		"postgres.url":      "DATABASE_URL",
		"server.port":       "APP_SERVER_PORT",
		"server.host":       "APP_SERVER_HOST",
		"jwt.accessSecret":  "APP_JWT_ACCESSSECRET",
		"jwt.refreshSecret": "APP_JWT_REFRESHSECRET",
		"jwt.accessExpire":  "APP_JWT_ACCESSEXPIRE",
		"jwt.refreshExpire": "APP_JWT_REFRESHEXPIRE",
		"redis.url":         "REDIS_URL",
		"redis.password":    "REDIS_PASSWORD",
		"redis.db":          "REDIS_DB",
	}

	for configKey, envKey := range envBindings {
		if err := viper.BindEnv(configKey, envKey); err != nil {
			log.Fatalf("Failed to bind environment variable: %v", err)
		}
	}
}

func validateConfig(config *Config) error {
	if err := config.DB.Validate(); err != nil {
		return fmt.Errorf("PostgreSQL config validation failed: %v", err)
	}
	if err := config.Redis.Validate(); err != nil {
		return fmt.Errorf("Redis config validation failed: %v", err)
	}
	if err := config.JWT.Validate(); err != nil { // Валидация JWTConfig
		return fmt.Errorf("JWT config validation failed: %v", err)
	}
	if err := config.Server.Validate(); err != nil {
		return fmt.Errorf("Server config validation failed: %v", err)
	}
	return nil
}
