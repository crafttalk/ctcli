package release

import (
	"fmt"
	"testing"
)

func TestListReleases(t *testing.T) {
	releases, err := GetReleases("/home/lkmfwe/ctcli");
	if err != nil {
		t.Errorf("folder was not found!")
	}
	fmt.Println(releases)
}

func TestGetCurrentReleaseInfo(t *testing.T) {
	releases, err := GetCurrentReleaseInfo("/home/lkmfwe/ctcli");
	if err != nil {
		t.Errorf("release info was not found!")
	}
	fmt.Println(releases)
}