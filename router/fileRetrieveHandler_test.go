package router

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type FileRetrieveHandlerTestSuite struct {
	suite.Suite
}

func TestFileRetrieveHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(FileRetrieveHandlerTestSuite))
}

func (suite *FileRetrieveHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func (suite *FileRetrieveHandlerTestSuite) TestRetrieveFile() {
	mockFop := new(mockFileOp)
	mockFop.On("Read", "test/tempFile").Return("Any content can go inside", nil)

	handler := NewHandler("test", mockFop, nil)
	router := InitRouter(handler)

	req, _ := http.NewRequest("GET", "/file?path=tempFile", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	suite.Equal(200, resp.Code)
	suite.Equal("{\"code\":200,\"msg\":\"Any content can go inside\"}", resp.Body.String())
}

func (suite *FileRetrieveHandlerTestSuite) TestRetrieveFile_should_give_500_when_file_not_able_to_read() {
	mockFop := new(mockFileOp)
	mockFop.On("Read", "test/anyPath").Return("", errors.New("any error message"))

	handler := NewHandler("test", mockFop, new(mockFileStatic))
	router := InitRouter(handler)

	req, _ := http.NewRequest("GET", "/file?path=anyPath", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	suite.Equal(500, resp.Code)
	suite.Equal("{\"code\":500,\"msg\":\"Unable to read file, got error: any error message\"}", resp.Body.String())
}

func (suite *FileRetrieveHandlerTestSuite) TestRetrieveFile_should_give_400_when_target_path_not_provided() {
	handler := NewHandler("test", nil, nil)
	router := InitRouter(handler)

	req, _ := http.NewRequest("GET", "/file?path=", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	suite.Equal(400, resp.Code)
	suite.Equal("{\"code\":400,\"msg\":\"You need to provide target path\"}", resp.Body.String())
}