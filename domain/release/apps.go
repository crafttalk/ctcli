package release

import (
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"os"
	"path"
	"time"
)

func GetAppNamesFromFolders(appFolders []os.FileInfo) []string {
	var appNames []string
	for _, appFolder := range appFolders {
		appNames = append(appNames, appFolder.Name())
	}
	return appNames
}

func CheckIfAppsInReleaseFolder(appsFromArgs, appsInReleaseFolder []string) ([]string, []string) {
	var appsExist []string
	var appsNotExist []string
	for _, appFromArgs := range appsFromArgs {
		appExist := false
		for _, appInReleaseFolder := range appsInReleaseFolder {
			if appFromArgs == appInReleaseFolder {
				appsExist = append(appsExist, appFromArgs)
				appExist = true
				break
			}
		}
		if !appExist {
			appsNotExist = append(appsNotExist, appFromArgs)
		}
	}
	return appsExist, appsNotExist
}

func GetReleaseAppsList(rootDir string) ([]string, error) {
	apps, err := ioutil.ReadDir(GetCurrentReleaseAppsFolder(rootDir))
	if err != nil {
		return nil, err
	}

	result := []string{}
	for _, app := range apps {
		if app.IsDir() {
			result = append(result, app.Name())
		}
	}
	return result, nil
}

func CreateReleaseInfo(rootDir string) (ReleaseMeta, error) {
	apps, err := GetReleaseAppsList(rootDir)

	releaseInfo := ReleaseMeta{}

	if err != nil {
		return ReleaseMeta{}, err
	}
	for _, app := range apps {
		versionPath := path.Join(GetCurrentReleaseAppFolder(rootDir, app), "version.json")
		version, err := GetVersionJsonFromFile(versionPath)
		if err != nil {
			return ReleaseMeta{}, err
		}

		releaseInfo.AppVersions = append(releaseInfo.AppVersions, version)
	}

	releaseInfo.PreviousRelease = "" // Todo!
	releaseInfo.Id = uuid.NewV4().String()
	releaseInfo.CreatedAt = time.Now()
	return releaseInfo, nil
}