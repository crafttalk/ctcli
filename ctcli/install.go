package ctcli

import (
	"ctcli/util"
	"fmt"
	"log"
	"os"
	"path"
)

func extractBlobs() {
	//# skopeo_image_path = os.path.join(tmp_abs_path, f"skopeo-{name}")
	//# runc_bundle_path = os.path.join(tmp_abs_path, f"runc-bundle-{name}")
	//#
	//# run_umoci = [
	//#     "./umoci.amd64",
	//#     "unpack",
	//#     "--rootless",
	//#     "--image",
	//#     f"{skopeo_image_path}:latest",
	//#     runc_bundle_path
	//# ]
	//# output = subprocess.run(run_umoci)
	//# if output.returncode != 0:
	//#     print(output)
	//#     raise OSError(f"Failed to create runc bundle from docker image {image} to {container_tmp_dir}")
	//#
	//# shutil.rmtree(skopeo_image_path)
	//# shutil.move(os.path.join(runc_bundle_path, "rootfs"), os.path.join(container_tmp_dir, "rootfs"))
	//# shutil.move(os.path.join(runc_bundle_path, "config.json"), os.path.join(container_tmp_dir, "config.json"))
	//# shutil.rmtree(runc_bundle_path)-
}

func Install(rootDir string, packagePath string) error {
	if !util.PathExists(packagePath) {
		return fmt.Errorf("couldn't find package %s", packagePath)
	}
	if !util.PathExists(rootDir) {
		return fmt.Errorf("root dir %s doesn't exists", rootDir)
	}

	log.Printf("Root dir: %s", rootDir)
	_ = os.MkdirAll(rootDir, os.ModePerm)

	tempFolder := path.Join(rootDir, "tmp")
	log.Printf("Extracting package %s to %s", packagePath, tempFolder)
	//_ = os.RemoveAll(tempFolder)
	//_ = os.MkdirAll(tempFolder, os.ModePerm)
	//
	//err := util.ExtractTarGz(packagePath, tempFolder)
	//if err != nil {
	//	return err
	//}

	return nil
}