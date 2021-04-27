package cmd

import (
	"ctcli/domain"
	"github.com/spf13/cobra"
	"log"
	"path/filepath"
)

var installCmd = &cobra.Command{
	Use: "install [path to package]",
	Short: "install a release",
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

		err = domain.Install(rootDir, packagePath)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init()  {
	rootCmd.AddCommand(installCmd)
}
