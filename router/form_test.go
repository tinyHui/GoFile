package router

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type FormTestSuite struct {
	suite.Suite
	c    *gin.Context
}

func TestFormTestSuite(t *testing.T) {
	suite.Run(t, new(FormTestSuite))
}

func (suite *FormTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	suite.c, _ = gin.CreateTestContext(resp)
}

func (suite *FormTestSuite) TestParseUploadFileRequest_should_parse_post_form() {
	payload := "------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"file\"; filename=\"example\"\r\nContent-Type: false\r\n\r\n\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"path\"\r\n\r\n/testpath\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW--"

	req := httptest.NewRequest(
		"POST", "/file",
		strings.NewReader(payload))
	req.Header.Add("content-type", "multipart/form-data;boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")

	suite.c.Request = req

	f, err := parseUploadFileRequest(suite.c)

	suite.Equal("/testpath", f.filePath)
	suite.Equal("example", f.file.Filename)
	suite.Nil(err)
}

func (suite *FormTestSuite) TestParseUploadFileRequest_should_use_empty_string_when_filePath_not_provided() {
	payload := "------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"file\"; filename=\"example\"\r\nContent-Type: false\r\n\r\n\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW--"

	req := httptest.NewRequest(
		"POST", "/file",
		strings.NewReader(payload))
	req.Header.Add("content-type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	suite.c.Request = req

	f, err := parseUploadFileRequest(suite.c)

	suite.Equal("", f.filePath)
	suite.Equal("example", f.file.Filename)
	suite.Nil(err)
}

func (suite *FormTestSuite) TestParseUploadFileRequest_should_give_400_if_given_filePath_got_error() {
	payload := "Content-Disposition: form-data; name=\"file\"; filename=\"\"\r\nContent-Type: false\r\n\r\n\r\n"

	req := httptest.NewRequest(
		"POST", "/file",
		strings.NewReader(payload))
	req.Header.Add("content-type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	suite.c.Request = req

	_, err := parseUploadFileRequest(suite.c)

	suite.Error(err)
	suite.Equal("multipart: NextPart: EOF", err.Error())
}

func (suite *FormTestSuite) TestParseFileRequest_should_parse_json_body() {
	payload := "{\n\t\"filePath\": \"/abc\",\n\t\"fileName\": \"file_name\",\n\t\"fileContent\": \"content content\"\n}"

	req := httptest.NewRequest(
		"PUT", "/file",
		strings.NewReader(payload))
	req.Header.Add("Content-Type", "application/json")

	suite.c.Request = req

	content, err := parseFileRequest(suite.c)

	suite.NoError(err)
	suite.Equal("/abc", content.FilePath)
	suite.Equal("file_name", content.FileName)
	suite.Equal("content content", content.FileContent)
}

func (suite *FormTestSuite) TestParseFileRequest_should_give_empty_value_if_field_missing() {
	payload := "{\n\t\"filePath\": \"/abc\",\n\t\"fileContent\": \"content content\"\n}"

	req := httptest.NewRequest(
		"PUT", "/file",
		strings.NewReader(payload))
	req.Header.Add("Content-Type", "application/json")

	suite.c.Request = req

	content, err := parseFileRequest(suite.c)

	suite.NoError(err)
	suite.Equal("/abc", content.FilePath)
	suite.Equal("", content.FileName)
	suite.Equal("content content", content.FileContent)
}

func (suite *FormTestSuite) TestParseFileRequest_should_return_error_if_not_able_to_parse_body() {
	payload := ""

	req := httptest.NewRequest(
		"PUT", "/file",
		strings.NewReader(payload))
	req.Header.Add("Content-Type", "application/json")

	suite.c.Request = req

	content, err := parseFileRequest(suite.c)

	suite.Error(err)
	suite.Equal("not able to parse request body", err.Error())

	suite.Empty(content)
}

func (suite *FormTestSuite) TestParseTargetPath_should_return_parameter() {
	q := url.Values{}
	q.Add("path", "anyPath")

	path, err := parseTargetPath(q)
	suite.NoError(err)
	suite.Equal("anyPath", path)
}

func (suite *FormTestSuite) TestParseTargetPath_should_return_empty_string_and_result_if_not_got_parameter() {
	q := url.Values{}

	path, err := parseTargetPath(q)
	suite.Error(err)
	suite.Equal("target path parameter not provided", err.Error())
	suite.Equal("", path)
}
