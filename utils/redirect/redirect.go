package utilsRedirect

import (
	"strings"

	utilsConfigFile "github.com/parinyapt/prinflix_backend/utils/config_file"
)

const (
	ProviderLine   = "Line"
	ProviderGoogle = "Google"
	ProviderApple  = "Apple"
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

func GenerateAppOAuthConnectRedirectUrl(platform string, success bool) (returnPath string) {
	returnPath = utilsConfigFile.GetAppBaseURL()
	if success {
		returnPath += utilsConfigFile.GetAppRedirectPath(utilsConfigFile.AppConnectOAuthSuccessPagePath)
	} else {
		returnPath += utilsConfigFile.GetAppRedirectPath(utilsConfigFile.AppConnectOAuthFailPagePath)
	}

	returnPath = strings.Replace(returnPath, ":provider", platform, -1)

	return returnPath
}

func GenerateOAuthLoginRedirectUrl(platform string, success bool, code string) (returnPath string) {
	returnPath = utilsConfigFile.GetAppBaseURL()
	if success {
		returnPath += utilsConfigFile.GetAppRedirectPath(utilsConfigFile.AppLoginOAuthSuccessPagePath)
		returnPath = strings.Replace(returnPath, ":auth_temp_code", code, -1)
	} else {
		returnPath += utilsConfigFile.GetAppRedirectPath(utilsConfigFile.AppLoginOAuthFailPagePath)
	}

	returnPath = strings.Replace(returnPath, ":provider", platform, -1)

	return returnPath
}
