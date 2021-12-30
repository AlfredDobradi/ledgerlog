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
