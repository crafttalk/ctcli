package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use: "start [app]",
	Short: "start a service",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("not implemented")
	},
}

func init()  {
	rootCmd.AddCommand(startCmd)
}