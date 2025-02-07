package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/kiper0808/s3/internal/domain"
)

type storageRepository struct {
	db *sqlx.DB
}

func newStorageRepository(db *sqlx.DB) *storageRepository {
	return &storageRepository{
		db: db,
	}
}

func (r *storageRepository) Create(ctx context.Context, address *domain.Storage) error {
	return nil
}

func (r *storageRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Storage, error) {
	return &domain.Storage{}, nil
}
