package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tinyhui/GoFile/fileop"
)

func (h *handler) getStatic(c *gin.Context) {
	q := c.Request.URL.Query()
	targetPath, err := parseTargetPath(q)

	if targetPath == "" || err != nil {
		targetPath = h.storageRoot
	} else {
		targetPath = fileop.PathJoin(h.storageRoot, targetPath)
	}

	result := h.fstatic.GetStaticResult(targetPath)

	c.JSON(
		http.StatusOK,
		result,
	)
}
