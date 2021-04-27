package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use: "upgrade",
	Short: "upgrade a release",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("not implemented")
	},
}

func init()  {
	rootCmd.AddCommand(upgradeCmd)
}