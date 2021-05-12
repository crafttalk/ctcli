package runc

import (
	"ctcli/domain/release"
	"github.com/valyala/fastjson"
	"os/exec"
	"time"
)

func GetStatus(rootDir string, appName string) string {
	runcPath := release.GetCurrentReleaseRuncPath(rootDir)
	out, err := exec.Command(
		runcPath,
		"state",
		appName,
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
		"state",
		appName,
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
		"delete",
		appName,
	).Run()
	if err != nil {
		return err
	}
	return nil
}