package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ffardo/user-crud/interfaces/mocks"
	"github.com/ffardo/user-crud/models"
	"github.com/ffardo/user-crud/routes"
	"github.com/ffardo/user-crud/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	gin.SetMode("test")
	us := new(mocks.UserService)

	uc := UserController{
		us,
	}

	router := routes.InitRouter(&uc, "test_key")

	address := "3197 Woodrow Way"
	hashed_password := "736563726574e3b0c44298fc1c"
	name := "John Doe"
	email := "joe25@mailprovider.com"

	birthDate := "1970-01-31"

	user_uuid := "d035e79d-ffe9-4ebf-b665-747353b3ea40"
	user_bd, _ := time.Parse("2006-01-02", birthDate)

	user := models.User{
		UUID:      uuid.MustParse(user_uuid),
		BirthDate: user_bd,
		Name:      name,
		Email:     email,
		Password:  hashed_password,
		Address:   address,
	}

	us.On(
		"CreateUser",
		"John Doe",
		"1970-01-31",
		"joe25@mailprovider.com",
		address,
		"secret",
	).Return(user, nil)

	new_user_request := UserRequest{
		Name:      name,
		BirthDate: birthDate,
		Email:     email,
		Address:   address,
		Password:  "secret",
	}

	marshalled, _ := json.Marshal(new_user_request)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/users/", bytes.NewReader(marshalled))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", "test_key")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var res UserResponse

	err := json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(t, err, nil)
	assert.Equal(t, res.Name, name)
	assert.Equal(t, res.Email, email)
	assert.Equal(t, res.Address, address)
	assert.Equal(t, res.Password, hashed_password)
	assert.Equal(t, res.BirthDate, birthDate)
	assert.Equal(t, res.UUID, user_uuid)

}

func TestCreateUserUnauthorized(t *testing.T) {
	gin.SetMode("test")
	us := new(mocks.UserService)

	uc := UserController{
		us,
	}

	router := routes.InitRouter(&uc, "test_key")

	address := "3197 Woodrow Way"
	name := "John Doe"
	email := "joe25@mailprovider.com"

	birthDate := "1970-01-31"

	new_user_request := UserRequest{
		Name:      name,
		BirthDate: birthDate,
		Email:     email,
		Address:   address,
		Password:  "secret",
	}

	marshalled, _ := json.Marshal(new_user_request)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/users/", bytes.NewReader(marshalled))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", "invalid key")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

}

func runCreateTestWithError(t *testing.T, serviceError error, statusCode int) {
	gin.SetMode("test")
	us := new(mocks.UserService)

	uc := UserController{
		us,
	}

	router := routes.InitRouter(&uc, "test_key")

	address := "3197 Woodrow Way"
	name := "John Doe"
	email := "joe25@mailprovider.com"

	birthDate := "1970-01-31"

	us.On(
		"CreateUser",
		"John Doe",
		"1970-01-31",
		"joe25@mailprovider.com",
		address,
		"secret",
	).Return(models.User{}, serviceError)

	new_user_request := UserRequest{
		Name:      name,
		BirthDate: birthDate,
		Email:     email,
		Address:   address,
		Password:  "secret",
	}

	marshalled, _ := json.Marshal(new_user_request)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/users/", bytes.NewReader(marshalled))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", "test_key")
	router.ServeHTTP(w, req)
	assert.Equal(t, statusCode, w.Code)

	var res ErrorResponse

	err := json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(t, err, nil)
	assert.Equal(t, res.Error, serviceError.Error())

}

func TestCreateUserInvalidDateFormat(t *testing.T) {
	runCreateTestWithError(t, services.ErrInvalidDateFormat, http.StatusBadRequest)
}

func TestCreateUserInvalidEmailFormat(t *testing.T) {
	runCreateTestWithError(t, services.ErrInvalidEmailFormat, http.StatusBadRequest)
}

func TestCreateUserExistingEmail(t *testing.T) {
	runCreateTestWithError(t, services.ErrEmailRegistered, http.StatusBadRequest)
}

func TestGetUser(t *testing.T) {
	gin.SetMode("test")
	us := new(mocks.UserService)

	uc := UserController{
		us,
	}

	router := routes.InitRouter(&uc, "test_key")

	address := "3197 Woodrow Way"
	name := "John Doe"
	email := "joe25@mailprovider.com"
	user_uuid := "d035e79d-ffe9-4ebf-b665-747353b3ea40"

	birthDate := "1970-01-31"

	user_bd, _ := time.Parse("2006-01-02", birthDate)
	hashed_password := "736563726574e3b0c44298fc1c"

	user := models.User{
		UUID:      uuid.MustParse(user_uuid),
		BirthDate: user_bd,
		Name:      name,
		Email:     email,
		Password:  hashed_password,
		Address:   address,
	}

	us.On(
		"GetUser",
		user_uuid,
	).Return(user, nil)

	url := fmt.Sprintf("/api/users/%s", user_uuid)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", "test_key")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var res UserResponse

	err := json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(t, err, nil)
	assert.Equal(t, res.Name, name)
	assert.Equal(t, res.Email, email)
	assert.Equal(t, res.Address, address)
	assert.Equal(t, res.Password, hashed_password)
	assert.Equal(t, res.BirthDate, birthDate)
	assert.Equal(t, res.UUID, user_uuid)

}

func TestGetUserUnauthorized(t *testing.T) {
	gin.SetMode("test")
	us := new(mocks.UserService)

	uc := UserController{
		us,
	}

	router := routes.InitRouter(&uc, "test_key")
	user_uuid := "d035e79d-ffe9-4ebf-b665-747353b3ea40"

	url := fmt.Sprintf("/api/users/%s", user_uuid)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", "some_other_key")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

}

func runGetUserTestWithError(t *testing.T, serviceError error, statusCode int) {
	gin.SetMode("test")
	us := new(mocks.UserService)

	uc := UserController{
		us,
	}

	router := routes.InitRouter(&uc, "test_key")

	user_uuid := "d035e79d-ffe9-4ebf-b665-747353b3ea40"

	us.On(
		"GetUser",
		user_uuid,
	).Return(models.User{}, serviceError)

	url := fmt.Sprintf("/api/users/%s", user_uuid)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", "test_key")
	router.ServeHTTP(w, req)
	assert.Equal(t, statusCode, w.Code)

	var res ErrorResponse

	err := json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(t, err, nil)
	assert.Equal(t, res.Error, serviceError.Error())

}

func TestGetUserInvalidUUID(t *testing.T) {
	runGetUserTestWithError(t, services.ErrInvalidUuidFormat, http.StatusBadRequest)
}

func TestGetUserNotFound(t *testing.T) {
	runGetUserTestWithError(t, services.ErrUserNotFound, http.StatusNotFound)
}

func TestUpdateUser(t *testing.T) {
	gin.SetMode("test")
	us := new(mocks.UserService)

	uc := UserController{
		us,
	}

	router := routes.InitRouter(&uc, "test_key")

	address := "3197 Woodrow Way"
	hashed_password := "736563726574e3b0c44298fc1c"
	name := "John Doe"
	email := "joe25@mailprovider.com"

	birthDate := "1970-01-31"

	user_uuid := "d035e79d-ffe9-4ebf-b665-747353b3ea40"
	user_bd, _ := time.Parse("2006-01-02", birthDate)

	user := models.User{
		UUID:      uuid.MustParse(user_uuid),
		BirthDate: user_bd,
		Name:      name,
		Email:     email,
		Password:  hashed_password,
		Address:   address,
	}

	user_request := UserRequest{
		Name:      name,
		BirthDate: birthDate,
		Email:     email,
		Address:   address,
		Password:  "secret",
	}

	params := map[string]string{
		"name":       user_request.Name,
		"email":      user_request.Email,
		"birth_date": user_request.BirthDate,
		"address":    user_request.Address,
		"password":   user_request.Password,
	}

	us.On(
		"UpdateUser",
		user_uuid,
		params,
	).Return(user, nil)

	marshalled, _ := json.Marshal(user_request)
	url := fmt.Sprintf("/api/users/%s", user_uuid)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewReader(marshalled))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", "test_key")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var res UserResponse

	err := json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(t, err, nil)
	assert.Equal(t, res.Name, name)
	assert.Equal(t, res.Email, email)
	assert.Equal(t, res.Address, address)
	assert.Equal(t, res.Password, hashed_password)
	assert.Equal(t, res.BirthDate, birthDate)
	assert.Equal(t, res.UUID, user_uuid)

}

func TestUpdateUserUnauthorized(t *testing.T) {
	gin.SetMode("test")
	us := new(mocks.UserService)

	uc := UserController{
		us,
	}

	router := routes.InitRouter(&uc, "test_key")

	address := "3197 Woodrow Way"
	name := "John Doe"
	email := "joe25@mailprovider.com"

	birthDate := "1970-01-31"

	user_uuid := "d035e79d-ffe9-4ebf-b665-747353b3ea40"

	user_request := UserRequest{
		Name:      name,
		BirthDate: birthDate,
		Email:     email,
		Address:   address,
		Password:  "secret",
	}

	marshalled, _ := json.Marshal(user_request)
	url := fmt.Sprintf("/api/users/%s", user_uuid)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewReader(marshalled))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", "some_other_key")

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

}

func runUpdateTestWithError(t *testing.T, serviceError error, statusCode int) {
	gin.SetMode("test")
	us := new(mocks.UserService)

	uc := UserController{
		us,
	}

	router := routes.InitRouter(&uc, "test_key")

	address := "3197 Woodrow Way"
	name := "John Doe"
	email := "joe25@mailprovider.com"

	birthDate := "1970-01-31"

	user_uuid := "d035e79d-ffe9-4ebf-b665-747353b3ea40"

	user_request := UserRequest{
		Name:      name,
		BirthDate: birthDate,
		Email:     email,
		Address:   address,
		Password:  "secret",
	}

	params := map[string]string{
		"name":       user_request.Name,
		"email":      user_request.Email,
		"birth_date": user_request.BirthDate,
		"address":    user_request.Address,
		"password":   user_request.Password,
	}

	us.On(
		"UpdateUser",
		user_uuid,
		params,
	).Return(models.User{}, serviceError)

	marshalled, _ := json.Marshal(user_request)
	url := fmt.Sprintf("/api/users/%s", user_uuid)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewReader(marshalled))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", "test_key")
	router.ServeHTTP(w, req)
	assert.Equal(t, statusCode, w.Code)

	var res ErrorResponse

	err := json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(t, err, nil)
	assert.Equal(t, res.Error, serviceError.Error())

}

func TestUpdateUserInvalidUUID(t *testing.T) {
	runUpdateTestWithError(t, services.ErrInvalidUuidFormat, http.StatusBadRequest)
}

func TestUpdateUserInvalidDateFormat(t *testing.T) {
	runUpdateTestWithError(t, services.ErrInvalidDateFormat, http.StatusBadRequest)
}

func TestUpdateUserInvalidEmail(t *testing.T) {
	runUpdateTestWithError(t, services.ErrInvalidEmailFormat, http.StatusBadRequest)
}

func TestUpdateUserExistingEmail(t *testing.T) {
	runUpdateTestWithError(t, services.ErrEmailRegistered, http.StatusBadRequest)
}

func TestUpdateUserNotFound(t *testing.T) {
	runUpdateTestWithError(t, services.ErrUserNotFound, http.StatusNotFound)
}

func TestDeleteUser(t *testing.T) {
	gin.SetMode("test")
	us := new(mocks.UserService)

	uc := UserController{
		us,
	}

	router := routes.InitRouter(&uc, "test_key")

	user_uuid := "d035e79d-ffe9-4ebf-b665-747353b3ea40"

	us.On(
		"DeleteUser",
		user_uuid,
	).Return(nil)

	url := fmt.Sprintf("/api/users/%s", user_uuid)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", "test_key")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

}

func TestDeleteUserUnauthorized(t *testing.T) {
	gin.SetMode("test")
	us := new(mocks.UserService)

	uc := UserController{
		us,
	}

	router := routes.InitRouter(&uc, "test_key")

	user_uuid := "d035e79d-ffe9-4ebf-b665-747353b3ea40"

	us.On(
		"DeleteUser",
		user_uuid,
	).Return(nil)

	url := fmt.Sprintf("/api/users/%s", user_uuid)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", "some_other_key")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

}

func runDeleteUserTestWithError(t *testing.T, serviceError error, statusCode int) {
	gin.SetMode("test")
	us := new(mocks.UserService)

	uc := UserController{
		us,
	}

	router := routes.InitRouter(&uc, "test_key")

	user_uuid := "d035e79d-ffe9-4ebf-b665-747353b3ea40"

	us.On(
		"DeleteUser",
		user_uuid,
	).Return(serviceError)

	url := fmt.Sprintf("/api/users/%s", user_uuid)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", "test_key")
	router.ServeHTTP(w, req)
	assert.Equal(t, statusCode, w.Code)

	var res ErrorResponse

	err := json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(t, err, nil)
	assert.Equal(t, res.Error, serviceError.Error())

}

func TestDeleteUserInvalidUuidFormat(t *testing.T) {
	runDeleteUserTestWithError(t, services.ErrInvalidUuidFormat, http.StatusBadRequest)
}

func TestDeleteUserNotFound(t *testing.T) {
	runDeleteUserTestWithError(t, services.ErrUserNotFound, http.StatusNotFound)
}
