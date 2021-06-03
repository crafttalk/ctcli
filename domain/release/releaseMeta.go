package release

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

type AppVersion struct {
	AppName string `json:"appName"`
	Image string `json:"image"`
	ImageSha string `json:"imageSha"`
	ImageTag string `json:"imageTag"`
	CommitSha string `json:"commit"`
	BuiltAt time.Time `json:"builtAt"`
}

type ReleaseMeta struct {
	Id string `json:"id"`
	PreviousRelease string `json:"baseRelease"`
	CreatedAt time.Time `json:"createdAt"`
	AppVersions []AppVersion `json:"appVersions"`
}

func GetReleaseInfoFromJsonFile(filePath string) (ReleaseMeta, error) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return ReleaseMeta{}, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var releaseMeta ReleaseMeta
	err = json.Unmarshal(byteValue, &releaseMeta)
	if err != nil {
		return ReleaseMeta{}, err
	}

	return releaseMeta, nil
}

func GetVersionJsonFromFile(versionJsonFilePath string) (AppVersion, error) {
	jsonFile, err := os.Open(versionJsonFilePath)
	if err != nil {
		return AppVersion{}, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var releaseMeta AppVersion
	err = json.Unmarshal(byteValue, &releaseMeta)
	if err != nil {
		return AppVersion{}, err
	}

	return releaseMeta, nil
}