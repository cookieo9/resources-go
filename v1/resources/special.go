package resources

import (
	"bitbucket.org/kardianos/osext"
	"fmt"
	"go/build"
	"io"
	"path/filepath"
)

// An AutoBundle is a wrapper around a niladic
// bundle generator function (eg: OpenCurrentPackage)
// which calls the generator prior to all Bundle
// operations to receive a fresh bundle to operate
// on. It is used to implement the behaviour of the
// Default search sequence w.r.t. opening the bundle
// of the current package at the time of the call.
type AutoBundle func() (Bundle, error)

// Get generates a bundle, and calls Get on it.
func (ab AutoBundle) Get(path string) (Resource, error) {
	b, e := ab()
	if e != nil {
		return nil, e
	}
	return b.Get(path)
}

// Glob generates a bundle, and calls Glob on it.
func (ab AutoBundle) Glob(path string) ([]Resource, error) {
	b, e := ab()
	if e != nil {
		return nil, e
	}
	return b.Glob(path)
}

// List generates a bundle, and calls List on it.
func (ab AutoBundle) List() []Resource {
	b, e := ab()
	if e != nil {
		return nil
	}
	return b.List()
}

// Gets a string representation of this bundle.
func (ab AutoBundle) String() string {
	b, e := ab()
	if e != nil {
		return e.Error()
	}
	return fmt.Sprintf("AutoBundle:%v", b)
}

var (
	ExecutableDirectoryBundle Bundle       // Directory containing the executable
	ExecutableZipBundle       Bundle       // Executable file as zip bundle (possibly nil)
	WorkingDirectoryBundle    Bundle       // Working directory bundle
	CurrentPackageBundle      Bundle       // The current package bundle
	GopathBundle              NoListBundle // The gopath source directories as an un-List()-able bundle
)

// DefaultBundle holds the bundle to use by package functions
// resources.Get, resources.Glob, and resources.List.
//
// It implements the following search-path:
//	- The current working directory
//	- The directory containing the executable
//	- The package source directory of the calling code
//	- The GOPATH source directories
//	- The executable treated as a ZipFile
var DefaultBundle BundleSequence

func init() {
	exepath, _ := ExecutablePath()

	ExecutableZipBundle, _ = OpenZip(exepath)
	ExecutableDirectoryBundle = OpenFS(filepath.Dir(exepath))
	WorkingDirectoryBundle = AutoBundle(func() (Bundle, error) { return OpenFS("."), nil })
	CurrentPackageBundle = AutoBundle(OpenCurrentPackage)

	gopathdirs := build.Default.SrcDirs()
	gopathbundle := make(BundleSequence, len(gopathdirs))
	for i, dir := range gopathdirs {
		gopathbundle[i] = OpenFS(dir)
	}
	GopathBundle = NoListBundle{gopathbundle}

	DefaultBundle = BundleSequence{
		WorkingDirectoryBundle,
		ExecutableDirectoryBundle,
		CurrentPackageBundle,
		GopathBundle,
		ExecutableZipBundle,
	}
}

// ExecutablePath returns a system-native path to the currently running
// executable.
//
// It is implemented using the bitbucket.org/kardianos/osext package.
func ExecutablePath() (string, error) {
	return osext.Executable()
}

// List runs DefaultBundle.List()
func List() []Resource {
	return DefaultBundle.List()
}

// Get runs DefaultBundle.Get()
func Get(path string) (Resource, error) {
	return DefaultBundle.Get(path)
}

// Glob runs DefaultBundle.Glob()
func Glob(pattern string) ([]Resource, error) {
	return DefaultBundle.Glob(pattern)
}

// Open opens a file in the default search path.
func Open(path string) (io.ReadCloser, error) {
	res, err := DefaultBundle.Get(path)
	if err != nil {
		return nil, err
	}
	return res.Open()
}

// A NoListBundle wraps another bundle, and
// replaces the List() method to return nil at
// all times.
type NoListBundle struct {
	Bundle
}

// Always returns nil (AKA: an empty list of resources.)
func (nlb NoListBundle) List() []Resource {
	return nil
}
