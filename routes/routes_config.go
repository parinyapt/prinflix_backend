package routes

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"github.com/parinyapt/prinflix_backend/logger"
)

func configCors(router *gin.Engine) {
	config := cors.DefaultConfig()

	// Set Allow Origins
	config.AllowAllOrigins = true
	// config.AllowOrigins = []string{
	// 	"https://prinpt.com",
	// }

	// Set Allow Methods
	config.AllowMethods = []string{
		"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS",
	}

	// Set Allow Headers
	config.AllowHeaders = []string{
		"Origin", "Content-Length", "Content-Type", "Accept-Language", "Authorization",
	}

	// Set Allow Credentials
	config.AllowCredentials = true

	// Set Max Age
	config.MaxAge = 60 * time.Minute

	router.Use(cors.New(config))
}

func configRateLimit(router *gin.Engine) {
	limiter := rate.NewLimiter(1000, 1)
	router.Use(func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		c.Next()
	})
}

func configApi(router *gin.Engine, port string) *http.Server {
	// Set Max Memory for multipart forms to 5MB
	router.MaxMultipartMemory = 5 << 20

	s := &http.Server{
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if port == "" {
		logger.Fatal("Failed to run server because port not config")
	} else {
		s.Addr = ":" + port
		logger.Info("Running on PORT : " + port)
	}

	return s
}
