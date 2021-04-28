package util

import (
	"io"
	"os"
	"path"
)

func PathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateDirIfNotExist(rootDir string, dirToCreate string) {
	absPath := path.Join(rootDir, dirToCreate)
	if !PathExists(absPath) {
		_ = os.MkdirAll(absPath, os.ModePerm)
	}
}

func CopyFile(pathFrom, pathTo string) error {
	in, err := os.Open(pathFrom)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(pathTo)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}