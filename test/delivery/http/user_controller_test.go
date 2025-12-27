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

	request := httptest.NewRequest(http.MethodPost, "/api/v1/session", strings.NewReader(string(bodyJson)))
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

	request := httptest.NewRequest(http.MethodPost, "/api/v1/session", strings.NewReader(string(bodyJson)))
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

	request := httptest.NewRequest(http.MethodPost, "/api/v1/session", strings.NewReader(string(bodyJson)))
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

func TestRefresh(t *testing.T) {
	test.ClearAll()
	registerUser()

	// 1. Login to get refresh token
	loginReq := model.LoginUserRequest{
		Email:    "alice@mail.com",
		Password: "123456",
	}
	bodyJson, _ := json.Marshal(loginReq)
	reqFn := func(method, url string, body io.Reader) *http.Request {
		r := httptest.NewRequest(method, url, body)
		r.Header.Set("Content-Type", "application/json")
		return r
	}

	respLogin, _ := test.App.Test(reqFn(http.MethodPost, "/api/v1/session", strings.NewReader(string(bodyJson))))
	bytesLogin, _ := io.ReadAll(respLogin.Body)
	var loginBody model.WebResponse[model.LoginResponse]
	json.Unmarshal(bytesLogin, &loginBody)
	refreshToken := loginBody.Data.RefreshToken

	// 2. Refresh Success
	reqRefresh := reqFn(http.MethodPut, "/api/v1/session", nil)
	reqRefresh.Header.Set("Authorization", "Bearer "+refreshToken)

	resp, err := test.App.Test(reqRefresh)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	var responseBody model.WebResponse[model.LoginResponse]
	err = json.Unmarshal(bytes, &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, responseBody.Ok)
	assert.NotEmpty(t, responseBody.Data.AccessToken)

	// 3. Refresh Invalid Token
	reqInvalid := reqFn(http.MethodPut, "/api/v1/session", nil)
	reqInvalid.Header.Set("Authorization", "Bearer invalid-token")
	respInvalid, _ := test.App.Test(reqInvalid)
	assert.Equal(t, http.StatusUnauthorized, respInvalid.StatusCode)
}
