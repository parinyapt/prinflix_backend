package modelHandler

type RequestUpdateProfileHandler struct {
	Name string `json:"name" validate:"required"`
}

type RequestUpdatePasswordHandler struct {
	CurrentPassword string `json:"current_password" validate:"required,min=8,max=100"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=100"`
}
