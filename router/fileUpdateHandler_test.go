package router

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type FileUpdateTestSuite struct {
	suite.Suite
}

func (suite *FileUpdateTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func TestFileUpdateTestSuite(t *testing.T) {
	suite.Run(t, new(FileUpdateTestSuite))
}

func (suite *FileUpdateTestSuite) TestUpdateFile() {
	mockFop := new(mockFileOp)
	mockFop.On(
		"Write",
		"test/test/file_name",
		"content content",
		true,
	).Return(nil)

	handler := NewHandler("test", mockFop, nil)
	router := InitRouter(handler)

	payload := "{\n\t\"filePath\": \"/test\",\n\t\"fileName\": \"file_name\",\n\t\"fileContent\": \"content content\"\n}"

	req, _ := http.NewRequest("PUT", "/file",
		strings.NewReader(payload))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	suite.Equal(200, resp.Code)
	suite.Equal("{\"code\":200,\"msg\":\"File updated successfully\"}", resp.Body.String())
}

func (suite *FileUpdateTestSuite) TestUpdateFile_should_give_500_when_can_not_create_file() {
	mockFop := new(mockFileOp)
	mockFop.On(
		"Write",
		"test/test/file_name",
		"content content",
		true,
	).Return(errors.New("any error"))

	handler := NewHandler("test", mockFop, nil)
	router := InitRouter(handler)

	payload := "{\n\t\"filePath\": \"/test\",\n\t\"fileName\": \"file_name\",\n\t\"fileContent\": \"content content\"\n}"

	req, _ := http.NewRequest("PUT", "/file",
		strings.NewReader(payload))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	suite.Equal(500, resp.Code)
	suite.Equal("{\"code\":500,\"msg\":\"Unable to update file, got error: any error\"}", resp.Body.String())
}

func (suite *FileUpdateTestSuite) TestUpdateFile_should_give_400_when_request_body_not_provided() {
	handler := NewHandler("test", nil, nil)
	router := InitRouter(handler)

	req, _ := http.NewRequest("PUT", "/file",
		strings.NewReader(""))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	suite.Equal(400, resp.Code)
	suite.Equal("{\"code\":400,\"msg\":\"not able to parse request body\"}", resp.Body.String())
}

func (suite *FileUpdateTestSuite) TestUpdateFile_should_give_400_when_target_path_not_provided() {
	handler := NewHandler("test", nil, nil)
	router := InitRouter(handler)

	payload := "{\n\t\"fileName\": \"file_name\",\n\t\"fileContent\": \"content content\"\n}"

	req, _ := http.NewRequest("PUT", "/file",
		strings.NewReader(payload))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	suite.Equal(400, resp.Code)
	suite.Equal("{\"code\":400,\"msg\":\"File name or path is not provided\"}", resp.Body.String())
}
