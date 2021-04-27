package release

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const BACKUP_FOLDER_NAME = "backups"
const CURRENT_RELEASE_FOLDER_NAME = "current-release"
const RELEASE_INFO_JSON = "release-info.json"

type RootDirMeta struct {
	currentReleaseFolder string
	logsFolder string
	configFolder string

}

func GetReleases(rootDirPath string) ([]ReleaseMeta, error) {
	backupsFolder := path.Join(rootDirPath, BACKUP_FOLDER_NAME)

	if _, err := os.Stat(backupsFolder); os.IsNotExist(err) {
		return nil, err
	}

	files, err := ioutil.ReadDir(backupsFolder)
	if err != nil {
		return nil, err
	}

	result := []ReleaseMeta{}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".json") {
			releaseMeta, err := GetReleaseInfoFromJsonFile(path.Join(backupsFolder, f.Name()))
			if err == nil {
				result = append(result, releaseMeta)
			}
		}
	}

	return result, nil
}

func GetCurrentReleaseInfo(rootDirPath string) (ReleaseMeta, error) {
	releaseInfoPath := path.Join(rootDirPath, CURRENT_RELEASE_FOLDER_NAME, RELEASE_INFO_JSON)
	if _, err := os.Stat(releaseInfoPath); os.IsNotExist(err) {
		return ReleaseMeta{}, err
	}

	return GetReleaseInfoFromJsonFile(releaseInfoPath)

}