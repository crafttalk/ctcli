package cmd

import (
	"ctcli/domain/ctcliDir"
	"ctcli/domain/lifetime"
	"ctcli/util"
	"path/filepath"

	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec <app> <command>",
	Short: "exec inside a service filesystem",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		rootFlag := cmd.Flag("root")
		rootDir, err := filepath.Abs(rootFlag.Value.String())
		if err != nil {
			return err
		}
		if err := ctcliDir.OkIfIsARootDir(rootDir); err != nil {
			return err
		}
		fn := util.MirrorStdoutToFile(ctcliDir.GetCtcliLogFilePath(rootDir))
		defer fn()

		if err := lifetime.Exec(rootDir, args); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
