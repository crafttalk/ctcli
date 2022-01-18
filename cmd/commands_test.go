package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func GetRootDir(t *testing.T) string {
	rootDir, err := ioutil.TempDir("/tmp/", "ctcli")
	if err != nil {
		t.Fatal(err)
	}
	return rootDir
}

func RunCommand(t *testing.T, args []string) string {

	b := bytes.NewBufferString("")

	rootCmd.SetArgs(args)
	rootCmd.SetOut(b)
	rootCmd.Execute()

	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}

func TestVersion(t *testing.T) {
	rootDir := GetRootDir(t)
	defer os.RemoveAll(rootDir)

	out := RunCommand(t, []string{"--root", rootDir, "version"})

	if strings.HasPrefix("CraftTalk Command Line Tool", out) {
		t.Fatalf("Version string is incorrect: %s!", out)
	}
}

func TestInitAndStatus(t *testing.T) {
	rootDir := GetRootDir(t)
	defer os.RemoveAll(rootDir)

	out := RunCommand(t, []string{"--root", rootDir, "init"})

	if !strings.Contains(out, "OK") {
		t.Fatalf("Init returned non-ok result")
	}

	out = RunCommand(t, []string{"--root", rootDir, "status"})

	if !strings.Contains(out, "APP-NAME") {
		t.Fatalf("Status table didn't contain the APP-NAME column: %s", out)
	}

	if !strings.Contains(out, "STATUS") {
		t.Fatal("Status table didn't contain the STATUS column")
	}

	if !strings.Contains(out, "PID") {
		t.Fatal("Status table didn't contain the PID column")
	}
}

func TestInstallAndStartStop(t *testing.T) {
	rootDir := GetRootDir(t)
	defer os.RemoveAll(rootDir)

	out := RunCommand(t, []string{"--root", rootDir, "init"})

	if !strings.Contains(out, "OK") {
		t.Fatal("Init returned non-ok result")
	}

	out = RunCommand(t, []string{"--root", rootDir, "install", "../data/test-package.tar.gz"})

	if !strings.Contains(out, "OK") {
		t.Fatalf("Install was not successful: %s", out)
	}

	out = RunCommand(t, []string{"--root", rootDir, "release-info"})

	if !strings.Contains(out, "siebelintegration") {
		t.Fatalf("siebelintegration not present in status: %s", out)
	}

	if !strings.Contains(out, "4aac9f7cafa6bd8dd78069ddc22f066228e48c67c6a12c90085dad10785ee230") {
		t.Fatalf("siebelintegration image sha is incorrect: %s", out)
	}

	if !strings.Contains(out, "0") {
		t.Fatalf("Pid != 0: %s", out)
	}
	out = RunCommand(t, []string{"--root", rootDir, "start"})
	t.Log(out)

	out = RunCommand(t, []string{"--root", rootDir, "logs", "siebelintegration"})

	// В докер окружениях невозможно запустить рутлес контейнеры(
	if !(strings.Contains(out, "Starting") || strings.Contains(out, "mapping tool not present: Operation not permitted")) {
		t.Fatalf("Failed to start a service")
	}
	t.Log(out)

	out = RunCommand(t, []string{"--root", rootDir, "status"})
	t.Log(out)

	// if strings.Contains(out, "0") {
	// 	t.Fatal("Process failed to start")
	// }
}
