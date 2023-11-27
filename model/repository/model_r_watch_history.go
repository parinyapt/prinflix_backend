package modelRepository

import (
	"github.com/google/uuid"
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
)

type ParamCreateWatchHistory struct {
	WatchHistoryUUID uuid.UUID
	AccountUUID      uuid.UUID
	MovieUUID        uuid.UUID
}

type ParamUpdateWatchHistory struct {
	AccountUUID      uuid.UUID
	MovieUUID        uuid.UUID
	LatestTimeStamp int64
	IsEnd           bool
}

type ParamFetchOneWatchHistory struct {
	AccountUUID uuid.UUID
	MovieUUID   uuid.UUID
}

type ResultFetchOneWatchHistory struct {
	IsFound bool
	Data    *modelDatabase.WatchHistory
}