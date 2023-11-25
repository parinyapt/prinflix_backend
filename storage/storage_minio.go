package storage

import (
	"context"
	"os"
	"strconv"

	"github.com/parinyapt/prinflix_backend/logger"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func initializeConnectMinio(){
	connectSecure, err := strconv.ParseBool(os.Getenv("OBJECT_STORAGE_CONNECT_SSL"))
	if err != nil {
		logger.Fatal("Failed to parse OBJECT_STORAGE_CONNECT_SSL", logger.Field("error", err))
	}
	// Initialize minio client object.
	minioClient, err := minio.New(os.Getenv("OBJECT_STORAGE_ENDPOINT"), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("OBJECT_STORAGE_ACCESS_KEY"), os.Getenv("OBJECT_STORAGE_SECRET_ACCESS_KEY"), ""),
		Secure: connectSecure,
	})
	if err != nil {
		logger.Fatal("Failed to connect minio object storage", logger.Field("error", err))
	}

	// Check if bucket already exists
	exists, errBucketExists := minioClient.BucketExists(context.Background(), os.Getenv("OBJECT_STORAGE_BUCKET_NAME"))
	if errBucketExists != nil {
		logger.Fatal("Failed to check bucket exist", logger.Field("error", errBucketExists))
	} 
	if !exists {
		logger.Fatal("Bucket not exist", logger.Field("bucket", os.Getenv("OBJECT_STORAGE_BUCKET_NAME")))
	}

	MinioClient = minioClient

	logger.Info("Initialize Minio Object Storage Success")
}