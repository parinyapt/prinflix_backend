package controller

import (
	"context"
	"os"
	"strings"

	"github.com/minio/minio-go/v7"
	modelController "github.com/parinyapt/prinflix_backend/model/controller"
	"github.com/parinyapt/prinflix_backend/storage"
	"github.com/pkg/errors"
)

func GetMovieVideoFile(param modelController.ParamGetMovieVideoFile) (returnData modelController.ReturnObjectDetail, err error) {
	object, err := storage.MinioClient.GetObject(context.Background(), os.Getenv("OBJECT_STORAGE_BUCKET_NAME"),
		storage.GenerateObjectPath(storage.MovieVideoFilePath, map[string]string{
			"movie_uuid": param.MovieUUID,
			"file_path":  strings.TrimLeft(param.FilePath,"/"),
		}),
		minio.GetObjectOptions{},
	)
	if err != nil {
		return returnData, errors.Wrap(err, "[Storage][GetMovieVideoFile()]->Fail to get object from minio")
	}

	returnData.Object = object
	returnData.Stat, err = object.Stat()
	if err != nil {
		return returnData, errors.Wrap(err, "[Storage][GetMovieVideoFile()]->Fail to get object stat")
	}

	return returnData, nil
}