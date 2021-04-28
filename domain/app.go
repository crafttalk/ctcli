package domain

type HealthCheck struct {
	Command []string `json:"command"`
	WaitFor int      `json:"waitFor"`
}

type AppPackageConfig struct {
	BaseImage   string      `json:"baseImage"`
	Healthcheck HealthCheck `json:"healthcheck"`
	LogsFolder  string      `json:"logsFolder"`
	Configs     []string    `json:"configs"`
}