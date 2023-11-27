package controller

import (
	modelController "github.com/parinyapt/prinflix_backend/model/controller"
	"github.com/parinyapt/prinflix_backend/repository"
	"github.com/pkg/errors"
)

func (receiver ControllerReceiverArgument) GetAllMovieCategory() (returnData modelController.ReturnGetAllMovieCategory, err error) {
	repoInstance := repository.NewRepository(receiver.databaseTX)

	repoData, repoErr := repoInstance.FetchManyMovieCategory()
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][GetAllMovieCategory()]->Fail to fetch many movie category")
	}
	if !repoData.IsFound {
		returnData.IsNotFound = true
		return returnData, nil
	}

	for _, data := range repoData.Data {
		returnData.Data = append(returnData.Data, modelController.ReturnGetAllMovieCategoryData{
			CategoryID:   data.ID,
			CategoryName: data.Name,
		})
	}

	return returnData, nil
}
