package cmd

import (
	"ctcli/domain/ctcliDir"
	"ctcli/domain/release"
	"fmt"
	"github.com/fatih/color"
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
		if err := ctcliDir.OkIfIsARootDir(rootDir); err != nil {
			cmd.PrintErr(err)
			return
		}

		releaseInfo, err := release.GetCurrentReleaseInfo(rootDir)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		nameColor := color.New(color.FgBlue)
		valueColor := color.New(color.FgGreen)

		subNameColor := color.New(color.FgYellow)
		subValueColor := color.New(color.FgCyan)

		appVersionsString := ""
		for _, appVersion := range releaseInfo.AppVersions {
			appVersionsString += subNameColor.Sprintf("  app: ")
			appVersionsString += color.New(color.FgRed).Sprintln(appVersion.AppName)
			appVersionsString += subNameColor.Sprintf("    image: ")
			appVersionsString += subValueColor.Sprintln(appVersion.Image)
			appVersionsString += subNameColor.Sprintf("    built at: ")
			appVersionsString += subValueColor.Sprintln(appVersion.BuiltAt)
			appVersionsString += subNameColor.Sprintf("    commit: ")
			appVersionsString += subValueColor.Sprintln(appVersion.CommitSha)
			appVersionsString += subNameColor.Sprintf("    imageSha: ")
			appVersionsString += subValueColor.Sprintln(appVersion.ImageSha)
		}

		nameColor.Printf("id: ")
		valueColor.Println(releaseInfo.Id)

		nameColor.Printf("prev release: ")
		valueColor.Println(releaseInfo.PreviousRelease)

		nameColor.Printf("created at: ")
		valueColor.Println(releaseInfo.CreatedAt)

		color.Blue("app versions:")
		fmt.Print(appVersionsString)
	},
}

func init()  {
	rootCmd.AddCommand(releaseInfoCmd)
}