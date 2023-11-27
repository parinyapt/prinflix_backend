package modelRepository

import (
	"time"

	"github.com/google/uuid"
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
)

type ParamCreateWatchSession struct {
	SessionUUID uuid.UUID
	AccountUUID uuid.UUID
	MovieUUID   uuid.UUID
	ExpiredAt   time.Time
}

type ResultFetchOneWatchSession struct {
	IsFound bool
	Data    *modelDatabase.WatchSession
}
