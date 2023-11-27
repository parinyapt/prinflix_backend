package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/parinyapt/prinflix_backend/handler"
	"github.com/parinyapt/prinflix_backend/storage"
)

func InitStorageAPI(router *gin.RouterGroup) {
	// watchSessionProtect := router.Group("").Use(middleware.WatchSessionCheck)
	// {
	// 	watchSessionProtect.GET(string(storage.MovieVideoFileRoutePath), handler.GetMovieVideoFileHandler)
	// }
	router.GET(string(storage.MovieVideoFileRoutePath), handler.GetMovieVideoFileHandler)
	router.GET(string(storage.MovieThumbnailRoutePath), handler.GetMovieThumbnailHandler)
}
