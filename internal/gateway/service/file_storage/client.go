package file_storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/google/uuid"
	"github.com/kiper0808/api/internal/gateway/config"
)

type fileStorageClient struct {
	httpClient *http.Client
	cfg        config.FileStorage
}

//go:generate mockgen -destination=mocks/mock_file_storage.go -package=mocks github.com/kiper0808/api/internal/gateway/service/file_storage Client
type Client interface {
	Upload(ctx context.Context, partReader io.Reader, hostname string, chunkID uuid.UUID) error
	Download(ctx context.Context, hostname string, fileID uuid.UUID, writer io.Writer) error
	GetMetrics(ctx context.Context, hostname string) ([]byte, error)
}

func NewClient(client *http.Client, cfg config.FileStorage) (*fileStorageClient, error) {
	return &fileStorageClient{
		httpClient: client,
		cfg:        cfg,
	}, nil
}

const (
	ApiPath     string = "/api/v1/files"
	MetricsPath string = "/api/v1/metrics"
)

type UploadResponse struct {
	Uuid *uuid.UUID          `json:"uuid,omitempty"`
	Meta *MetaUploadResponse `json:"meta,omitempty"`
}

type MetaUploadResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Upload загружает файл с использованием multipart/form-data
func (c *fileStorageClient) Upload(ctx context.Context, partReader io.Reader, hostname string, chunkID uuid.UUID) error {
	// Используем io.Pipe для стриминга без буферизации в память
	pipeReader, pipeWriter := io.Pipe()
	writer := multipart.NewWriter(pipeWriter)

	// Горутина для записи multipart данных
	go func() {
		defer pipeWriter.Close()
		defer writer.Close()

		// Добавляем файл в multipart
		part, err := writer.CreateFormFile("file", chunkID.String())
		if err != nil {
			pipeWriter.CloseWithError(fmt.Errorf("не удалось создать FormFile: %w", err))
			return
		}

		// Копируем содержимое файла в multipart-часть
		_, err = io.Copy(part, partReader)
		if err != nil {
			pipeWriter.CloseWithError(fmt.Errorf("не удалось записать файл в multipart: %w", err))
			return
		}

		// Добавляем objectName как поле формы
		err = writer.WriteField("objectName", chunkID.String())
		if err != nil {
			pipeWriter.CloseWithError(fmt.Errorf("не удалось записать поле objectName: %w", err))
			return
		}
	}()

	// Создаём HTTP-запрос
	url := "http://" + hostname + ApiPath
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, pipeReader)
	if err != nil {
		return fmt.Errorf("не удалось создать запрос: %w", err)
	}

	// Устанавливаем заголовки
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", c.cfg.Secret)

	// Отправляем запрос
	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(response.Body)
		return fmt.Errorf("wrong status: %d body: %v host: %v", response.StatusCode, string(bodyBytes), hostname)
	}

	// Дочитываем и отбрасываем тело ответа
	io.Copy(io.Discard, response.Body)

	fmt.Println("Файл успешно загружен!")
	return nil
}

func (c *fileStorageClient) Download(ctx context.Context, hostname string, fileID uuid.UUID, writer io.Writer) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+hostname+ApiPath+"/"+fileID.String(), nil)
	if err != nil {
		return fmt.Errorf("cant create request: %w", err)
	}
	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("cant do request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("cant get metrics from %s http code is: %d", hostname, response.StatusCode)
	}

	// Стримим напрямую в writer без загрузки в память
	_, err = io.Copy(writer, response.Body)
	if err != nil {
		return fmt.Errorf("cant copy response body: %w", err)
	}

	return nil
}

func (c *fileStorageClient) GetMetrics(ctx context.Context, hostname string) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+hostname+MetricsPath, nil)
	if err != nil {
		return nil, fmt.Errorf("cant create request: %w", err)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("cant do request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cant get metrics from %s http code is: %d", hostname, response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %w", err)
	}

	return body, nil
}
