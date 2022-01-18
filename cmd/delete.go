package cmd

import (
	"ctcli/domain/ctcliDir"
	"ctcli/util"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a release",
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
		currentReleaseDir := ctcliDir.GetCurrentReleaseDir(rootDir)
		if err := util.RemoveContentOfFolder(currentReleaseDir); err != nil {
			return err
		}
		cmd.Printf(color.HiGreenString("Current release was deleted\n"))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
