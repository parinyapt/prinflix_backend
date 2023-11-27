package controller

import (
	"os"
	"time"

	PTGUjwt "github.com/parinyapt/golang_utils/jwt/v1"
	modelController "github.com/parinyapt/prinflix_backend/model/controller"
	modelRepository "github.com/parinyapt/prinflix_backend/model/repository"
	"github.com/parinyapt/prinflix_backend/repository"
	utilsUUID "github.com/parinyapt/prinflix_backend/utils/uuid"
	"github.com/pkg/errors"
)

const (
	WatchSessionExpiredIn = time.Hour * 12
)

func (receiver ControllerReceiverArgument) CreateWatchSession(param modelController.ParamAccountUUIDandMovieUUID) (returnData modelController.ReturnCreateWatchSession, err error) {
	accountUUIDparse, err := utilsUUID.ParseUUIDfromString(param.AccountUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][CreateFavoriteMovie()]->Fail to parse account uuid")
	}

	movieUUIDparse, err := utilsUUID.ParseUUIDfromString(param.MovieUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][CreateFavoriteMovie()]->Fail to parse movie uuid")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)

	sessionUUID := utilsUUID.GenerateUUIDv4()
	repoErr := repoInstance.CreateWatchSession(modelRepository.ParamCreateWatchSession{
		SessionUUID: sessionUUID,
		AccountUUID: accountUUIDparse,
		MovieUUID:   movieUUIDparse,
		ExpiredAt:   time.Now().Add(WatchSessionExpiredIn),
	})
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][CreateWatchSession()]->Fail to create watch session")
	}

	returnData.SessionUUID = sessionUUID
	returnData.ExpiredAt = time.Now().Add(WatchSessionExpiredIn)

	return returnData, nil
}

func (receiver ControllerReceiverArgument) CheckWatchSession(sessionUUID string) (returnData modelController.ReturnCheckWatchSession, err error) {
	sessionUUIDparse, err := utilsUUID.ParseUUIDfromString(sessionUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][CheckWatchSession()]->Fail to parse session uuid")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)

	repoData, repoErr := repoInstance.FetchOneWatchSessionBySessionUUID(sessionUUIDparse)
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][CheckWatchSession()]->Fail to fetch one watch session by session uuid")
	}

	if !repoData.IsFound {
		returnData.IsNotFound = true
		return returnData, nil
	}

	if repoData.Data.ExpiredAt.Before(time.Now()) {
		returnData.IsExpired = true
		return returnData, nil
	}

	returnData.AccountUUID = repoData.Data.AccountUUID
	returnData.MovieUUID = repoData.Data.MovieUUID

	return returnData, nil
}

func (receiver ControllerReceiverArgument) DeleteAllWatchSessionByAccountUUID(accountUUID string) (err error) {
	accountUUIDparse, err := utilsUUID.ParseUUIDfromString(accountUUID)
	if err != nil {
		return errors.Wrap(err, "[Controller][DeleteAllWatchSessionByAccountUUID()]->Fail to parse account uuid")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)

	_, repoErr := repoInstance.DeleteWatchSessionByAccountUUID(accountUUIDparse)
	if repoErr != nil {
		return errors.Wrap(repoErr, "[Controller][DeleteAllWatchSessionByAccountUUID()]->Fail to delete watch session by account uuid")
	}

	return nil
}

func GenerateWatchSessionToken(param modelController.ParamGenerateWatchSessionToken) (returnData modelController.ReturnGenerateWatchSessionToken, err error) {
	jwtWatchSessionTokenClaim := modelController.ClaimWatchSessionToken{
		SessionUUID: param.SessionUUID,
	}
	watchSessionToken, err := PTGUjwt.Sign(PTGUjwt.JwtSignConfig{
		SignKey:       os.Getenv("JWT_SIGN_KEY_WATCH_SESSION_TOKEN"),
		AppName:       os.Getenv("APP_NAME"),
		ExpireTime:    param.ExpiredAt,
		IssuedTime:    time.Now(),
		NotBeforeTime: time.Now(),
	}, jwtWatchSessionTokenClaim)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][GenerateWatchSessionToken()]->Fail to generate watch session token")
	}

	returnData.WatchSessionToken = watchSessionToken

	return returnData, nil
}

func ValidateWatchSessionToken(watchSessionToken string) (returnData modelController.ReturnValidateWatchSessionToken, err error) {
	claims, isExpireOrNotValidYet, err := PTGUjwt.Validate(watchSessionToken, PTGUjwt.JwtValidateConfig{
		SignKey: os.Getenv("JWT_SIGN_KEY_WATCH_SESSION_TOKEN"),
	})
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][ValidateWatchSessionToken()]->Fail to validate watch session token")
	}

	if isExpireOrNotValidYet {
		returnData.IsExpired = true
		return returnData, nil
	}

	returnData.SessionUUID = claims.(map[string]interface{})["SessionUUID"].(string)

	return returnData, nil
}
