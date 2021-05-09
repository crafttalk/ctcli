package domain

import (
	"ctcli/domain/ctcliDir"
	"ctcli/domain/moving"
	"ctcli/domain/packaging"
	"ctcli/domain/release"
	"ctcli/util"
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"strings"
)

func AddAppDiffToText(diffText *string, newAppVersion, currentAppVersion release.AppVersion) error {
	appKeyColor := color.New(color.FgHiCyan)
	appValueColor := color.New(color.FgRed)
	oldKeyColor := color.New(color.FgYellow)
	newKeyColor := color.New(color.FgGreen)
	keyValueColor := color.New(color.FgBlue)

	*diffText += appKeyColor.Sprintf("  app: ")
	*diffText += appValueColor.Sprintln(newAppVersion.AppName)
	*diffText += oldKeyColor.Sprintf("    old image: ")
	*diffText += keyValueColor.Sprintln(currentAppVersion.Image)
	*diffText += newKeyColor.Sprintf("    new image: ")
	*diffText += keyValueColor.Sprintln(newAppVersion.Image)
	*diffText += oldKeyColor.Sprintf("    old imageSha: ")
	*diffText += keyValueColor.Sprintln(currentAppVersion.ImageSha)
	*diffText += newKeyColor.Sprintf("    new imageSha: ")
	*diffText += keyValueColor.Sprintln(newAppVersion.ImageSha)
	return nil
}

func MakeReleasesDifferenceRow(newReleaseAppVersions, currentReleaseAppVersions []release.AppVersion) (string, error) {
	diffText := ""
	for _, newAppVersion := range newReleaseAppVersions {
		appExistInCurrentRelease := false
		for _, currentAppVersion := range currentReleaseAppVersions {
			if newAppVersion.AppName == currentAppVersion.AppName {
				appExistInCurrentRelease = true
				if err := AddAppDiffToText(&diffText, newAppVersion, currentAppVersion); err != nil {
					return "", err
				}
				break
			}
		}
		if !appExistInCurrentRelease {
			if err := AddAppDiffToText(&diffText, newAppVersion, release.AppVersion{}); err != nil {
				return "", err
			}
		}
	}
	return diffText, nil
}

func ShowReleasesDifference(newReleaseAppVersions, currentReleaseAppVersions []release.AppVersion) error {
	diffText, err := MakeReleasesDifferenceRow(newReleaseAppVersions, currentReleaseAppVersions)
	if err != nil {
		return err
	}
	fmt.Printf("%s %s", color.New(color.FgHiMagenta).Sprintln("Services to be upgraded:"), diffText)
	return nil
}

func AskDoUpgrade() bool {
	fmt.Println("Do Upgrade? (Y/N)")
	var answer string
	var doUpgrade bool
	userAnswered := false

	for !userAnswered {
		fmt.Fscan(os.Stdin, &answer)
		answer = strings.ToLower(answer)
		if answer == "y" || answer == "yes" {
			userAnswered = true
			doUpgrade = true
		} else if answer == "n" || answer == "no" {
			userAnswered = true
			doUpgrade = false
		} else {
			fmt.Println("You need to answer \"y\"(yes) or \"n\"(no)")
		}
	}
	return doUpgrade
}

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

	ctcliDir.RecreateTempDir(rootDir)

	if err := util.ExtractTarGz(packagePath, tempFolder); err != nil {
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

	if err := ShowReleasesDifference(newReleaseInfo.AppVersions, currentReleaseInfo.AppVersions); err != nil {
		return err
	}

	doUpgrade := AskDoUpgrade()
	if doUpgrade {
		if err := moving.LoadRelease(rootDir, tempFolder); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("upgrade was cancelled by user")
		// TODO: Delete tmp folder
	}
	return nil
}
