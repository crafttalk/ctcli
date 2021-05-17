package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const VERSION int = 3
const COMMIT string = "%%commit_hash%%"

var versionCmd = &cobra.Command{
	Use: "version",
	Short: "print the version number of CraftTalk CLI",
	Run: func(cmd *cobra.Command, args []string) {
		blue := color.New(color.FgBlue)
		green := color.New(color.FgGreen)
		yellow := color.New(color.FgYellow)

		blue.Printf("CraftTalk Command Line Tool ")
		green.Printf("v%d", VERSION)
		blue.Printf(" -- ")
		yellow.Println(COMMIT)
	},
}

func init()  {
	rootCmd.AddCommand(versionCmd)
}
