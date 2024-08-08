package cmd

import (
	"ctcli/domain/ctcliDir"
	"ctcli/domain/lifetime"
	"ctcli/util"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:   "restart [app]",
	Short: "restart a service",
	Run: func(cmd *cobra.Command, args []string) {
		rootFlag := cmd.Flag("root")
		rootDir, err := filepath.Abs(rootFlag.Value.String())
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		
		disableFlag := cmd.Flag("disable")
		
		isDisableWriteLogsString := disableFlag.Value.String()
		isDisableWriteLogs, err := strconv.ParseBool(isDisableWriteLogsString)
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

		if err := lifetime.StopApps(rootDir, args); err != nil {
			cmd.PrintErr(err)
			return
		}
		if err := lifetime.StartApps(rootDir, args, isDisableWriteLogs); err != nil {
			cmd.PrintErr(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)
	restartCmd.Flags().BoolP("disable", "d", false, "Disable write to stdout-stderr.log")
}
