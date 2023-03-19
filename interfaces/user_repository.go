package interfaces

import (
	"github.com/ffardo/user-crud/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetUserByUUID(uuid.UUID) (models.User, error)
	UserExistsWithEmail(string) (bool, error)
	UserExistsWithEmailAndNotUuid(string, uuid.UUID) (bool, error)
	CreateUser(models.User) (models.User, error)
	UpdateUser(models.User) (models.User, error)
	DeleteUser(uuid.UUID) error
}
