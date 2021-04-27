package cmd

import (
	"crafttalk-cli/util"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path"
	"path/filepath"
)

var installCmd = &cobra.Command{
	Use: "install [path to package]",
	Short: "Install a release",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		packagePath, err := filepath.Abs(args[0])
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		if _, err := os.Stat(packagePath); os.IsNotExist(err) {
			log.Fatalf("Couldn't find package %s", packagePath)
			return
		}

		rootFlag := cmd.Flag("root")
		rootDir, err := filepath.Abs(rootFlag.Value.String())
		if err != nil {
			cmd.PrintErr(err)
			return
		}


		log.Printf("Root dir: %s", rootDir)
		_ = os.MkdirAll(rootDir, os.ModePerm)

		tempFolder := path.Join(rootDir, "tmp")
		log.Printf("Extracting package %s to %s", packagePath, tempFolder)
		_ = os.RemoveAll(tempFolder)
		_ = os.MkdirAll(tempFolder, os.ModePerm)

		err = util.ExtractTarGz(packagePath, tempFolder)
		if err != nil {
			cmd.PrintErr(err)
			return
		}
	},
}

func init()  {
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	installCmd.Flags().String("root", workDir, "root of the installation")
	rootCmd.AddCommand(installCmd)
}
