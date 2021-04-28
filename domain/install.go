package domain

import (
	"ctcli/domain/ctcliDir"
	"ctcli/domain/packaging"
	"ctcli/domain/release"
	"ctcli/util"
	"encoding/json"
	"fmt"
	"github.com/otiai10/copy"
	"io/ioutil"
	"log"
	"os"
	"path"
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

		fmt.Println(appPackageConfig)
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

func GetCurrentCurrentReleaseInfo (currentReleaseInfoPath string) (packaging.PackageMeta, error) {
	file, err := os.Open(currentReleaseInfoPath)
	if err != nil {
		return packaging.PackageMeta{}, err
	}
	defer file.Close()
	bytesFromFile, err := ioutil.ReadAll(file)
	if err != nil {
		return packaging.PackageMeta{}, err
	}
	var meta packaging.PackageMeta
	if err := json.Unmarshal(bytesFromFile, &meta); err != nil {
		return packaging.PackageMeta{}, err
	}
	return meta, nil
}

func MakeCurrentReleaseInfoFile (rootDir string) error {
	//currentReleaseInfoPath := release.GetCurrentReleaseInfoPath(rootDir)
	//meta, err := GetCurrentCurrentReleaseInfo(currentReleaseInfoPath)
	//if err != nil {
	//	return err
	//}
	//currentReleaseInfo := meta.ReleaseMeta
	return nil

}

func Install(rootDir string, packagePath string) error {
	if !util.PathExists(packagePath) {
		return fmt.Errorf("couldn't find package %s", packagePath)
	}
	if !util.PathExists(rootDir) {
		return fmt.Errorf("root dir %s doesn't exists", rootDir)
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

	// ПОТОМ:
	// make backup from /current-release to -> backups/<release-name>.tar.gz
	// make backup manifest file -> backups/<release-name>.json

	if err := copyBinariesToRelease(rootDir, tempFolder); err != nil {
		return err
	}
	if err := copyPackagesToRelease(rootDir, tempFolder); err != nil {
		return err
	}

	//  tmp/meta.json -> current-release/

	return nil
}
