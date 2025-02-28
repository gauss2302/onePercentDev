package app

import (
	"fmt"
	"onepercentdev_server/config"
)

type App struct {
	config *config.Config
}

func NewApp(*App, error) {

	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	return &App{
		config: cfg,
	}, nil
}
