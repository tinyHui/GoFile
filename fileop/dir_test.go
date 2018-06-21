package fileop

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirWalk(t *testing.T) {
	t.Run("should list all files under dir", func(t *testing.T) {
		dummies := generateDummyNames("all")
		createDummies(dummies)
		defer os.RemoveAll("all")

		count := 0
		fileDirs := make(map[string]bool)

		for _, name := range dummies.fileDirs {
			fileDirs[name] = false
		}

		files := make(chan osFile)
		gw := new(sync.WaitGroup)
		gw.Add(1)

		go dirWalk("all", files, gw)

		for file := range files {
			count += 1
			fileDirs[file.f.Name()] = true
		}

		gw.Wait()

		assert.Equal(t, 30, count)

		for name, covered := range fileDirs {
			if !covered {
				assert.Fail(t, fmt.Sprintf("%s not covered", name))
			}
		}
	})

	t.Run("should list the file if given dir is a file", func(t *testing.T) {
		os.Create("anyFile")
		defer os.Remove("anyFile")

		count := 0
		fileDirs := make(map[string]bool)
		fileDirs["anyFile"] = false

		files := make(chan osFile)
		gw := new(sync.WaitGroup)
		gw.Add(1)

		go dirWalk("anyFile", files, gw)

		for file := range files {
			count += 1
			fileDirs[file.f.Name()] = true
		}

		gw.Wait()

		assert.Equal(t, 1, count)

		for name, covered := range fileDirs {
			if !covered {
				assert.Fail(t, fmt.Sprintf("%s not covered", name))
			}
		}
	})

	t.Run("should give empty channel when given dir not exist", func(t *testing.T) {
		count := 0

		files := make(chan osFile)
		gw := new(sync.WaitGroup)
		gw.Add(1)

		go dirWalk("notexistdir", files, gw)
		gw.Wait()

		assert.Equal(t, 0, count)
	})
}
