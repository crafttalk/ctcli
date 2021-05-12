package release

import "os"

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
