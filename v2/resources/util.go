package resources

import (
	"os"
	"os/exec"
	"path/filepath"
)

// ExecutablePath returns a system-native path to the currently running
// executable.
//
// If the value of os.Args[0] has been tampered with, this function may
// give innacurate results.
func ExecutablePath() (string, error) {
	exepath, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	return filepath.Abs(filepath.Clean(exepath))
}
