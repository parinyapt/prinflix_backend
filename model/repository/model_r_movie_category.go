package modelRepository

import (
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
)

type ResultFetchOneMovieCategory struct {
	IsFound bool
	Data    *modelDatabase.MovieCategory
}

type ResultFetchManyMovieCategory struct {
	IsFound bool
	Data    []modelDatabase.MovieCategory
}