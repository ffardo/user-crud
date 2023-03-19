package services

import (
	"errors"
	"testing"
	"time"

	"github.com/ffardo/user-crud/interfaces/mocks"
	"github.com/ffardo/user-crud/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	userRepository := new(mocks.UserRepository)
	user_uuid := "d035e79d-ffe9-4ebf-b665-747353b3ea40"

	u := uuid.MustParse(user_uuid)

	user := models.User{
		UUID:      uuid.MustParse(user_uuid),
		BirthDate: time.Now(),
		Name:      "John Doe",
		Email:     "joe25@mailprovider.com",
		Password:  "some_password",
		Address:   "3197 Woodrow Way",
	}

	userRepository.On("GetUserByUUID", u).Return(user, nil)

	userService := UserService{userRepository}

	expectedResult := user

	actualResult, _ := userService.GetUser(user_uuid)

	assert.Equal(t, expectedResult, actualResult)
}

func TestGetMissingUser(t *testing.T) {
	userRepository := new(mocks.UserRepository)

	user := models.User{}

	user_uuid := "d035e79d-ffe9-4ebf-b665-747353b3ea40"
	u := uuid.MustParse(user_uuid)

	userRepository.On("GetUserByUUID", u).Return(user, errors.New("could not find document"))

	userService := UserService{userRepository}

	_, err := userService.GetUser(user_uuid)

	assert.Equal(t, err, ErrUserNotFound)
}

func TestGetUserWithInvalidUUID(t *testing.T) {
	userRepository := new(mocks.UserRepository)

	userService := UserService{userRepository}

	_, err := userService.GetUser("not an uuid")

	assert.Equal(t, err, ErrInvalidUuidFormat)
}

func TestCreateUser(t *testing.T) {
	userRepository := new(mocks.UserRepository)
	userService := UserService{userRepository}
	user_uuid := "d035e79d-ffe9-4ebf-b665-747353b3ea40"
	user_bd, _ := time.Parse("2006-01-02", "1970-01-31")

	address := "3197 Woodrow Way"
	email := "joe25@mailprovider.com"
	password := "736563726574e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	name := "John Doe"

	user := models.User{
		BirthDate: user_bd,
		Name:      name,
		Email:     email,
		Password:  password,
		Address:   address,
	}

	created_user := models.User{
		UUID:      uuid.MustParse(user_uuid),
		BirthDate: user_bd,
		Name:      name,
		Email:     email,
		Password:  password,
		Address:   address,
	}

	userRepository.On("UserExistsWithEmail", "joe25@mailprovider.com").Return(false, nil)

	userRepository.On(
		"CreateUser", user,
	).Return(created_user, nil)

	ret_user, err := userService.CreateUser(
		"John Doe", "1970-01-31", "joe25@mailprovider.com", "3197 Woodrow Way", "secret",
	)

	assert.Equal(t, err, nil)
	assert.Equal(t, ret_user, created_user)

}

func TestCreateUserWithInvalidBirthDate(t *testing.T) {
	userRepository := new(mocks.UserRepository)
	userService := UserService{userRepository}

	userRepository.On("UserExistsWithEmail", "joe25@mailprovider.com").Return(false, nil)

	_, err := userService.CreateUser(
		"John Doe", "1970-13-1", "joe25@mailprovider.com", "3197 Woodrow Way", "secret",
	)

	userRepository.AssertNotCalled(t, "CreateUser")

	assert.Equal(t, err, ErrInvalidDateFormat)

}

func TestCreateUserWithInvalidEmail(t *testing.T) {
	userRepository := new(mocks.UserRepository)
	userService := UserService{userRepository}

	userRepository.On("UserExistsWithEmail", "joe25@mailprovider.com").Return(false, nil)

	_, err := userService.CreateUser(
		"John Doe", "1970-01-02", "not_valid_email", "3197 Woodrow Way", "secret",
	)

	userRepository.AssertNotCalled(t, "CreateUser")

	assert.Equal(t, err, ErrInvalidEmailFormat)

}

func TestCreateUserWithExistingEmail(t *testing.T) {
	userRepository := new(mocks.UserRepository)
	userService := UserService{userRepository}

	userRepository.On("UserExistsWithEmail", "joe25@mailprovider.com").Return(true, nil)

	_, err := userService.CreateUser(
		"John Doe", "1970-01-01", "joe25@mailprovider.com", "3197 Woodrow Way", "secret",
	)

	userRepository.AssertNotCalled(t, "CreateUser")

	assert.Equal(t, err, ErrEmailRegistered)
}

func TestUpdateUser(t *testing.T) {
	userRepository := new(mocks.UserRepository)
	userService := UserService{userRepository}
	user_uuid := "d035e79d-ffe9-4ebf-b665-747353b3ea40"

	address := "3197 Woodrow Way"
	email := "joe25@mailprovider.com"
	password := "4d79206e65772070617373776f7264e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	name := "John Doe"

	original_user := models.User{
		UUID:      uuid.MustParse(user_uuid),
		BirthDate: time.Now(),
		Name:      name,
		Email:     email,
		Password:  password,
		Address:   address,
	}

	modified_user := models.User{
		UUID:      original_user.UUID,
		BirthDate: original_user.BirthDate,
		Name:      "John Nobody",
		Email:     email,
		Password:  password,
		Address:   address,
	}

	userRepository.On("GetUserByUUID", uuid.MustParse(user_uuid)).Return(original_user, nil)
	userRepository.On("UpdateUser", modified_user).Return(modified_user, nil)

	params := map[string]string{
		"name":     "John Nobody",
		"password": "My new password",
	}

	update_user, err := userService.UpdateUser(user_uuid, params)

	assert.Equal(t, err, nil)
	assert.Equal(t, modified_user.Password, update_user.Password)
}

func TestUpdateUserNotFound(t *testing.T) {
	userRepository := new(mocks.UserRepository)
	userService := UserService{userRepository}
	user_uuid := "d035e79d-ffe9-4ebf-b665-747353b3ea40"

	u := uuid.MustParse(user_uuid)

	userRepository.On("GetUserByUUID", u).Return(models.User{}, errors.New("document not found"))

	params := map[string]string{
		"name":     "John Nobody",
		"password": "My new password",
	}

	_, err := userService.UpdateUser(user_uuid, params)
	assert.Equal(t, err, ErrUserNotFound)

}

func TestUpdateUserWithInvalidUUID(t *testing.T) {
	userRepository := new(mocks.UserRepository)

	userService := UserService{userRepository}

	params := map[string]string{
		"name":       "John Nobody",
		"password":   "My new password",
		"birth_date": "1970-13-32",
	}

	_, err := userService.UpdateUser("not an uuid", params)

	assert.Equal(t, err, ErrInvalidUuidFormat)
}

func TestUpdateUserWithInvalidDateFormat(t *testing.T) {
	userRepository := new(mocks.UserRepository)
	userService := UserService{userRepository}
	user_uuid := "d035e79d-ffe9-4ebf-b665-747353b3ea40"
	user := models.User{
		UUID:      uuid.MustParse(user_uuid),
		BirthDate: time.Now(),
		Name:      "John Doe",
		Email:     "joe25@mailprovider.com",
		Password:  "some_password",
		Address:   "3197 Woodrow Way",
	}

	userRepository.On("GetUserByUUID", uuid.MustParse(user_uuid)).Return(user, nil)

	params := map[string]string{
		"name":       "John Nobody",
		"password":   "My new password",
		"birth_date": "1970-13-32",
	}

	_, err := userService.UpdateUser(user_uuid, params)

	assert.Equal(t, err, ErrInvalidDateFormat)

}

func TestUpdateUserWithEmailCollision(t *testing.T) {
	userRepository := new(mocks.UserRepository)
	userService := UserService{userRepository}
	user_uuid := "d035e79d-ffe9-4ebf-b665-747353b3ea40"

	user := models.User{
		UUID:      uuid.MustParse(user_uuid),
		BirthDate: time.Now(),
		Name:      "John Doe",
		Email:     "joe25@mailprovider.com",
		Password:  "some_password",
		Address:   "3197 Woodrow Way",
	}

	userRepository.On("GetUserByUUID", uuid.MustParse(user_uuid)).Return(user, nil)
	userRepository.On("UserExistsWithEmailAndNotUuid", "different_email@mailprovider.com", uuid.MustParse(user_uuid)).Return(true, nil)

	params := map[string]string{
		"name":     "John Nobody",
		"password": "My new password",
		"email":    "different_email@mailprovider.com",
	}

	_, err := userService.UpdateUser(user_uuid, params)

	assert.Equal(t, err, ErrEmailRegistered)

}

func TestDeleteUser(t *testing.T) {
	userRepository := new(mocks.UserRepository)
	userService := UserService{userRepository}
	user_uuid := "d035e79d-ffe9-4ebf-b665-747353b3ea40"

	u := uuid.MustParse(user_uuid)

	user := models.User{
		UUID:      uuid.MustParse(user_uuid),
		BirthDate: time.Now(),
		Name:      "John Doe",
		Email:     "joe25@mailprovider.com",
		Password:  "some_password",
		Address:   "3197 Woodrow Way",
	}

	userRepository.On("GetUserByUUID", u).Return(user, nil)

	userRepository.On("DeleteUser", u).Return(nil)

	err := userService.DeleteUser(user_uuid)
	assert.Equal(t, err, nil)

}

func TestDeleteUserWithInvalidUUID(t *testing.T) {
	userRepository := new(mocks.UserRepository)

	userService := UserService{userRepository}

	err := userService.DeleteUser("not an uuid")

	assert.Equal(t, err, ErrInvalidUuidFormat)
}

func TestDeleteUserNotFound(t *testing.T) {
	userRepository := new(mocks.UserRepository)
	userService := UserService{userRepository}
	user_uuid := "d035e79d-ffe9-4ebf-b665-747353b3ea40"

	u := uuid.MustParse(user_uuid)

	userRepository.On("GetUserByUUID", u).Return(models.User{}, errors.New("document not found"))

	err := userService.DeleteUser(user_uuid)
	assert.Equal(t, ErrUserNotFound, err)

}
