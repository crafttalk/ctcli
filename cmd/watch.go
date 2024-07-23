package cmd

import (
	"ctcli/domain/ctcliDir"
	"ctcli/domain/release"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"
)

var PIDFile = "/tmp/daemonize.pid"

func savePID(pid int) {
	file, err := os.Create(PIDFile)
	if err != nil {
			log.Printf("Unable to create pid file : %v\n", err)
			os.Exit(1)
	}

	defer file.Close()

	_, err = file.WriteString(strconv.Itoa(pid))

	if err != nil {
			log.Printf("Unable to create pid file : %v\n", err)
			os.Exit(1)
	}

	file.Sync()
}

var watchCmd = &cobra.Command{
	Use:   "watch [start|stop]",
	Short: "watch a log files",
	Run: func(cmd *cobra.Command, args []string) {
		rootFlag := cmd.Flag("root")
		maxFileSizeFlag := cmd.Flag("size")
		daysFlag := cmd.Flag("days")
		
		switch args[0] {
			case "main":
				days, err := strconv.ParseInt(daysFlag.Value.String(), 8, 64)
				if err != nil {
					cmd.PrintErr(err)
					return
				}

				maxFileSize, err := strconv.ParseInt(maxFileSizeFlag.Value.String(), 8, 64)
				if err != nil {
					cmd.PrintErr(err)
					return
				}

				rootDir, err := filepath.Abs(rootFlag.Value.String())
				if err != nil {
					cmd.PrintErr(err)
					return
				}
				if err := ctcliDir.OkIfIsARootDir(rootDir); err != nil {
					cmd.PrintErr(err)
					return
				}
				ch := make(chan os.Signal, 1)
				signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)

				go func() {
					_ = <- ch
					signal.Stop(ch)

					os.Remove(PIDFile)

					os.Exit(0)
				}()

				appsPath := release.GetCurrentReleaseAppsFolder(rootDir)
				appFolders, _ := ioutil.ReadDir(appsPath)
				appNames := release.GetAppNamesFromFolders(appFolders)
				
				if maxFileSize > 0 {
					go ctcliDir.CheckFilesSize(rootDir, appNames, maxFileSize)
				}
		
				if days > 0 {
					go ctcliDir.ArchiveLifeTime(rootDir, appNames, days)
				}
		
				select {}
			case "start":
				rootDir, err := filepath.Abs(rootFlag.Value.String())
				if err != nil {
					cmd.PrintErr(err)
					return
				}
				if err := ctcliDir.OkIfIsARootDir(rootDir); err != nil {
					cmd.PrintErr(err)
					return
				}
				if _, err := os.Stat(PIDFile); err == nil {
					fmt.Println("Already running or /tmp/daemonize.pid file exist.")
					os.Exit(1)
				}
				
				cmdTest := exec.Command(os.Args[0], "watch", "main", "--size", maxFileSizeFlag.Value.String(), "--days", daysFlag.Value.String())
				cmdTest.Start()
				fmt.Println("Daemon process ID is : ", cmdTest.Process.Pid)
				savePID(cmdTest.Process.Pid)
				os.Exit(0)
			case "stop":
				if _, err := os.Stat(PIDFile); err == nil {
					data, err := ioutil.ReadFile(PIDFile)
					if err != nil {
						fmt.Println("Not running")
						os.Exit(1)
					}
					ProcessID, err := strconv.Atoi(string(data))

					if err != nil {
						fmt.Println("Unable to read and parse process id found in ", PIDFile)
						os.Exit(1)
					}

					process, err := os.FindProcess(ProcessID)

					if err != nil {
						fmt.Printf("Unable to find process ID [%v] with error %v \n", ProcessID, err)
						os.Exit(1)
					}
					os.Remove(PIDFile)

					fmt.Printf("Killing process ID [%v] now.\n", ProcessID)
					err = process.Kill()

					if err != nil {
						fmt.Printf("Unable to kill process ID [%v] with error %v \n", ProcessID, err)
						os.Exit(1)
					} else {
						fmt.Printf("Killed process ID [%v]\n", ProcessID)
						os.Exit(0)
					}

				} else {
					fmt.Println("Not running.")
					os.Exit(1)
				}
			}
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
	watchCmd.Flags().Int("size", 0, "Set max .log file size after which it will be archived")
	watchCmd.Flags().Int("days", 0, "Set the number of days in rotation for archived logs")
}