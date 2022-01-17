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
	Run: func(cmd *cobra.Command, args []string) {
		rootFlag := cmd.Flag("root")
		rootDir, err := filepath.Abs(rootFlag.Value.String())
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		fn := util.MirrorStdoutToFile(ctcliDir.GetCtcliLogFilePath(rootDir))
		defer fn()
		ctcliDir.Init(rootDir)
		cmd.Printf("%s", color.HiGreenString("OK\n"))
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
