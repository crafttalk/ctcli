package packaging

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
)

func chmodPackageBinaries(packageDir string) error {
	if err := os.Chmod(GetPackageUmociPath(packageDir), 0775); err != nil {
		return err
	}
	if err := os.Chmod(GetPackageRuncPath(packageDir), 0775); err != nil {
		return err
	}

	return nil
}

func extractBlobs(umociPath string, containerTmpPath string, name string) error {
	var skopeoImagePath = path.Join(containerTmpPath, SkopeoFolder)
	var runcBundlePath = path.Join(containerTmpPath, "runc-bundle")

	cmd := exec.Command(
		umociPath,
		"unpack",
		"--rootless",
		"--image",
		fmt.Sprintf("%s:%s", skopeoImagePath, name),
		runcBundlePath,
	)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("original error: %s\nArgs:%s\nStdout:\n%s\n", err, cmd.Args, stdoutStderr)
	}

	if err := os.RemoveAll(skopeoImagePath); err != nil {
		return err
	}
	if err := os.Rename(path.Join(runcBundlePath, RuncRootFs), path.Join(containerTmpPath, RuncRootFs)); err != nil {
		return err
	}
	if err := os.Rename(path.Join(runcBundlePath, RuncConfigJson), path.Join(containerTmpPath, RuncConfigJson)); err != nil {
		return err
	}
	if err := os.RemoveAll(runcBundlePath); err != nil {
		return err
	}
	return nil
}

func PreparePackage(packageFolder string) error {
	if err := chmodPackageBinaries(packageFolder); err != nil {
		return err
	}

	apps, err := GetPackageAppsList(packageFolder)
	if err != nil {
		return err
	}

	for _, app := range apps {
		var containerTmpPath = path.Join(GetPackageAppsFolder(packageFolder), app)
		log.Printf("Extracting blob %s", containerTmpPath)
		if err := extractBlobs(GetPackageUmociPath(packageFolder), containerTmpPath, app); err != nil {
			return err
		}
	}

	return nil
}

func GetPackageAppsList(packageFolder string) ([]string, error) {
	appsPath := path.Join(packageFolder, PackageContentFolder, AppsFolder)
	apps, err := ioutil.ReadDir(appsPath)
	if err != nil {
		return nil, err
	}

	result := []string{}
	for _, app := range apps {
		result = append(result, app.Name())
	}
	return result, nil
}