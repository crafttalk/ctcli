package cmd

import (
	"ctcli/domain"
	"ctcli/domain/ctcliDir"
	"ctcli/util"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "make a backup of current release",
	RunE: func(cmd *cobra.Command, args []string) error {
		rootFlag := cmd.Flag("root")
		rootDir, err := filepath.Abs(rootFlag.Value.String())
		if err != nil {
			return err
		}

		backupDataFlag := cmd.Flag("ignore-data")
		backupData := false
		if backupDataFlag.Value.String() == "false" {
			backupData = true
		}

		fn := util.MirrorStdoutToFile(ctcliDir.GetCtcliLogFilePath(rootDir))
		defer fn()
		if err := domain.MakeABackup(rootDir, backupData); err != nil {
			return err
		}
		cmd.Printf("%s", color.HiGreenString("backup was made\n"))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
	backupCmd.Flags().Bool("ignore-data", false, "Do not include data folder into backup")
}
