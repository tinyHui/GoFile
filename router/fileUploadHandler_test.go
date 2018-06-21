package router

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"github.com/tinyhui/GoFile/fileop"
)

type FileUploadHandlerTestSuite struct {
	suite.Suite
	fop fileop.FileOp
}

func (suite *FileUploadHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.fop = fileop.NewFileOp()
}

func TestFileUploadHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(FileUploadHandlerTestSuite))
}

func (suite *FileUploadHandlerTestSuite) TestCreateFile_should_create_file_in_exist_folder() {
	folderDir := "test"
	os.Mkdir(folderDir, os.ModeDir|os.ModePerm)
	defer os.RemoveAll("test")

	os.Mkdir("test/exist", os.ModeDir|os.ModePerm)
	defer func() {
		os.Remove("test/exist/tempFile")
		os.Remove("test/exist")
	}()

	handler := NewHandler(folderDir, suite.fop, nil)
	router := InitRouter(handler)

	os.Create("tempFile")
	defer os.Remove("tempFile")

	payload := "------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"file\"; filename=\"tempFile\"\r\nContent-Type: false\r\n\r\n\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"path\"\r\n\r\n/exist\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW--"

	req, _ := http.NewRequest("POST", "/file",
		strings.NewReader(payload))
	req.Header.Add("content-type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	suite.Equal(200, resp.Code)
	suite.Equal("{\"code\":200,\"msg\":\"File tempFile uploaded successfully\"}", resp.Body.String())

	uploaded, _ := suite.fop.Exists("tempFile")
	suite.True(uploaded)
}

func (suite *FileUploadHandlerTestSuite) TestCreateFile_should_create_file_and_folder_when_target_folder_not_exist() {
	folderDir := "test"
	os.Mkdir(folderDir, os.ModeDir|os.ModePerm)
	defer os.RemoveAll("test")

	handler := NewHandler(folderDir, suite.fop, nil)
	router := InitRouter(handler)

	os.Create("tempFile")
	defer os.Remove("tempFile")

	payload := "------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"file\"; filename=\"tempFile\"\r\nContent-Type: false\r\n\r\n\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"path\"\r\n\r\n/notexist\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW--"

	req, _ := http.NewRequest("POST", "/file",
		strings.NewReader(payload))
	req.Header.Add("content-type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	suite.Equal(200, resp.Code)
	suite.Equal("{\"code\":200,\"msg\":\"File tempFile uploaded successfully\"}", resp.Body.String())

	uploaded, _ := suite.fop.Exists("tempFile")
	suite.True(uploaded)

	folderCreated, _ := suite.fop.Exists("test/notexist")
	suite.True(folderCreated)
}

func (suite *FileUploadHandlerTestSuite) TestCreateFile_should_give_400_when_request_wrong() {
	handler := NewHandler("", suite.fop, nil)
	router := InitRouter(handler)

	req, _ := http.NewRequest("POST", "/file",
		strings.NewReader(""))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	suite.Equal(400, resp.Code)
	suite.Equal("{\"code\":400,\"msg\":\"Get file err: request Content-Type isn't multipart/form-data\"}", resp.Body.String())
}

func (suite *FileUploadHandlerTestSuite) TestCreateFile_should_give_400_when_file_already_exist() {
	folderDir := "test"
	os.Mkdir(folderDir, os.ModeDir|os.ModePerm)
	defer os.RemoveAll(folderDir)

	handler := NewHandler(folderDir, suite.fop, nil)
	router := InitRouter(handler)

	os.Create("tempFile")
	defer os.Remove("tempFile")

	payload := "------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"file\"; filename=\"tempFile\"\r\nContent-Type: false\r\n\r\n\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW--"

	req, _ := http.NewRequest("POST", "/file",
		strings.NewReader(payload))
	req.Header.Add("content-type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	suite.Equal(400, resp.Code)
	suite.Equal("{\"code\":400,\"msg\":\"File already exist\"}", resp.Body.String())
}

func (suite *FileUploadHandlerTestSuite) TestCreateFile_should_give_400_when_not_able_to_create_folder() {
	folderDir := "test"
	os.Mkdir(folderDir, os.ModeDir)
	defer os.RemoveAll("test")

	handler := NewHandler(folderDir, suite.fop, nil)
	router := InitRouter(handler)

	os.Create("tempFile")
	defer os.Remove("tempFile")

	payload := "------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"path\"\r\n\r\nanyfolder\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"file\"; filename=\"tempFile\"\r\nContent-Type: false\r\n\r\n\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW--"

	req, _ := http.NewRequest("POST", "/file",
		strings.NewReader(payload))
	req.Header.Add("content-type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	suite.Equal(400, resp.Code)
	suite.Equal("{\"code\":400,\"msg\":\"Not able to create folder: stat test/anyfolder: permission denied\"}", resp.Body.String())
}
