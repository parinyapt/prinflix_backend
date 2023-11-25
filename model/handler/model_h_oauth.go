package modelHandler

type ResponseRequestConnectOAuth struct {
	AuthURL string `json:"auth_url"`
}

type QueryParamOAuthCallback struct {
	Code  string `form:"code" validate:"required"`
	State string `form:"state" validate:"required,uuid"`
}
