package cmd

import (
	"ctcli/domain"
	"ctcli/domain/ctcliDir"
	"ctcli/util"
	"log"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade <path to package>",
	Short: "upgrade a release",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		rootFlag := cmd.Flag("root")
		rootDir, err := filepath.Abs(rootFlag.Value.String())
		if err != nil {
			return err
		}
		packagePath, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}
		if err := ctcliDir.OkIfIsARootDir(rootDir); err != nil {
			return err
		}
		fn := util.MirrorStdoutToFile(ctcliDir.GetCtcliLogFilePath(rootDir))
		defer fn()
		err = domain.Upgrade(rootDir, packagePath)
		if err != nil {
			log.Fatal(err)
		}
		color.Green("OK\n")

		return err
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
