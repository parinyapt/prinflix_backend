package handler

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	PTGUvalidator "github.com/parinyapt/golang_utils/validator/v1"
	"github.com/parinyapt/prinflix_backend/controller"
	"github.com/parinyapt/prinflix_backend/logger"
	modelController "github.com/parinyapt/prinflix_backend/model/controller"
	modelHandler "github.com/parinyapt/prinflix_backend/model/handler"
)

func GetMovieVideoFileHandler(c *gin.Context) {
	var uriParam modelHandler.UriParamGetMovieVideoFile

	if err := c.ShouldBindUri(&uriParam); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	isValidatePass, _, validatorError := PTGUvalidator.Validate(uriParam)
	if validatorError != nil {
		logger.Error("[Handler][GetMovieVideoFileHandler()]->Error Validate Data", logger.Field("error", validatorError.Error()))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !isValidatePass {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if uriParam.MovieUUID != c.GetString("WATCHSESSION_MOVIE_UUID") {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	movieVideo, err := controller.GetMovieVideoFile(modelController.ParamGetMovieVideoFile{
		MovieUUID: uriParam.MovieUUID,
		FilePath:  uriParam.FilePath,
	})
	if err != nil {
		logger.Error("[Handler][GetMovieVideoFileHandler()]->Error GetMovieVideoFile()", logger.Field("error", err.Error()))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Header("Content-Type", "application/x-mpegURL")
	c.Header("Content-Length", strconv.FormatInt(movieVideo.Stat.Size, 10))

	if _, err := io.Copy(c.Writer, movieVideo.Object); err != nil {
		logger.Error("[Handler][GetMovieVideoFileHandler()]->Error Copy Object", logger.Field("error", err.Error()))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
}

func GetTestMovieVideoFileHandler(c *gin.Context) {
	var uriParam modelHandler.UriParamGetMovieVideoFile

	if err := c.ShouldBindUri(&uriParam); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	isValidatePass, _, validatorError := PTGUvalidator.Validate(uriParam)
	if validatorError != nil {
		logger.Error("[Handler][GetMovieVideoFileHandler()]->Error Validate Data", logger.Field("error", validatorError.Error()))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !isValidatePass {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	movieVideo, err := controller.GetMovieVideoFile(modelController.ParamGetMovieVideoFile{
		MovieUUID: uriParam.MovieUUID,
		FilePath:  uriParam.FilePath,
	})
	if err != nil {
		logger.Error("[Handler][GetMovieVideoFileHandler()]->Error GetMovieVideoFile()", logger.Field("error", err.Error()))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Header("Content-Type", "application/x-mpegURL")
	c.Header("Content-Length", strconv.FormatInt(movieVideo.Stat.Size, 10))

	if _, err := io.Copy(c.Writer, movieVideo.Object); err != nil {
		logger.Error("[Handler][GetMovieVideoFileHandler()]->Error Copy Object", logger.Field("error", err.Error()))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
}

func GetMovieThumbnailHandler(c *gin.Context) {
	var uriParam modelHandler.UriParamGetMovieThumbnail

	if err := c.ShouldBindUri(&uriParam); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	isValidatePass, _, validatorError := PTGUvalidator.Validate(uriParam)
	if validatorError != nil {
		logger.Error("[Handler][GetMovieThumbnailHandler()]->Error Validate Data", logger.Field("error", validatorError.Error()))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !isValidatePass {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	movieThumbnail, err := controller.GetMovieThumbnail(modelController.ParamGetMovieThumbnail{
		MovieUUID: uriParam.MovieUUID,
	})
	if err != nil {
		movieThumbnail, err = controller.GetMovieThumbnailNotFound()
		if err != nil {
			logger.Error("[Handler][GetMovieThumbnailHandler()]->Error GetMovieThumbnail() and GetMovieThumbnailNotFound()", logger.Field("error", err.Error()))
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
	}

	c.Header("Content-Type", "image/jpeg")
	c.Header("Content-Length", strconv.FormatInt(movieThumbnail.Stat.Size, 10))

	if _, err := io.Copy(c.Writer, movieThumbnail.Object); err != nil {
		logger.Error("[Handler][GetMovieThumbnailHandler()]->Error Copy Object", logger.Field("error", err.Error()))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
}
