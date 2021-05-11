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
		if err := StopApp(appName, runcPath); err != nil {
			color.Red(fmt.Sprintf("error while stoping %s app, error: %s", appName, err))
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
