package repository

import (
	"context"
	"fmt"
	"github.com/kiper0808/s3/internal/gateway/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type chunkRepository struct {
	db *sqlx.DB
}

func newChunkRepository(db *sqlx.DB) *chunkRepository {
	return &chunkRepository{
		db: db,
	}
}

func (r *chunkRepository) Create(ctx context.Context, chunk *domain.Chunk) error {
	result, err := r.db.ExecContext(ctx, "INSERT INTO chunk (id, file_id, part, storage_hostname) VALUES (uuid_to_bin(?), uuid_to_bin(?), ?, ?)", chunk.ID, chunk.FileID, chunk.Part, chunk.StorageHostname)
	if err != nil {
		return fmt.Errorf("insert chunk err: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected err: %w", err)
	}
	if affected != 1 {
		return fmt.Errorf("expected 1 row affected, got %d", affected)
	}

	return nil
}

func (r *chunkRepository) GetAllByFileID(ctx context.Context, fileID uuid.UUID) ([]domain.Chunk, error) {
	var chunks []domain.Chunk
	if err := r.db.SelectContext(ctx, &chunks, "select id, file_id, part, storage_hostname, created_at from chunk where file_id = uuid_to_bin(?)", fileID); err != nil {
		return nil, fmt.Errorf("select err: %w", err)
	}
	return chunks, nil
}
