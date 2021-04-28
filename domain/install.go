package domain

import (
	"ctcli/domain/ctcliDir"
	"ctcli/domain/packaging"
	"ctcli/domain/release"
	"ctcli/util"
	"encoding/json"
	"fmt"
	"github.com/otiai10/copy"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

func GetAppConfig(path string) (AppPackageConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return AppPackageConfig{}, err
	}
	defer file.Close()
	bytesFromFile, err := ioutil.ReadAll(file)
	if err != nil {
		return AppPackageConfig{}, err
	}
	var config AppPackageConfig
	if err := json.Unmarshal(bytesFromFile, &config); err != nil {
		return AppPackageConfig{}, err
	}
	return config, nil
}

func copyBinariesToRelease(rootDir string, packagePath string) error {
	if err := copy.Copy(packaging.GetPackageRuncPath(packagePath), release.GetCurrentReleaseRuncPath(rootDir)); err != nil {
		return err
	}
	if err := copy.Copy(packaging.GetPackageAppsFolder(packagePath), release.GetCurrentReleaseAppsFolder(rootDir)); err != nil {
		return err
	}
	return nil
}

func copyPackagesToRelease(rootDir string, packagePath string) error {
	apps, err := packaging.GetPackageAppsList(packagePath)
	if err != nil {
		return err
	}

	for _, app := range apps {
		appPackageConfigPath := release.GetCurrentReleasePackageConfigPath(rootDir, app)
		appPackageConfig, err := GetAppConfig(appPackageConfigPath)
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

func createReleaseInfoFile(rootDir string, tempFolder string) error {
	apps, err := packaging.GetPackageAppsList(tempFolder)

	releaseInfo := release.ReleaseMeta{}

	if err != nil {
		return err
	}
	for _, app := range apps {
		versionPath := packaging.GetAppVersionFilePath(tempFolder, app)
		version, err := packaging.GetVersionJsonFromFile(versionPath)
		if err != nil {
			return err
		}

		releaseInfo.AppVersions = append(releaseInfo.AppVersions, version)
	}

	releaseInfo.PreviousRelease = "" // Todo!
	releaseInfo.Id = uuid.NewV4().String()
	releaseInfo.CreatedAt = time.Now()

	content, _ := json.MarshalIndent(releaseInfo, "", " ")
	err = ioutil.WriteFile(release.GetCurrentReleaseInfoPath(rootDir), content, 0644)
	return err
}

func Install(rootDir string, packagePath string) error {
	if !util.PathExists(packagePath) {
		return fmt.Errorf("couldn't find package %s", packagePath)
	}
	if !util.PathExists(rootDir) {
		return fmt.Errorf("root dir %s doesn't exists", rootDir)
	}
	if err := ctcliDir.OkIfIsARootDir(rootDir); err != nil {
		return err
	}
	if util.PathExists(release.GetCurrentReleaseInfoPath(rootDir)) {
		return fmt.Errorf("current release already installed. Maybe you intended to use upgrade?")
	}

	log.Printf("Installing in root dir: %s", rootDir)

	tempFolder := ctcliDir.GetTempDir(rootDir)
	log.Printf("Extracting package %s to %s", packagePath, tempFolder)

	ctcliDir.DeleteTempDir(rootDir)
	ctcliDir.CreateTempDir(rootDir)

	err := util.ExtractTarGz(packagePath, tempFolder)
	if err != nil {
		return err
	}

	if err := packaging.PreparePackage(tempFolder); err != nil {
		return err
	}

	// TODO:
	// make backup from /current-release to -> backups/<release-name>.tar.gz
	// make backup manifest file -> backups/<release-name>.json

	if err := copyBinariesToRelease(rootDir, tempFolder); err != nil {
		return err
	}
	if err := copyPackagesToRelease(rootDir, tempFolder); err != nil {
		return err
	}
	if err := createReleaseInfoFile(rootDir, tempFolder); err != nil {
		return err
	}

	return nil
}
