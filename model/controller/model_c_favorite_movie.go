package modelController

type ParamFavoriteMovie struct {
	AccountUUID string
	MovieUUID   string
}

type ReturnGetAllFavoriteMovie struct {
	IsNotFound bool
	Data       []ReturnGetManyMovieData
}