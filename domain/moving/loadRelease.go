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
			if util.PathExists(absContainerConfigPath) {
				configPath := path.Join(ctcliDir.GetAppConfigDir(rootDir, app), configPath)
				_ = os.MkdirAll(path.Dir(configPath), os.ModePerm)

				if err := util.CopyFile(absContainerConfigPath, configPath); err != nil {
					return err
				}
			}
		}

		if appPackageConfig.LogsFolder != "" {
			absContainerLogPath := path.Join(release.GetCurrentReleaseAppRootfsFolder(rootDir, app), appPackageConfig.LogsFolder)
			if util.PathExists(absContainerLogPath) {
				_ = os.RemoveAll(absContainerLogPath)
			}

			rootDirLogPath := ctcliDir.GetAppLogsDir(rootDir, app)
			_ = os.MkdirAll(rootDirLogPath, os.ModePerm)
			if err := os.Symlink(rootDirLogPath, absContainerLogPath); err != nil {
				log.Printf("can not make symlink from %s to %s", rootDirLogPath, absContainerLogPath)
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
	util.RemoveContentOfFolder(ctcliDir.GetCurrentReleaseDir(rootDir))
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
