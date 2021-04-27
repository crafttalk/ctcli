package cmd

import (
	"ctcli/domain/release"
	"fmt"
	"github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"
	"path/filepath"
)

var releaseInfoCmd = &cobra.Command{
	Use: "release-info",
	Short: "get current release info",
	Run: func(cmd *cobra.Command, args []string) {
		rootFlag := cmd.Flag("root")
		rootDir, err := filepath.Abs(rootFlag.Value.String())
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		releaseInfo, err := release.GetCurrentReleaseInfo(rootDir)
		if err != nil {
			cmd.PrintErr(err)
			return
		}



		fmt.Printf(
			"%s%s\n%s%s\n%s%s\n%s%s\n",
			runewidth.FillRight("branch:", 16),
			releaseInfo.Branch,
			runewidth.FillRight("commit:", 16),
			releaseInfo.CommitSha,
			runewidth.FillRight("built at:", 16),
			releaseInfo.BuildDate.String(),
			runewidth.FillRight("installed at:", 16),
			releaseInfo.InstalledDate.String(),
		)
	},
}

func init()  {
	rootCmd.AddCommand(releaseInfoCmd)
}