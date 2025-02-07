package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
)

func (h *Handler) serviceIdentityMiddleware(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header != "secret" {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
