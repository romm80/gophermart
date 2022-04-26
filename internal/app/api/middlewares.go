package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *API) authMiddleware(c *gin.Context) {
	token := c.GetHeader(authHeader)
	if token == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	login, err := a.Services.ParseToken(token)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Set("user", login)
	c.Next()
}
