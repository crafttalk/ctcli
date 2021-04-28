package packaging

import "path"

const (
	PackageContentFolder = "package"
	AppsFolder = "apps"
	UmociBinary = "umoci.amd64"
	RuncBinary = "runc.amd64"
	SkopeoFolder = "skopeo"
	RuncConfigJson = "config.json"
	RuncRootFs = "rootfs"
	MetaFile = "meta.json"
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

