package repository

import (
	"github.com/google/uuid"
	PTGUdata "github.com/parinyapt/golang_utils/data/v1"
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
	modelRepository "github.com/parinyapt/prinflix_backend/model/repository"
	utilsDatabase "github.com/parinyapt/prinflix_backend/utils/database"
	"github.com/pkg/errors"
)

func (receiver RepositoryReceiverArgument) CreateAccount(param modelRepository.ParamCreateAccount) (err error) {
	resultDB := receiver.databaseTX.Create(&modelDatabase.Account{
		UUID:          param.UUID,
		Name:          param.Name,
		Email:         param.Email,
		EmailVerified: param.EmailVerified,
		Password:      param.Password,
		Status:        param.Status,
		Role:          param.Role,
	})
	if resultDB.Error != nil {
		return errors.Wrap(resultDB.Error, "[Repository][CreateAccount()]->"+errorDatabaseQueryFailed)
	}

	return nil
}

func (receiver RepositoryReceiverArgument) FetchOneAccountByEmail(email string) (result modelRepository.ResultFetchOneAccount, err error) {
	resultDB := receiver.databaseTX.Where(&modelDatabase.Account{Email: email}).Limit(1).Find(&result.Data)
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][FetchOneAccountByEmail()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}

func (receiver RepositoryReceiverArgument) FetchOneAccountByUUID(accountUUID uuid.UUID) (result modelRepository.ResultFetchOneAccount, err error) {
	resultDB := receiver.databaseTX.Where(&modelDatabase.Account{UUID: accountUUID}).Limit(1).Find(&result.Data)
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][FetchOneAccountByUUID()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}

func (receiver RepositoryReceiverArgument) UpdateAccountByUUID(accountUUID uuid.UUID, param modelRepository.ParamUpdateAccount) (result modelRepository.ResultIsFoundOnly, err error) {
	resultDB := receiver.databaseTX.Model(&modelDatabase.Account{UUID: accountUUID}).Updates(&modelDatabase.Account{
		Name:          param.Name,
		Email:         param.Email,
		EmailVerified: param.EmailVerified,
		Password:      param.Password,
		Status:        param.Status,
		Image:         param.Image,
		Role:          param.Role,
	})
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][UpdateAccountByUUID()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}

func (receiver RepositoryReceiverArgument) DeleteAccountByUUID(accountUUID uuid.UUID) (result modelRepository.ResultIsFoundOnly, err error) {
	resultDB := receiver.databaseTX.Where(&modelDatabase.Account{UUID: accountUUID}).Delete(&modelDatabase.Account{})
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][DeleteAccountByUUID()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}

func (receiver RepositoryReceiverArgument) FetchOneAccountMostViewCategory(accountUUID uuid.UUID) (result modelRepository.ResultFetchOneAccountMostViewCategory, err error) {
	sqlCommand := `WITH temp_data AS (
			SELECT 
				@table_watch_history.watch_history_movie_uuid AS movie_uuid,
				@table_watch_history.watch_history_account_uuid AS account_uuid,
				@table_watch_history.watch_history_updated_at AS updated_at,
				ROW_NUMBER() OVER (PARTITION BY @table_watch_history.watch_history_account_uuid ORDER BY @table_watch_history.watch_history_updated_at DESC) AS rownumber
			FROM @table_watch_history
			ORDER BY rownumber ASC
		)
		SELECT
			@table_movie.movie_movie_category_id AS category_id,
			COUNT(@table_movie.movie_movie_category_id) AS category_count,
			temp_data.updated_at
		FROM temp_data
		INNER JOIN @table_movie ON temp_data.movie_uuid = @table_movie.movie_uuid
		WHERE temp_data.account_uuid = @account_uuid
		GROUP BY @table_movie.movie_movie_category_id
		ORDER BY category_count DESC, temp_data.updated_at DESC
		LIMIT 1;`
	sqlCommand = PTGUdata.ReplaceString(sqlCommand, map[string]string{
		"@table_watch_history": utilsDatabase.GenerateTableName("watch_history"),
		"@table_movie":         utilsDatabase.GenerateTableName("movie"),
	})
	resultDB := receiver.databaseTX.Raw(sqlCommand, map[string]interface{}{
		"account_uuid": accountUUID,
	}).Scan(&result.Data)
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][FetchOneAccountMostViewCategory()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}
