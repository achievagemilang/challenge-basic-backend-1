package usecase

import (
	"context"
	"time"

	"challenge-backend-1/internal/entity"
	"challenge-backend-1/internal/gateway/messaging"
	"challenge-backend-1/internal/model"
	"challenge-backend-1/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	DB             *gorm.DB
	Log            *zap.SugaredLogger
	Validate       *validator.Validate
	UserRepository *repository.UserRepository
	UserProducer   *messaging.UserProducer
	Config         *viper.Viper
}

func NewUserUseCase(db *gorm.DB, logger *zap.SugaredLogger, validate *validator.Validate,
	userRepository *repository.UserRepository, userProducer *messaging.UserProducer, config *viper.Viper,
) *UserUseCase {
	return &UserUseCase{
		DB:             db,
		Log:            logger,
		Validate:       validate,
		UserRepository: userRepository,
		UserProducer:   userProducer,
		Config:         config,
	}
}

func (c *UserUseCase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.LoginResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body  : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := c.UserRepository.FindByEmail(tx, user, request.Email); err != nil {
		c.Log.Warnf("Failed find user by email : %+v", err)
		return nil, fiber.NewError(fiber.StatusUnauthorized, "incorrect username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.Log.Warnf("Failed to compare user password with bcrype hash : %+v", err)
		return nil, fiber.NewError(fiber.StatusUnauthorized, "incorrect username or password")
	}

	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Second * 20).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(c.Config.GetString("jwt.secret")))
	if err != nil {
		c.Log.Warnf("Failed to generate access token : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	refreshToken := uuid.New().String()

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return &model.LoginResponse{
		User: model.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
