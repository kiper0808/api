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
	fileStorageClient file_storage.Client
}

func newStorageService(httpClient *client.Client,
	logger *zap.Logger,
	storageRepository repository.Storage,
	chunkRepository repository.Chunk,
	fileStorageClient file_storage.Client,
) *serviceStorage {
	return &serviceStorage{
		storageRepository: storageRepository,
		chunkRepository:   chunkRepository,
		httpClient:        httpClient,
		logger:            logger,
		fileStorageClient: fileStorageClient,
	}
}

func (s *serviceStorage) AddStorage(ctx context.Context, storage *domain.Storage) error {
	exist, err := s.storageRepository.IsExist(ctx, storage.Hostname)
	if err != nil {
		return fmt.Errorf("cant check exist storage by hostname: %w", err)
	}
	if exist {
		return ErrStorageAlreadyExists
	}

	v7, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("cant generate uuid: %w", err)
	}
	storage.ID = v7

	return s.storageRepository.Create(ctx, storage)
}

type File struct {
	ID uuid.UUID `json:"id"`
}

// UploadFile загружает файл, разделяя его на части
func (s *serviceStorage) UploadFile(ctx context.Context, file *multipart.FileHeader) (*File, error) {
	fileID := uuid.New()

	// Определяем количество чанков в зависимости от размера файла
	chunks := s.calculateChunksCount(file.Size)

	// Получаем хранилища с метриками
	storages, err := s.getStoragesWithMetrics(ctx, chunks)
	if err != nil {
		return nil, fmt.Errorf("cant get storages: %w", err)
	}

	if len(storages) < chunks {
		return nil, fmt.Errorf("not enough storages available")
	}

	s.logger.Info("chunks", zap.Any("chunks", chunks), zap.Any("file_size", file.Size))

	// Определяем размер чанка
	partSize := (file.Size + int64(chunks) - 1) / int64(chunks)
	if partSize < 1 {
		partSize = 1
	}

	// Открываем файл, чтобы получить доступ к его содержимому
	fileReader, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("cant open file: %w", err)
	}
	defer fileReader.Close()

	var wg sync.WaitGroup
	errCh := make(chan error, chunks)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Читаем чанки и загружаем их
	for i := 0; i < chunks; i++ {
		if i >= len(storages) {
			return nil, fmt.Errorf("not enough storages available")
		}

		hostname := storages[i].Hostname
		chunkID := uuid.New()
		partReader := io.LimitReader(fileReader, partSize) // Ограничиваем чтение чанка

		wg.Add(1)
		go func(hostname string, partReader io.Reader, chunkID uuid.UUID, i int) {
			defer wg.Done()

			// Загружаем чанк
			if err := s.fileStorageClient.Upload(ctx, partReader, hostname, chunkID); err != nil {
				errCh <- fmt.Errorf("can't upload chunk: %w", err)
				cancel()
				return
			}

			// Сохраняем метаданные в БД
			if err := s.chunkRepository.Create(ctx, &domain.Chunk{
				ID:              chunkID,
				FileID:          fileID,
				Part:            i,
				StorageHostname: hostname,
			}); err != nil {
				errCh <- fmt.Errorf("can't save chunk info: %w", err)
				cancel()
				return
			}
		}(hostname, partReader, chunkID, i)
	}

	// Ждём завершения всех горутин
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// Обрабатываем ошибки
	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}

	return &File{ID: fileID}, nil
}

func (s *serviceStorage) calculateChunksCount(fileSize int64) int {
	// Расчет количества чанков на основе размера файла
	if fileSize <= 1*1024*1024*1024 { // 1 GB
		return 6
	} else if fileSize <= 3*1024*1024*1024 { // 3 GB
		return 7
	} else if fileSize <= 5*1024*1024*1024 { // 5 GB
		return 8
	} else if fileSize <= 7*1024*1024*1024 { // 7 GB
		return 9
	}

	return 10
}

type fileChunk struct {
	ChunkID uuid.UUID `json:"chunk_id"`
	Data    []byte    `json:"data"`
	Part    int       `json:"part"`
}

func (s *serviceStorage) DownloadFile(ctx context.Context, fileID uuid.UUID, writer io.Writer) error {
	// Получаем все чанки для данного файла
	fileChunks, err := s.chunkRepository.GetAllByFileID(ctx, fileID)
	if err != nil {
		return fmt.Errorf("can't get all chunks: %w", err)
	}

	// Сортируем чанки по порядку (Part)
	sort.Slice(fileChunks, func(i, j int) bool {
		return fileChunks[i].Part < fileChunks[j].Part
	})

	// Скачиваем чанки последовательно с prefetch
	// Начинаем скачивать chunk N+1 пока пишем chunk N
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	type chunkPipe struct {
		pr  *io.PipeReader
		err error
	}

	// Функция для скачивания чанка в pipe
	downloadChunk := func(chunk domain.Chunk) chunkPipe {
		pr, pw := io.Pipe()
		go func() {
			defer pw.Close()
			err := s.fileStorageClient.Download(ctx, chunk.StorageHostname, chunk.ID, pw)
			if err != nil {
				pw.CloseWithError(err)
			}
		}()
		return chunkPipe{pr: pr, err: nil}
	}

	// Начинаем скачивать первый чанк
	var currentChunk *chunkPipe
	if len(fileChunks) > 0 {
		cp := downloadChunk(fileChunks[0])
		currentChunk = &cp
	}

	// Обрабатываем чанки по порядку с prefetch
	for i := 0; i < len(fileChunks); i++ {
		// Начинаем скачивать следующий чанк (prefetch)
		var nextChunk *chunkPipe
		if i+1 < len(fileChunks) {
			cp := downloadChunk(fileChunks[i+1])
			nextChunk = &cp
		}

		// Пишем текущий чанк в writer
		_, err := io.Copy(writer, currentChunk.pr)
		currentChunk.pr.Close()
		if err != nil {
			cancel()
			return fmt.Errorf("can't write chunk to response (part: %d): %w", i, err)
		}

		// Переходим к следующему чанку
		currentChunk = nextChunk
	}

	return nil
}

type FileStorageData interface {
	UsagePercentage() float64
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
		return 0
	}
	return (m.UsedBytes / total) * 100
}

const metricMinioSystemDriveFreeBytes = "minio_system_drive_free_bytes"
const metricMinioSystemDriveUsedBytes = "minio_system_drive_used_bytes"

func (s *serviceStorage) getStoragesWithMetrics(ctx context.Context, chunks int) ([]StorageData, error) {
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
			data, err := s.GetMetrics(ctx, &storage)
			if err != nil {
				errCh <- fmt.Errorf("get metrics err: %w", err)
				cancel()
				return
			}
			mu.Lock()
			storageData = append(storageData, *data)
			mu.Unlock()
		}(ctx, storage)
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

	sort.Slice(storageData, func(i, j int) bool {
		return storageData[i].UsagePercentage() < storageData[j].UsagePercentage()
	})

	if len(storageData) < 6 {
		return nil, fmt.Errorf("not enough storages available")
	}

	return storageData[:chunks], nil
}

func (s *serviceStorage) GetMetrics(ctx context.Context, storage *domain.Storage) (*StorageData, error) {
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
