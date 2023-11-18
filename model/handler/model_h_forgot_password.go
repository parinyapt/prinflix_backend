package modelHandler

type RequestForgotPassword struct {
	Email    string `json:"email" validate:"required,email"`
}

type UriParamCheckForgotPasswordSession struct {
	SessionID string `uri:"session_id" validate:"required,base64url"`
}

type RequestResetPassword struct {
	SessionID string `json:"session_id" validate:"required,base64url"`
	Password  string `json:"password" validate:"required,min=8,max=100"`
}