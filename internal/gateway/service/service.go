package service

import (
	"context"
	"github.com/kiper0808/api/internal/gateway/config"
	"github.com/kiper0808/api/internal/gateway/domain"
	"github.com/kiper0808/api/internal/gateway/repository"
	"github.com/kiper0808/api/internal/gateway/service/file_storage"
	"mime/multipart"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/kiper0808/api/pkg/http"
)

type Services struct {
	Storage Storage
}

type Deps struct {
	Logger            *zap.Logger
	Repos             *repository.Repositories
	HttpClient        *http.Client
	Config            *config.Config
	FileStorageClient fileStorage.Client
}

type Storage interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader) (*File, error)
	DownloadFile(ctx context.Context, id uuid.UUID) ([]byte, error)
	getStoragesWithMetrics(ctx context.Context) ([]StorageData, error)
	getMetrics(ctx context.Context, storage *domain.Storage) (*StorageData, error)
}

func NewServices(deps *Deps) *Services {
	return &Services{
		Storage: newStorageService(deps.HttpClient,
			deps.Logger,
			deps.Repos.Storage,
			deps.Repos.Chunk,
			deps.FileStorageClient,
		),
	}
}
