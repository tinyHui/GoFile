package fileop

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileOp interface {
	Exists(path string) (bool, error)
	Delete(path string) error
	Read(path string) (string, error)
	Write(path string, content string, override ...bool) error
}

type fileOp struct {
}

func NewFileOp() *fileOp {
	return &fileOp{}
}

func (o *fileOp) Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func (o *fileOp) Delete(path string) error {
	exists, _ := o.Exists(path)
	if exists {
		err := os.Remove(path)
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("file not exists")
	}
}

func (o *fileOp) Read(path string) (string, error) {
	exists, _ := o.Exists(path)
	if exists {
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			return "", err
		}
		return string(buf), nil
	} else {
		return "", errors.New("file not exists")
	}
}

func (o *fileOp) Write(path string, content string, override ...bool) error {
	if len(override) == 0 {
		override = []bool{false}
	}

	exists, err := o.Exists(path)
	if err != nil {
		return err
	}
	if exists && !override[0] {
		return errors.New("file already exists")
	}

	parentPath := filepath.Dir(path)
	exists, err = o.Exists(parentPath)
	if err != nil {
		return err
	}
	if !exists {
		err := os.MkdirAll(parentPath, os.ModeDir|os.ModePerm)
		if err != nil {
			return err
		}
	}

	f, err := os.Create(path)
	defer f.Close()

	if err != nil {
		return err
	}

	f.WriteString(content)

	return nil
}
