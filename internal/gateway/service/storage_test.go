package service

import (
	"context"
	"github.com/kiper0808/api/internal/gateway/domain"
	clientMocks "github.com/kiper0808/api/internal/gateway/service/file_storage/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestUsagePercentage(t *testing.T) {
	testTable := []struct {
		name      string
		freeBytes float64
		usedBytes float64
		result    float64
	}{
		{
			name:      "OK",
			freeBytes: 5000,
			usedBytes: 3000,
			result:    37.5,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			storageData := &StorageData{
				FreeBytes: testCase.freeBytes,
				UsedBytes: testCase.usedBytes,
			}
			assert.Equal(t, testCase.result, storageData.UsagePercentage())
		})
	}

}

func TestGetMetrics(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		c := gomock.NewController(t)
		defer c.Finish()

		mockClient := clientMocks.NewMockClient(c)

		mockClient.EXPECT().
			GetMetrics(gomock.Any(), "storage1").
			Return([]byte("minio_system_drive_free_bytes 5000\nminio_system_drive_used_bytes 3000"), nil)

		service := &serviceStorage{
			fileStorageClient: mockClient,
		}

		storage := &domain.Storage{
			Hostname: "storage1",
		}

		result, err := service.GetMetrics(context.Background(), storage)

		assert.NoError(t, err)
		assert.Equal(t, &StorageData{
			ID:        storage.ID,
			Hostname:  "storage1",
			FreeBytes: 5000,
			UsedBytes: 3000,
		}, result)
	})
}
