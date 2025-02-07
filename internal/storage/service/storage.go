package service

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"mime/multipart"

	"github.com/kiper0808/api/internal/storage/service/minio_storage"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/kiper0808/api/internal/storage/config"
	"github.com/kiper0808/api/internal/storage/repository"
	client "github.com/kiper0808/api/pkg/http"
)

type serviceStorage struct {
	storageRepository  repository.Storage
	httpClient         *client.Client
	logger             *zap.Logger
	minioStorageClient minio_storage.Client
	minioStorageConfig *config.MinioStorage
}

func newStorageService(httpClient *client.Client,
	logger *zap.Logger,
	storageRepository repository.Storage,
	minioStorageClient minio_storage.Client,
	minioStorageConfig *config.MinioStorage,
) *serviceStorage {
	return &serviceStorage{
		storageRepository:  storageRepository,
		httpClient:         httpClient,
		logger:             logger,
		minioStorageClient: minioStorageClient,
		minioStorageConfig: minioStorageConfig,
	}
}

func (s *serviceStorage) UploadFile(ctx context.Context, file *multipart.FileHeader, objectID uuid.UUID) error {
	err := s.minioStorageClient.Upload(ctx, file, objectID)
	if err != nil {
		return fmt.Errorf("file storage client upload err: %w", err)
	}
	return nil
}

func (s *serviceStorage) DownloadFile(ctx context.Context, fileID uuid.UUID) (*minio.Object, error) {
	//request, err := http.NewRequestWithContext(ctx, http.MethodGet, pdfUrl, http.NoBody)
	//if err != nil {
	//	return nil, fmt.Errorf("create request failed: %w", err)
	//}
	//
	//response, err := http.DefaultClient.Do(request)
	//if err != nil {
	//	return nil, fmt.Errorf("do request failed: %w", err)
	//}
	//
	//defer response.Body.Close()
	//
	//pdfBytes, err := io.ReadAll(response.Body)
	//if err != nil {
	//	return nil, fmt.Errorf("read response body failed: %w", err)
	//}
	//
	//fileUUID, err := s.fileStorageClient.Upload(ctx, pdfBytes)
	//if err != nil {
	//	metrics.IncFileStorageUploadCounter(metrics.StatusFailed)
	//	return nil, fmt.Errorf("upload to file storage failed: %w", err)
	//}
	//
	//metrics.IncFileStorageUploadCounter(metrics.StatusSuccessful)
	//return fileUUID, nil
	object, err := s.minioStorageClient.Download(ctx, fileID)
	if err != nil {
		return nil, fmt.Errorf("file storage client download err: %w", err)
	}

	return object, nil
}

func (s *serviceStorage) Metrics(ctx context.Context) ([]byte, error) {
	return s.minioStorageClient.GetMetrics(ctx)
}
