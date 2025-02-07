package service

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/kiper0808/api/internal/gateway/domain"
	"github.com/kiper0808/api/internal/gateway/repository"
	"github.com/kiper0808/api/internal/gateway/service/file_storage"
	"io"
	"mime/multipart"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"

	client "github.com/kiper0808/api/pkg/http"
)

type serviceStorage struct {
	storageRepository repository.Storage
	chunkRepository   repository.Chunk
	httpClient        *client.Client
	logger            *zap.Logger
	fileStorageClient fileStorage.Client
}

const chunks = 6

func newStorageService(httpClient *client.Client,
	logger *zap.Logger,
	storageRepository repository.Storage,
	chunkRepository repository.Chunk,
	fileStorageClient fileStorage.Client,
) *serviceStorage {
	return &serviceStorage{
		storageRepository: storageRepository,
		chunkRepository:   chunkRepository,
		httpClient:        httpClient,
		logger:            logger,
		fileStorageClient: fileStorageClient,
	}
}

type File struct {
	ID uuid.UUID `json:"id"`
}

// UploadFile загружает файл, разделяя его на части
func (s *serviceStorage) UploadFile(ctx context.Context, file *multipart.FileHeader) (*File, error) {
	fileID := uuid.New()

	storages, err := s.getStoragesWithMetrics(ctx)
	if err != nil {
		return nil, fmt.Errorf("cant get storages: %w", err)
	}

	data, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("can't open file: %w", err)
	}
	defer data.Close()

	fileData, err := io.ReadAll(data)
	if err != nil {
		return nil, fmt.Errorf("can't read file: %w", err)
	}

	fileSize := len(fileData)
	partSize := fileSize / chunks
	if partSize == 0 {
		partSize = 1
	}

	var wg sync.WaitGroup
	errCh := make(chan error, chunks)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for i := 0; i < chunks; i++ {
		start := i * partSize
		end := start + partSize
		if i == chunks-1 {
			end = fileSize
		}

		partCopy := append([]byte(nil), fileData[start:end]...)
		hostname := storages[i].Hostname
		wg.Add(1)

		go func(hostname string, partCopy []byte, i int) {
			defer wg.Done()
			chunkID := uuid.New()

			if err := s.fileStorageClient.Upload(ctx, partCopy, hostname, chunkID); err != nil {
				errCh <- fmt.Errorf("cant upload chunk: %w", err)
				cancel()
				return
			}

			err := s.chunkRepository.Create(ctx, &domain.Chunk{
				ID:              chunkID,
				FileID:          fileID,
				Part:            i,
				StorageHostname: hostname,
			})
			if err != nil {
				errCh <- fmt.Errorf("cant save chunk info: %w", err)
				cancel()
				return
			}
		}(hostname, partCopy, i)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}

	return &File{
		ID: fileID,
	}, nil
}

type fileChunk struct {
	ChunkID uuid.UUID `json:"chunk_id"`
	Data    []byte    `json:"data"`
	Part    int       `json:"part"`
}

func (s *serviceStorage) DownloadFile(ctx context.Context, fileID uuid.UUID) ([]byte, error) {
	fileChunks, err := s.chunkRepository.GetAllByFileID(ctx, fileID)
	if err != nil {
		return nil, fmt.Errorf("cant get all chunks: %w", err)
	}

	var chunksWithFiles []*fileChunk
	for _, chunk := range fileChunks {
		data, err := s.fileStorageClient.Download(ctx, chunk.StorageHostname, chunk.ID)
		if err != nil {
			return nil, fmt.Errorf("cant download file: %w", err)
		}
		chunksWithFiles = append(chunksWithFiles, &fileChunk{
			ChunkID: chunk.ID,
			Data:    data,
			Part:    chunk.Part,
		})
	}

	sort.Slice(chunksWithFiles, func(i, j int) bool {
		return chunksWithFiles[i].Part < chunksWithFiles[j].Part
	})

	var resultFile []byte

	for _, chunkWithFile := range chunksWithFiles {
		resultFile = append(resultFile, chunkWithFile.Data...)
	}

	return resultFile, nil
}

type StorageData struct {
	ID        uuid.UUID
	Hostname  string
	FreeBytes float64
	UsedBytes float64
}

func (m *StorageData) UsagePercentage() float64 {
	total := m.FreeBytes + m.UsedBytes
	if total == 0 {
		return 0 // Если диска нет, считаем 0% использования
	}
	return (m.UsedBytes / total) * 100
}

const metricMinioSystemDriveFreeBytes = "minio_system_drive_free_bytes"
const metricMinioSystemDriveUsedBytes = "minio_system_drive_used_bytes"

func (s *serviceStorage) getStoragesWithMetrics(ctx context.Context) ([]StorageData, error) {
	storages, err := s.storageRepository.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all storages err: %w", err)
	}

	var storageData []StorageData
	var wg sync.WaitGroup
	var mu sync.Mutex

	errCh := make(chan error, len(storages))
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, storage := range storages {
		wg.Add(1)
		go func(ctx context.Context, storage domain.Storage) {
			defer wg.Done()
			data, err := s.getMetrics(ctx, &storage)
			if err != nil {
				errCh <- fmt.Errorf("get metrics err: %w", err)
				cancel() // Прерываем остальные горутины
				return
			}
			mu.Lock()
			storageData = append(storageData, *data)
			mu.Unlock()
		}(ctx, storage)
	}

	// Ждём завершения всех горутин
	go func() {
		wg.Wait()
		close(errCh) // Закрываем канал после завершения всех горутин
	}()

	for err := range errCh {
		if err != nil {
			return nil, err // Немедленный выход при первой ошибке
		}
	}

	sort.Slice(storageData, func(i, j int) bool {
		return storageData[i].UsagePercentage() < storageData[j].UsagePercentage()
	})

	if len(storageData) < chunks {
		return nil, fmt.Errorf("not enough storages available")
	}

	return storageData[:chunks], nil
}

func (s *serviceStorage) getMetrics(ctx context.Context, storage *domain.Storage) (*StorageData, error) {
	body, err := s.fileStorageClient.GetMetrics(ctx, storage.Hostname)
	if err != nil {
		return nil, fmt.Errorf("get metrics err: %w", err)
	}

	storageDiskMetrics := &StorageData{
		ID:       storage.ID,
		Hostname: storage.Hostname,
	}

	scanner := bufio.NewScanner(bytes.NewReader(body))
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) == 2 {
			metricName := strings.TrimSpace(parts[0])

			if slices.Contains([]string{metricMinioSystemDriveFreeBytes, metricMinioSystemDriveUsedBytes}, metricName) {
				metricValue, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
				if err != nil {
					return nil, fmt.Errorf("parse float err: %w", err)
				}

				if metricName == metricMinioSystemDriveFreeBytes {
					storageDiskMetrics.FreeBytes = metricValue
				} else if metricName == metricMinioSystemDriveUsedBytes {
					storageDiskMetrics.UsedBytes = metricValue
				}
			}
		}
	}
	return storageDiskMetrics, nil
}
