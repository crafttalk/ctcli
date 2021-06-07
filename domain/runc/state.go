package runc

import (
	"crypto/md5"
	"ctcli/domain/ctcliDir"
	"ctcli/domain/release"
	"fmt"
	"github.com/valyala/fastjson"
	"os/exec"
	"time"
)

func GetContainerName(rootDir string, appName string) string {
	return fmt.Sprintf("%x-%s", md5.Sum([]byte(rootDir)), appName)
}

func GetStatus(rootDir string, appName string) string {
	runcPath := release.GetCurrentReleaseRuncPath(rootDir)
	out, err := exec.Command(
		runcPath,
		"--root",
		ctcliDir.GetRuncRoot(rootDir),
		"state",
		GetContainerName(rootDir, appName),
	).Output()
	if err != nil {
		return "unknown"
	}
	jsonResult := fastjson.MustParse(string(out))
	return string(jsonResult.GetStringBytes("status"))
}

func GetPid(rootDir string, appName string) int {
	runcPath := release.GetCurrentReleaseRuncPath(rootDir)
	out, err := exec.Command(
		runcPath,
		"--root",
		ctcliDir.GetRuncRoot(rootDir),
		"state",
		GetContainerName(rootDir, appName),
	).Output()
	if err != nil {
		return 0
	}
	jsonResult := fastjson.MustParse(string(out))
	return jsonResult.GetInt("pid")
}

func WaitUntilNotRunning(rootDir string, appName string) {
	for i := 0; i < 10; i++ {
		status := GetStatus(rootDir, appName)
		if status != "running" {
			break
		}
		time.Sleep(time.Second)
	}
}

func DeleteContainer(rootDir string, appName string) error {
	runcPath := release.GetCurrentReleaseRuncPath(rootDir)
	err := exec.Command(
		runcPath,
		"--root",
		ctcliDir.GetRuncRoot(rootDir),
		"delete",
		GetContainerName(rootDir, appName),
	).Run()
	if err != nil {
		return err
	}
	return nil
}