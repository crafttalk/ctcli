package ctcliDir

import (
	"ctcli/util"
	"os"
	"path"
	"path/filepath"
	"time"
)

func FilePathWalkDir(root string) []string {
    var files []string
    filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err == nil {
			if !info.IsDir() {
				files = append(files, path)
			}
		}
		return nil
    })
    return files
}

func ArchivesOfAppLifeTime(rootDir string, app string, days int64) {
	archivesPath := path.Join(GetAppLogsDir(rootDir, app), "archives")
	util.CreateDirIfNotExist(GetAppLogsDir(rootDir, app), "archives")
	
	files := FilePathWalkDir(archivesPath)
	for _, file := range files {
		fileStat, err := os.Stat(file)
		if err == nil {
			diffHours := time.Now().Sub(fileStat.ModTime()).Hours()

			if int64(diffHours / 24) >= days {
				_ = os.Remove(file)
			}
		}
	}
}

func ArchiveLifeTime(rootDir string, apps []string, days int64)  {
	for {
		for _, app := range apps {
			ArchivesOfAppLifeTime(rootDir, app, days)
		}
		
		time.Sleep(time.Duration(10) * time.Minute)
	}
}