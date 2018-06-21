package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"github.com/tinyhui/GoFile/fileop"
)

type FileStaticHandlerTestSuite struct {
	suite.Suite
}

func (suite *FileStaticHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func TestFileStaticHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(FileStaticHandlerTestSuite))
}

func (suite *FileStaticHandlerTestSuite) TestGetStatic() {
	mockFstatic := new(mockFileStatic)

	result := fileop.StaticResult{
		TotalFileNumber: 100,
		WordLengthStd:   1.2323223,
	}
	mockFstatic.On("GetStaticResult", "test/anyPath").Return(result)

	handler := NewHandler("test", nil, mockFstatic)
	router := InitRouter(handler)

	req, _ := http.NewRequest("GET", "/static?path=/anyPath", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	mockFstatic.AssertCalled(suite.T(), "GetStaticResult", "test/anyPath")

	suite.Equal(200, resp.Code)
	suite.Equal("{\"TotalFileNumber\":100,\"CharNumAvg\":0,\"CharNumStd\":0,\"WordLengthAvg\":0,\"WordLengthStd\":1.2323223,\"TotalBytes\":0}",
		resp.Body.String())
}

func (suite *FileStaticHandlerTestSuite) TestGetStatic_should_do_static_on_root_if_no_parameter_provided() {
	mockFstatic := new(mockFileStatic)

	result := fileop.StaticResult{
		TotalFileNumber: 50,
		WordLengthStd:   1.2323223,
		CharNumAvg:      3,
	}
	mockFstatic.On("GetStaticResult", "test").Return(result)

	handler := NewHandler("test", nil, mockFstatic)
	router := InitRouter(handler)

	req, _ := http.NewRequest("GET", "/static", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	mockFstatic.AssertCalled(suite.T(), "GetStaticResult", "test")

	suite.Equal(200, resp.Code)
	suite.Equal("{\"TotalFileNumber\":50,\"CharNumAvg\":3,\"CharNumStd\":0,\"WordLengthAvg\":0,\"WordLengthStd\":1.2323223,\"TotalBytes\":0}",
		resp.Body.String())
}
