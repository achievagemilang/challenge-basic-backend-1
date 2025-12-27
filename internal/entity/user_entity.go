package entity

import (
	"time"
)

// User is a struct that represents a user entity
type User struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Email     string    `gorm:"column:email"`
	Password  string    `gorm:"column:password"`
	Name      string    `gorm:"column:name"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (u *User) TableName() string {
	return "users"
}
