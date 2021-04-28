package packaging

import (
	"ctcli/domain/release"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

const (
	PackageContentFolder = "package"
	AppsFolder = "apps"
	UmociBinary = "umoci.amd64"
	RuncBinary = "runc.amd64"
	SkopeoFolder = "skopeo"
	RuncConfigJson = "config.json"
	RuncRootFs = "rootfs"
	AppVersionFile = "version.json"
)

func GetPackageRuncPath(packageDir string) string {
	return path.Join(packageDir, PackageContentFolder, RuncBinary)
}

func GetPackageUmociPath(packageDir string) string {
	return path.Join(packageDir, PackageContentFolder, UmociBinary)
}

func GetPackageAppsFolder(packageDir string) string {
	return path.Join(packageDir, PackageContentFolder, AppsFolder)
}

func GetPackageAppFolder(packageDir string, app string) string {
	return path.Join(GetPackageAppsFolder(packageDir), app)
}

func GetAppVersionFilePath(packageDir string, app string) string {
	return path.Join(GetPackageAppFolder(packageDir, app), AppVersionFile)
}

func GetVersionJsonFromFile(versionJsonFilePath string) (release.AppVersion, error) {
	jsonFile, err := os.Open(versionJsonFilePath)
	if err != nil {
		return release.AppVersion{}, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var releaseMeta release.AppVersion
	err = json.Unmarshal(byteValue, &releaseMeta)
	if err != nil {
		return release.AppVersion{}, err
	}

	return releaseMeta, nil
}

