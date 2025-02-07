package service

import (
	"context"
	"github.com/minio/minio-go/v7"
	"mime/multipart"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/kiper0808/api/internal/storage/service/minio_storage"

	"github.com/kiper0808/api/pkg/http"

	"github.com/kiper0808/api/internal/storage/config"
	"github.com/kiper0808/api/internal/storage/repository"
)

type Services struct {
	Storage Storage
}

type Deps struct {
	Logger            *zap.Logger
	Repos             *repository.Repositories
	HttpClient        *http.Client
	Config            *config.Config
	FileStorageClient minio_storage.Client
}

type Storage interface {
	Metrics(ctx context.Context) ([]byte, error)
	UploadFile(ctx context.Context, file *multipart.FileHeader, objectID uuid.UUID) error
	DownloadFile(ctx context.Context, id uuid.UUID) (*minio.Object, error)
}

func NewServices(deps *Deps) *Services {
	return &Services{
		Storage: newStorageService(deps.HttpClient,
			deps.Logger,
			deps.Repos.Storage,
			deps.FileStorageClient,
			&deps.Config.MinioStorage,
		),
	}
}
