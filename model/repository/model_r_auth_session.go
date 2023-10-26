package modelRepository

import (
	"time"

	"github.com/google/uuid"
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
)

type ParamCreateAuthSession struct {
	SessionUUID uuid.UUID
	AccountUUID uuid.UUID
	ExpiredAt   time.Time
}

type ResultFetchOneAuthSession struct {
	IsFound bool
	Data    *modelDatabase.AuthSession
}