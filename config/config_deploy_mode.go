package config

import (
	"flag"
	"log"
	"os"
)

func InitializeDeployModeFlag() {
	DeployModeFlag := flag.String("mode", "development", "deploy mode (development, production)")
	flag.Parse()

	if (*DeployModeFlag == "development") || (*DeployModeFlag == "production") {
		os.Setenv("DEPLOY_MODE", *DeployModeFlag)
		log.Printf("Deploy Mode : %s \n", os.Getenv("DEPLOY_MODE"))
	} else {
		log.Fatalf("Please set deploy mode to 'development' or 'production'")
		return
	}
}