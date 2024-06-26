package modelHandler

import "time"

type RequestLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RequestRegister struct {
	Name     string `json:"name" validate:"required,alpha_space,min=2,max=200"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type ResponseAccessToken struct {
	TokenType           string        `json:"token_type"`
	AccessToken         string        `json:"access_token"`
	AccessTokenExpireIn time.Duration `json:"expire_in"`
}

type ResponseVerifyToken struct {
	Name          string                   `json:"name"`
	Email         string                   `json:"email"`
	EmailVerified bool                     `json:"email_verified"`
	ImageStatus   bool                     `json:"have_image"`
	ImageURL      string                   `json:"image_url"`
	Status        string                   `json:"status"`
	Role          string                   `json:"role"`
	OAuth         ResponseVerifyTokenOAuth `json:"oauth"`
	SessionUUID   string                   `json:"session_id"`
}

type ResponseVerifyTokenOAuth struct {
	Google struct {
		IsConnected bool   `json:"connected"`
		Name        string `json:"name,omitempty"`
		Picture     string `json:"picture,omitempty"`
	} `json:"google"`
	Line struct {
		IsConnected bool   `json:"connected"`
		Name        string `json:"name,omitempty"`
		Picture     string `json:"picture,omitempty"`
	} `json:"line"`
	Apple struct {
		IsConnected bool   `json:"connected"`
		Name        string `json:"name,omitempty"`
		Picture     string `json:"picture,omitempty"`
	} `json:"apple"`
}

type RequestInternalOAuthLogin struct {
	UserID string `json:"user_id" validate:"required"`
}

type RequestExchangeCodeToAuthToken struct {
	Code string `json:"code" validate:"required,base64url"`
}