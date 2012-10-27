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

// An ErrNotFound is returned when a Resource cannot be found.
type ErrNotFound struct {
	Path string // The path which couldn't be found
}

func (enf ErrNotFound) Error() string {
	return "Resource Not Found: " + enf.Path
}

// IsNotFound returns true if the error given is an error representing
// a Resource that was not found.
func IsNotFound(e error) bool {
	_, ok := e.(*ErrNotFound)
	return ok
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
		return (*ErrEscapeRoot)(&path)
	}
	if len(clean) >= 1 && clean[0] == '/' {
		return (*ErrNotRelative)(&path)
	}
	return nil
}

// An ErrEscapeRoot is an error returned when a Resource path
// escapes the root directory of a bundle.
type ErrEscapeRoot string

func (eer ErrEscapeRoot) Error() string {
	return "Path escapes root: " + string(eer)
}

// An ErrNotRelative error is returned when a Resource path
// is not a relative path.
type ErrNotRelative string

func (enr ErrNotRelative) Error() string {
	return "Path not relative: " + string(enr)
}
