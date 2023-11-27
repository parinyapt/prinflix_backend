package modelController

type ReturnGetAllFavoriteMovie struct {
	IsNotFound bool
	Data       []ReturnGetManyMovieData
}