package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/parinyapt/prinflix_backend/handler"
	"github.com/parinyapt/prinflix_backend/middleware"
)

func InitMovieAPI(router *gin.RouterGroup) {
	r := router.Group("/movie").Use(middleware.GetHeaderAuthorizationToken, middleware.AuthWithAccessToken)
	{
		middlewareUser := middleware.NewMiddleware(middleware.MiddlewareReceiverArgument{})
		user := r.Use(middlewareUser.CheckAccount)
		user.GET("category", handler.GetMovieCategoryListHandler)
		user.GET("", handler.GetMovieListHandler)
		user.GET("/:movie_uuid", handler.GetMovieDetailHandler)
		// user.GET("favorite", handler.GetFavoriteMovieListHandler)
		// user.GET("watch", handler.GetWatchMovieListHandler)
		// user.POST("/:movie_uuid/favorite", handler.AddFavoriteMovieHandler)
		// user.DELETE("/:movie_uuid/favorite", handler.RemoveFavoriteMovieHandler)
		// user.POST("/:movie_uuid/watch", handler.RequestWatchMovieHandler)
	}
}