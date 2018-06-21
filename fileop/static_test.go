package fileop

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type dummies struct {
	fileDirs   []string
	folderDirs []string
}

func TestFileStatic(t *testing.T) {
	t.Run("should get static result for multilayer file structure", func(t *testing.T) {
		names := generateDummyNames("all")
		createDummies(names)
		defer os.RemoveAll("all")

		fs := NewFileStatic()
		result := fs.GetStaticResult("all")

		assert.Equal(t, 30, int(result.TotalFileNumber))
		assert.Equal(t, 7*30, int(result.TotalBytes))
		assert.Equal(t, 3, int(result.WordLengthAvg))
		assert.Equal(t, 0, int(result.WordLengthStd))
		assert.Equal(t, 6, int(result.CharNumAvg))
		assert.Equal(t, 0, int(result.CharNumStd))
	})
}

func generateDummyNames(base string) dummies {
	var d dummies

	d.folderDirs = append(d.folderDirs, base)

	for i := 1; i <= 10; i++ {
		if i < 5 {
			d.folderDirs = append(d.folderDirs, fmt.Sprintf("all/folder%d", i))
			for j := 1; j <= 5; j++ {
				d.fileDirs = append(d.fileDirs, fmt.Sprintf("all/folder%d/file%d", i, j))
			}
		}

		d.fileDirs = append(d.fileDirs, fmt.Sprintf("all/file%d", i))
	}

	return d
}

func createDummies(d dummies) {
	// for folders
	for _, name := range d.folderDirs {
		os.Mkdir(name, os.ModeDir|os.ModePerm)
	}

	// for files
	for _, name := range d.fileDirs {
		f, _ := os.Create(name)
		f.WriteString("abc bcd")
	}
}

type mockFileInfo struct {
	mock.Mock
}

func (i *mockFileInfo) Name() string {
	args := i.Called()
	return args.String(0)
}

func (i *mockFileInfo) Size() int64 {
	args := i.Called()
	return int64(args.Int(0))
}
func (i *mockFileInfo) Mode() os.FileMode {
	return 0666
}
func (i *mockFileInfo) ModTime() time.Time {
	return time.Now()
}

func (i *mockFileInfo) IsDir() bool {
	args := i.Called()
	return args.Bool(0)
}
func (i *mockFileInfo) Sys() interface{} {
	return nil
}
