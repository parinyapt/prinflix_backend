package modelController

import (
	"time"

	"github.com/google/uuid"
)

type ReturnCreateWatchSession struct {
	SessionUUID uuid.UUID
	ExpiredAt   time.Time
}

type ReturnCheckWatchSession struct {
	IsNotFound  bool
	IsExpired   bool
	AccountUUID uuid.UUID
	MovieUUID   uuid.UUID
}

type ParamGenerateWatchSessionToken struct {
	SessionUUID string
	ExpiredAt   time.Time
}

type ReturnGenerateWatchSessionToken struct {
	WatchSessionToken string
}

type ReturnValidateWatchSessionToken struct {
	IsExpired   bool
	SessionUUID string
}

type ClaimWatchSessionToken struct {
	SessionUUID string
}
