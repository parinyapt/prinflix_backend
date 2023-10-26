package controller

import (
	"os"
	"time"

	modelController "github.com/parinyapt/prinflix_backend/model/controller"
	"github.com/parinyapt/prinflix_backend/repository"

	PTGUjwt "github.com/parinyapt/golang_utils/jwt/v1"
	PTGUpassword "github.com/parinyapt/golang_utils/password/v1"

	"github.com/pkg/errors"
)

func (receiver ControllerReceiverArgument) CheckLogin(param modelController.ParamCheckLogin) (returnData modelController.ReturnCheckLogin, err error) {
	repoInstance := repository.NewRepository(receiver.databaseTX)

	repoData, repoErr := repoInstance.FetchOneAccountByEmail(param.Email)
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][CreateAccount()]->Fail to fetch one account by email")
	}

	if !repoData.IsFound {
		returnData.IsNotFound = true
		return returnData, nil
	}

	isPasswordMatch := PTGUpassword.VerifyHashPassword(param.Password, repoData.Data.Password)
	if !isPasswordMatch {
		returnData.IsPasswordNotMatch = true
		return returnData, nil
	}

	returnData.AccountUUID = repoData.Data.UUID

	return returnData, nil
}

func GenerateAccessToken(param modelController.ParamGenerateAccessToken) (returnData modelController.ReturnGenerateAccessToken, err error) {
	jwtAccessTokenClaim := modelController.ClaimAuthToken{
		SessionID: param.SessionUUID,
	}
	accessToken, err := PTGUjwt.Sign(PTGUjwt.JwtSignConfig{
		SignKey:       os.Getenv("JWT_SIGN_KEY_ACCESS_TOKEN"),
		AppName:       os.Getenv("APP_NAME"),
		ExpireTime:    param.ExpiredAt,
		IssuedTime:    time.Now(),
		NotBeforeTime: time.Now(),
	}, jwtAccessTokenClaim)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][CreateAuthSession()]->Fail to generate access token")
	}

	returnData.TokenType = "Bearer"
	returnData.AccessToken = accessToken

	return returnData, nil
}
