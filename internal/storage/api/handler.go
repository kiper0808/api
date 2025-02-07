package http

import (
	"time"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	v1 "github.com/kiper0808/s3/internal/api/http/v1"

	_ "github.com/kiper0808/s3/docs"
	"github.com/kiper0808/s3/internal/config"
	"github.com/kiper0808/s3/internal/service"
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
		sentrygin.New(sentrygin.Options{}),
		otelgin.Middleware(cfg.EnvName),
	)

	if cfg.Server.SwaggerEnabled {
		router.GET("/swagger/api/*any", ginSwagger.WrapHandler(swaggerfiles.NewHandler(), ginSwagger.InstanceName("api")))
	}

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.services, h.logger, h.apiVersion)
	api := router.Group("/api")
	handlerV1.Init(api)
}
