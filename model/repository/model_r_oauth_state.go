package modelRepository

import (
	"github.com/google/uuid"
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
)

type ParamCreateOauthState struct {
	UUID     uuid.UUID
	Provider string
}

type ResultFetchOneOauthState struct {
	IsFound bool
	Data    *modelDatabase.OauthState
}
