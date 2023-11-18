package modelController

import (
	"time"

	"github.com/google/uuid"
)

type ParamCheckLogin struct {
	Email    string
	Password string
}

type ReturnCheckLogin struct {
	IsNotFound         bool
	IsPasswordNotMatch bool
	AccountUUID        uuid.UUID
}

type ParamGenerateAccessToken struct {
	SessionUUID string
	ExpiredAt   time.Time
}

type ClaimAuthToken struct {
	SessionID string
}

type ReturnGenerateAccessToken struct {
	TokenType   string
	AccessToken string
}

type ReturnValidateAccessToken struct {
	IsExpired   bool
	SessionUUID string
}
