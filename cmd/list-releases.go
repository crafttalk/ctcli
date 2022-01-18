package cmd

import (
	"ctcli/domain/ctcliDir"
	"ctcli/domain/release"
	"ctcli/util"
	"path/filepath"
	"sort"

	"github.com/kataras/tablewriter"
	"github.com/spf13/cobra"
)

func getAppsInRelease(appVersions []release.AppVersion) string {
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
	Use:   "list-releases",
	Short: "list all installed releases",
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

		releases, err := release.GetReleases(rootDir)
		if err != nil {
			return err
		}

		sort.Slice(releases, func(i, j int) bool {
			return releases[i].CreatedAt.Unix() > releases[j].CreatedAt.Unix()
		})

		table := tablewriter.NewWriter(cmd.OutOrStdout())
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
		table.SetOutput(cmd.OutOrStdout())
		table.Render()

		return nil
	},
}

func init() {
	// Работает криво, поэтому закоменчено
	// rootCmd.AddCommand(listReleasesCmd)
}
