package cmd

import (
	"ctcli/domain"
	"ctcli/domain/ctcliDir"
	"ctcli/util"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install <path to package>",
	Short: "install a release",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		packagePath, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}

		rootFlag := cmd.Flag("root")
		rootDir, err := filepath.Abs(rootFlag.Value.String())
		if err != nil {
			return err
		}
		fn := util.MirrorStdoutToFile(ctcliDir.GetCtcliLogFilePath(rootDir))
		defer fn()

		err = domain.Install(rootDir, packagePath)
		if err != nil {
			return err
		}
		cmd.Printf("%s", color.HiGreenString("OK\n"))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
