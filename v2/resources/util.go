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

// IsNotFound returns true if the error given is an error representing
// a Resource that was not found.
func IsNotFound(e error) bool {
	return e == ErrNotFound
}

// CheckPath() returns nil if given a valid path. Valid paths are
// forward slash delimeted, relative paths, which don't escape the
// base-level directory.
//
// Otherwise it returns one of the following error types:
//  - ErrEscapeRoot: if the path leaves the base directory
//  - ErrNotRelative: if the path is not a relative path
func CheckPath(path string) error {
	clean := filepath.Clean(path)
	if len(clean) >= 2 && clean[:2] == ".." {
		return ErrEscapeRoot
	}
	if len(clean) >= 1 && clean[0] == '/' {
		return ErrNotRelative
	}
	return nil
}
