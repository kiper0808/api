package healthcheck

import (
	"github.com/gin-gonic/gin"
	"github.com/heptiolabs/healthcheck"
)

type HealthCheck struct {
	healthcheck.Handler
}

func NewHealthCheck() *HealthCheck {
	healthCheck := &HealthCheck{
		Handler: healthcheck.NewHandler(),
	}
	baseChecker := NewBaseChecker()
	baseChecker.Start()
	healthCheck.Handler.AddReadinessCheck("traffic", baseChecker.Check)

	return healthCheck
}

func (hc *HealthCheck) LiveHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		hc.LiveEndpoint(c.Writer, c.Request)
	}
}

func (hc *HealthCheck) ReadyHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		hc.ReadyEndpoint(c.Writer, c.Request)
	}
}
