package fileop

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExists(t *testing.T) {
	t.Run("should return false for non exist file", func(t *testing.T) {
		fop := NewFileOp()
		r, _ := fop.Exists("arandomfilewillneverexist")
		assert.False(t, r)
	})

	t.Run("should detect exist for folder", func(t *testing.T) {
		os.Mkdir("test.temp", os.ModeDir|os.ModeTemporary)
		defer os.Remove("test.temp")

		fop := NewFileOp()
		r, _ := fop.Exists("test.temp")
		assert.True(t, r)
	})

	t.Run("should detect exist for file", func(t *testing.T) {
		f, _ := os.Create("test.temp")
		defer f.Close()
		defer os.Remove("test.temp")

		fop := NewFileOp()
		r, _ := fop.Exists("test.temp")
		assert.True(t, r)
	})
}

func TestDelete(t *testing.T) {
	t.Run("should return false if file not exist", func(t *testing.T) {
		fop := NewFileOp()
		err := fop.Delete("arandomfilewillneverexist")

		assert.Error(t, err)
		assert.Equal(t, "file not exists", err.Error())
	})

	t.Run("should return true if file deleted", func(t *testing.T) {
		os.Create("anyFile")
		defer os.Remove("anyFile")

		fop := NewFileOp()

		exists, _ := fop.Exists("anyFile")
		assert.True(t, exists)

		err := fop.Delete("anyFile")

		assert.NoError(t, err)

		exists, _ = fop.Exists("anyFile")
		assert.False(t, exists)
	})
}

func TestRead(t *testing.T) {
	t.Run("should return empty content if file not exist", func(t *testing.T) {
		fop := NewFileOp()

		content, err := fop.Read("arandomfilewillneverexist")

		assert.Equal(t, "", content)
		assert.Error(t, err)
		assert.Equal(t, "file not exists", err.Error())
	})

	t.Run("should return file content if file read successful", func(t *testing.T) {
		f, _ := os.Create("anyFile")
		f.WriteString("file content")
		defer os.Remove("anyFile")
		defer f.Close()

		fop := NewFileOp()

		content, err := fop.Read("anyFile")

		assert.NoError(t, err)
		assert.Equal(t, "file content", content)
	})
}

func TestWrite(t *testing.T) {
	t.Run("should write file", func(t *testing.T) {
		fop := NewFileOp()

		fop.Write("test.tmp", "any content here")
		defer os.Remove("test.tmp")

		buf, err := ioutil.ReadFile("./test.tmp")
		content := string(buf)

		assert.NoError(t, err)
		assert.Equal(t, "any content here", content)
	})

	t.Run("should write file under non-exist folder", func(t *testing.T) {
		fop := NewFileOp()

		fop.Write("nonexist/secondlevel/test.tmp", "any content here")
		defer os.RemoveAll("nonexist")

		buf, err := ioutil.ReadFile("nonexist/secondlevel/test.tmp")
		content := string(buf)

		assert.NoError(t, err)
		assert.Equal(t, "any content here", content)
	})

	t.Run("should refuse to override file by default", func(t *testing.T) {
		f, _ := os.Create("test.tmp")
		f.WriteString("initial content")
		defer os.Remove("test.tmp")
		defer f.Close()

		fop := NewFileOp()

		err := fop.Write("test.tmp", "any content here")
		assert.Error(t, err)
		assert.Equal(t, "file already exists", err.Error())

		buf, err := ioutil.ReadFile("test.tmp")
		content := string(buf)

		assert.Equal(t, "initial content", content)
	})

	t.Run("should override file if require", func(t *testing.T) {
		f, _ := os.Create("test.tmp")
		f.WriteString("initial content")
		defer os.Remove("test.tmp")
		defer f.Close()

		fop := NewFileOp()

		err := fop.Write("test.tmp", "any content here", true)
		assert.NoError(t, err)

		buf, err := ioutil.ReadFile("test.tmp")
		content := string(buf)

		assert.Equal(t, "any content here", content)
	})
}
