package modelRepository

import (
	"github.com/google/uuid"
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
)

type ParamFavoriteMovie struct {
	AccountUUID uuid.UUID
	MovieUUID   uuid.UUID
}

type ResultFetchOneFavoriteMovie struct {
	IsFound bool
	Data    *modelDatabase.FavoriteMovie
}

type ResultFetchManyFavoriteMovie struct {
	IsFound bool
	Data    []DBResultFetchMovie
}
