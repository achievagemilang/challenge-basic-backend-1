package usecase

import (
	"context"

	"challenge-backend-1/internal/entity"
	"challenge-backend-1/internal/model"
	"challenge-backend-1/internal/security"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	DB             *gorm.DB
	Log            *zap.SugaredLogger
	Validate       *validator.Validate
	UserRepository UserRepository
	UserProducer   UserProducer
	TokenProvider  security.TokenProvider
}

func NewUserUseCase(db *gorm.DB, logger *zap.SugaredLogger, validate *validator.Validate,
	userRepository UserRepository, userProducer UserProducer,
	tokenProvider security.TokenProvider,
) *UserUseCase {
	return &UserUseCase{
		DB:             db,
		Log:            logger,
		Validate:       validate,
		UserRepository: userRepository,
		UserProducer:   userProducer,
		TokenProvider:  tokenProvider,
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

	accessToken, err := c.TokenProvider.GenerateAccessToken(user)
	if err != nil {
		c.Log.Warnf("Failed to generate access token : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	refreshToken, err := c.TokenProvider.GenerateRefreshToken(user)
	if err != nil {
		c.Log.Warnf("Failed to generate refresh token : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return &model.LoginResponse{
		User: &model.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (c *UserUseCase) Refresh(ctx context.Context, refreshToken string) (*model.LoginResponse, error) {
	claims, err := c.TokenProvider.ValidateToken(refreshToken)
	if err != nil {
		c.Log.Warnf("Invalid refresh token : %+v", err)
		return nil, fiber.NewError(fiber.StatusUnauthorized, "ERR_INVALID_REFRESH_TOKEN", "invalid refresh token")
	}

	if (*claims)["type"] != "refresh" {
		c.Log.Warnf("Invalid token type : %v", (*claims)["type"])
		return nil, fiber.NewError(fiber.StatusUnauthorized, "ERR_INVALID_REFRESH_TOKEN", "invalid refresh token")
	}

	var userID int64
	switch v := (*claims)["sub"].(type) {
	case float64:
		userID = int64(v)
	case string:
		return nil, fiber.NewError(fiber.StatusUnauthorized, "ERR_INVALID_REFRESH_TOKEN", "invalid refresh token")
	default:
		return nil, fiber.NewError(fiber.StatusUnauthorized, "ERR_INVALID_REFRESH_TOKEN", "invalid refresh token")
	}

	user := new(entity.User)
	if err := c.UserRepository.FindById(c.DB.WithContext(ctx), user, userID); err != nil {
		c.Log.Warnf("User not found for refresh token : %+v", err)
		return nil, fiber.NewError(fiber.StatusUnauthorized, "ERR_INVALID_REFRESH_TOKEN", "invalid refresh token")
	}

	accessToken, err := c.TokenProvider.GenerateAccessToken(user)
	if err != nil {
		c.Log.Warnf("Failed to generate access token : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return &model.LoginResponse{
		AccessToken: accessToken,
	}, nil
}
