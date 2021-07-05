package lifetime

import (
	"ctcli/domain/ctcliDir"
	"ctcli/domain/release"
	"ctcli/domain/runc"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func Exec(rootDir string, args[] string) error {
	runcPath := release.GetCurrentReleaseRuncPath(rootDir)
	appName := args[0]
	if !release.CurrentReleaseAppExists(rootDir, appName) {
		return fmt.Errorf("App \"%s\" not exists\n", appName)
	}

	execArgs := args[1:]
	execCommandArgs := append([]string{runc.GetContainerName(rootDir, appName)}, execArgs...)

	appStatus := runc.GetStatus(rootDir, appName)
	log.Printf("%s", appStatus)

	if appStatus == "unknown" {
		cmd := runc.CreateContainer(rootDir, appName)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	commandArgs := append([]string{
		"--root",
		ctcliDir.GetRuncRoot(rootDir),
		"exec",}, execCommandArgs...)
	cmd := exec.Command(runcPath, commandArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	if appStatus == "unknown" {
		if err := runc.DeleteContainer(rootDir, appName); err != nil {
			return err
		}
	}

	return nil
}