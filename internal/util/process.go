package util

import (
	"bufio"
	"fmt"
	"os"
)

// ParentIsShell returns true if the parent process is a shell.
// If this could not be determined, it will return false.
func ParentIsShell() bool {
	ppid := os.Getppid()

	procExePath := fmt.Sprintf("/proc/%d/exe", ppid)
	parentExePath, err := os.Readlink(procExePath)
	if err != nil {
		return false
	}

	shellsFile, err := os.Open("/etc/shells")
	if err != nil {
		return false
	}
	defer shellsFile.Close()

	scanner := bufio.NewScanner(shellsFile)
	for scanner.Scan() {
		line := scanner.Text()
		if line == parentExePath {
			return true
		}
	}

	if err := scanner.Err(); err != nil {
		return false
	}

	return false
}
