package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var rollbackCmd = &cobra.Command{
	Use: "rollback",
	Short: "rollback a release",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("not implemented")
	},
}

func init()  {
	rootCmd.AddCommand(rollbackCmd)
}