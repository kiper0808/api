package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kiper0808/api/internal/gateway/domain"

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
	v7, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("cant create uuid for storage: %w", err)
	}
	_, err = r.db.ExecContext(ctx, "insert into storage (id, hostname) values (uuid_to_bin(?), ?)",
		v7, address.Hostname)
	if err != nil {
		return fmt.Errorf("create storage: %w", err)
	}
	return nil

}

func (r *storageRepository) GetAll(ctx context.Context) ([]domain.Storage, error) {
	var storages []domain.Storage
	if err := r.db.SelectContext(ctx, &storages, "select id, hostname, created_at from storage"); err != nil {
		return nil, err
	}
	return storages, nil
}

func (r *storageRepository) GetByHostname(ctx context.Context, hostname string) (*domain.Storage, error) {
	var storage domain.Storage
	if err := r.db.GetContext(ctx, &storage, "select id, hostname, created_at from storage where hostname = ?", hostname); err != nil {
		return nil, fmt.Errorf("get storage by hostname: %w", err)
	}
	return &storage, nil
}

func (r *storageRepository) IsExist(ctx context.Context, hostname string) (bool, error) {
	var count int
	if err := r.db.GetContext(ctx, &count, "select count(*) from storage where hostname = ?", hostname); err != nil {
		return false, fmt.Errorf("check storage exist: %w", err)
	}
	return count > 0, nil
}
