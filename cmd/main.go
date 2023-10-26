package main

import (
	"os"

	"github.com/parinyapt/prinflix_backend/config"
	"github.com/parinyapt/prinflix_backend/database"
	"github.com/parinyapt/prinflix_backend/logger"
	"github.com/parinyapt/prinflix_backend/routes"
)

func main() {
	config.InitializeDeployModeFlag()

	logger.InitializeLogger(os.Getenv("DEPLOY_MODE"))
	config.InitializeConfig()
	database.InitializeDatabase()
	routes.InitializeRoutes()
}