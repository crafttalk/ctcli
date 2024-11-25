package cmd

import (
	"strconv"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const VERSION float64 = 7.4
const COMMIT string = "%%commit_hash%%"

var versionCmd = &cobra.Command{
	Use: "version",
	Short: "print the version number of CraftTalk CLI",
	Run: func(cmd *cobra.Command, args []string) {
		blue := color.New(color.FgHiBlue)
		green := color.New(color.FgHiGreen)
		yellow := color.New(color.FgHiYellow)

		blue.Printf("CraftTalk Command Line Tool ")
		green.Printf("v" + strconv.FormatFloat(VERSION, 'f', -1, 64))
		blue.Printf(" -- ")
		yellow.Println(COMMIT)
	},
}

func init()  {
	rootCmd.AddCommand(versionCmd)
}
