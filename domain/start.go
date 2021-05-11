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

func StartApps(rootDir string) error {
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
		appPath := release.GetCurrentReleaseAppFolder(rootDir, appName)
		if err := StartApp(appName, appPath, runcPath); err != nil {
			color.Red(fmt.Sprintf("error while starting %s app, error: %s", appName, err))
		}
	}
	return nil
}

func StartApp(appName, appPath, runcPath string) error {
	cmd := exec.Command(
		runcPath,
		"run",
		"--bundle",
		appPath,
		appName,
	)

	if err := cmd.Start(); err != nil {
		return err
	}
	// TODO Search how to do it without sleep
	time.Sleep(4 * time.Second)

	if err := cmd.Process.Release(); err != nil {
		return err
	}
	color.Green(fmt.Sprintf("%s application started", appName))
	return nil
}
