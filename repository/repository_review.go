package repository

import (
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
	modelRepository "github.com/parinyapt/prinflix_backend/model/repository"
	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

func (receiver RepositoryReceiverArgument) UpsertReview(param modelRepository.ParamUpsertReview) (result modelRepository.ResultIsFoundOnly, err error) {
	resultDB := receiver.databaseTX.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{
				Name: "AccountUUID",
			},
			{
				Name: "MovieUUID",
			},
		}, // key colume
		DoUpdates: clause.AssignmentColumns([]string{
			"review_rating",
			"review_updated_at",
		}), // column needed to be updated
	}).Create(&modelDatabase.Review{
		AccountUUID: param.AccountUUID,
		MovieUUID:   param.MovieUUID,
		Rating:      param.ReviewRating,
	})
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][UpsertReview()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}

func (receiver RepositoryReceiverArgument) DeleteReview(param modelRepository.ParamDeleteReview) (result modelRepository.ResultIsFoundOnly, err error) {
	resultDB := receiver.databaseTX.Where(&modelDatabase.Review{
		AccountUUID: param.AccountUUID,
		MovieUUID:   param.MovieUUID,
	}).Delete(&modelDatabase.Review{})
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][DeleteReview()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}
