package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
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
	Owner   uuid.UUID `json:"owner"`
	Message string    `json:"message"`
}

type SendPostResponse struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

type Post struct {
	ID        uuid.UUID `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	User      User      `json:"user"`
	Message   string    `json:"message"`
}

type PostDisplay struct {
	ID                uuid.UUID
	Timestamp         time.Time
	UserID            uuid.UUID
	UserPreferredName string
	UserFingerprint   string
	Message           string
}

type LedgerEntry struct {
	Timestamp time.Time
	ID        uuid.UUID
	Prev      uuid.UUID
	Subject   uuid.UUID
	Content   []byte
}

type User struct {
	ID            uuid.UUID     `db:"id"`
	Email         string        `db:"email" json:"-"`
	PreferredName string        `db:"preferred_name"`
	PublicKey     ssh.PublicKey `db:"public_key"`
	CreatedAt     time.Time     `db:"created_at" json:"-"`
	UpdatedAt     time.Time     `db:"updated_at" json:"-"`
}
