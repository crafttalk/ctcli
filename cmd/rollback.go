package cmd

import (
	"ctcli/domain/ctcliDir"
	"github.com/spf13/cobra"
	"path/filepath"
)

var rollbackCmd = &cobra.Command{
	Use: "rollback",
	Short: "rollback a release",
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
	},
}

func init()  {
	//rootCmd.AddCommand(rollbackCmd)
}