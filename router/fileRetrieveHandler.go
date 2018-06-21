package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tinyhui/GoFile/fileop"
)

func (h *handler) retrieveFile(c *gin.Context) {
	q := c.Request.URL.Query()
	targetPath, err := parseTargetPath(q)

	if targetPath == "" || err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code": http.StatusBadRequest,
				"msg":  "You need to provide target path",
			},
		)
		return
	}

	targetPath = fileop.PathJoin(h.storageRoot, targetPath)

	content, err := h.fop.Read(targetPath)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"code": http.StatusInternalServerError,
				"msg":  fmt.Sprintf("Unable to read file, got error: %s", err.Error()),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"code": http.StatusOK,
			"msg":  content,
		},
	)
}
