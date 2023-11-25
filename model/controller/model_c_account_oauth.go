package modelController

import "github.com/google/uuid"

type ParamCreateAccountOAuth struct {
	AccountUUID string
	Provider    string
	UserID      string
	UserEmail   string
	UserPicture string
}

type ParamCheckAccountOAuth struct {
	UserID      string
	AccountUUID string
}

type ReturnCheckAccountOAuth struct {
	IsNotFound  bool
	AccountUUID uuid.UUID
}

type ParamDeleteAccountOAuth struct {
	AccountUUID string
	Provider    string
}
