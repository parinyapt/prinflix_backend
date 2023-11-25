package modelHandler

type UriParamGetMovieVideoFile struct {
	MovieUUID string `uri:"movie_uuid" validate:"required,uuid"`
	FilePath  string `uri:"file_path" validate:"required,max=500"`
}
