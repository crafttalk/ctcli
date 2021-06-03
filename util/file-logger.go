package util

import (
	"fmt"
	"github.com/fatih/color"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func MirrorStdoutToFile(logFile string) func() {
	// open file read/write | create if not exist | clear file at open if exists
	f, _ := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(f, "\n=== %s ===\n", time.Now())
	fmt.Fprintf(f, "Working directory: %s\n", wd)
	fmt.Fprintf(f, "> %s\n", strings.Join(os.Args, " "))

	// save existing stdout | MultiWriter writes to saved stdout and file
	out := os.Stdout
	mw := io.MultiWriter(out, f)

	// get pipe reader and writer | writes to pipe writer come out pipe reader
	r, w, _ := os.Pipe()

	// replace stdout,stderr with pipe writer | all writes to stdout, stderr will go through pipe instead (fmt.print, log)
	os.Stdout = w
	os.Stderr = w
	color.Output = mw
	color.Error = mw

	// writes with log.Print should also write to mw
	log.SetOutput(mw)

	//create channel to control exit | will block until all copies are finished
	exit := make(chan bool)


	go func() {
		// copy all reads from pipe to multiwriter, which writes to stdout and file
		_,_ = io.Copy(mw, r)
		// when r or w is closed copy will finish and true will be sent to channel
		exit <- true
	}()

	// function to be deferred in main until program exits
	return func() {
		// close writer then block on exit channel | this will let mw finish writing before the program exits
		_ = w.Close()
		<-exit
		// close file after all writes have finished
		_ = f.Close()
	}
}