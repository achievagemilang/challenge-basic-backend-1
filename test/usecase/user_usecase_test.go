package usecase_test

import (
	"context"
	"errors"
	"testing"

	"challenge-backend-1/internal/entity"
	"challenge-backend-1/internal/model"
	"challenge-backend-1/internal/usecase"
	"challenge-backend-1/test/mocks"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestUserUseCase_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockProducer := mocks.NewMockUserProducer(ctrl)
	mockToken := mocks.NewMockTokenProvider(ctrl)

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	logger := zap.NewNop().Sugar()
	validate := validator.New()

	uc := usecase.NewUserUseCase(gormDB, logger, validate, mockRepo, mockProducer, mockToken)

	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectCommit()

		user := &entity.User{
			ID:       1,
			Email:    "test@example.com",
			Password: string(hashedPassword),
			Name:     "Test User",
		}

		mockRepo.EXPECT().
			FindByEmail(gomock.Any(), gomock.Any(), "test@example.com").
			DoAndReturn(func(db *gorm.DB, u *entity.User, email string) error {
				*u = *user
				return nil
			})

		mockToken.EXPECT().GenerateAccessToken(user).Return("access_token", nil)
		mockToken.EXPECT().GenerateRefreshToken(user).Return("refresh_token", nil)

		req := &model.LoginUserRequest{
			Email:    "test@example.com",
			Password: password,
		}

		resp, err := uc.Login(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "access_token", resp.AccessToken)
		assert.Equal(t, "refresh_token", resp.RefreshToken)
	})

	t.Run("InvalidCredentials_EmailNotFound", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectRollback()

		mockRepo.EXPECT().
			FindByEmail(gomock.Any(), gomock.Any(), "wrong@example.com").
			Return(errors.New("record not found"))

		req := &model.LoginUserRequest{
			Email:    "wrong@example.com",
			Password: password,
		}

		resp, err := uc.Login(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)

		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			assert.Equal(t, fiber.StatusUnauthorized, fiberErr.Code)
		}
	})

	t.Run("InvalidCredentials_WrongPassword", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectRollback()

		user := &entity.User{
			ID:       1,
			Email:    "test@example.com",
			Password: string(hashedPassword),
		}

		mockRepo.EXPECT().
			FindByEmail(gomock.Any(), gomock.Any(), "test@example.com").
			DoAndReturn(func(db *gorm.DB, u *entity.User, email string) error {
				*u = *user
				return nil
			})

		req := &model.LoginUserRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		resp, err := uc.Login(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)

		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			assert.Equal(t, fiber.StatusUnauthorized, fiberErr.Code)
		}
	})

	t.Run("TokenGenerationError", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectRollback()

		user := &entity.User{
			ID:       1,
			Email:    "test@example.com",
			Password: string(hashedPassword),
		}

		mockRepo.EXPECT().
			FindByEmail(gomock.Any(), gomock.Any(), "test@example.com").
			DoAndReturn(func(db *gorm.DB, u *entity.User, email string) error {
				*u = *user
				return nil
			})

		mockToken.EXPECT().GenerateAccessToken(user).Return("", errors.New("token error"))

		req := &model.LoginUserRequest{
			Email:    "test@example.com",
			Password: password,
		}

		resp, err := uc.Login(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			assert.Equal(t, fiber.StatusInternalServerError, fiberErr.Code)
		}
	})
}

func TestUserUseCase_Refresh(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockProducer := mocks.NewMockUserProducer(ctrl)
	mockToken := mocks.NewMockTokenProvider(ctrl)

	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	logger := zap.NewNop().Sugar()
	validate := validator.New()

	uc := usecase.NewUserUseCase(gormDB, logger, validate, mockRepo, mockProducer, mockToken)

	t.Run("Success", func(t *testing.T) {
		refreshToken := "valid_refresh_token"
		claims := jwt.MapClaims{
			"sub":  float64(1),
			"type": "refresh",
		}
		user := &entity.User{ID: 1}

		mockToken.EXPECT().ValidateToken(refreshToken).Return(&claims, nil)
		mockRepo.EXPECT().FindById(gomock.Any(), gomock.Any(), int64(1)).DoAndReturn(func(db *gorm.DB, u *entity.User, id any) error {
			*u = *user
			return nil
		})
		mockToken.EXPECT().GenerateAccessToken(user).Return("new_access_token", nil)

		resp, err := uc.Refresh(context.Background(), refreshToken)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "new_access_token", resp.AccessToken)
	})

	t.Run("InvalidToken", func(t *testing.T) {
		refreshToken := "invalid_token"

		mockToken.EXPECT().ValidateToken(refreshToken).Return(nil, errors.New("invalid"))

		resp, err := uc.Refresh(context.Background(), refreshToken)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("InvalidTokenType", func(t *testing.T) {
		refreshToken := "invalid_type_token"
		claims := jwt.MapClaims{
			"sub":  float64(1),
			"type": "access", // Incorrect type
		}

		mockToken.EXPECT().ValidateToken(refreshToken).Return(&claims, nil)

		resp, err := uc.Refresh(context.Background(), refreshToken)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}
