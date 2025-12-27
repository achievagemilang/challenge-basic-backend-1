package test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"challenge-backend-1/internal/entity"
	"challenge-backend-1/internal/model"
	"challenge-backend-1/test"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func registerUser() *entity.User {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	user := &entity.User{
		Email:     "alice@mail.com",
		Password:  string(passwordHash),
		Name:      "Alice",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	test.Db.Create(user)
	return user
}

func TestLogin(t *testing.T) {
	test.ClearAll()

	user := registerUser()

	requestBody := model.LoginUserRequest{
		Email:    "alice@mail.com",
		Password: "123456",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/session", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := test.App.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.LoginResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.True(t, responseBody.Ok)
	assert.Equal(t, user.ID, responseBody.Data.User.ID)
	assert.Equal(t, user.Email, responseBody.Data.User.Email)
	assert.NotEmpty(t, responseBody.Data.AccessToken)
	assert.NotEmpty(t, responseBody.Data.RefreshToken)
}

func TestLoginWrongEmail(t *testing.T) {
	test.ClearAll()

	registerUser()

	requestBody := model.LoginUserRequest{
		Email:    "wrong@mail.com",
		Password: "123456",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/session", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := test.App.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse)
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.False(t, responseBody.Ok)
	assert.Equal(t, "ERR_INVALID_CREDS", responseBody.Err)
}

func TestLoginWrongPassword(t *testing.T) {
	test.ClearAll()

	registerUser()

	requestBody := model.LoginUserRequest{
		Email:    "alice@mail.com",
		Password: "wrong",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/session", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := test.App.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse)
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.False(t, responseBody.Ok)
	assert.Equal(t, "ERR_INVALID_CREDS", responseBody.Err)
}
