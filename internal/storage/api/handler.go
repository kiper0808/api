package http

import (
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	v1 "github.com/kiper0808/api/internal/storage/api/http/v1"

	_ "github.com/kiper0808/api/docs"
	"github.com/kiper0808/api/internal/storage/config"
	"github.com/kiper0808/api/internal/storage/service"
)

type Handler struct {
	services        *service.Services
	logger          *zap.Logger
	minioStorageCfg *config.MinioStorage
	apiVersion      string
}

func NewHandler(services *service.Services,
	logger *zap.Logger,
	minioStorageCfg *config.MinioStorage,
	apiVersion string,
) *Handler {
	return &Handler{
		services:        services,
		logger:          logger,
		minioStorageCfg: minioStorageCfg,
		apiVersion:      apiVersion,
	}
}

func (h *Handler) Init(cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(
		gin.Recovery(),
		ginzap.Ginzap(h.logger, time.RFC3339, true),
	)

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.services, h.logger, h.apiVersion)
	api := router.Group("/api")
	handlerV1.Init(api)
}
