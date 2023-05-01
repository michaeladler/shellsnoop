package proc

import (
	"os"
	"syscall"
)

type PID int

func IsRunning(pid PID) bool {
	process, err := os.FindProcess(int(pid))
	if err != nil {
		return false
	}

	err = process.Signal(syscall.Signal(0))
	return err == nil
}
