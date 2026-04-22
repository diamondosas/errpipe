package utils

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
)

func IsInstalled(name string) (string, bool) {
	path, err := exec.LookPath(name)
	if err != nil {
		return "", false
	}
	return path, true
}

// IsRunning returns all PIDs whose name or cmdline contains `name`.
// The second return value is true if at least one was found.
func IsRunning(name string) ([]int32, bool) {
	var pids []int32
	nameLower := strings.ToLower(name)

	procs, err := process.Processes()
	if err != nil {
		fmt.Println(err)
		return nil, false
	}

	for _, p := range procs {
		n, _ := p.Name()
		n = strings.ToLower(strings.TrimSuffix(n, ".exe"))
		if n == nameLower {
			pids = append(pids, p.Pid)
			continue
		}

		cmdline, err := p.Cmdline()
		if err != nil {
			continue
		}
		if strings.Contains(strings.ToLower(cmdline), nameLower) {
			pids = append(pids, p.Pid)
		}
	}

	return pids, len(pids) > 0
}