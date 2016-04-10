package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/odacremolbap/fsisolate"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

var debug bool
var root string

func init() {

	flag.BoolVarP(&debug, "debug", "d", false, "enable debug messages")
	flag.StringVarP(&root, "root", "r", "", "directory to place the new root")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] IMAGE COMMAND [command args ...]\n", os.Args[0])
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
	exec := flag.Arg(1)
	args := flag.Args()[2:]

	// if no root is informed, use a generated temp dir
	if root == "" {
		root, err = ioutil.TempDir("", "isocli")
		if err != nil {
			log.Fatal(err)
		}
		log.Debugf("Root will be located at '%s'", root)
	}

	// prepare process
	chrootProc, err := fsisolate.PrepareChrootedProcess(image, root, exec, args)
	if err != nil {
		log.Fatal(err)
	}

	if err = chrootProc.SandboxExec(exec, args...); err != nil {
		log.Fatal(err)
	}

}
