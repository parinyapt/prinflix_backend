package modelRepository

import (
	"github.com/google/uuid"
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
)

type ParamCreateAccount struct {
	UUID          uuid.UUID
	Name          string
	Email         string
	EmailVerified bool
	Password      string
	Status        string
	Role          string
}

type ResultFetchOneAccount struct {
	IsFound bool
	Data    *modelDatabase.Account
}

type ParamUpdateAccount struct {
	Name          string
	Email         string
	EmailVerified bool
	Password      string
	Status        string
	Image         bool
	Role          string
}
