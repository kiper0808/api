package v1

import (
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initMetricsRoutes(api *gin.RouterGroup) {
	api.GET("/metrics", h.metrics)
}

// @Summary Метрики хранилища
// @Tags metrics
// @Description Метрики хранилища
// @ModuleID metrics
// @Produce  json
// @Success 200
// @Failure 400
// @Router /metrics [get]
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
