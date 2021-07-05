package release

import (
	"ctcli/domain/ctcliDir"
	"ctcli/util"
	"path"
)
const (
	ReleaseInfoJson = "release-info.json"
	RuncBinary = "runc.amd64"
	AppsFolder = "apps"
	PackageConfigJson = "package-config.json"
	RuncConfigJson = "config.json"
	RootFsFolder = "rootfs"
	ReleasesFolder = "releases"
)

func GetCurrentReleaseRuncPath(rootDir string) string {
	return path.Join(ctcliDir.GetCurrentReleaseDir(rootDir), RuncBinary)
}

func GetCurrentReleaseInfoPath(rootDir string) string {
	return path.Join(rootDir, ctcliDir.CurrentReleaseDir, ReleaseInfoJson)
}

func GetCurrentReleaseAppsFolder(rootDir string) string {
	return path.Join(ctcliDir.GetCurrentReleaseDir(rootDir), AppsFolder)
}

func GetCurrentReleaseAppFolder(rootDir string, app string) string {
	return path.Join(GetCurrentReleaseAppsFolder(rootDir), app)
}

func CurrentReleaseAppExists(rootDir string, app string) bool {
	return util.PathExists(GetCurrentReleaseAppFolder(rootDir, app)) &&
		   util.PathExists(GetCurrentReleaseRuncConfigPath(rootDir, app))
}

func GetCurrentReleasePackageConfigPath(rootDir string, app string) string {
	return path.Join(GetCurrentReleaseAppFolder(rootDir, app), PackageConfigJson)
}

func GetCurrentReleaseRuncConfigPath(rootDir string, app string) string {
	return path.Join(GetCurrentReleaseAppFolder(rootDir, app), RuncConfigJson)
}

func GetCurrentReleaseAppRootfsFolder(rootDir string, app string) string {
	return path.Join(GetCurrentReleaseAppFolder(rootDir, app), RootFsFolder)
}

func GetReleasesFolder(rootDir string) string {
	return path.Join(rootDir, ReleasesFolder)
}