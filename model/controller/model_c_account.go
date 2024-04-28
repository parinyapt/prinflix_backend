package modelController

import "github.com/google/uuid"

type ParamCreateAccount struct {
	Name               string
	Email              string
	Password           string
	EmailVerifyApprove bool
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
	PasswordHash  string
}

type ParamUpdateAccount struct {
	Name          string
	EmailVerified bool
	Password      string
	Image         bool
}
