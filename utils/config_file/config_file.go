package utilsConfigFile

import "github.com/spf13/viper"

const (
	ConnectOAuthSuccessPagePath = "connect_oauth.success"
	ConnectOAuthFailPagePath    = "connect_oauth.fail"
	EmailVerifySuccessPagePath  = "email_verification.success"
	EmailVerifyFailPagePath     = "email_verification.fail"
	ResetPasswordPagePath       = "reset_password_page"
)

func GetConfigDomain() []string {
	return viper.GetStringSlice("domain")
}

func GetFrontendBaseURL() string {
	return viper.GetString("base_url")
}

func GetRedirectPagePath(page string) string {
	return viper.GetString("path." + page)
}
