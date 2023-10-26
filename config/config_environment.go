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

	DatabaseTablePrefix        string `required:"true" split_words:"true"`
	DatabasePostgresqlHost     string `required:"true" split_words:"true"`
	DatabasePostgresqlPort     int    `required:"true" split_words:"true"`
	DatabasePostgresqlUsername string `required:"true" split_words:"true"`
	DatabasePostgresqlPassword string `required:"true" split_words:"true"`
	DatabasePostgresqlDbname   string `required:"true" split_words:"true"`

	Oauth2GoogleRedirectUrl   string `required:"true" split_words:"true"`
	Oauth2FacebookRedirectUrl string `required:"true" split_words:"true"`
	Oauth2AppleRedirectUrl    string `required:"true" split_words:"true"`

	ObjectStorageEndpoint        string `required:"true" split_words:"true"`
	ObjectStorageAccessKey       string `required:"true" split_words:"true"`
	ObjectStorageSecretAccessKey string `required:"true" split_words:"true"`
	ObjectStorageBucketName      string `required:"true" split_words:"true"`

	JwtSignKeyAccessToken  string `required:"true" split_words:"true"`
	JwtSignKeyRefreshToken string `required:"true" split_words:"true"`

	EncryptionKeyAuthTempCode string `required:"true" split_words:"true"`
	EncryptionKeyStorage      string `required:"true" split_words:"true"`
}

func initializeEnvironmentVariableCheck() {
	var envSpec envSpecification
	err := envconfig.Process("", &envSpec)
	if err != nil {
		logger.Fatal(err.Error())
	}
}
