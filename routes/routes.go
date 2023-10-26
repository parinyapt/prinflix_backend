package routes

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/parinyapt/prinflix_backend/logger"
)

func InitializeRoutes() {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	
	router := gin.Default()
	//config
	configCors(router)
	configRateLimit(router)
	s := configApi(router, os.Getenv("PORT"))

	//setup all api route
	configApiRoutes(router)

	// Gracefully shutdown the server
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("ListenAndServe Fail", logger.Field("error", err))
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	logger.Info("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 60 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown fail", logger.Field("error", err))
	}

	logger.Info("Server shutting down completely")
}
