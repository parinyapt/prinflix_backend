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
