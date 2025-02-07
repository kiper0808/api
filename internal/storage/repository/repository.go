package repository

import (
	"context"
	"database/sql"

	"go.uber.org/zap"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/kiper0808/api/internal/storage/domain"
)

type Repositories struct {
	Db
	Storage Storage
}

func NewRepositories(db *sqlx.DB, logger *zap.Logger) *Repositories {
	return &Repositories{
		Storage: newStorageRepository(db),
	}
}

type Db interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
}

type Storage interface {
	Create(ctx context.Context, storage *domain.Storage) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Storage, error)
}
