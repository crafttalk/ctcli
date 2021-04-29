package cmd

import (
	"ctcli/domain"
	"ctcli/domain/ctcliDir"
	"github.com/spf13/cobra"
	"log"
	"path/filepath"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "upgrade a release",
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
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
