package domain

import (
	"ctcli/domain/ctcliDir"
	"ctcli/domain/moving"
	"ctcli/domain/packaging"
	"ctcli/domain/release"
	"ctcli/util"
	"fmt"
	"log"
)

func Install(rootDir string, packagePath string) error {
	if !util.PathExists(packagePath) {
		return fmt.Errorf("couldn't find package %s", packagePath)
	}
	if !util.PathExists(rootDir) {
		return fmt.Errorf("root dir %s doesn't exists", rootDir)
	}
	if err := ctcliDir.OkIfIsARootDir(rootDir); err != nil {
		return err
	}
	if util.PathExists(release.GetCurrentReleaseInfoPath(rootDir)) {
		return fmt.Errorf("current release already installed. Maybe you intended to use upgrade?")
	}

	log.Printf("Installing in root dir: %s", rootDir)

	tempFolder := ctcliDir.GetTempDir(rootDir)
	log.Printf("Extracting package %s to %s", packagePath, tempFolder)

	ctcliDir.RecreateTempDir(rootDir)

	err := util.ExtractTarGz(packagePath, tempFolder)
	if err != nil {
		return err
	}

	if err := packaging.PreparePackage(tempFolder); err != nil {
		return err
	}

	if err := moving.LoadRelease(rootDir, tempFolder); err != nil {
		return err
	}

	log.Printf("Cleaning up tmp folder\n")
	ctcliDir.DeleteTempDir(rootDir)

	return nil
}
