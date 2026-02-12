package service

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/kiper0808/api/internal/gateway/config"
	"github.com/kiper0808/api/internal/gateway/domain"
	"github.com/kiper0808/api/internal/gateway/repository"
	"github.com/kiper0808/api/internal/gateway/service/file_storage"
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
	FileStorageClient file_storage.Client
}

//go:generate mockgen -destination=mocks/mock_storage.go -package=mocks github.com/kiper0808/api/internal/gateway/service Storage
type Storage interface {
	AddStorage(ctx context.Context, storage *domain.Storage) error
	UploadFile(ctx context.Context, file *multipart.FileHeader) (*File, error)
	DownloadFile(ctx context.Context, id uuid.UUID, writer io.Writer) error
	getStoragesWithMetrics(ctx context.Context, chunks int) ([]StorageData, error)
	GetMetrics(ctx context.Context, storage *domain.Storage) (*StorageData, error)
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
