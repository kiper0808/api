package v1

import (
	"bytes"
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) initStorageRoutes(api *gin.RouterGroup) {
	api.POST("/files", h.serviceIdentityMiddleware, h.uploadFile)
	api.GET("/files/:id", h.downloadFile)
}

type uploadFileResponse struct {
	ID uuid.UUID `json:"id" binding:"required"`
}

// @Summary Загрузка файла
// @Tags files
// @Description Загрузка файла
// @ModuleID files
// @Accept multipart/form-data
// @Produce  json
// @Param X-Idempotency-Key header string true "Idempotency Key"
// @Param file formData file true "File to upload"
// @Success 201 {object} uploadFileResponse
// @Failure 400 {object} ErrorStruct
// @Router /files [post]
// @Security Bearer
func (h *Handler) uploadFile(c *gin.Context) {
	ctx := c.Request.Context()

	// Получаем ключ идемпотентности
	idempotencyKey := c.GetHeader("X-Idempotency-Key")
	if idempotencyKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing idempotency key"})
		return
	}

	// Проверка, использовался ли этот ключ (эмуляция)
	// В реальном приложении здесь должна быть проверка в БД или Redis
	if idempotencyKey == "used-key" {
		c.JSON(http.StatusConflict, gin.H{"error": "Duplicate request"})
		return
	}

	// Получаем файл
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.services.Storage.UploadFile(ctx, file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// @Summary Скачивание файла
// @Tags files
// @Description Скачивание файла
// @ModuleID files
// @Accept  json
// @Produce  json
// @Param	id	path		string		true	"ID файла"
// @Success 200
// @Failure 400 {object} ErrorStruct
// @Router /files/{id} [get]
// @Security Bearer
func (h *Handler) downloadFile(c *gin.Context) {
	ctx := c.Request.Context()
	fileIDStr := c.Param("id")
	fileID, err := uuid.Parse(fileIDStr)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	data, err := h.services.Storage.DownloadFile(ctx, fileID)
	if err != nil {
		h.logger.Error("cant download file", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
	}

	reader := bytes.NewReader(data)

	// Отправляем файл потоково
	c.DataFromReader(http.StatusOK, reader.Size(), "application/octet-stream", reader, nil)
}
