package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/tinyhui/GoFile/fileop"
)

type HandlerTestSuite struct {
	suite.Suite
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func (suite *HandlerTestSuite) TestHealCheckHandler() {
	handler := NewHandler("", nil, nil)
	router := InitRouter(handler)

	req, _ := http.NewRequest("GET", "/healthcheck", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	suite.Equal(200, resp.Code)
	suite.Equal("{\"code\":200,\"msg\":\"OK\"}", resp.Body.String())
}

type mockFileOp struct {
	mock.Mock
}

func (o *mockFileOp) Exists(path string) (bool, error) {
	args := o.Called(path)
	return args.Bool(0), args.Error(1)
}

func (o *mockFileOp) Delete(path string) error {
	args := o.Called(path)
	return args.Error(0)

}

func (o *mockFileOp) Read(path string) (string, error) {
	args := o.Called(path)
	return args.String(0), args.Error(1)
}

func (o *mockFileOp) Write(path string, content string, override ...bool) error {
	args := o.Called(path, content, override[0])
	return args.Error(0)
}

type mockFileStatic struct {
	mock.Mock
}

func (s *mockFileStatic) GetStaticResult(rootDir string) fileop.StaticResult {
	args := s.Called(rootDir)
	return args.Get(0).(fileop.StaticResult)
}
