package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/mail"
	"time"

	"github.com/ffardo/user-crud/interfaces"
	"github.com/ffardo/user-crud/models"
	"github.com/google/uuid"
)

var ErrUserNotFound = errors.New("Could not find user")
var ErrEmailRegistered = errors.New("Email already registered")
var ErrInvalidDateFormat = errors.New("Invalid date format")
var ErrInvalidUuidFormat = errors.New("Invalid uuid format")
var ErrInvalidEmailFormat = errors.New("Invalid email format")

type UserService struct {
	interfaces.UserRepository
}

func (s UserService) GetUser(user_uuid string) (models.User, error) {
	u, err := uuid.Parse(user_uuid)

	if err != nil {
		return models.User{}, ErrInvalidUuidFormat
	}
	user, err := s.UserRepository.GetUserByUUID(u)

	if err != nil {
		return models.User{}, ErrUserNotFound
	}

	return user, err
}

func (s UserService) CreateUser(name, birthDate, email, address, password string) (models.User, error) {
	bd, err := time.Parse("2006-01-02", birthDate)
	if err != nil {
		return models.User{}, ErrInvalidDateFormat
	}

	_, err = mail.ParseAddress(email)
	if err != nil {
		return models.User{}, ErrInvalidEmailFormat
	}

	exists, _ := s.UserRepository.UserExistsWithEmail(email)

	if exists == true {
		return models.User{}, ErrEmailRegistered
	}

	c := sha256.New()
	h := c.Sum([]byte(password))

	p := hex.EncodeToString(h)

	user := models.User{
		Name:      name,
		BirthDate: bd,
		Email:     email,
		Address:   address,
		Password:  p,
	}

	user, err = s.UserRepository.CreateUser(user)

	return user, err
}

func (s UserService) UpdateUser(user_uuid string, params map[string]string) (models.User, error) {
	u, err := uuid.Parse(user_uuid)

	if err != nil {
		return models.User{}, ErrInvalidUuidFormat
	}

	user, err := s.UserRepository.GetUserByUUID(u)

	if err != nil {
		return models.User{}, ErrUserNotFound
	}

	name, ok := params["name"]
	if ok {
		user.Name = name
	}

	bd, ok := params["birth_date"]
	if ok {
		t, err := time.Parse("2006-01-01", bd)

		if err != nil {
			return user, ErrInvalidDateFormat
		}
		user.BirthDate = t
	}

	password, ok := params["password"]
	if ok {
		s := sha256.New()
		h := s.Sum([]byte(password))
		p := hex.EncodeToString(h)
		user.Password = p
	}

	address, ok := params["address"]
	if ok {
		user.Address = address
	}

	email, ok := params["email"]
	if ok {
		_, err = mail.ParseAddress(email)
		if err != nil {
			return models.User{}, ErrInvalidEmailFormat
		}

		exists, _ := s.UserRepository.UserExistsWithEmailAndNotUuid(email, u)

		if exists {
			return models.User{}, ErrEmailRegistered
		}
		user.Email = email
	}

	user, err = s.UserRepository.UpdateUser(user)

	return user, err
}

func (s UserService) DeleteUser(user_uuid string) error {
	u, err := uuid.Parse(user_uuid)

	if err != nil {
		return ErrInvalidUuidFormat
	}

	_, err = s.UserRepository.GetUserByUUID(u)

	if err != nil {
		return ErrUserNotFound
	}
	return s.UserRepository.DeleteUser(u)

}
