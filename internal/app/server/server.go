package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func New(handlers *gin.Engine) (*http.Server, error) {
	return &http.Server{
		Addr:    CFG.Host,
		Handler: handlers,
	}, nil
}
