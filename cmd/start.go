package cmd

import (
	"ctcli/domain/ctcliDir"
	"ctcli/domain/lifetime"
	"ctcli/util"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [app]",
	Short: "start a service",
	Run: func(cmd *cobra.Command, args []string) {
		rootFlag := cmd.Flag("root")
		disableFlag := cmd.Flag("disable")
		
		isDisableWriteLogsString := disableFlag.Value.String()
		isDisableWriteLogs, err := strconv.ParseBool(isDisableWriteLogsString)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		rootDir, err := filepath.Abs(rootFlag.Value.String())
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		if err := ctcliDir.OkIfIsARootDir(rootDir); err != nil {
			cmd.PrintErr(err)
			return
		}

		if isDisableWriteLogs == false {
			fn := util.MirrorStdoutToFile(ctcliDir.GetCtcliLogFilePath(rootDir))
			defer fn()
		}

		if err := lifetime.StartApps(rootDir, args, isDisableWriteLogs); err != nil {
			cmd.PrintErr(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().BoolP("disable", "d", false, "Disable write to stdout-stderr.log")
}
