package util

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/fatih/color"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func ExtractTarGz(pathToGzip string, pathToExtractTo string) error {
	gzipStream, err := os.Open(pathToGzip)
	if err != nil {
		log.Fatalf("ExtractTarGz: Open file descriptor failed: %s", err.Error())
		return err
	}

	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		log.Fatal("ExtractTarGz: NewReader failed")
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("ExtractTarGz: Next() failed: %s", err.Error())
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(path.Join(pathToExtractTo, header.Name), 0755); err != nil {
				log.Fatalf("ExtractTarGz: Mkdir() failed: %s", err.Error())
				return err
			}
		case tar.TypeReg:
			filePath := path.Join(pathToExtractTo, header.Name)
			_ = os.MkdirAll(path.Dir(filePath), os.ModePerm)

			outFile, err := os.OpenFile(path.Join(pathToExtractTo, header.Name), os.O_RDWR | os.O_CREATE, os.FileMode(header.Mode))
			if err != nil {
				log.Fatalf("ExtractTarGz: Create() failed: %s", err.Error())
				return err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				log.Fatalf("ExtractTarGz: Copy() failed: %s", err.Error())
				return err
			}
			if err := outFile.Close(); err != nil {
				log.Fatalf("ExtractTarGz: Close() failed: %s", err.Error())
				return err
			}

		default:
			log.Fatalf(
				"ExtractTarGz: uknown type: %s in %s",
				header.Typeflag,
				header.Name)
			return err
		}

	}
	return nil
}

func ArchiveTarGz(writer io.Writer, srcs ...string) error {

	gzw := gzip.NewWriter(writer)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	for ind, src := range srcs {
		// ensure the src actually exists before trying to tar it
		if _, err := os.Stat(src); err != nil {
			return fmt.Errorf("Unable to tar files - %v", err.Error())
		}

		srcDir := filepath.Dir(src)
		color.HiGreen(fmt.Sprintf("starting to archive %s (%d of %d)", strings.Replace(src, srcDir, "", -1), ind+1, len(srcs)))
		err := filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !fi.Mode().IsRegular() {
				return nil
			}

			header, err := tar.FileInfoHeader(fi, fi.Name())
			if err != nil {
				return err
			}

			header.Name = strings.TrimPrefix(strings.Replace(file, srcDir, "", -1), string(filepath.Separator))
			if err := tw.WriteHeader(header); err != nil {
				return err
			}

			f, err := os.Open(file)
			if err != nil {
				return err
			}

			if _, err := io.Copy(tw, f); err != nil {
				return err
			}
			f.Close()

			return nil
		})
		if err != nil {
			return err
		}

	}
	return nil
}
