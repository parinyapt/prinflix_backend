package config

import (
	"os"

	"github.com/parinyapt/prinflix_backend/logger"
)

func InitializeConfig() {
	if os.Getenv("DEPLOY_MODE") == "development" {
		initializeEnvironmentFile()
	}
	if os.Getenv("DEPLOY_MODE") == "production" {
		initializeSetGinReleaseMode()
	}
	initializeEnvironmentVariableCheck()
	initializeGlobalTimezone()

	logger.Info("Initialize Config Success")
}