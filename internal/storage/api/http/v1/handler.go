package v1

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/kiper0808/api/internal/storage/service"
)

// @title Karma8 Storage Service
// @version 1.0
// @description API for Karma8 StorageService

// @BasePath /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

type Handler struct {
	services   *service.Services
	logger     *zap.Logger
	apiVersion string
}

func NewHandler(services *service.Services, logger *zap.Logger, apiVersion string) *Handler {
	return &Handler{
		services:   services,
		logger:     logger,
		apiVersion: apiVersion,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	h.initStorageRoutes(v1)
	h.initMetricsRoutes(v1)
}
