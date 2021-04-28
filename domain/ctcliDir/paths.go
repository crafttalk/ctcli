package ctcliDir

import "path"

const (
	TempDir           = "tmp"
	ConfigDir         = "config"
	CurrentReleaseDir = "current-release"
	LogsDir           = "logs"
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

func GetAppConfigDir(rootDir string, app string) string {
	return path.Join(GetConfigDir(rootDir), app)
}

func GetAppLogsDir(rootDir string, app string) string {
	return path.Join(GetLogsDir(rootDir), app)
}