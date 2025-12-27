//go:generate mockgen -source=interfaces.go -destination=../../test/mocks/usecase_mocks.go -package=mocks
package usecase

import (
	"challenge-backend-1/internal/entity"
	"challenge-backend-1/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(db *gorm.DB, user *entity.User, email string) error
	FindById(db *gorm.DB, user *entity.User, id any) error
}

type UserProducer interface {
	Send(event *model.UserEvent) error
}
