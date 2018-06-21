package fileop

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type osFile struct {
	f     *os.File
	stats os.FileInfo
}

func dirWalk(dir string, files chan osFile, gw *sync.WaitGroup) {
	defer gw.Done()
	defer close(files)

	f, err := os.Open(dir)
	if err != nil {
		return
	}
	fs, err := f.Stat()
	if err != nil {
		return
	}

	if !fs.IsDir() {
		files <- osFile{
			f:     f,
			stats: fs,
		}
	} else {
		walkerWait := new(sync.WaitGroup)
		walkerWait.Add(1)
		go execDirWalk(dir, files, walkerWait)

		walkerWait.Wait()
	}
}

func execDirWalk(dir string, files chan osFile, walkerWait *sync.WaitGroup) {
	defer walkerWait.Done()

	fnames, err := filepath.Glob(fmt.Sprintf("%s/*", dir))

	if err != nil {
	}

	for _, fname := range fnames {
		f, err := os.Open(fname)
		if err != nil {
			continue
		}

		fs, err := f.Stat()
		if err != nil {
			continue
		}

		if fs.IsDir() {
			walkerWait.Add(1)
			execDirWalk(f.Name(), files, walkerWait)
		} else {
			files <- osFile{
				f:     f,
				stats: fs,
			}
		}
	}
}
