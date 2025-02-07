package repository

import (
	"context"
	"database/sql"
	domain2 "github.com/kiper0808/s3/internal/gateway/domain"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repositories struct {
	Db
	Storage Storage
	Chunk   Chunk
}

func NewRepositories(db *sqlx.DB, redisClient redis.UniversalClient, logger *zap.Logger) *Repositories {
	return &Repositories{
		Storage: newStorageRepository(db),
		Chunk:   newChunkRepository(db),
	}
}

type Db interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
}

type Storage interface {
	Create(ctx context.Context, storage *domain2.Storage) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain2.Storage, error)
	GetAll(ctx context.Context) ([]domain2.Storage, error)
}

type Chunk interface {
	Create(ctx context.Context, storage *domain2.Chunk) error
	GetAllByFileID(ctx context.Context, fileID uuid.UUID) ([]domain2.Chunk, error)
}
