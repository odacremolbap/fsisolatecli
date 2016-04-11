package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/odacremolbap/fsisolate"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

var debug bool
var root string
var delayAfterExited int64

func init() {

	flag.BoolVarP(&debug, "debug", "d", false, "enable debug messages")
	flag.StringVarP(&root, "root", "r", "", "directory to place the new root")
	flag.Int64VarP(&delayAfterExited, "postdelay", "t", 0, "ammount of time in seconds before exiting")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: isocli [OPTIONS] IMAGE COMMAND [command args ...]\n")
		fmt.Fprintf(os.Stderr, "\nA naive chroot wrapper\n")
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
		os.Exit(2)
	}

	flag.Parse()

	if debug {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {

	// check existence for IMAGE nad COMMAND args
	if len(flag.Args()) < 2 {
		flag.Usage()
		os.Exit(2)
	}
	var err error
	image := flag.Arg(0)
	command := flag.Arg(1)
	args := flag.Args()[2:]

	// if no root is informed, use a generated temp dir
	// TODO if image parameter is directory, this tempdir is not going to be used. Check that
	if root == "" {
		root, err = ioutil.TempDir("", "isocli")
		if err != nil {
			log.Fatal(err)
		}
		log.Debugf("Root will be located at '%s'", root)
	}

	// prepare process
	chrootProc, err := fsisolate.PrepareChrootedProcess(image, root)
	if err != nil {
		log.Fatal(err)
	}

	// launch commands and signals processor
	go inputProc(chrootProc)

	// execute process
	if err = chrootProc.SandboxExec(command, args...); err != nil {
		log.Fatal(err)
	}

	// wait for the process to finish
	if err = chrootProc.Wait(); err != nil {
		log.Fatal(err)
	}

	// delay after exited
	if delayAfterExited > 0 {
		time.Sleep(time.Duration(delayAfterExited) * time.Second)
	}

}

func printMetaInfo(message string, args ...interface{}) {
	message = "[ISOCLI]" + message + "\n"
	fmt.Printf(message, args...)
}
