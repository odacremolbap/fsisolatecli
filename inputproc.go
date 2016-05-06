package main

import (
	"bufio"
	"fmt"
	"os"
	"syscall"

	"github.com/odacremolbap/fsisolate"
)

func inputProc(chrootProc *fsisolate.ChrootedProcess) {
	for {

		// TODO capture like getch
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadByte()

		if err != nil {
			continue
		}

		switch text {
		case 's', 'S':
			// get status
			pid, e := chrootProc.GetPID()
			if err != nil {
				fmt.Println(e)
				continue
			}

			state := chrootProc.GetState()

			if state == fsisolate.Finished {
				exitStatus, e := chrootProc.GetExitStatus()
				if err != nil {
					fmt.Println(e)
					continue
				}
				printMetaInfo(" PID: %d STATE: %s EXIT-STATUS: %d", pid, state, exitStatus)

			} else {
				printMetaInfo(" PID: %d STATE: %s", pid, state)
			}
		case 'h', 'H':
			err = chrootProc.SendSignal(syscall.SIGHUP)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case 'i', 'I':
			err = chrootProc.SendSignal(syscall.SIGINT)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case 'k', 'K':
			err = chrootProc.SendSignal(syscall.SIGKILL)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case 'u', 'U':
			err = chrootProc.SendSignal(syscall.SIGUSR1)
			if err != nil {
				fmt.Println(err)
				continue
			}

		}

	}
}
