package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UUID      uuid.UUID
	BirthDate time.Time
	Name      string
	Email     string
	Password  string
	Address   string
	Created   time.Time
	Updated   time.Time
}
