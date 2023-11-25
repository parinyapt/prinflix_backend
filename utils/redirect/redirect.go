package utilsRedirect

import (
	"strings"

	utilsConfigFile "github.com/parinyapt/prinflix_backend/utils/config_file"
)

const (
	ProviderLine = "Line"
	ProviderGoogle = "Google"
)

func GenerateOAuthConnectRedirectUrl(platform string, success bool) (returnPath string) {
	returnPath = utilsConfigFile.GetFrontendBaseURL()
	if success {
		returnPath += utilsConfigFile.GetRedirectPagePath(utilsConfigFile.ConnectOAuthSuccessPagePath)
	} else {
		returnPath += utilsConfigFile.GetRedirectPagePath(utilsConfigFile.ConnectOAuthFailPagePath)
	}

	returnPath = strings.Replace(returnPath, ":provider", platform, -1)

	return returnPath
}