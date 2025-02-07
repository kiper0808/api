package minio_storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/kiper0808/api/internal/storage/config"
)

type minioStorageClient struct {
	httpClient  *http.Client
	minioClient *minio.Client
	config      *config.MinioStorage
}

type Client interface {
	GetMetrics(ctx context.Context) ([]byte, error)
	Upload(ctx context.Context, file *multipart.FileHeader, objectID uuid.UUID) error
	Download(ctx context.Context, fileID uuid.UUID) (*minio.Object, error)
}

func NewClient(cfg *config.MinioStorage, httpClient *http.Client) (*minioStorageClient, error) {
	minioClient, err := minio.New(cfg.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("cant create minio client: %w", err)
	}

	return &minioStorageClient{
		httpClient:  httpClient,
		minioClient: minioClient,
		config:      cfg,
	}, nil
}

func (c *minioStorageClient) Upload(ctx context.Context, fileHeader *multipart.FileHeader, objectID uuid.UUID) error {
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("cant read file: %w", err)
	}

	defer file.Close()

	contentType := "application/octet-stream"

	_, err = c.minioClient.PutObject(ctx, c.config.Bucket, objectID.String(), file, fileHeader.Size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return fmt.Errorf("minio client fput object err: %w", err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectID.String(), fileHeader.Size)
	return nil
}

func (c *minioStorageClient) Download(ctx context.Context, fileID uuid.UUID) (*minio.Object, error) {
	object, err := c.minioClient.GetObject(ctx, c.config.Bucket, fileID.String(), minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("cant get object from minio: %w", err)
	}
	return object, nil
}

func (c *minioStorageClient) GetMetrics(ctx context.Context) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+c.config.Host+"/minio/metrics/v3/system", nil)
	if err != nil {
		return nil, fmt.Errorf("cant create request: %w", err)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("cant do request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cant get metrics from %s http code is: %d", c.config.Host, response.StatusCode)
	}

	// Чтение тела ответа в строку
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %w", err)
	}

	return body, nil
}
