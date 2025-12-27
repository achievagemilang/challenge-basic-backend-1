package repository

import (
	"challenge-backend-1/internal/entity"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.User]
	Log *zap.SugaredLogger
}

func NewUserRepository(log *zap.SugaredLogger) *UserRepository {
	return &UserRepository{
		Log: log,
	}
}

func (r *UserRepository) FindByToken(db *gorm.DB, user *entity.User, token string) error {
	return db.Where("token = ?", token).First(user).Error
}
