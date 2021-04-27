package cmd

import (
	"ctcli/domain/release"
	"github.com/kataras/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"sort"
)

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

		releases, err := release.GetReleases(rootDir)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		sort.Slice(releases, func(i, j int) bool {
			return releases[i].BuildDate.Unix() > releases[j].BuildDate.Unix()
		})

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"current", "branch", "commit", "build-date"})
		tableContent := [][]string{}
		for _, release := range releases {
			row := []string{ " ", release.Branch, release.CommitSha, release.BuildDate.String() }
			tableContent = append(tableContent, row)
		}

		if releaseInfo, err := release.GetCurrentReleaseInfo(rootDir); err == nil {
			currentReleaseRow := []string{ "*", releaseInfo.Branch, releaseInfo.CommitSha, releaseInfo.BuildDate.String() }
			table.Append(currentReleaseRow)
		}

		table.AppendBulk(tableContent)
		table.Render()
	},
}

func init()  {
	rootCmd.AddCommand(listReleasesCmd)
}