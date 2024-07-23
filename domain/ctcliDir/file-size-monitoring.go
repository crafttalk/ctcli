package ctcliDir

import (
	"ctcli/util"
	"fmt"
	"os"
	"path"
	"sync"
	"time"
)

var rw sync.RWMutex

func getFileSize(path string) (int64, error) {
	file, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return file.Size(), nil
}

func checkFileSize(rootDir string, app string, maxSize int64) {
	logFilePath := GetAppStdoutLogFilePath(rootDir, app)
	fileSize, _ := getFileSize(logFilePath)

	if fileSize > maxSize * 1024 * 1024 {
		rw.Lock()
		defer rw.Unlock()
		archiveFileName := fmt.Sprintf("%s.tar.gz", time.Now().UTC().Format("2006-01-02_15-04-05"))
		util.CreateDirIfNotExist(GetAppLogsDir(rootDir, app), "archives")
		err := util.ArchiveTarGz(path.Join(GetAppLogsDir(rootDir, app), "archives", archiveFileName), GetAppLogsDir(rootDir, app) + "/stdout-stderr.log")
		if err == nil {
			os.Truncate(logFilePath, 0)
		}
	}
}

func CheckFilesSize(rootDir string, apps []string, maxSize int64) {
	for {
		for _, app := range apps {
			checkFileSize(rootDir, app, maxSize)
		}
		
		time.Sleep(time.Duration(500) * time.Millisecond)
	}
}

