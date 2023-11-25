package modelHandler

type RequestUpdateProfile struct {
	Name string `json:"name" validate:"required,alpha_space,min=2,max=200"`
}

type RequestUpdatePassword struct {
	CurrentPassword string `json:"current_password" validate:"required,min=8,max=100"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=100"`
}
