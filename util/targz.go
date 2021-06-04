package util

import (
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

func ExtractTarGz(archivePath string, pathToExtractTo string) error {
	args := append([]string{ "-xzvf", archivePath, "-C", pathToExtractTo })

	log.Printf("tar %s\n", strings.Join(args, " "))

	cmd := exec.Command(
		"tar",
		args...
	)

	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	if err := cmd.Process.Release(); err != nil {
		return err
	}
	return nil
}

func ArchiveTarGz(archivePath string, srcs ...string) error {
	sourceDir := path.Dir(srcs[0])
	for i, src := range srcs {
		srcs[i] = strings.Replace(src, sourceDir + "/", "", -1)
	}

	args := append([]string{ "-czvf", archivePath, "-C", sourceDir }, srcs...)

	log.Printf("tar %s\n", strings.Join(args, " "))

	cmd := exec.Command(
		"tar",
		args...
	)

	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	if err := cmd.Process.Release(); err != nil {
		return err
	}
	return nil
}
