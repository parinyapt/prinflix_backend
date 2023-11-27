package modelController

import "github.com/minio/minio-go/v7"

type ReturnObjectDetail struct {
	Object *minio.Object
	Stat   minio.ObjectInfo
}

type ParamGetMovieVideoFile struct {
	MovieUUID string
	FilePath  string
}

type ParamGetMovieThumbnail struct {
	MovieUUID string
}