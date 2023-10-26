package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/parinyapt/prinflix_backend/logger"
)

func initializeEnvironmentFile() {
	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatal("Failed to load environment file", logger.Field("error", err))
	}
}

type envSpecification struct {
	Port       int    `required:"true" split_words:"true"`
	TZ         string `required:"true" split_words:"true"`
	AppName    string `required:"true" split_words:"true"`
	AppBaseUrl string `required:"true" split_words:"true"`

	DatabaseTablePrefix     string `required:"true" split_words:"true"`
	DatabaseMariadbHost     string `required:"true" split_words:"true"`
	DatabaseMariadbPort     string `split_words:"true"`
	DatabaseMariadbUsername string `required:"true" split_words:"true"`
	DatabaseMariadbPassword string `required:"true" split_words:"true"`
	DatabaseMariadbDbname   string `required:"true" split_words:"true"`
}

func initializeEnvironmentVariableCheck() {
	var envSpec envSpecification
	err := envconfig.Process("", &envSpec)
	if err != nil {
		logger.Fatal(err.Error())
	}
}
