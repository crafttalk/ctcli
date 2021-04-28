package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var rootCmd = &cobra.Command{
	Use: "ctcli",
	Short: "CraftTalk CLI helps with managing CraftTalk installations",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute()  {
	//rootCmd.SetArgs([]string{ "--root", "/home/lkmfwe/ctcli", "install", "/home/lkmfwe/Programming/FSharp/opbot/packaging/package/crafttalk-opbot-release-2021-03-22-6-commit.tar.gz" })
	//rootCmd.SetArgs([]string{ "--root", "/home/lkmfwe/ctcli", "release-info" })

	workDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	rootCmd.PersistentFlags().String("root", workDir, "root of the installation")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}