package cmd

import (
	"ctcli/domain"
	"ctcli/domain/ctcliDir"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"log"
	"path/filepath"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade <path to package>",
	Short: "upgrade a release",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rootFlag := cmd.Flag("root")
		rootDir, err := filepath.Abs(rootFlag.Value.String())
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		packagePath, err := filepath.Abs(args[0])
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		if err := ctcliDir.OkIfIsARootDir(rootDir); err != nil {
			cmd.PrintErr(err)
			return
		}
		err = domain.Upgrade(rootDir, packagePath)
		if err != nil {
			log.Fatal(err)
		}
		color.Green("OK\n")
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
