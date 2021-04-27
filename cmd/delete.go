package cmd

import "github.com/spf13/cobra"

var deleteCmd = &cobra.Command{
	Use: "delete",
	Short: "Delete a release",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init()  {
	rootCmd.AddCommand(deleteCmd)
}