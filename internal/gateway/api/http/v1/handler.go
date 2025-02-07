package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kiper0808/s3/internal/gateway/service"
	"go.uber.org/zap"
)

// @title Karma8 Test FileStorage
// @version 1.0
// @description API for Karma8 Test FileStorage

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
}
