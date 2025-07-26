// @title Banner API
// @version 1.0
// @description Simple banner click tracking API
// @host localhost:8080
// @BasePath /
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/acom21/banner-api/internal/handlers"
	"github.com/acom21/banner-api/internal/repository"
	"github.com/acom21/banner-api/internal/service"
	"github.com/acom21/banner-api/pkg/config"
	"github.com/acom21/banner-api/pkg/database"
	"github.com/acom21/banner-api/pkg/http"
	"github.com/acom21/banner-api/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	if err := logger.Init(); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.Log.Sync()

	cfg := config.Load()

	db, err := database.New(cfg.DB)
	if err != nil {
		logger.Log.Fatal("Failed to connect to database", zap.Error(err))
	}

	repo := repository.New(db)
	bannerService := service.New(repo)
	handler := handlers.New(bannerService)

	router := http.New(handler)
	server := http.NewServer(router)

	go func() {
		logger.Log.Info("Starting server", zap.String("address", cfg.Address()))
		if err := server.Start(cfg.Address()); err != nil {
			logger.Log.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Log.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Log.Info("Closing database connection...")
	db.Close()

	logger.Log.Info("Server exited")
}
