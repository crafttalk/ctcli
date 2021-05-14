package domain

import (
	"ctcli/domain/ctcliDir"
	"ctcli/domain/release"
	"ctcli/util"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"os"
	"path"
	"time"
)

const (
	BackupInfoFile = "backup-info.json"
)

type BackupInfo struct {
	Uid         string              `json:"uid"`
	BackupDate  string              `json:"backupDate"`
	ReleaseInfo release.ReleaseMeta `json:"releaseInfo"`
}

func MakeBackupInfoFile(backupInfoFilePath string, backupInfo BackupInfo) error {
	json, err := json.MarshalIndent(backupInfo, "", " ")
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(backupInfoFilePath, json, 0644); err != nil {
		return err
	}
	return nil
}

func MakeABackup(rootDir string) error {
	if err := ctcliDir.OkIfIsARootDir(rootDir); err != nil {
		return err
	}

	currentReleasePath := ctcliDir.GetCurrentReleaseDir(rootDir)
	dataPath := ctcliDir.GetDataDir(rootDir)
	configPath := ctcliDir.GetConfigDir(rootDir)

	releaseInfo, err := release.GetCurrentReleaseInfo(rootDir)
	if err != nil {
		return err
	}

	backupUid := uuid.NewV4().String()
	backupDate := time.Now().Format("02-01-2006 15:04:05 MST")

	backupInfo := BackupInfo{
		Uid:         backupUid,
		BackupDate:  backupDate,
		ReleaseInfo: releaseInfo,
	}

	releasesPath := ctcliDir.GetReleasesDir(rootDir)

	currentBackupFolderPath := path.Join(releasesPath, backupDate)
	if err := os.Mkdir(currentBackupFolderPath, os.ModePerm); err != nil {
		return err
	}
	backupInfoFilePath := path.Join(currentBackupFolderPath, BackupInfoFile)
	if err := MakeBackupInfoFile(backupInfoFilePath, backupInfo); err != nil {
		return err
	}

	tarfile, err := os.Create(path.Join(currentBackupFolderPath, backupUid+".tar.gz"))
	defer tarfile.Close()

	if err = util.ArchiveTarGz(tarfile, currentReleasePath, dataPath, configPath); err != nil {
		return err
	}
	return nil
}
