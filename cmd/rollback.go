package cmd

import (
	"ctcli/domain"
	"ctcli/domain/ctcliDir"
	"ctcli/util"
	"path/filepath"

	"github.com/spf13/cobra"
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback <path-to-backup>",
	Short: "rollback a release",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		rootFlag := cmd.Flag("root")
		rootDir, err := filepath.Abs(rootFlag.Value.String())
		if err != nil {
			return err
		}
		if err := ctcliDir.OkIfIsARootDir(rootDir); err != nil {
			return err
		}
		backupPath, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}
		fn := util.MirrorStdoutToFile(ctcliDir.GetCtcliLogFilePath(rootDir))
		defer fn()

		if err := domain.Rollback(rootDir, backupPath); err != nil {
			cmd.PrintErr(err)
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rollbackCmd)
}
