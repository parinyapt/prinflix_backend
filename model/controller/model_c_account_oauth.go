package modelController

import "github.com/google/uuid"

type ParamCreateAccountOAuth struct {
	AccountUUID string
	Provider    string
	UserID      string
	UserName    string
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
	Name        string
	Email       string
	Picture     string
}

type ParamDeleteAccountOAuth struct {
	AccountUUID string
	Provider    string
}
