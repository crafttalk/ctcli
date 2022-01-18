package cmd

import (
	"ctcli/domain/ctcliDir"
	"ctcli/util"
	"os"
	"path/filepath"
	"strconv"

	"github.com/fatih/color"
	"github.com/hpcloud/tail"
	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs <app>",
	Short: "show stdout and stderr of an app",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		rootFlag := cmd.Flag("root")
		rootDir, err := filepath.Abs(rootFlag.Value.String())
		if err != nil {
			return err
		}
		if err := ctcliDir.OkIfIsARootDir(rootDir); err != nil {
			return err
		}
		fn := util.MirrorStdoutToFile(ctcliDir.GetCtcliLogFilePath(rootDir))
		defer fn()

		appName := args[0]
		logFilePath := ctcliDir.GetAppStdoutLogFilePath(rootDir, appName)

		followFlag := cmd.Flag("follow")
		follow := true
		if followFlag.Value.String() == "false" {
			follow = false
		}

		tailFlag := cmd.Flag("tail")
		tailAmount, err := strconv.ParseInt(tailFlag.Value.String(), 10, 32)
		if err != nil {
			cmd.PrintErr(err)
		}

		config := tail.Config{Follow: follow}
		if tailAmount > 0 {
			config.Location = &tail.SeekInfo{-tailAmount, os.SEEK_END}
		}
		t, err := tail.TailFile(logFilePath, config)
		if err != nil {
			cmd.PrintErr(err)
		}

		for line := range t.Lines {
			if follow {
				cmd.Printf("[%s] %s\n", color.HiBlueString("%s", line.Time), line.Text)
			} else {
				cmd.Println(line.Text)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
	logsCmd.Flags().BoolP("follow", "f", false, "Follow logs")
	logsCmd.Flags().Int32("tail", 0, "Tail logs")
}
