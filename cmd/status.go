package cmd

import (
	"ctcli/domain/ctcliDir"
	"ctcli/domain/release"
	"ctcli/domain/runc"
	"ctcli/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/kataras/tablewriter"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "shows current status",
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

		runcPath := release.GetCurrentReleaseRuncPath(rootDir)
		if !util.PathExists(runcPath) {
			cmd.PrintErr("there is no runc in current-relase folder")
		}
		appsPath := release.GetCurrentReleaseAppsFolder(rootDir)
		appFolders, err := ioutil.ReadDir(appsPath)
		if err != nil {
			cmd.PrintErr(err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"app-name", "status", "pid"})
		tableContent := [][]string{}

		for _, appFolder := range appFolders {
			appName := appFolder.Name()
			status := runc.GetStatus(rootDir, appName)
			if status == "unknown" {
				status = "-"
			}
			pid := runc.GetPid(rootDir, appName)
			row := []string{
				appName,
				status,
				strconv.Itoa(pid),
			}
			tableContent = append(tableContent, row)
		}

		table.AppendBulk(tableContent)
		table.SetOutput(cmd.OutOrStdout())
		table.Render()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
