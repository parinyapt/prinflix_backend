package modelRepository

import (
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
)

type ResultFetchOneMovie struct {
	IsFound bool
	Data    *modelDatabase.Movie
}
