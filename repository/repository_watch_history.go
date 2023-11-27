package repository

import (
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
	modelRepository "github.com/parinyapt/prinflix_backend/model/repository"
	"github.com/pkg/errors"
)

func (receiver RepositoryReceiverArgument) CreateWatchHistory(param modelRepository.ParamCreateWatchHistory) (err error) {
	resultDB := receiver.databaseTX.Create(&modelDatabase.WatchHistory{
		UUID:        param.WatchHistoryUUID,
		AccountUUID: param.AccountUUID,
		MovieUUID:   param.MovieUUID,
	})
	if resultDB.Error != nil {
		return errors.Wrap(resultDB.Error, "[Repository][CreateWatchHistory()]->"+errorDatabaseQueryFailed)
	}

	return nil
}

func (receiver RepositoryReceiverArgument) UpdateWatchHistoryLatestTimeStamp(param modelRepository.ParamUpdateWatchHistory) (result modelRepository.ResultIsFoundOnly, err error) {
	resultDB := receiver.databaseTX.Model(&modelDatabase.WatchHistory{}).Where(&modelDatabase.WatchHistory{AccountUUID: param.AccountUUID, MovieUUID: param.MovieUUID}).Updates(map[string]interface{}{
		"watch_history_latest_timestamp": param.LatestTimeStamp,
	})
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][UpdateWatchHistoryLatestTimeStamp()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}

func (receiver RepositoryReceiverArgument) UpdateWatchHistoryIsEnd(param modelRepository.ParamUpdateWatchHistory) (result modelRepository.ResultIsFoundOnly, err error) {
	resultDB := receiver.databaseTX.Model(&modelDatabase.WatchHistory{}).Where(&modelDatabase.WatchHistory{AccountUUID: param.AccountUUID, MovieUUID: param.MovieUUID}).Updates(map[string]interface{}{
		"watch_history_is_end": param.IsEnd,
	})
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][UpdateWatchHistoryIsEnd()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}

func (receiver RepositoryReceiverArgument) FetchOneWatchHistory(param modelRepository.ParamFetchOneWatchHistory) (result modelRepository.ResultFetchOneWatchHistory, err error) {
	resultDB := receiver.databaseTX.Where(&modelDatabase.WatchHistory{
		AccountUUID: param.AccountUUID,
		MovieUUID:   param.MovieUUID,
	}).Limit(1).Find(&result.Data)
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][FetchOneWatchHistory()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}
