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
var beforeDelay int64
var afterDelay int64

func init() {

	flag.BoolVarP(&debug, "debug", "d", false, "enable debug messages")
	flag.StringVarP(&root, "root", "r", "", "directory to place the new root")
	flag.Int64VarP(&beforeDelay, "beforedelay", "b", 2, "ammount of time in seconds before executing process")
	flag.Int64VarP(&afterDelay, "afterdelay", "a", 0, "ammount of time in seconds after finishing process")
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

	// if no root is informed and image is not a directory
	// use a generated temp dir
	r, err := os.Stat(image)
	if root == "" && (err != nil || !r.IsDir()) {
		root, err = ioutil.TempDir("", "isocli")
		if err != nil {
			log.Fatal(err)
		}
		log.Debugf("Root will be located at '%s'", root)
	}

	// prepare process
	chrootProc, err := fsisolate.Prepare(image, root)
	if err != nil {
		log.Fatal(err)
	}
	printMetaInfo("When running the process enter:")
	printMetaInfo("\t's + enter' to get process state")
	printMetaInfo("\t'i + enter' to send SIGINT")
	printMetaInfo("\t'h + enter' to send SIGHUP")
	printMetaInfo("\t'k + enter' to send SIGKILL")
	printMetaInfo("\t'u + enter' to send SIGUSR1")

	// launch commands and signals processor
	go inputProc(chrootProc)

	// delay before execution
	if beforeDelay > 0 {
		time.Sleep(time.Duration(beforeDelay) * time.Second)
	}

	// execute process
	if err = chrootProc.Exec(command, args...); err != nil {
		log.Fatal(err)
	}

	// wait for the process to finish
	if err = chrootProc.Wait(); err != nil {
		// this error is related to the chrooted process
		log.Error(err)
	}

	// delay after exited
	if afterDelay > 0 {
		time.Sleep(time.Duration(afterDelay) * time.Second)
	}

}

func printMetaInfo(message string, args ...interface{}) {
	message = "[ISOCLI]" + message + "\n"
	fmt.Printf(message, args...)
}
