package modelController

import (
	"time"

	"github.com/google/uuid"
)

type ReturnCreateAuthSession struct {
	SessionUUID uuid.UUID
	ExtiredIn   time.Duration
	ExpiredAt   time.Time
}
