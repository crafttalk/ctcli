package cmd

import (
	"ctcli/domain/ctcliDir"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"path/filepath"
)

var logsCmd = &cobra.Command{
	Use: "logs <app>",
	Short: "show stdout and stderr of an app",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rootFlag := cmd.Flag("root")
		rootDir, err := filepath.Abs(rootFlag.Value.String())
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		if err := ctcliDir.OkIfIsARootDir(rootDir); err != nil {
			cmd.PrintErr(err)
			return
		}

		appName := args[0]
		logFilePath := ctcliDir.GetAppStdoutLogFilePath(rootDir, appName)
		b, err := ioutil.ReadFile(logFilePath)
		if err != nil {
			cmd.PrintErr(err)
		}

		fmt.Print(string(b))
	},
}

func init()  {
	rootCmd.AddCommand(logsCmd)
}