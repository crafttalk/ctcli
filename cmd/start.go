package cmd

import (
	"ctcli/domain/ctcliDir"
	"ctcli/domain/lifetime"
	"ctcli/util"
	"github.com/spf13/cobra"
	"path/filepath"
)

var startCmd = &cobra.Command{
	Use:   "start [app]",
	Short: "start a service",
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
		fn := util.MirrorStdoutToFile(ctcliDir.GetCtcliLogFilePath(rootDir))
		defer fn()
		if err := lifetime.StartApps(rootDir, args); err != nil {
			cmd.PrintErr(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
