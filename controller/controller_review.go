package controller

import (
	modelController "github.com/parinyapt/prinflix_backend/model/controller"
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
	modelRepository "github.com/parinyapt/prinflix_backend/model/repository"
	"github.com/parinyapt/prinflix_backend/repository"
	utilsUUID "github.com/parinyapt/prinflix_backend/utils/uuid"
	"github.com/pkg/errors"
)

func (receiver ControllerReceiverArgument) CreateUpdateReviewMovie(param modelController.ParamCreateUpdateReviewMovie) (err error) {
	accountUUIDparse, err := utilsUUID.ParseUUIDfromString(param.AccountUUID)
	if err != nil {
		return errors.Wrap(err, "[Controller][CreateUpdateReviewMovie()]->Fail to parse account uuid")
	}
	movieUUIDparse, err := utilsUUID.ParseUUIDfromString(param.MovieUUID)
	if err != nil {
		return errors.Wrap(err, "[Controller][CreateUpdateReviewMovie()]->Fail to parse movie uuid")
	}

	if param.ReviewRating < modelDatabase.ReviewRatingBad || param.ReviewRating > modelDatabase.ReviewRatingGood {
		return errors.New("[Controller][CreateUpdateReviewMovie()]->Review rating must be between 1 and 3")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)

	_, repoErr := repoInstance.UpsertReview(modelRepository.ParamUpsertReview{
		AccountUUID:   accountUUIDparse,
		MovieUUID:     movieUUIDparse,
		ReviewRating:  param.ReviewRating,
	})
	if repoErr != nil {
		return errors.Wrap(repoErr, "[Controller][CreateUpdateReviewMovie()]->Fail to upsert review")
	}

	return nil
}

func (receiver ControllerReceiverArgument) DeleteReviewMovie(param modelController.ParamDeleteReviewMovie) (err error) {
	accountUUIDparse, err := utilsUUID.ParseUUIDfromString(param.AccountUUID)
	if err != nil {
		return errors.Wrap(err, "[Controller][DeleteReviewMovie()]->Fail to parse account uuid")
	}
	movieUUIDparse, err := utilsUUID.ParseUUIDfromString(param.MovieUUID)
	if err != nil {
		return errors.Wrap(err, "[Controller][DeleteReviewMovie()]->Fail to parse movie uuid")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)

	_, repoErr := repoInstance.DeleteReview(modelRepository.ParamDeleteReview{
		AccountUUID: accountUUIDparse,
		MovieUUID:   movieUUIDparse,
	})
	if repoErr != nil {
		return errors.Wrap(repoErr, "[Controller][DeleteReviewMovie()]->Fail to delete review by account uuid and movie uuid")
	}

	return nil
}
