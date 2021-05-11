package appConfig

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type HealthCheck struct {
	Command []string `json:"command"`
	WaitFor int      `json:"waitFor"`
}

type AppPackageConfig struct {
	BaseImage   string      `json:"baseImage"`
	Healthcheck HealthCheck `json:"healthcheck"`
	LogsFolder  string      `json:"logsFolder"`
	Configs     []string    `json:"configs"`
	Data 		[]string	`json:"data"`
}

func GetAppConfig(path string) (AppPackageConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return AppPackageConfig{}, err
	}
	defer file.Close()
	bytesFromFile, err := ioutil.ReadAll(file)
	if err != nil {
		return AppPackageConfig{}, err
	}
	var config AppPackageConfig
	if err := json.Unmarshal(bytesFromFile, &config); err != nil {
		return AppPackageConfig{}, err
	}
	return config, nil
}
