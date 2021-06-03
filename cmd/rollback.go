package cmd

import (
	"ctcli/domain"
	"ctcli/domain/ctcliDir"
	"ctcli/util"
	"github.com/spf13/cobra"
	"path/filepath"
)

var rollbackCmd = &cobra.Command{
	Use: "rollback <path-to-backup>",
	Short: "rollback a release",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rootFlag := cmd.Flag("root")
		rootDir, err := filepath.Abs(rootFlag.Value.String())
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		if err := ctcliDir.OkIfIsARootDir(rootDir); err != nil {
			cmd.PrintErr(err)
			return
		}
		backupPath, err := filepath.Abs(args[0])
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		fn := util.MirrorStdoutToFile(ctcliDir.GetCtcliLogFilePath(rootDir))
		defer fn()

		if err := domain.Rollback(rootDir, backupPath); err != nil {
			cmd.PrintErr(err)
			return
		}
	},
}

func init()  {
	rootCmd.AddCommand(rollbackCmd)
}