package domain

import (
	"time"

	"github.com/google/uuid"
)

type Storage struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Hostname  string    `json:"hostname" db:"hostname"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
