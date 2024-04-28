package modelHandler

type ResponseRequestConnectOAuth struct {
	AuthURL string `json:"auth_url"`
}

type QueryParamOAuthCallback struct {
	Code  string `form:"code" validate:"required"`
	State string `form:"state" validate:"required,uuid"`
}

type RequestAppleCallback struct {
	Code  string `form:"code" validate:"required"`
	State string `form:"state" validate:"required,uuid"`
	User  string `form:"user"`
}

type RequestAppleCallbackUser struct {
	Name struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	} `json:"name"`
}