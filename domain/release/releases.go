package release

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func GetReleases(rootDir string) ([]ReleaseMeta, error) {
	releasesFolder := GetReleasesFolder(rootDir)

	if _, err := os.Stat(releasesFolder); os.IsNotExist(err) {
		return nil, err
	}

	files, err := ioutil.ReadDir(releasesFolder)
	if err != nil {
		return nil, err
	}

	result := []ReleaseMeta{}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".json") {
			releaseMeta, err := GetReleaseInfoFromJsonFile(path.Join(releasesFolder, f.Name()))
			if err == nil {
				result = append(result, releaseMeta)
			}
		}
	}

	return result, nil
}

func GetCurrentReleaseInfo(rootDir string) (ReleaseMeta, error) {
	releaseInfoPath := GetCurrentReleaseInfoPath(rootDir)
	if _, err := os.Stat(releaseInfoPath); os.IsNotExist(err) {
		return ReleaseMeta{}, err
	}

	return GetReleaseInfoFromJsonFile(releaseInfoPath)

}