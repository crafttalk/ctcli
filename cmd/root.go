package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ctcli",
	Short: "CraftTalk CLI helps with managing CraftTalk installations",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	rootCmd.PersistentFlags().String("root", workDir, "root of the installation")
}

func Execute() {
	//rootCmd.SetArgs([]string{ "--root", "/home/lkmfwe/ctcli", "install", "/home/lkmfwe/Programming/FSharp/opbot/packaging/package/crafttalk-opbot-release-2021-03-22-6-commit.tar.gz" })
	//rootCmd.SetArgs([]string{ "--root", "/home/lkmfwe/ctcli", "backup" })

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
