package modelController

import "github.com/google/uuid"

type ParamCreateAccount struct {
	Name     string
	Email    string
	Password string
}

type ReturnCreateAccount struct {
	IsExist bool
	UUID    uuid.UUID
}

type ParamGetAccountInfo struct {
	AccountUUID string
	Email       string
}

type ReturnGetAccountInfo struct {
	IsNotFound bool

	AccountUUID   uuid.UUID
	Name          string
	Email         string
	EmailVerified bool
	Status        string
	Image         bool
	Role          string
}

type ParamUpdateAccount struct {
	AccountUUID   string
	Name          string
	EmailVerified bool
	Password      string
	Image         bool
}
