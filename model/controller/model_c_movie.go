package modelController

type ReturnGetAllMovieCategory struct {
	IsNotFound bool
	Data       []ReturnGetAllMovieCategoryData
}

type ReturnGetAllMovieCategoryData struct {
	CategoryID   uint   `json:"id"`
	CategoryName string `json:"name"`
}
