package domain

import (
	"ctcli/domain/ctcliDir"
	"ctcli/domain/packaging"
	"ctcli/domain/release"
	"ctcli/util"
	"fmt"
	"github.com/kataras/tablewriter"
	"log"
	"os"
)

const (
	AppColumn         = "app"
	OldImageColumn    = "old image"
	OldImageShaColumn = "old imageSha"
	NewImageColumn    = "new image"
	NewImageShaColumn = "new imageSha"
)

func Upgrade(rootDir, packagePath string) error {
	if !util.PathExists(packagePath) {
		return fmt.Errorf("couldn't find package %s", packagePath)
	}
	if !util.PathExists(rootDir) {
		return fmt.Errorf("root dir %s doesn't exists", rootDir)
	}
	if err := ctcliDir.OkIfIsARootDir(rootDir); err != nil {
		return err
	}

	tempFolder := ctcliDir.GetTempDir(rootDir)
	log.Printf("Extracting package %s to %s", packagePath, tempFolder)

	ctcliDir.DeleteTempDir(rootDir)
	ctcliDir.CreateTempDir(rootDir)

	err := util.ExtractTarGz(packagePath, tempFolder)
	if err != nil {
		return err
	}

	if err := packaging.PreparePackage(tempFolder); err != nil {
		return err
	}
	currentReleaseInfo, err := release.GetCurrentReleaseInfo(rootDir)
	if err != nil {
		return err
	}
	newReleaseInfo, err := packaging.CreateReleaseInfo(tempFolder)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	// Todo: Table looks like shit! Need to remake without using tablewriter
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{AppColumn, OldImageColumn, OldImageShaColumn, NewImageColumn, NewImageShaColumn})
	tableContent := [][]string{}
	for _, newAppversion := range newReleaseInfo.AppVersions {
		appExistInCurrentRelease := false
		for _, currentAppVersion := range currentReleaseInfo.AppVersions {
			if newAppversion.AppName == currentAppVersion.AppName {
				appExistInCurrentRelease = true
				row := []string{
					newAppversion.AppName,
					currentAppVersion.Image,
					currentAppVersion.ImageSha,
					newAppversion.Image,
					newAppversion.ImageSha,
				}
				tableContent = append(tableContent, row)
				break
			}
		}
		if !appExistInCurrentRelease {
			row := []string{
				newAppversion.AppName,
				"* ",
				"* ",
				newAppversion.Image,
				newAppversion.ImageSha,
			}
			tableContent = append(tableContent, row)
		}
	}
	table.AppendBulk(tableContent)
	table.Render()
	return nil
}
