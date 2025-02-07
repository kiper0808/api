package http

import (
	"github.com/kiper0808/api/internal/gateway/api/http/v1"
	"github.com/kiper0808/api/internal/gateway/config"
	"github.com/kiper0808/api/internal/gateway/service"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	_ "github.com/kiper0808/api/docs"
)

type Handler struct {
	services   *service.Services
	logger     *zap.Logger
	apiVersion string
}

func NewHandler(services *service.Services,
	logger *zap.Logger,
	apiVersion string,
) *Handler {
	return &Handler{
		services:   services,
		logger:     logger,
		apiVersion: apiVersion,
	}
}

func (h *Handler) Init(cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(
		gin.Recovery(),
		ginzap.Ginzap(h.logger, time.RFC3339, true),
	)

	if cfg.Server.SwaggerEnabled {
		router.GET("/swagger/gateway/*any", ginSwagger.WrapHandler(swaggerfiles.NewHandler(), ginSwagger.InstanceName("gateway")))
	}

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.services, h.logger, h.apiVersion)
	api := router.Group("/api")
	handlerV1.Init(api)
}
