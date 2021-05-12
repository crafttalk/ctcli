package domain

import (
	"ctcli/domain/release"
	"ctcli/util"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"os/exec"
	"time"
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
		if err := StopApp(appToStop, runcPath); err != nil {
			color.Red(fmt.Sprintf("error while stoping %s app, error: %s", appToStop, err))
		}
	}
	return nil
}

func StopApp(appName, runcPath string) error {
	cmd := exec.Command(
		runcPath,
		"kill",
		appName,
		"SIGTERM",
	)
	if err := cmd.Run(); err != nil {
		return err
	}

	time.Sleep(time.Second)
	cmd = exec.Command(
		runcPath,
		"delete",
		appName,
	)
	if err := cmd.Run(); err != nil {
		return err
	}
	color.Green(fmt.Sprintf("%s application stoped", appName))
	return nil
}
