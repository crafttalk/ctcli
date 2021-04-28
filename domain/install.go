package domain

import (
	"ctcli/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
)

type HealthCheck struct {
	command []string
	waitFor int
}

type AppPackageConfig struct {
	baseImage   string
	healthcheck HealthCheck
	logsFolder  string
	configs     []string
}

/// tmp_abs_path example: /home/lkmfwe/ctcli/tmp/package/apps/agents
func extractBlobs(umociPath string, containerTmpPath string, name string) error {
	var skopeoImagePath = path.Join(containerTmpPath, "skopeo")
	var runcBundlePath = path.Join(containerTmpPath, "runc-bundle")

	cmd := exec.Command(
		umociPath,
		"unpack",
		"--rootless",
		"--image",
		fmt.Sprintf("%s:%s", skopeoImagePath, name),
		runcBundlePath,
	)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("original error: %s\nArgs:%s\nStdout:\n%s\n", err, cmd.Args, stdoutStderr)
	}

	if err := os.RemoveAll(skopeoImagePath); err != nil {
		return err
	}
	if err := os.Rename(path.Join(runcBundlePath, "rootfs"), path.Join(containerTmpPath, "rootfs")); err != nil {
		return err
	}
	if err := os.Rename(path.Join(runcBundlePath, "config.json"), path.Join(containerTmpPath, "config.json")); err != nil {
		return err
	}
	if err := os.RemoveAll(runcBundlePath); err != nil {
		return err
	}
	return nil
}

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
	config := new(AppPackageConfig)
	if err := json.Unmarshal(bytesFromFile, config); err != nil {
		return AppPackageConfig{}, err
	}
	return *config, nil
}

func Install(rootDir string, packagePath string) error {
	if !util.PathExists(packagePath) {
		return fmt.Errorf("couldn't find package %s", packagePath)
	}
	if !util.PathExists(rootDir) {
		return fmt.Errorf("root dir %s doesn't exists", rootDir)
	}

	log.Printf("Root dir: %s", rootDir)
	_ = os.MkdirAll(rootDir, os.ModePerm)

	tempFolder := path.Join(rootDir, "tmp")
	log.Printf("Extracting package %s to %s", packagePath, tempFolder)

	_ = os.RemoveAll(tempFolder)
	_ = os.MkdirAll(tempFolder, os.ModePerm)
	// We need to defer delete tempFolder HERE, immediately after we made it

	err := util.ExtractTarGz(packagePath, tempFolder)
	if err != nil {
		return err
	}

	umociPath := path.Join(tempFolder, "package", "umoci.amd64")
	runcPath := path.Join(tempFolder, "package", "runc.amd64")

	if err := os.Chmod(umociPath, 0775); err != nil {
		return err
	}
	if err := os.Chmod(runcPath, 0775); err != nil {
		return err
	}

	folders, err := ioutil.ReadDir(fmt.Sprintf("%s/package/apps", tempFolder))
	if err != nil {
		return err
	}
	for _, folder := range folders {
		var containerTmpPath = fmt.Sprintf("%s/package/apps/%s", tempFolder, folder.Name())
		log.Printf("Extracting blob %s", containerTmpPath)
		if err := extractBlobs(tempFolder, containerTmpPath, folder.Name()); err != nil {
			return err
		}
	}

	tempPackagePath := path.Join(tempFolder, "package")
	currentReleasePath := path.Join(rootDir, "current-release")

	// Need to move apps from current release to backups HERE, before copying tmp to current release

	// /tmp/package/runc.amd64 -> /current-release/runc.amd64
	if err := os.Rename(path.Join(tempPackagePath, "runc.amd64"), path.Join(currentReleasePath, "runc.amd64")); err != nil {
		return err
	}
	// /tmp/package/apps -> /current-release/apps
	if err := os.Rename(path.Join(tempPackagePath, "apps"), path.Join(currentReleasePath, "apps")); err != nil { // /tmp/package/apps DELETING!!!!!
		return err
	}

	// if not exists -> create /config
	// if not exists -> create /logs
	// if not exists -> create /releases
	util.CreateDirIfNotExist(rootDir, "config")
	util.CreateDirIfNotExist(rootDir, "logs")
	util.CreateDirIfNotExist(rootDir, "releases")

	// for each app in package:
	// copy container configs to /config/{appName}/config/config.json, /config/{appName}/bin/NLog.config (e.g)
	// make symlink from logFolder to /logs/{appName}
	// ln -s /home/lkmfwe/ctcli/logs/agents /home/lkmfwe/ctcli/current-release/apps/agents/rootfs/runtime/logs

	for _, app := range folders {
		appPackageConfigPath := path.Join(currentReleasePath, "apps", app.Name(), "package-config.json")
		appPackageConfig, err := GetAppConfig(appPackageConfigPath)
		if err != nil {
			return err
		}
		for _, configPath := range appPackageConfig.configs {
			absContainerConfigPath := path.Join(currentReleasePath, "apps", app.Name(), "rootfs", configPath)
			if util.PathExists(absContainerConfigPath) {
				rootDirConfigPath := path.Join(rootDir, "config", app.Name(), configPath)
				if err := util.CopyFile(absContainerConfigPath, rootDirConfigPath); err != nil {
					return err
				}
			}
		}
		absContainerLogPath := path.Join(currentReleasePath, "apps", app.Name(), "rootfs", appPackageConfig.logsFolder)
		if util.PathExists(absContainerLogPath) {
			rootDirLogPath := path.Join(rootDir, "logs", app.Name())
			if err := os.Symlink(absContainerLogPath, rootDirLogPath); err != nil {
				log.Printf("can not make symlink from %s to %s", rootDirLogPath, absContainerLogPath)
				return err
			}
		}
	}

	// defer удаление папки /tmp

	// ПОТОМ:
	// make backup from /current-release to -> backups/<release-name>.tar.gz
	// make backup manifest file -> backups/<release-name>.json

	return nil
}
