package router

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tinyhui/GoFile/fileop"
)

func (h *handler) createFile(c *gin.Context) {
	form, err := parseUploadFileRequest(c)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code": http.StatusBadRequest,
				"msg":  fmt.Sprintf("Get file err: %s", err.Error()),
			},
		)
		return
	}

	filePath := form.filePath
	file := form.file

	folderDir := fileop.PathJoin(h.storageRoot, filePath)
	exists, err := h.fop.Exists(folderDir)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code": http.StatusBadRequest,
				"msg":  fmt.Sprintf("Not able to create folder: %s", err.Error()),
			},
		)
		return
	}

	if !exists {
		os.MkdirAll(folderDir, os.ModeDir|os.ModePerm)
	}

	filePath = fileop.PathJoin(h.storageRoot, filePath, file.Filename)

	exists, _ = h.fop.Exists(filePath)
	if exists {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code": http.StatusBadRequest,
				"msg":  "File already exist",
			},
		)
		return
	}

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code": http.StatusBadRequest,
				"msg":  fmt.Sprintf("Upload file err: %s", err.Error()),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"code": http.StatusOK,
			"msg":  fmt.Sprintf("File %s uploaded successfully", file.Filename),
		},
	)
}
