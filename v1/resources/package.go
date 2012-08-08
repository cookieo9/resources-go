package resources

import (
	"go/build"
	"path/filepath"
	"runtime"
	"strings"
)

// Check a package to see if it's a valid place to look for resources
// Anywhere in GOROOT, and in the resources package are not valid
// locations
func checkPackage(pkg *build.Package) bool {
	if strings.Index(pkg.Dir, runtime.GOROOT()) != -1 {
		return false
	}

	var thispkg *build.Package
	_, sfile, _, _ := runtime.Caller(0)
	if p, err := build.ImportDir(filepath.Dir(sfile), build.FindOnly); err != nil {
		panic(err)
	} else {
		thispkg = p
	}

	if pkg.Dir == thispkg.Dir {
		return false
	}
	return true
}

// Opens the source directory of the current package as a Bundle.
// The current package is the package of the code calling
// OpenCurrentPackage() (as determined by runtime.Caller())
func OpenCurrentPackage() (Bundle, error) {
	// Keep calling runtime.Caller with increasing values until we are no longer in
	// this package
	for i := 1; ; i++ {
		_, sfile, _, _ := runtime.Caller(i)
		if p, err := build.ImportDir(filepath.Dir(sfile), build.FindOnly); err == nil {
			if checkPackage(p) {
				return OpenFS(p.Dir), nil
			}
		} else {
			return nil, err
		}
	}
	panic("Shouldn't Get Here!")
}

// OpenPackagePath returns a FSBundle which accesses files
// in the source directory of the package named by the given
// import path.
func OpenPackage(import_path string) (Bundle, error) {
	pkg, err := build.Import(import_path, "", build.FindOnly)
	if err != nil {
		return nil, err
	}
	return OpenFS(pkg.Dir), nil
}
