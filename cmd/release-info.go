package cmd

import (
	"ctcli/domain/ctcliDir"
	"ctcli/domain/release"
	"ctcli/util"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var releaseInfoCmd = &cobra.Command{
	Use:   "release-info",
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
		fn := util.MirrorStdoutToFile(ctcliDir.GetCtcliLogFilePath(rootDir))
		defer fn()

		releaseInfo, err := release.GetCurrentReleaseInfo(rootDir)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		nameColor := color.New(color.FgHiBlue)
		valueColor := color.New(color.FgHiGreen)

		subNameColor := color.New(color.FgHiYellow)
		subValueColor := color.New(color.FgHiCyan)

		appVersionsString := ""
		for _, appVersion := range releaseInfo.AppVersions {
			appVersionsString += subNameColor.Sprintf("  app: ")
			appVersionsString += color.New(color.FgHiRed).Sprintln(appVersion.AppName)
			appVersionsString += subNameColor.Sprintf("    image: ")
			appVersionsString += subValueColor.Sprintln(appVersion.Image)
			appVersionsString += subNameColor.Sprintf("    built at: ")
			appVersionsString += subValueColor.Sprintln(appVersion.BuiltAt)
			appVersionsString += subNameColor.Sprintf("    commit: ")
			appVersionsString += subValueColor.Sprintln(appVersion.CommitSha)
			appVersionsString += subNameColor.Sprintf("    imageSha: ")
			appVersionsString += subValueColor.Sprintln(appVersion.ImageSha)
		}

		cmd.Printf("%s", nameColor.Sprintf("id: "))
		cmd.Printf("%s\n", valueColor.Sprintf(releaseInfo.Id))

		cmd.Printf("%s", nameColor.Sprintf("prev release: "))
		cmd.Printf("%s\n", valueColor.Sprintf(releaseInfo.PreviousRelease))

		cmd.Printf("%s", nameColor.Sprintf("created at: "))
		cmd.Printf("%s\n", valueColor.Sprint(releaseInfo.CreatedAt))

		cmd.Print(color.HiBlueString("app versions:\n"))
		cmd.Print(appVersionsString)
	},
}

func init() {
	rootCmd.AddCommand(releaseInfoCmd)
}
