package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"challenge-backend-1/internal/config"
)

// @title Challenge Backend 1
// @version 1.0.0
// @description Challenge Backend 1

// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	validate := config.NewValidator(viperConfig)
	app := config.NewFiber(viperConfig)
	producer := config.NewKafkaProducer(viperConfig, log)

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		Config:   viperConfig,
		Producer: producer,
	})

	webPort := viperConfig.GetInt("web.port")

	go func() {
		if err := app.Listen(fmt.Sprintf(":%d", webPort)); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Close SQL DB
	sqlDB, err := db.DB()
	if err != nil {
		log.Errorf("Failed to get SQL DB: %v", err)
	} else {
		if err := sqlDB.Close(); err != nil {
			log.Errorf("Failed to close SQL DB: %v", err)
		} else {
			log.Info("SQL DB closed")
		}
	}

	log.Info("Server exited")
}
