package test

import (
	"challenge-backend-1/internal/config"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var App *fiber.App

var Db *gorm.DB

var ViperConfig *viper.Viper

var Log *zap.SugaredLogger

var Validate *validator.Validate

func init() {
	ViperConfig = config.NewViper()
	Log = config.NewLogger(ViperConfig)
	Validate = config.NewValidator(ViperConfig)
	App = config.NewFiber(ViperConfig)
	Db = config.NewDatabase(ViperConfig, Log)
	producer := config.NewKafkaProducer(ViperConfig, Log)

	config.Bootstrap(&config.BootstrapConfig{
		DB:       Db,
		App:      App,
		Log:      Log,
		Validate: Validate,
		Config:   ViperConfig,
		Producer: producer,
	})
}
