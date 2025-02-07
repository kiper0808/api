package repository

import (
	"context"
	"github.com/kiper0808/s3/internal/gateway/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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

func (r *storageRepository) GetAll(ctx context.Context) ([]domain.Storage, error) {
	var storages []domain.Storage
	if err := r.db.SelectContext(ctx, &storages, "select id, hostname, created_at from storage"); err != nil {
		return nil, err
	}
	return storages, nil
}
