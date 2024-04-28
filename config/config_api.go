package config

import (
	"github.com/parinyapt/prinflix_backend/logger"
	utilsConfigFile "github.com/parinyapt/prinflix_backend/utils/config_file"
	"github.com/spf13/viper"
)

var CorsAllowOrigins []string = []string{"https://appleid.apple.com"}

func initializeAPIConfigFile() {
	viper.SetConfigName("api_config") // name of config file (without extension)
	viper.SetConfigType("json")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./config-file")      // optionally look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		logger.Fatal("Failed to load api_config.json file", logger.Field("error", err))
	}

	CorsAllowOrigins = append(CorsAllowOrigins, utilsConfigFile.GetConfigDomain()...)
}
