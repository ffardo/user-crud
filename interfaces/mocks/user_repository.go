package mocks

import (
	"github.com/ffardo/user-crud/models"
	"github.com/google/uuid"
	mock "github.com/stretchr/testify/mock"
)

type UserRepository struct {
	mock.Mock
}

func (u UserRepository) GetUserByUUID(user_uuid uuid.UUID) (models.User, error) {
	ret := u.Called(user_uuid)

	var user models.User
	if rf, ok := ret.Get(0).(func(uuid.UUID) models.User); ok {
		user = rf(user_uuid)
	} else {
		user = ret.Get(0).(models.User)
	}

	var err error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		err = rf(user_uuid)
	} else {
		err = ret.Error(1)
	}

	return user, err
}

func (u UserRepository) UserExistsWithEmail(email string) (bool, error) {
	ret := u.Called(email)

	var exists bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		exists = rf(email)
	} else {
		exists = ret.Get(0).(bool)
	}

	var err error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		err = rf(email)
	} else {
		err = ret.Error(1)
	}

	return exists, err
}

func (u UserRepository) UserExistsWithEmailAndNotUuid(email string, user_uuid uuid.UUID) (bool, error) {
	ret := u.Called(email, user_uuid)

	var exists bool
	if rf, ok := ret.Get(0).(func(string, uuid.UUID) bool); ok {
		exists = rf(email, user_uuid)
	} else {
		exists = ret.Get(0).(bool)
	}

	var err error
	if rf, ok := ret.Get(1).(func(string, uuid.UUID) error); ok {
		err = rf(email, user_uuid)
	} else {
		err = ret.Error(1)
	}

	return exists, err
}

func (u UserRepository) CreateUser(user models.User) (models.User, error) {
	ret := u.Called(user)

	var mockedUser models.User
	if rf, ok := ret.Get(0).(func(models.User) models.User); ok {
		mockedUser = rf(user)
	} else {
		mockedUser = ret.Get(0).(models.User)
	}

	var err error
	if rf, ok := ret.Get(1).(func(models.User) error); ok {
		err = rf(user)
	} else {
		err = ret.Error(1)
	}

	return mockedUser, err
}

func (u UserRepository) UpdateUser(user models.User) (models.User, error) {
	ret := u.Called(user)

	var mockedUser models.User
	if rf, ok := ret.Get(0).(func(models.User) models.User); ok {
		mockedUser = rf(user)
	} else {
		mockedUser = ret.Get(0).(models.User)
	}

	var err error
	if rf, ok := ret.Get(1).(func(models.User) error); ok {
		err = rf(user)
	} else {
		err = ret.Error(1)
	}

	return mockedUser, err
}

func (u UserRepository) DeleteUser(user_uuid uuid.UUID) error {
	ret := u.Called(user_uuid)
	var err error
	if rf, ok := ret.Get(0).(func(uuid.UUID) error); ok {
		err = rf(user_uuid)
	} else {
		err = ret.Error(0)
	}
	return err
}
