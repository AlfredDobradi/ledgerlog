package models

import (
	"time"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	Email     string `json:"email"`
	PublicKey string `json:"public_key"`
	Name      string `json:"name"`
}

type RegisterResponse struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

type SendPostRequest struct {
	Owner   string `json:"-"`
	Message string `json:"message"`
}

type SendPostResponse struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

type Post struct {
	Timestamp time.Time `json:"timestamp"`
	Email     string    `json:"email"`
	Message   string    `json:"message"`
}

type LedgerEntry struct {
	Timestamp time.Time
	ID        uuid.UUID
	Prev      uuid.UUID
	Subject   uuid.UUID
	Content   []byte
}

type User struct {
	ID            uuid.UUID `db:"id"`
	Email         string    `db:"email"`
	PreferredName string    `db:"preferred_name"`
	PublicKey     string    `db:"public_key"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
