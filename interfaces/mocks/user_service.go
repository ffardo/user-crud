package mocks

import (
	"github.com/ffardo/user-crud/models"
	mock "github.com/stretchr/testify/mock"
)

type UserService struct {
	mock.Mock
}

func (s UserService) GetUser(user_uuid string) (models.User, error) {
	ret := s.Called(user_uuid)
	var user models.User

	if rf, ok := ret.Get(0).(func(string) models.User); ok {
		user = rf(user_uuid)
	} else {
		user = ret.Get(0).(models.User)
	}

	var err error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		err = rf(user_uuid)
	} else {
		err = ret.Error(1)
	}

	return user, err
}

func (s UserService) CreateUser(name, birthDate, email, address, password string) (models.User, error) {
	ret := s.Called(name, birthDate, email, address, password)
	var user models.User

	if rf, ok := ret.Get(0).(func(string, string, string, string, string) models.User); ok {
		user = rf(name, birthDate, email, address, password)
	} else {
		user = ret.Get(0).(models.User)
	}

	var err error
	if rf, ok := ret.Get(1).(func(string, string, string, string, string) error); ok {
		err = rf(name, birthDate, email, address, password)
	} else {
		err = ret.Error(1)
	}

	return user, err
}

func (s UserService) UpdateUser(user_uuid string, params map[string]string) (models.User, error) {
	ret := s.Called(user_uuid, params)
	var user models.User

	if rf, ok := ret.Get(0).(func(string, map[string]string) models.User); ok {
		user = rf(user_uuid, params)
	} else {
		user = ret.Get(0).(models.User)
	}

	var err error
	if rf, ok := ret.Get(1).(func(string, map[string]string) error); ok {
		err = rf(user_uuid, params)
	} else {
		err = ret.Error(1)
	}

	return user, err
}

func (s UserService) DeleteUser(user_uuid string) error {
	ret := s.Called(user_uuid)

	var err error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		err = rf(user_uuid)
	} else {
		err = ret.Error(0)
	}

	return err
}
