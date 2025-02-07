package domain

import (
	"time"

	"github.com/google/uuid"
)

type Chunk struct {
	ID              uuid.UUID `json:"id" db:"id"`
	FileID          uuid.UUID `json:"file_id" db:"file_id"`
	Part            int       `json:"part" db:"part"`
	StorageHostname string    `json:"storage_hostname" db:"storage_hostname"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}
