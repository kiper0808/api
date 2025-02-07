package v1

import (
	"fmt"
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
// @Failure 400
// @Router /files [post]
// @Security Bearer
func (h *Handler) uploadFile(c *gin.Context) {
	ctx := c.Request.Context()

	// Получаем файл
	file, err := c.FormFile("file")
	if err != nil {
		h.logger.Error("failed to get file", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	objectName := c.PostForm("objectName")
	objectID, err := uuid.Parse(objectName)
	if err != nil {
		h.logger.Error("failed to parse uuid", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = h.services.Storage.UploadFile(ctx, file, objectID)
	if err != nil {
		h.logger.Error("failed to upload file", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
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
// @Failure 400
// @Router /files/{id} [get]
// @Security Bearer
func (h *Handler) downloadFile(c *gin.Context) {
	ctx := c.Request.Context()
	fileIDStr := c.Param("id")
	fileID, err := uuid.Parse(fileIDStr)
	if err != nil {
		h.logger.Error("failed to parse uuid", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	object, err := h.services.Storage.DownloadFile(ctx, fileID)
	if err != nil {
		h.logger.Error("cant download file", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	defer object.Close()

	objectStat, err := object.Stat()
	if err != nil {
		h.logger.Error("cant download file", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Устанавливаем заголовки для скачивания
	c.Header("Content-Disposition", "attachment; filename="+"test")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", fmt.Sprintf("%d", objectStat.Size))

	// Отправляем файл потоково
	c.DataFromReader(http.StatusOK, objectStat.Size, "application/octet-stream", object, nil)
}
