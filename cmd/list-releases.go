package cmd

import (
	"ctcli/domain/ctcliDir"
	"ctcli/domain/release"
	"github.com/kataras/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"sort"
)

func getAppsInRelease(appVersions []release.AppVersion)  string {
	result := ""
	for i, appVersion := range appVersions {
		if i > 0 {
			result += ", "
		}
		result += appVersion.Image
	}
	return result
}

var listReleasesCmd = &cobra.Command{
	Use: "list-releases",
	Short: "list all installed releases",
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

		releases, err := release.GetReleases(rootDir)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		sort.Slice(releases, func(i, j int) bool {
			return releases[i].CreatedAt.Unix() > releases[j].CreatedAt.Unix()
		})

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"current", "branch", "commit", "build-date", "apps"})
		tableContent := [][]string{}
		for _, release := range releases {
			row := []string{
				" ",
				release.Id,
				release.PreviousRelease,
				release.CreatedAt.String(),
				getAppsInRelease(release.AppVersions),
			}
			tableContent = append(tableContent, row)
		}

		if releaseInfo, err := release.GetCurrentReleaseInfo(rootDir); err == nil {
			currentReleaseRow := []string{
				"*",
				releaseInfo.Id,
				releaseInfo.PreviousRelease,
				releaseInfo.CreatedAt.String(),
			}
			table.Append(currentReleaseRow)
		}

		table.AppendBulk(tableContent)
		table.Render()
	},
}

func init()  {
	rootCmd.AddCommand(listReleasesCmd)
}