package api

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func (r *gzipWriter) Write(b []byte) (int, error) {
	return r.writer.Write(b)
}

func (a *API) gzipMiddleware(c *gin.Context) {
	if strings.Contains(c.Request.Header.Get("Content-Encoding"), "gzip") {
		gz, err := gzip.NewReader(c.Request.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		defer gz.Close()
		c.Request.Body = gz
	}
	if strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
		c.Header("Content-Encoding", "gzip")
		gz := gzip.NewWriter(c.Writer)
		defer gz.Close()
		c.Writer = &gzipWriter{
			ResponseWriter: c.Writer,
			writer:         gz,
		}
	}
	c.Next()
}

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
