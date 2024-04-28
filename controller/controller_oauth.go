package controller

import (
	"os"

	PTGUdata "github.com/parinyapt/golang_utils/data/v1"
	PTGUoauthApple "github.com/parinyapt/golang_utils/oauth/apple/v1"
	PTGUoauthGoogle "github.com/parinyapt/golang_utils/oauth/google/v1"
	PTGUoauthLine "github.com/parinyapt/golang_utils/oauth/line/v1"
	modelController "github.com/parinyapt/prinflix_backend/model/controller"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// OAuth Config
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

func lineOAuthConfigV2() *PTGUoauthLine.LineOAuthConfig {
	return &PTGUoauthLine.LineOAuthConfig{
		ClientID:     os.Getenv("OAUTH2_LINE_CLIENT_ID_V2"),
		ClientSecret: os.Getenv("OAUTH2_LINE_CLIENT_SECRET_V2"),
		RedirectURL:  os.Getenv("OAUTH2_LINE_REDIRECT_URL_V2"),
	}
}

func appleOAuthConfig() *PTGUoauthApple.AppleOAuthConfig {
	privatekey, err := PTGUoauthApple.GetApplePrivateKeyFromFile(os.Getenv("OAUTH2_APPLE_PRIVATE_KEY_FILE_PATH"))
	if err != nil {
		privatekey = ""
	}

	return &PTGUoauthApple.AppleOAuthConfig{
		ClientID:    os.Getenv("OAUTH2_APPLE_CLIENT_ID"),
		RedirectURL: os.Getenv("OAUTH2_APPLE_REDIRECT_URL"),
		TeamID:      os.Getenv("OAUTH2_APPLE_TEAM_ID"),
		KeyID:       os.Getenv("OAUTH2_APPLE_KEY_ID"),
		PrivateKey:  privatekey,
	}
}

func appleOAuthConfigV2() *PTGUoauthApple.AppleOAuthConfig {
	privatekey, err := PTGUoauthApple.GetApplePrivateKeyFromFile(os.Getenv("OAUTH2_APPLE_PRIVATE_KEY_FILE_PATH"))
	if err != nil {
		privatekey = ""
	}

	return &PTGUoauthApple.AppleOAuthConfig{
		ClientID:    os.Getenv("OAUTH2_APPLE_CLIENT_ID"),
		RedirectURL: os.Getenv("OAUTH2_APPLE_REDIRECT_URL_V2"),
		TeamID:      os.Getenv("OAUTH2_APPLE_TEAM_ID"),
		KeyID:       os.Getenv("OAUTH2_APPLE_KEY_ID"),
		PrivateKey:  privatekey,
	}
}
// OAuth Config

// OAuth Generate URL
func GenerateGoogleOAuthURL(state string, version int8) (url string) {
	var oauthConfig *oauth2.Config
	if version == 2 {
		oauthConfig = googleOAuthConfigV2()
	} else {
		oauthConfig = googleOAuthConfig()
	}

	return oauthConfig.AuthCodeURL(state)
}

func GenerateLineOAuthURL(state string, version int8) (url string) {
	var lineOauthConfig *PTGUoauthLine.LineOAuthConfig
	if version == 2 {
		lineOauthConfig = lineOAuthConfigV2()
	} else {
		lineOauthConfig = lineOAuthConfig()
	}
	oauthConfig := PTGUoauthLine.NewLineOAuth(lineOauthConfig)

	return oauthConfig.GenerateOAuthURL(PTGUoauthLine.OptionLineGenerateOAuthURL{
		Scopes: []string{
			"openid",
			"profile",
			"email",
		},
		State: state,
	})
}

func GenerateAppleOAuthURL(state string, version int8) (url string) {
	var appleOauthConfig *PTGUoauthApple.AppleOAuthConfig
	if version == 2 {
		appleOauthConfig = appleOAuthConfigV2()
	} else {
		appleOauthConfig = appleOAuthConfig()
	}
	appleOAuth := PTGUoauthApple.NewAppleOAuth(appleOauthConfig)

	return appleOAuth.GenerateOAuthURL(PTGUoauthApple.OptionAppleGenerateOAuthURL{
		ResponseType: []string{"code"},
		ResponseMode: "form_post",
		Scope:        []string{"name", "email"},
		State:        state,
	})
}
// OAuth Generate URL

// OAuth Get User Info
func GetGoogleOAuthUserInfo(code string, version int8) (returnData modelController.ReturnGetOAuthUserInfo, err error) {
	var googleOauthConfig *oauth2.Config
	if version == 2 {
		googleOauthConfig = googleOAuthConfigV2()
	} else {
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

func GetLineOAuthUserInfo(code string, version int8) (returnData modelController.ReturnGetOAuthUserInfo, err error) {
	var lineOauthConfig *PTGUoauthLine.LineOAuthConfig
	if version == 2 {
		lineOauthConfig = lineOAuthConfigV2()
	} else {
		lineOauthConfig = lineOAuthConfig()
	}
	lineOAuth := PTGUoauthLine.NewLineOAuth(lineOauthConfig)

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

func GetAppleOAuthUserInfo(code string, version int8) (returnData modelController.ReturnGetOAuthUserInfo, err error) {
	var appleOauthConfig *PTGUoauthApple.AppleOAuthConfig
	if version == 2 {
		appleOauthConfig = appleOAuthConfigV2()
	} else {
		appleOauthConfig = appleOAuthConfig()
	}
	appleOAuth := PTGUoauthApple.NewAppleOAuth(appleOauthConfig)

	token, err := appleOAuth.ValidateAuthorizationCode(code, PTGUoauthApple.PlatformWeb)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][GetAppleOAuthUserInfo()]->Fail to get token")
	}

	userInfo, err := PTGUoauthApple.GetIDTokenInfo(token.IDToken)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][GetAppleOAuthUserInfo()]->Fail to get user info")
	}

	returnData.UserID = userInfo.Subject
	if userInfo.Email != nil {
		returnData.Email = *userInfo.Email
	}

	return returnData, nil
}
