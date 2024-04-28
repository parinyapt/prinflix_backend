package controller

import (
	"os"

	PTGUdata "github.com/parinyapt/golang_utils/data/v1"
	PTGUoauthGoogle "github.com/parinyapt/golang_utils/oauth/google/v1"
	PTGUoauthLine "github.com/parinyapt/golang_utils/oauth/line/v1"
	modelController "github.com/parinyapt/prinflix_backend/model/controller"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func googleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  os.Getenv("OAUTH2_GOOGLE_REDIRECT_URL"),
		ClientID:     os.Getenv("OAUTH2_GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH2_GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func googleOAuthConfigV2() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  os.Getenv("OAUTH2_GOOGLE_REDIRECT_URL_V2"),
		ClientID:     os.Getenv("OAUTH2_GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH2_GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func lineOAuthConfig() *PTGUoauthLine.LineOAuthConfig {
	return &PTGUoauthLine.LineOAuthConfig{
		ClientID:     os.Getenv("OAUTH2_LINE_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH2_LINE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("OAUTH2_LINE_REDIRECT_URL"),
	}
}

func GenerateGoogleOAuthURL(state string) (url string) {
	oauthConfig := googleOAuthConfig()

	return oauthConfig.AuthCodeURL(state)
}

func GenerateGoogleOAuthURLV2(state string) (url string) {
	oauthConfig := googleOAuthConfigV2()

	return oauthConfig.AuthCodeURL(state)
}

func GenerateLineOAuthURL(state string) (url string) {
	oauthConfig := PTGUoauthLine.NewLineOAuth(lineOAuthConfig())

	return oauthConfig.GenerateOAuthURL(PTGUoauthLine.OptionLineGenerateOAuthURL{
		Scopes: []string{
			"openid",
			"profile",
			"email",
		},
		State: state,
	})
}

func GetGoogleOAuthUserInfo(code string, version int8) (returnData modelController.ReturnGetOAuthUserInfo, err error) {
	var googleOauthConfig *oauth2.Config
	if version == 2 {
		googleOauthConfig = googleOAuthConfigV2()
	}else{
		googleOauthConfig = googleOAuthConfig()
	}
	googleOAuth := PTGUoauthGoogle.NewGoogleOAuth(googleOauthConfig)


	accessToken, err := googleOAuth.GetAccessToken(code)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][GetGoogleOAuthUserInfo()]->Fail to get access token")
	}

	_, validateStatus, err := googleOAuth.GetTokenInfo(accessToken)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][GetGoogleOAuthUserInfo()]->Fail to get token info")
	}
	if !validateStatus.Aud || !validateStatus.Exp {
		return returnData, errors.Wrap(errors.New("Token invalid"), "[Controller][GetGoogleOAuthUserInfo()]->Fail to validate token")
	}

	userInfo, err := googleOAuth.GetUserInfo(accessToken)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][GetGoogleOAuthUserInfo()]->Fail to get user info")
	}

	returnData.UserID = userInfo.UserID
	returnData.Email = userInfo.Email
	returnData.Name = userInfo.Name
	returnData.Picture = userInfo.Picture

	return returnData, nil
}

func GetLineOAuthUserInfo(code string) (returnData modelController.ReturnGetOAuthUserInfo, err error) {
	lineOAuth := PTGUoauthLine.NewLineOAuth(lineOAuthConfig())

	tokenData, err := lineOAuth.GetToken(code)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][GetLineOAuthUserInfo()]->Fail to get token")
	}

	idTokenInfo, err := lineOAuth.GetIDTokenInfo(tokenData.IDToken)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][GetLineOAuthUserInfo()]->Fail to get id token info")
	}

	returnData.UserID = idTokenInfo.Sub
	returnData.Email = PTGUdata.PointerToStringValue(idTokenInfo.Email)
	returnData.Name = PTGUdata.PointerToStringValue(idTokenInfo.Name)
	returnData.Picture = PTGUdata.PointerToStringValue(idTokenInfo.Picture)

	return returnData, nil
}