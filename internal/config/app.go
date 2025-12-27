package config

import (
	"challenge-backend-1/internal/delivery/http"
	"challenge-backend-1/internal/delivery/http/middleware"
	"challenge-backend-1/internal/delivery/http/route"
	"challenge-backend-1/internal/gateway/messaging"
	"challenge-backend-1/internal/repository"
	"challenge-backend-1/internal/security"
	"challenge-backend-1/internal/usecase"

	"github.com/IBM/sarama"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Log      *zap.SugaredLogger
	Validate *validator.Validate
	Config   *viper.Viper
	Producer sarama.SyncProducer
}

func Bootstrap(config *BootstrapConfig) {
	// setup security
	tokenProvider := security.NewJwtTokenProvider(config.Config.GetString("jwt.secret"))
	authMiddleware := middleware.NewAuthMiddleware(tokenProvider, config.Log)

	// setup repositories
	userRepository := repository.NewUserRepository(config.Log)

	// setup producer
	var userProducer *messaging.UserProducer

	if config.Producer != nil {
		userProducer = messaging.NewUserProducer(config.Producer, config.Log)
	}

	// setup use cases
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository, userProducer, tokenProvider)

	// setup controller
	userController := http.NewUserController(userUseCase, config.Log)

	routeConfig := route.RouteConfig{
		App:            config.App,
		Log:            config.Log,
		UserController: userController,
		AuthMiddleware: authMiddleware,
	}
	routeConfig.Setup()
}
