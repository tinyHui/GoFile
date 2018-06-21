package router

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type FileDeleteHandlerTestSuite struct {
	suite.Suite
}

func (suite *FileDeleteHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func TestFileDeleteHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(FileDeleteHandlerTestSuite))
}

func (suite *FileDeleteHandlerTestSuite) TestDeleteFile() {
	mockFop := new(mockFileOp)
	mockFop.On("Delete", "test/tempFile").Return(nil)

	handler := NewHandler("test", mockFop, nil)
	router := InitRouter(handler)

	req, _ := http.NewRequest("DELETE", "/file?path=/tempFile", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	mockFop.AssertCalled(suite.T(), "Delete", "test/tempFile")

	suite.Equal(200, resp.Code)
	suite.Equal("{\"code\":200,\"msg\":\"File deleted\"}", resp.Body.String())
}

func (suite *FileDeleteHandlerTestSuite) TestDeleteFile_should_give_500_when_file_not_exist() {
	mockFop := new(mockFileOp)
	mockFop.On("Delete", "test/anyPath").Return(errors.New("any error"))

	handler := NewHandler("test", mockFop, nil)
	router := InitRouter(handler)

	req, _ := http.NewRequest("DELETE", "/file?path=anyPath", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	suite.Equal(500, resp.Code)
	suite.Equal("{\"code\":500,\"msg\":\"Unable to delete file, got error: any error\"}", resp.Body.String())
}

func (suite *FileDeleteHandlerTestSuite) TestDeleteFile_should_give_400_when_target_path_not_provided() {
	handler := NewHandler("test", nil, nil)
	router := InitRouter(handler)

	req, _ := http.NewRequest("DELETE", "/file?path=", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	suite.Equal(400, resp.Code)
	suite.Equal("{\"code\":400,\"msg\":\"You need to provide target path\"}", resp.Body.String())
}
