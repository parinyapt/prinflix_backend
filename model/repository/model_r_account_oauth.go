package modelRepository

import (
	"github.com/google/uuid"
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
)

type ParamCreateAccountOAuth struct {
	AccountUUID uuid.UUID
	Provider    string
	UserID      string
	UserEmail   string
	UserPicture string
}

type ParamFetchOneAccountOAuthByProviderAndUserID struct {
	Provider string
	UserID   string
}

type ParamFetchOneAccountOAuthByProviderAndAccountUUID struct {
	Provider    string
	AccountUUID uuid.UUID
}

type ResultFetchOneAccountOAuth struct {
	IsFound bool
	Data    *modelDatabase.AccountOAuth
}

type ResultFetchManyAccountOAuth struct {
	IsFound bool
	Data []modelDatabase.AccountOAuth
}

type ParamDeleteAccountOAuthByProviderAndAccountUUID struct {
	Provider    string
	AccountUUID uuid.UUID
}
