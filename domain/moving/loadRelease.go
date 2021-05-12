package moving

import (
	"bytes"
	"ctcli/domain/appConfig"
	"ctcli/domain/ctcliDir"
	"ctcli/domain/packaging"
	"ctcli/domain/release"
	"ctcli/util"
	"encoding/json"
	"fmt"
	"github.com/otiai10/copy"
	"github.com/valyala/fastjson"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func CopyBinariesToRelease(rootDir string, packagePath string) error {
	if err := copy.Copy(packaging.GetPackageRuncPath(packagePath), release.GetCurrentReleaseRuncPath(rootDir)); err != nil {
		return err
	}
	if err := copy.Copy(packaging.GetPackageAppsFolder(packagePath), release.GetCurrentReleaseAppsFolder(rootDir)); err != nil {
		return err
	}
	return nil
}

func createMountValue(pathFrom string, pathTo string) *fastjson.Value {
	result := fastjson.MustParse("{}")
	result.Set("destination", fastjson.MustParse(fmt.Sprintf(`"%s"`, pathTo)))
	result.Set("source", fastjson.MustParse(fmt.Sprintf(`"%s"`, pathFrom)))
	result.Set("type", fastjson.MustParse(`"bind"`))
	result.Set("options", fastjson.MustParse(`["noexec", "nosuid", "rbind"]`))
	return result
}

func configureRuncConfig(rootDir string, app string, config appConfig.AppPackageConfig, runcConfigPath string) error {
	jsonBytes, err := ioutil.ReadFile(runcConfigPath)
	if err != nil {
		return err
	}

	jsonValue, err := fastjson.ParseBytes(jsonBytes)
	if err != nil {
		return err
	}

	mounts, err := jsonValue.Get("mounts").Array()
	if err != nil {
		return err
	}

	jsonValue.Get("process").Set("terminal", fastjson.MustParse(`false`))
	if config.Entrypoint != nil {
		resultValue := fastjson.MustParse("[]")
		for i, cmdPart := range config.Entrypoint {
			resultValue.SetArrayItem(i, fastjson.MustParse(`"` + strings.Replace(cmdPart, `"`, `\"`, -1) + `"`))
		}
		jsonValue.Get("process").Set("args", resultValue)
	}

	mountsMap := map[string]*fastjson.Value{}

	for _, mount := range mounts {
		mountDst := string(mount.GetStringBytes("destination"))
		mountsMap[mountDst] = mount
	}

	if config.LogsFolder != "" {
		if _, contains := mountsMap[config.LogsFolder]; contains {
			delete(mountsMap, config.LogsFolder)
		}
		rootDirLogPath := ctcliDir.GetAppLogsDir(rootDir, app)
		_ = os.MkdirAll(rootDirLogPath, os.ModePerm)
		mountsMap[config.LogsFolder] = createMountValue(rootDirLogPath, config.LogsFolder)
	}
	for _, dataPath := range config.Data {
		if _, contains := mountsMap[dataPath]; contains {
			delete(mountsMap, dataPath)
		}

		rootDirDataPath := path.Join(ctcliDir.GetAppDataDir(rootDir, app), dataPath)
		_ = os.MkdirAll(rootDirDataPath, os.ModePerm)
		mountsMap[dataPath] = createMountValue(rootDirDataPath, dataPath)
	}
	for _, configPath := range config.Configs {
		if _, contains := mountsMap[configPath]; contains {
			delete(mountsMap, configPath)
		}

		absContainerConfigPath := path.Join(release.GetCurrentReleaseAppRootfsFolder(rootDir, app), configPath)
		rootDirConfigPath := path.Join(ctcliDir.GetAppConfigDir(rootDir, app), configPath)
		_ = os.MkdirAll(path.Dir(rootDirConfigPath), os.ModePerm)
		if !util.PathExists(rootDirConfigPath) {
			if err := util.CopyFile(absContainerConfigPath, rootDirConfigPath); err != nil {
				return err
			}
		}

		mountsMap[configPath] = createMountValue(rootDirConfigPath, configPath)
	}

	resultMountsJson := ""
	i := 0
	for _, mount := range mountsMap {
		if i > 0 {
			resultMountsJson += ","
		}
		resultMountsJson += mount.String()
		i++
	}

	jsonValue.Set("mounts", fastjson.MustParse("[" + resultMountsJson + "]"))
	resultJson := []byte(jsonValue.String())

	buf := bytes.Buffer{}
	if err := json.Indent(&buf, resultJson, "", "    "); err != nil {
		return err
	}

	if err := ioutil.WriteFile(runcConfigPath, buf.Bytes(), os.ModePerm); err != nil {
		return err
	}
	return nil
}

func CopyPackagesToRelease(rootDir string, packagePath string) error {
	apps, err := packaging.GetPackageAppsList(packagePath)
	if err != nil {
		return err
	}

	for _, app := range apps {
		appPackageConfigPath := release.GetCurrentReleasePackageConfigPath(rootDir, app)
		appPackageConfig, err := appConfig.GetAppConfig(appPackageConfigPath)
		if err != nil {
			return err
		}

		if err := configureRuncConfig(rootDir, app, appPackageConfig, release.GetCurrentReleaseRuncConfigPath(rootDir, app)); err != nil {
			return err
		}
	}
	return nil
}

func CreateReleaseInfoFile(rootDir string, tempFolder string) error {
	releaseInfo, err := packaging.CreateReleaseInfo(tempFolder)
	if err != nil {
		return err
	}
	content, _ := json.MarshalIndent(releaseInfo, "", " ")
	err = ioutil.WriteFile(release.GetCurrentReleaseInfoPath(rootDir), content, 0644)
	return err
}

func LoadRelease(rootDir, tempFolder string) error {
	if err := util.RemoveContentOfFolder(ctcliDir.GetCurrentReleaseDir(rootDir)); err != nil {
		return err
	}
	if err := CopyBinariesToRelease(rootDir, tempFolder); err != nil {
		return err
	}
	if err := CopyPackagesToRelease(rootDir, tempFolder); err != nil {
		return err
	}
	if err := CreateReleaseInfoFile(rootDir, tempFolder); err != nil {
		return err
	}
	return nil
}
