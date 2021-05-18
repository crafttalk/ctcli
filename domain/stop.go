package domain

import (
	"ctcli/domain/release"
	"ctcli/domain/runc"
	"ctcli/util"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"os/exec"
)

func StopApps(rootDir string, apps []string) error {
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
	var appsToStop []string
	if len(apps) > 0 {
		var _appsToStop, appsNotExistInReleaseFolder = release.CheckIfAppsInReleaseFolder(apps, appNames)
		appsToStop = _appsToStop
		if len(appsNotExistInReleaseFolder) > 0 {
			for _, notExistingApp := range appsNotExistInReleaseFolder {
				color.Red(fmt.Sprintf("app with name: %s is not installed", notExistingApp))
			}
		}
	} else {
		appsToStop = appNames
	}
	for _, appToStop := range appsToStop {
		if err := StopApp(rootDir, appToStop, runcPath); err != nil {
			color.Red(fmt.Sprintf("error while stopping %s app, error: %s", appToStop, err))
		}
	}
	return nil
}

func StopApp(rootDir, appName, runcPath string) error {
	status := runc.GetStatus(rootDir, appName)
	if status == "running" {
		cmd := exec.Command(
			runcPath,
			"kill",
			runc.GetContainerName(rootDir, appName),
			"SIGTERM",
		)
		if err := cmd.Run(); err != nil {
			return err
		}

		runc.WaitUntilNotRunning(rootDir, appName)

		status := runc.GetStatus(rootDir, appName)
		if status == "running" {
			cmd := exec.Command(
				runcPath,
				"kill",
				runc.GetContainerName(rootDir, appName),
				"SIGKILL",
			)
			if err := cmd.Run(); err != nil {
				return err
			}

			runc.WaitUntilNotRunning(rootDir, appName)
		}
	} else if status == "unknown" {
		color.Green(fmt.Sprintf("%s already stopped", appName))
		return nil
	}
	if err := runc.DeleteContainer(rootDir, appName); err != nil {
		return err
	}
	color.Green(fmt.Sprintf("%s stopped", appName))

	return nil
}
