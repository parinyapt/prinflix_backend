package modelHandler

type UriParamEmailVerify struct {
	Code string `uri:"code" validate:"required,base64url"`
}