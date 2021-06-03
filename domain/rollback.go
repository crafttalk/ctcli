package domain

import (
	"ctcli/domain/ctcliDir"
	"ctcli/domain/lifetime"
	"ctcli/util"
	"github.com/fatih/color"
	"log"
	"path"
)

func Rollback(rootDir, backupPath string) error {
	if err := lifetime.StopApps(rootDir, []string{}); err != nil {
		color.HiRedString("Error while stopping apps: %s\n", err)
		color.HiRedString("Continuing to rollback...\n")
	}

	ctcliDir.RecreateTempDir(rootDir)

	tempFolder := ctcliDir.GetTempDir(rootDir)
	log.Printf("Extracting backup %s to %s", backupPath, tempFolder)
	if err := util.ExtractTarGz(backupPath, tempFolder); err != nil {
		return err
	}

	if err := util.CopyDir(path.Join(tempFolder, "current-release"), ctcliDir.GetCurrentReleaseDir(rootDir)); err != nil {
		return err
	}

	if err := util.CopyDir(path.Join(tempFolder, "config"), ctcliDir.GetConfigDir(rootDir)); err != nil {
		return err
	}

	log.Printf("Cleaning up tmp folder\n")
	ctcliDir.DeleteTempDir(rootDir)

	return lifetime.StartApps(rootDir, []string{})
}
