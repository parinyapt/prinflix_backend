package main

import (
	"os"

	"github.com/parinyapt/prinflix_backend/config"
	"github.com/parinyapt/prinflix_backend/database"
	"github.com/parinyapt/prinflix_backend/logger"
	"github.com/parinyapt/prinflix_backend/routes"
	"github.com/parinyapt/prinflix_backend/storage"
)

func main() {
	config.InitializeDeployModeFlag()

	logger.InitializeLogger(os.Getenv("DEPLOY_MODE"))
	config.InitializeConfig()
	database.InitializeDatabase()
	storage.InitializeStorage()
	routes.InitializeRoutes()
}