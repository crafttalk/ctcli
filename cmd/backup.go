package cmd

import (
	"ctcli/domain"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"log"
	"path/filepath"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "making backup of current release",
	Run: func(cmd *cobra.Command, args []string) {
		rootFlag := cmd.Flag("root")
		rootDir, err := filepath.Abs(rootFlag.Value.String())
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		if err := domain.MakeABackup(rootDir); err != nil {
			log.Fatal(err)
		}
		color.Green("backup was made\n")
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
