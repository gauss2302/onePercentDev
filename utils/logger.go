package utils

import "go.uber.org/zap"

var logger *zap.Logger

func InitializeLogger(env string) {
	var err error
	if env == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
}

func GetLogger() *zap.Logger {
	if logger == nil {
		panic("Logger is not initialized")
	}
	return logger
}
