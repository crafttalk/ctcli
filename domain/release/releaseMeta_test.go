package release

import (
	"fmt"
	"testing"
)

func TestGetReleaseInfoFromJsonFile(t *testing.T) {
	releaseInfo, err := GetReleaseInfoFromJsonFile("/home/lkmfwe/ctcli/backups/some-release.json")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(releaseInfo)
}
