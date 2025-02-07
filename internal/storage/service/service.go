package service

import (
	"context"
	"github.com/minio/minio-go/v7"
	"mime/multipart"

	"github.com/google/uuid"
	"go.uber.org/zap"

	fileStorage "github.com/kiper0808/s3/internal/service/minio_storage"

	"github.com/kiper0808/s3/pkg/http"

	"github.com/kiper0808/s3/internal/config"
	"github.com/kiper0808/s3/internal/repository"
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
