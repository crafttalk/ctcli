package util

import (
	"io"
	"os"
	"path"
	"path/filepath"
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

func RemoveContentOfFolder(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
