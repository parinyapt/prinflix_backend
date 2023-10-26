package config

import "github.com/gin-gonic/gin"

func initializeSetGinReleaseMode() {
	gin.SetMode(gin.ReleaseMode)
}