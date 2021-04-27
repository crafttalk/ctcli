package cmd

import (
	"ctcli/ctcli"
	"github.com/spf13/cobra"
	"log"
	"os"
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

		rootFlag := cmd.Flag("root")
		rootDir, err := filepath.Abs(rootFlag.Value.String())
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		err = ctcli.Install(rootDir, packagePath)
		if err != nil {
			log.Fatal(err)
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
