package resources

import (
	"os"
	"os/exec"
	"path/filepath"
)

// ExecutablePath returns a system-native path to the currently running
// executable.
//
// BUG: This code uses a simple and portable technique to determine the
// path to the currently running executable. If os.Args[0] is tampered
// with either inside or outside the program, the executable might
// not be found. When issue 4057 in the Go Standard Library is resolved
// this function will be unnecessary, and the bug should no longer
// be present.
func ExecutablePath() (string, error) {
	exepath, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	return filepath.Abs(filepath.Clean(exepath))
}
