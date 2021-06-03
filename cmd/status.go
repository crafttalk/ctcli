package cmd

import (
	"ctcli/domain/ctcliDir"
	"ctcli/domain/release"
	"ctcli/domain/runc"
	"ctcli/util"
	"fmt"
	"github.com/kataras/tablewriter"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "shows current status",
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
		fn := util.MirrorStdoutToFile(ctcliDir.GetCtcliLogFilePath(rootDir))
		defer fn()

		runcPath := release.GetCurrentReleaseRuncPath(rootDir)
		if !util.PathExists(runcPath) {
			cmd.PrintErr(fmt.Errorf("there is no runc in current-relase folder"))
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
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
