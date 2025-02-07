package v1

import (
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initMetricsRoutes(api *gin.RouterGroup) {
	api.GET("/metrics", h.metrics)
}

// @Summary Статус хранилища
// @Tags metrics
// @Description Статус хранилища
// @ModuleID metrics
// @Accept
// @Produce  json
// @Param X-Idempotency-Key header string true "Idempotency Key"
// @Param file formData file true "File to upload"
// @Success 200
// @Failure 400 {object} ErrorStruct
// @Router /metrics [get]
// @Security Bearer
func (h *Handler) metrics(c *gin.Context) {
	ctx := c.Request.Context()

	data, err := h.services.Storage.Metrics(ctx)
	if err != nil {
		h.logger.Error("storage metrics", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, string(data))
}
