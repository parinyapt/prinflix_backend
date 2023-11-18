package modelHandler

type UriParamEmailVerifyHandler struct {
	Code string `uri:"code" validate:"required,base64url"`
}