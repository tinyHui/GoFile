package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter(handler *handler) http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())
	e.GET("/healthcheck", handler.healthCheck)
	e.POST("/file", handler.createFile)
	e.PUT("/file", handler.updateFile)
	e.GET("/file", handler.retrieveFile)
	e.DELETE("/file", handler.deleteFile)
	e.GET("/static", handler.getStatic)

	return e
}

