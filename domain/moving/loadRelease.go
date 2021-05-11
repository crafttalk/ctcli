package moving

import (
	"ctcli/domain/appConfig"
	"ctcli/domain/ctcliDir"
	"ctcli/domain/packaging"
	"ctcli/domain/release"
	"ctcli/util"
	"encoding/json"
	"github.com/otiai10/copy"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func CopyBinariesToRelease(rootDir string, packagePath string) error {
	if err := copy.Copy(packaging.GetPackageRuncPath(packagePath), release.GetCurrentReleaseRuncPath(rootDir)); err != nil {
		return err
	}
	if err := copy.Copy(packaging.GetPackageAppsFolder(packagePath), release.GetCurrentReleaseAppsFolder(rootDir)); err != nil {
		return err
	}
	return nil
}

func MakeConfigSymLink(absContainerConfigPath string, rootDirConfigPath string) error {
	if util.PathExists(absContainerConfigPath) {
		_ = os.MkdirAll(path.Dir(rootDirConfigPath), os.ModePerm)

		if err := util.CopyFile(absContainerConfigPath, rootDirConfigPath); err != nil {
			return err
		}
		_ = os.RemoveAll(absContainerConfigPath)
		if err := os.Symlink(rootDirConfigPath, absContainerConfigPath); err != nil {
			log.Printf("can not make symlink from %s to %s", rootDirConfigPath, absContainerConfigPath)
			return err
		}
	}
	return nil
}

func MakeContainerDirSymLink(containerDirPath string, rootDirPath string) error {
	if util.PathExists(containerDirPath) {
		_ = os.RemoveAll(containerDirPath)
	}

	_ = os.MkdirAll(rootDirPath, os.ModePerm)
	if err := os.Symlink(rootDirPath, containerDirPath); err != nil {
		log.Printf("can not make symlink from %s to %s", rootDirPath, containerDirPath)
		return err
	}
	return nil
}

func CopyPackagesToRelease(rootDir string, packagePath string) error {
	apps, err := packaging.GetPackageAppsList(packagePath)
	if err != nil {
		return err
	}

	for _, app := range apps {
		appPackageConfigPath := release.GetCurrentReleasePackageConfigPath(rootDir, app)
		appPackageConfig, err := appConfig.GetAppConfig(appPackageConfigPath)
		if err != nil {
			return err
		}

		for _, configPath := range appPackageConfig.Configs {
			absContainerConfigPath := path.Join(release.GetCurrentReleaseAppRootfsFolder(rootDir, app), configPath)
			rootDirConfigPath := path.Join(ctcliDir.GetAppConfigDir(rootDir, app), configPath)
			if err := MakeConfigSymLink(absContainerConfigPath, rootDirConfigPath); err != nil {
				return err
			}
		}

		if appPackageConfig.LogsFolder != "" {
			absContainerLogPath := path.Join(release.GetCurrentReleaseAppRootfsFolder(rootDir, app), appPackageConfig.LogsFolder)
			rootDirLogPath := ctcliDir.GetAppLogsDir(rootDir, app)
			if err := MakeContainerDirSymLink(absContainerLogPath, rootDirLogPath); err != nil {
				return err
			}
		}

		for _, dataPath := range appPackageConfig.Data {
			absContainerDataPath := path.Join(release.GetCurrentReleaseAppRootfsFolder(rootDir, app), dataPath)
			rootDirDataPath := path.Join(ctcliDir.GetAppDataDir(rootDir, app), dataPath)
			if err := MakeContainerDirSymLink(absContainerDataPath, rootDirDataPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func CreateReleaseInfoFile(rootDir string, tempFolder string) error {
	releaseInfo, err := packaging.CreateReleaseInfo(tempFolder)
	if err != nil {
		return err
	}
	content, _ := json.MarshalIndent(releaseInfo, "", " ")
	err = ioutil.WriteFile(release.GetCurrentReleaseInfoPath(rootDir), content, 0644)
	return err
}

func LoadRelease(rootDir, tempFolder string) error {
	if err := util.RemoveContentOfFolder(ctcliDir.GetCurrentReleaseDir(rootDir)); err != nil {
		return err
	}
	if err := CopyBinariesToRelease(rootDir, tempFolder); err != nil {
		return err
	}
	if err := CopyPackagesToRelease(rootDir, tempFolder); err != nil {
		return err
	}
	if err := CreateReleaseInfoFile(rootDir, tempFolder); err != nil {
		return err
	}
	return nil
}
