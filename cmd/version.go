package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

const VERSION int = 1
const COMMIT string = "%%commit_hash%%"

var versionCmd = &cobra.Command{
	Use: "version",
	Short: "Print the version number of CraftTalk CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("CraftTalk Command Line Tool v%d -- %s", VERSION, COMMIT)
	},
}

func init()  {
	rootCmd.AddCommand(versionCmd)
}
