package router

import (
	"errors"
	"mime/multipart"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

type uploadFileForm struct {
	filePath string
	file     *multipart.FileHeader
}

func parseUploadFileRequest(c *gin.Context) (*uploadFileForm, error) {
	filePath := strings.TrimSpace(c.PostForm("path"))
	file, err := c.FormFile("file")

	if err != nil {
		return nil, err
	}

	return &uploadFileForm{
		filePath: filePath,
		file:     file,
	}, nil
}

type fileRequestParam struct {
	FilePath    string `json:"filePath"`
	FileName    string `json:"fileName"`
	FileContent string `json:"fileContent"`
}

func parseFileRequest(c *gin.Context) (*fileRequestParam, error) {
	var param *fileRequestParam
	err := c.BindJSON(&param)
	if err != nil {
		return nil, errors.New("not able to parse request body")
	}

	return param, err
}

func parseTargetPath(q url.Values) (string, error) {
	paths := q["path"]
	if len(paths) < 1 {
		return "", errors.New("target path parameter not provided")
	}

	return strings.TrimSpace(paths[0]), nil
}