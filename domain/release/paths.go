package release

import (
	"ctcli/domain/ctcliDir"
	"path"
)
const (
	ReleaseInfoJson = "release-info.json"
	RuncBinary = "runc.amd64"
	AppsFolder = "apps"
	PackageConfigJson = "package-config.json"
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

func GetCurrentReleasePackageConfigPath(rootDir string, app string) string {
	return path.Join(GetCurrentReleaseAppFolder(rootDir, app), PackageConfigJson)
}

func GetCurrentReleaseAppRootfsFolder(rootDir string, app string) string {
	return path.Join(GetCurrentReleaseAppFolder(rootDir, app), RootFsFolder)
}

func GetReleasesFolder(rootDir string) string {
	return path.Join(rootDir, ReleasesFolder)
}