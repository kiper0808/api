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
// @Param file formData file true "File to upload"
// @Success 201 {object} uploadFileResponse
// @Failure 400
// @Router /files [post]
// @Security Bearer
func (h *Handler) uploadFile(c *gin.Context) {
	ctx := c.Request.Context()

	file, err := c.FormFile("file")
	if err != nil {
		h.logger.Error("form file failed", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	uploadResponse, err := h.services.Storage.UploadFile(ctx, file)
	if err != nil {
		h.logger.Error("upload file failed", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusCreated, &uploadFileResponse{
		ID: uploadResponse.ID,
	})
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
		h.logger.Error("parse file id failed", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
	}

	data, err := h.services.Storage.DownloadFile(ctx, fileID)
	if err != nil {
		h.logger.Error("cant download file", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
	}

	reader := bytes.NewReader(data)
	c.DataFromReader(http.StatusOK, reader.Size(), "application/octet-stream", reader, nil)
}
