package model

import (
	"time"

	"github.com/google/uuid"
)

// User represents the structure of the 'users' table.
type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
