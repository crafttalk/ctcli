package cmd

import (
	"ctcli/domain/ctcliDir"
	"ctcli/util"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "inits directory as a ctcli work directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		rootFlag := cmd.Flag("root")
		rootDir, err := filepath.Abs(rootFlag.Value.String())
		if err != nil {
			return err
		}
		fn := util.MirrorStdoutToFile(ctcliDir.GetCtcliLogFilePath(rootDir))
		defer fn()
		ctcliDir.Init(rootDir)
		cmd.Printf("%s", color.HiGreenString("OK\n"))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
