package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	PTGUvalidator "github.com/parinyapt/golang_utils/validator/v1"
	"github.com/parinyapt/prinflix_backend/controller"
	"github.com/parinyapt/prinflix_backend/database"
	"github.com/parinyapt/prinflix_backend/logger"
	modelHandler "github.com/parinyapt/prinflix_backend/model/handler"
	modelUtils "github.com/parinyapt/prinflix_backend/model/utils"
	utilsResponse "github.com/parinyapt/prinflix_backend/utils/response"
)

func GetMovieListHandler(c *gin.Context) {
	var queryParam modelHandler.QueryParamGetMovieList

	if err := c.ShouldBind(&queryParam); err != nil {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
		})
		return
	}

	isValidatePass, _, validatorError := PTGUvalidator.Validate(queryParam)
	if validatorError != nil {
		logger.Error("[Handler][GetTipsByMatchUUID()]->Error Validate Data", logger.Field("error", validatorError.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	if !isValidatePass {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
			Error:        "Invalid Parameter",
		})
		return
	}

	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusOK,
		Error:        queryParam,
	})
}

func GetMovieCategoryListHandler(c *gin.Context) {
	var response []modelHandler.ResponseGetMovieCategoryList

	controllerInstance := controller.NewController(database.DB)

	getAllMovieCategory, err := controllerInstance.GetAllMovieCategory()
	if err != nil {
		logger.Error("[Handler][GetMovieCategoryListHandler()]->Error GetAllMovieCategory()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	if getAllMovieCategory.IsNotFound {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusNotFound,
			Error:        "Movie Category Not Found",
		})
		return
	}

	for _, movieCategory := range getAllMovieCategory.Data {
		response = append(response, modelHandler.ResponseGetMovieCategoryList{
			CategoryID:   movieCategory.CategoryID,
			CategoryName: movieCategory.CategoryName,
		})
	}

	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusOK,
		Data:         response,
	})
}
