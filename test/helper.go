package test

import (
	"testing"

	"challenge-backend-1/internal/entity"

	"github.com/stretchr/testify/assert"
)

func ClearAll() {
	ClearUsers()
}

func ClearUsers() {
	err := Db.Where("id is not null").Delete(&entity.User{}).Error
	if err != nil {
		Log.Fatalf("Failed clear user data : %+v", err)
	}
}

func GetFirstUser(t *testing.T) *entity.User {
	user := new(entity.User)
	err := Db.First(user).Error
	assert.Nil(t, err)
	return user
}
