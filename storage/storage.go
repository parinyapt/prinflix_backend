package storage

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	PTGUcryptography "github.com/parinyapt/golang_utils/cryptography/aes/v1"
	"github.com/pkg/errors"
)

func InitializeStorage() {
	// Connect to Minio Storage
	initializeConnectMinio()
}

type ObjectPath string

const (
	// Minio Storage Path
	// ObjectNotFoundPath              ObjectPath = "global/object_not_found.jpg"
	// AccountProfileImageNotfoundPath ObjectPath = "global/profile_no_image.jpg"

	// AccountProfileImagePath ObjectPath = "account/:account_uuid/profile.jpg"
	// MovieThumbnailPath      ObjectPath = "movie/:movie_uuid/image/thumbnail.jpg"
	MovieVideoFilePath ObjectPath = "movie/:movie_uuid/video/:file_path"
)

type RoutePath string

const (
	// AccountProfileImageRoutePath RoutePath = "/profile/:account_uuid_encrypt"
	// MovieThumbnailRoutePath      RoutePath = "/movie/:movie_uuid/thumbnail"
	MovieVideoFileRoutePath RoutePath = "/movie/:movie_uuid/video/*file_path"
)

func GenerateObjectPath(path ObjectPath, param map[string]string) (returnPath string) {
	returnPath = string(path)
	for key, value := range param {
		returnPath = strings.Replace(returnPath, ":"+key, value, -1)
	}

	return returnPath
}

func GenerateRoutePath(version int, path RoutePath, param map[string]string) (returnPath string, err error) {
	returnPath = string(path)
	for key, value := range param {
		if strings.HasSuffix(key, "_encrypt") {
			valueEncrypt, err := PTGUcryptography.Encrypt(os.Getenv("ENCRYPTION_KEY_STORAGE"), value)
			if err != nil {
				return returnPath, errors.Wrap(err, "[Storage][GenerateRoutePath()]->Fail to encrypt value")
			}
			value = base64.URLEncoding.EncodeToString(valueEncrypt)
		}

		returnPath = strings.Replace(returnPath, ":"+key, value, -1)
	}

	returnPath = fmt.Sprintf("%s/storage/v%d%s", os.Getenv("APP_BASE_URL"), version, returnPath)

	return returnPath, nil
}
