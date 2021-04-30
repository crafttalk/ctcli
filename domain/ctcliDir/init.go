package ctcliDir

import (
	"ctcli/util"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func CreateRootDir(rootDir string) {
	_ = os.MkdirAll(rootDir, os.ModePerm)
}

func CreateTempDir(rootDir string) {
	util.CreateDirIfNotExist(rootDir, TempDir)
}

func DeleteTempDir(rootDir string) {
	_ = os.RemoveAll(GetTempDir(rootDir))
}

func RecreateTempDir(rootDir string) {
	DeleteTempDir(rootDir)
	CreateTempDir(rootDir)
}

func getInitFilePath(rootDir string) string {
	return path.Join(rootDir, InitFile)
}

func CreateInitFile(rootDir string) {
	ioutil.WriteFile(getInitFilePath(rootDir), []byte{}, 0644)
}

func OkIfIsARootDir(rootDir string) error {
	if !util.PathExists(getInitFilePath(rootDir)) {
		return fmt.Errorf("%s is not a ctcli rootDir. Use `ctcli --root %s init`", rootDir, rootDir)
	}
	return nil
}

func Init(rootDir string) {
	CreateRootDir(rootDir)
	util.CreateDirIfNotExist(rootDir, ConfigDir)
	util.CreateDirIfNotExist(rootDir, LogsDir)
	util.CreateDirIfNotExist(rootDir, ReleasesFolder)
	util.CreateDirIfNotExist(rootDir, CurrentReleaseDir)
	CreateInitFile(rootDir)
}
