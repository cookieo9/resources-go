package resources

import (
	"io"
	"os"
)

// A Resource represents a pointer to an unopened file in a bundle. Resources
// may be volatile references to files so there is no guarantee that a
// Resource returned by a bundle may exist at the time it's methods are called.
// Resources only represent files, not directories (ie: Open() makes sense).
//
// Path() returns the portable path to the file in its bundle.
//
// Open() opens the Resource for reading, it must be closed. If the error
// is not nil, the io.ReadCloser will be nil.
//
// Stat() returns a io.FileInfo structure for the referenced file, or an error
// if a problem occurred accessing that information.
type Resource interface {
	Path() string                 // unique path to this file in the bundle
	Open() (io.ReadCloser, error) // open the file for reading
	Stat() (os.FileInfo, error)   // get information from the file
}

// A Bundle represents a collection of Resources.
//
// Get takes a portable path, and returns the matching resource.
// If no file can be found, an error is returned which can be checked
// with IsNotExist.
//
// Glob takes a portable glob pattern (see path.Glob) and returns a list
// of matching resources, or an error if one occured.
//
// List returns a list of all visible resources in the bundle. It is
// possible that there are additional resources beyond what are listed
// and can thus be retrieved with Get(). This vagueness is intentional
// as it allows for:
//	- Files added to the filesystem after List() is called
//	- Invisible files / directories
//	- Alternate paths to the same file (but weren't listed)
//	- Non-Listable bundles (eg: HTTP Servers)
//
// Portable paths are paths in the UNIX style: directories separated
// by forward slashes "/".
//
// Directories are not considered a resources, and will be omitted from
// List() results, and Glob() matches.
//
// No processing beyond finding/matching should be done by a bundle
// and actual file operations should be done by methods of the Resource
// type.
type Bundle interface {
	Get(string) (Resource, error)    // get a single file
	Glob(string) ([]Resource, error) // find a set of files that match a pattern
	List() ([]Resource, error)       // list all *known* resources in a bundle
}

// A Bundle that can be closed after use will implement the
// BundleCloser interface
type BundleCloser interface {
	Bundle
	Close() error // close underlying resources
}

// A Bundle error is returned by a failing bundle operation
type bundleError struct {
	Op   string
	Path string
	Err  error
}

func (berr *bundleError) Error() string { return berr.Op + " " + berr.Path + ": " + berr.Err.Error() }

// IsNotExist returns true if the error is reporting a file can't be found by a bundle
func IsNotExist(err error) bool {
	if berr, ok := err.(*bundleError); ok {
		return IsNotExist(berr.Err)
	}
	return os.IsNotExist(err)
}

// A BundleSequence stores a list of bundles which are probed sequentially for
// files, and can be used to implement a search path.
type BundleSequence []Bundle

// Get checks the bundles in sequence and returns the first matching file.
// Nil Bundles are skipped.
func (bs BundleSequence) Get(path string) (Resource, error) {
	for _, b := range bs {
		if b == nil {
			continue
		}
		f, err := b.Get(path)
		if err != nil && !IsNotExist(err) {
			return nil, &bundleError{"get", path, err}
		}
		if f != nil {
			return f, nil
		}
	}
	return nil, &bundleError{"get", path, os.ErrNotExist}
}

// Glob finds all resources in all the sequence's bundles that match the given
// pattern. The order of the sub-bundles determines priority w.r.t. collisions.
func (bs BundleSequence) Glob(pattern string) ([]Resource, error) {
	file_list := make([]Resource, 0, 32)
	file_set := make(map[string]bool, 32)

	for _, bundle := range bs {
		if bundle == nil {
			continue
		}

		if matches, err := bundle.Glob(pattern); err != nil {
			return nil, err
		} else {
			for _, match := range matches {
				path := match.Path()
				if _, ok := file_set[path]; ok {
					continue
				}
				file_set[path] = true
				file_list = append(file_list, match)
			}
		}
	}
	return file_list, nil
}

// List returns all the resources in all the sub-bundles. Be careful
// if using this on DefaultBundle() since there may be a lot of files
// in the search path.
func (bs BundleSequence) List() ([]Resource, error) {
	set := make(map[string]bool)
	list := make([]Resource, 0, 32)

	for _, b := range bs {
		if b == nil {
			continue
		}

		sublist, err := b.List()
		if err != nil {
			return nil, err
		}

		for _, res := range sublist {
			p := res.Path()
			if _, ok := set[p]; ok {
				continue
			}
			set[p] = true
			list = append(list, res)
		}
	}

	return list, nil
}
