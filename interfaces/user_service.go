package interfaces

import "github.com/ffardo/user-crud/models"

type UserService interface {
	GetUser(string) (models.User, error)
	CreateUser(string, string, string, string, string) (models.User, error)
	UpdateUser(string, map[string]string) (models.User, error)
	DeleteUser(string) error
}
