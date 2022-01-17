package release

import (
	"ctcli/domain/ctcliDir"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestListReleases(t *testing.T) {
	rootDir, err := ioutil.TempDir("/tmp/", "ctcli")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(rootDir)

	ctcliDir.Init(rootDir)

	releases, err := GetReleases(rootDir)
	if err != nil {
		t.Errorf("folder was not found!")
	}
	fmt.Println(releases)
}
