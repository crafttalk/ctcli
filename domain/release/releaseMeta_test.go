package release

import (
	"fmt"
	"testing"
)

func TestGetReleaseInfoFromJsonFile(t *testing.T) {
	releaseInfo, err := GetReleaseInfoFromJsonFile("../../data/test-package.json")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(releaseInfo)
}
