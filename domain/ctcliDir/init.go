package ctcliDir

import (
	"ctcli/util"
	"os"
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

func Init(rootDir string) {
	CreateRootDir(rootDir)
	util.CreateDirIfNotExist(rootDir, ConfigDir)
	util.CreateDirIfNotExist(rootDir, LogsDir)
	util.CreateDirIfNotExist(rootDir, ReleasesFolder)
	util.CreateDirIfNotExist(rootDir, CurrentReleaseDir)
}