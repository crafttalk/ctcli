package ctcliDir

import (
	"fmt"
	"path"
	"time"
)

const (
	TempDir           = "tmp"
	ConfigDir         = "config"
	CurrentReleaseDir = "current-release"
	LogsDir           = "logs"
	DataDir			  = "data"
	ReleasesFolder    = "releases"
	InitFile 		  = ".ctcli-init"
)

func GetTempDir(rootDir string) string {
	return path.Join(rootDir, TempDir)
}

func GetCurrentReleaseDir(rootDir string) string {
	return path.Join(rootDir, CurrentReleaseDir)
}

func GetConfigDir(rootDir string) string {
	return path.Join(rootDir, ConfigDir)
}

func GetLogsDir(rootDir string) string {
	return path.Join(rootDir, LogsDir)
}

func GetDataDir(rootDir string) string {
	return path.Join(rootDir, DataDir)
}

func GetAppConfigDir(rootDir string, app string) string {
	return path.Join(GetConfigDir(rootDir), app)
}

func GetAppDataDir(rootDir string, app string) string {
	return path.Join(GetDataDir(rootDir), app)
}

func GetAppLogsDir(rootDir string, app string) string {
	return path.Join(GetLogsDir(rootDir), app)
}

func GetAppStdoutLogFilePath(rootDir, app string) string {
	return path.Join(GetAppLogsDir(rootDir, app), "stdout-stderr.log")
}

func GetNewArchiveStdoutLogFilePath(rootDir, app string) string {
	return path.Join(
		GetAppLogsDir(rootDir, app),
		"stdout-stderr-arch",
		fmt.Sprintf("%s.log", time.Now().UTC().Format(time.RFC3339)))
}