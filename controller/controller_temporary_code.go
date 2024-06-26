package controller

import (
	"encoding/base64"
	"os"
	"time"

	PTGUcryptography "github.com/parinyapt/golang_utils/cryptography/aes/v1"
	modelController "github.com/parinyapt/prinflix_backend/model/controller"
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
	modelRepository "github.com/parinyapt/prinflix_backend/model/repository"
	"github.com/parinyapt/prinflix_backend/repository"
	utilsUUID "github.com/parinyapt/prinflix_backend/utils/uuid"

	"github.com/pkg/errors"
)

const (
	TemporaryCodeTypeEmailVerificationExpiredIn = time.Minute * 60 * 24
	TemporaryCodeTypePasswordResetExpiredIn     = time.Minute * 15
	TemporaryCodeTypeOAuthStateExpiredIn        = time.Minute * 5
	TemporaryCodeTypeAuthTokenCodeExpiredIn     = time.Minute * 5
)

func (receiver ControllerReceiverArgument) CreateTemporaryCode(param modelController.ParamTemporaryCode) (returnData modelController.ReturnCreateTemporaryCode, err error) {
	accountUUIDparse, err := utilsUUID.ParseUUIDfromString(param.AccountUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][CreateTemporaryCode()]->Fail to parse account uuid")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)

	codeUUID := utilsUUID.GenerateUUIDv4()
	repoErr := repoInstance.CreateTemporaryCode(modelRepository.ParamCreateTemporaryCode{
		UUID:        codeUUID,
		AccountUUID: accountUUIDparse,
		Type:        param.Type,
	})
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][CreateTemporaryCode()]->Fail to create temporary code")
	}

	returnData.CodeUUID = codeUUID

	return returnData, nil
}

func (receiver ControllerReceiverArgument) CheckTemporaryCode(param modelController.ParamCheckTemporaryCode) (returnData modelController.ReturnCheckTemporaryCode, err error) {
	codeUUIDparse, err := utilsUUID.ParseUUIDfromString(param.CodeUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][CheckTemporaryCode()]->Fail to parse code uuid")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)

	repoData, repoErr := repoInstance.FetchOneTemporaryCodeByUUID(codeUUIDparse)
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][CheckTemporaryCode()]->Fail to fetch one temporary code by uuid")
	}

	if !repoData.IsFound {
		returnData.IsNotFound = true
		return returnData, nil
	}

	if repoData.Data.Type != param.Type {
		returnData.IsNotFound = true
		return returnData, nil
	}

	if repoData.Data.Type == modelDatabase.TemporaryCodeTypeEmailVerification {
		if time.Now().After(repoData.Data.CreatedAt.Add(TemporaryCodeTypeEmailVerificationExpiredIn)) {
			returnData.IsExpired = true
			return returnData, nil
		}
	}
	if repoData.Data.Type == modelDatabase.TemporaryCodeTypePasswordReset {
		if time.Now().After(repoData.Data.CreatedAt.Add(TemporaryCodeTypePasswordResetExpiredIn)) {
			returnData.IsExpired = true
			return returnData, nil
		}
	}
	if repoData.Data.Type == modelDatabase.TemporaryCodeTypeOAuthStateLine || repoData.Data.Type == modelDatabase.TemporaryCodeTypeOAuthStateGoogle || repoData.Data.Type == modelDatabase.TemporaryCodeTypeOAuthStateApple {
		if time.Now().After(repoData.Data.CreatedAt.Add(TemporaryCodeTypeOAuthStateExpiredIn)) {
			returnData.IsExpired = true
			return returnData, nil
		}
	}
	if repoData.Data.Type == modelDatabase.TemporaryCodeTypeAuthTokenCode {
		if time.Now().After(repoData.Data.CreatedAt.Add(TemporaryCodeTypeAuthTokenCodeExpiredIn)) {
			returnData.IsExpired = true
			return returnData, nil
		}
	}

	returnData.CodeUUID = repoData.Data.UUID
	returnData.AccountUUID = repoData.Data.AccountUUID

	return returnData, nil
}

func (receiver ControllerReceiverArgument) DeleteTemporaryCode(param modelController.ParamTemporaryCode) (err error) {
	accountUUIDparse, err := utilsUUID.ParseUUIDfromString(param.AccountUUID)
	if err != nil {
		return errors.Wrap(err, "[Controller][DeleteTemporaryCode()]->Fail to parse account uuid")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)

	_, repoErr := repoInstance.DeleteTemporaryCodeByAccountUUIDAndType(modelRepository.ParamDeleteTemporaryCodeByAccountUUIDAndType{
		AccountUUID: accountUUIDparse,
		Type:        param.Type,
	})
	if repoErr != nil {
		return errors.Wrap(repoErr, "[Controller][DeleteTemporaryCode()]->Fail to delete temporary code")
	}

	return nil
}

func EncryptTemporaryCode(codeUUID string) (codeUUIDEncryptBase64 string, err error) {
	codeUUIDEncrypt, err := PTGUcryptography.Encrypt(os.Getenv("ENCRYPTION_KEY_TEMP_CODE"), codeUUID)
	if err != nil {
		return codeUUIDEncryptBase64, errors.Wrap(err, "[Controller][EncryptTemporaryCode()]->Fail to encrypt code uuid")
	}
	codeUUIDEncryptBase64 = base64.URLEncoding.EncodeToString(codeUUIDEncrypt)

	return codeUUIDEncryptBase64, nil
}

func DecryptTemporaryCode(codeUUIDEncryptBase64 string) (codeUUID string, err error) {
	codeUUIDEncrypt, err := base64.URLEncoding.DecodeString(codeUUIDEncryptBase64)
	if err != nil {
		return codeUUID, errors.Wrap(err, "[Controller][DecryptTemporaryCode()]->Fail to decode base64")
	}

	codeUUID, err = PTGUcryptography.Decrypt(os.Getenv("ENCRYPTION_KEY_TEMP_CODE"), codeUUIDEncrypt)
	if err != nil {
		return codeUUID, errors.Wrap(err, "[Controller][DecryptTemporaryCode()]->Fail to decrypt code uuid")
	}

	return codeUUID, nil
}
