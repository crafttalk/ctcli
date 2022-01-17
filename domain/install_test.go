package domain

import (
	"ctcli/domain/ctcliDir"
	"ctcli/domain/release"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestGetCurrentReleaseInfo(t *testing.T) {
	rootDir, err := ioutil.TempDir("/tmp/", "ctcli")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(rootDir)

	ctcliDir.Init(rootDir)

	if err := Install(rootDir, "../data/test-package.tar.gz"); err != nil {
		t.Errorf("Install failed %s", err)
	}

	releases, err := release.GetCurrentReleaseInfo(rootDir)
	if err != nil {
		t.Errorf("release info was not found!")
	}
	fmt.Println(releases)
}
