package modelController

import "github.com/google/uuid"

type ParamCreateAccount struct {
	Name     string
	Email    string
	Password string
}

type ReturnCreateAccount struct {
	IsExist bool
	UUID uuid.UUID
}
