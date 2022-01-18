package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const VERSION int = 7
const COMMIT string = "%%commit_hash%%"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the version number of CraftTalk CLI",
	RunE: func(cmd *cobra.Command, args []string) error {
		blue := color.New(color.FgHiBlue)
		green := color.New(color.FgHiGreen)
		yellow := color.New(color.FgHiYellow)

		cmd.Printf(
			"%s %s %s %s\n",
			blue.Sprintf("CraftTalk Command Line Tool"),
			green.Sprintf("v%d", VERSION),
			blue.Sprintf("--"),
			yellow.Sprintf(COMMIT),
		)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
