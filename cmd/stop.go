package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use: "stop [app]",
	Short: "start a service",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("not implemented")
	},
}

func init()  {
	rootCmd.AddCommand(stopCmd)
}