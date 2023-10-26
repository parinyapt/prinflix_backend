package config

import (
	"os"

	"github.com/parinyapt/golang_utils/timezone/v1"
	
	"github.com/parinyapt/prinflix_backend/logger"
)

func initializeGlobalTimezone() {
	// Global TimeZone Setup
	if err := PTGUtimezone.GlobalTimezoneSetup(os.Getenv("TZ")); err != nil {
		logger.Fatal("Fail to set Global Timezone", logger.Field("error", err.Error()))
	}
}