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

func (r *UserRepository) FindByEmail(db *gorm.DB, user *entity.User, email string) error {
	return db.Where("email = ?", email).First(user).Error
}
