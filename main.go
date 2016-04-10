package main

import (
	"bufio"
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

func init() {

	flag.BoolVarP(&debug, "debug", "d", false, "enable debug messages")
	flag.StringVarP(&root, "root", "r", "", "directory to place the new root")
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

	go func() {
		for {

			reader := bufio.NewReader(os.Stdin)
			text, err := reader.ReadByte()

			if err != nil {
				continue
			}

			switch text {
			case 's', 'S':
				// get status
				pid, err := chrootProc.GetPID()
				if err != nil {
					fmt.Println(err)
					continue
				}

				exited, err := chrootProc.GetExited()
				if err != nil {
					fmt.Println(err)
					continue
				}

				if exited {
					exitStatus, err := chrootProc.GetExitStatus()
					if err != nil {
						fmt.Println(err)
						continue
					}
					printMetaInfo(" PID: %d EXITED: %t EXIT-STATUS: %d", pid, exited, exitStatus)

				} else {
					printMetaInfo(" PID: %d EXITED: %t", pid, exited)
				}
			case 'q', 'Q':
				err = chrootProc.SendSignal(os.Interrupt)
				if err != nil {
					fmt.Println(err)
					continue
				}
			}

		}
	}()

	if err = chrootProc.SandboxExec(command, args...); err != nil {
		log.Fatal(err)
	}

	if err = chrootProc.Wait(); err != nil {
		log.Fatal(err)
	}

	time.Sleep(5 * time.Second)

}

func printMetaInfo(message string, args ...interface{}) {
	message = "[ISOCLI]" + message + "\n"
	fmt.Printf(message, args...)
}
