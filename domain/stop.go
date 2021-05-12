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

func StopApps(rootDir string) error {
	runcPath := release.GetCurrentReleaseRuncPath(rootDir)
	if !util.PathExists(runcPath) {
		return fmt.Errorf("there is no runc in current-relase folder")
	}
	appsPath := release.GetCurrentReleaseAppsFolder(rootDir)
	appFolders, err := ioutil.ReadDir(appsPath)
	if err != nil {
		return err
	}
	for _, appFolder := range appFolders {
		appName := appFolder.Name()
		if err := StopApp(rootDir, appName, runcPath); err != nil {
			color.Red(fmt.Sprintf("error while stopping %s app, error: %s", appName, err))
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
			appName,
			"SIGTERM",
		)
		if err := cmd.Run(); err != nil {
			return err
		}

		runc.WaitUntilNotRunning(rootDir, appName)
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
