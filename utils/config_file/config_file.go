package utilsConfigFile

import "github.com/spf13/viper"

const (
	ConnectOAuthSuccessPagePath = "connect_oauth.success"
	ConnectOAuthFailPagePath    = "connect_oauth.fail"
	EmailVerifySuccessPagePath  = "email_verification.success"
	EmailVerifyFailPagePath     = "email_verification.fail"
	ResetPasswordPagePath       = "reset_password_page"

	AppLoginOAuthSuccessPagePath = "oauth.success"
	AppLoginOAuthFailPagePath    = "oauth.fail"
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

func GetAppBaseURL() string {
	return viper.GetString("application.base_url")
}

func GetAppRedirectPath(page string) string {
	return viper.GetString("application.redirect_path." + page)
}