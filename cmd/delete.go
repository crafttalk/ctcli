package cmd

import (
	"ctcli/domain/ctcliDir"
	"ctcli/util"
	"github.com/spf13/cobra"
	"path/filepath"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a release",
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
		currentReleaseDir := ctcliDir.GetCurrentReleaseDir(rootDir)
		if err := util.RemoveContentOfFolder(currentReleaseDir); err != nil {
			cmd.PrintErr(err)
			return
		}
		cmd.Println("current release was deleted")
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
