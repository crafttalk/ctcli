package domain

import (
	"ctcli/domain/ctcliDir"
	"ctcli/domain/release"
	"ctcli/domain/runc"
	"ctcli/util"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

func StartApps(rootDir string, apps []string) error {
	runcPath := release.GetCurrentReleaseRuncPath(rootDir)
	if !util.PathExists(runcPath) {
		return fmt.Errorf("there is no runc in current-relase folder")
	}
	appsPath := release.GetCurrentReleaseAppsFolder(rootDir)
	appFolders, err := ioutil.ReadDir(appsPath)
	if err != nil {
		return err
	}
	appNames := release.GetAppNamesFromFolders(appFolders)
	var appsToStart []string
	if len(apps) > 0 {
		var _appsToStart, appsNotExistInReleaseFolder = release.CheckIfAppsInReleaseFolder(apps, appNames)
		appsToStart = _appsToStart
		if len(appsNotExistInReleaseFolder) > 0 {
			for _, notExistingApp := range appsNotExistInReleaseFolder {
				color.Red(fmt.Sprintf("app with name: %s is not installed", notExistingApp))
			}
		}
	} else {
		appsToStart = appNames
	}
	for _, appName := range appsToStart {
		appPath := release.GetCurrentReleaseAppFolder(rootDir, appName)
		if err := StartApp(rootDir, appName, appPath, runcPath); err != nil {
			color.Red(fmt.Sprintf("error while starting %s app, error: %s", appName, err))
		}
	}
	return nil
}

func StartApp(rootDir, appName, appPath, runcPath string) error {
	appStatus := runc.GetStatus(rootDir, appName)
	if appStatus == "running" {
		color.Green("%s already running", appName)
		return nil
	} else if appStatus == "stopped" {
		if err := runc.DeleteContainer(rootDir, appName); err != nil {
			return err
		}
	}

	cmd := exec.Command(
		runcPath,
		"create",
		"--bundle",
		appPath,
		appName,
	)

	logFilePath := ctcliDir.GetAppStdoutLogFilePath(rootDir, appName)
	_ = os.MkdirAll(path.Dir(logFilePath), os.ModePerm)
	if util.PathExists(logFilePath) {
		archiveLogFilePath := ctcliDir.GetNewArchiveStdoutLogFilePath(rootDir, appName)
		_ = os.MkdirAll(path.Dir(archiveLogFilePath), os.ModePerm)
		_ = os.Rename(logFilePath, archiveLogFilePath)
	}

	stdout, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	cmd.Stdout = stdout
	cmd.Stderr = stdout
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command(
		runcPath,
		"start",
		appName,
	)
	defer stdout.Close()

	if err := cmd.Run(); err != nil {
		return err
	}

	if err := cmd.Process.Release(); err != nil {
		return err
	}
	color.Green(fmt.Sprintf("%s started", appName))
	return nil
}
