package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tinyhui/GoFile/fileop"
)

func (h *handler) updateFile(c *gin.Context) {
	param, err := parseFileRequest(c)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code": http.StatusBadRequest,
				"msg":  err.Error(),
			},
		)
		return
	}

	fileName := param.FileName
	filePath := param.FilePath
	fileContent := param.FileContent

	if fileName == "" || filePath == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code": http.StatusBadRequest,
				"msg":  "File name or path is not provided",
			},
		)
		return
	}

	fileDir := fileop.PathJoin(h.storageRoot, filePath, fileName)

	err = h.fop.Write(fileDir, fileContent, true)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"code": http.StatusInternalServerError,
				"msg":  fmt.Sprintf("Unable to update file, got error: %s", err.Error()),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"code": http.StatusOK,
			"msg":  "File updated successfully",
		},
	)
}
