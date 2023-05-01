package storage

import (
	"sync"

	"github.com/michaeladler/shellsnoop/internal/log"
	"github.com/michaeladler/shellsnoop/internal/proc"
)

var (
	data  = make(map[proc.PID]string, 0)
	mutex sync.RWMutex
)

func Add(pid int, cmdline string) {
	mutex.Lock()
	data[proc.PID(pid)] = cmdline
	mutex.Unlock()
}

func Get(pid int) string {
	var cmdline string
	var ok bool
	p := proc.PID(pid)
	mutex.RLock()
	cmdline, ok = data[p]
	mutex.RUnlock()
	if !ok {
		return ""
	}
	return cmdline
}

func GC() {
	logger := log.Logger

	// make copy of pids
	mutex.RLock()
	toCheck := make([]proc.PID, 0, len(data))
	for k := range data {
		toCheck = append(toCheck, k)
	}
	mutex.RUnlock()

	// check if pids are still alive
	toRemove := make([]proc.PID, 0, len(toCheck))
	for _, pid := range toCheck {
		if !proc.IsRunning(pid) {
			toRemove = append(toRemove, pid)
			continue
		}
	}

	// remove dead pids
	if len(toRemove) > 0 {
		logger.Info("Removing stale entries", "entries", toRemove)
	}
	mutex.Lock()
	for _, pid := range toRemove {
		delete(data, pid)
	}
	mutex.Unlock()
}
