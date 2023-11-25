package controller

import (
	"os"

	PTGUoauthLine "github.com/parinyapt/golang_utils/oauth/line/v1"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GenerateGoogleOAuthURL(state string) (url string) {
	oauthConfig := &oauth2.Config{
		RedirectURL:  os.Getenv("OAUTH2_GOOGLE_REDIRECT_URL"),
		ClientID:     os.Getenv("OAUTH2_GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH2_GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return oauthConfig.AuthCodeURL(state)
}

func GenerateLineOAuthURL(state string) (url string) {
	oauthConfig := PTGUoauthLine.NewLineOAuth(&PTGUoauthLine.LineOAuthConfig{
		ClientID:     os.Getenv("OAUTH2_LINE_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH2_LINE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("OAUTH2_LINE_REDIRECT_URL"),
	})

	return oauthConfig.GenerateOAuthURL(PTGUoauthLine.OptionLineGenerateOAuthURL{
		Scopes: []string{
			"openid",
			"profile",
			"email",
		},
		State: state,
	})
}
