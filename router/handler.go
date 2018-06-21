package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tinyhui/GoFile/fileop"
)

type Handler interface {
	healthCheck(c *gin.Context)
}

type handler struct {
	storageRoot string
	fop         fileop.FileOp
	fstatic     fileop.FileStatic
}

func NewHandler(storageRoot string, fop fileop.FileOp, fstatic fileop.FileStatic) *handler {
	return &handler{
		storageRoot: storageRoot,
		fop:         fop,
		fstatic:     fstatic,
	}
}

func (h *handler) healthCheck(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"code": http.StatusOK,
			"msg":  "OK",
		},
	)
}
