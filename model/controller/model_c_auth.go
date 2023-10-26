package modelController

import "github.com/google/uuid"

type ParamCheckLogin struct {
	Email    string
	Password string
}

type ReturnCheckLogin struct {
	IsNotFound         bool
	IsPasswordNotMatch bool
	UUID               uuid.UUID
}
