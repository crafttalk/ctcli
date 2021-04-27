package util

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"os"
	"path"
)

func ExtractTarGz (pathToGzip string, pathToExtractTo string) error {
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
			outFile, err := os.Create(path.Join(pathToExtractTo, header.Name))
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