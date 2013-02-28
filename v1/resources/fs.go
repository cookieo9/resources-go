package resources

import (
	"io"
	"os"
	"path/filepath"
)

// A FSFile represents a file in the filesystem.
type fsResource struct {
	path   string
	bundle *fsBundle
}

// Open opens the file for reading, it is opened using
// os.Open, so the io.ReadCloser is an os.File.
func (f *fsResource) Open() (io.ReadCloser, error) {
	return os.Open(f.FullPath())
}

// Stat returns the os.FileInfo structure describing the file, if
// there is an error, it will be of type *os.PathError.
func (f *fsResource) Stat() (os.FileInfo, error) {
	return os.Stat(f.FullPath())
}

func (f *fsResource) String() string {
	return "FSFile(" + f.path + ")"
}

// Path returns the bundle-local path to the file.
func (f *fsResource) Path() string {
	return f.path
}

// FullPath returns the complete path (not necesarily an absolute path)
// to the file on the filesystem in system native format.
func (f *fsResource) FullPath() string {
	return filepath.Clean(
		filepath.Join(
			f.bundle.Root,
			filepath.FromSlash(f.path),
		))
}

// A FSBundle represents a bundle located in a folder
// in the filesystem.
type fsBundle struct {
	Root string
}

// OpenFS opens a FSBundle with the given root location
// in the filesystem. The path is expected to be in
// system native format.
func OpenFS(root string) Bundle {
	return &fsBundle{filepath.Clean(root)}
}

// get returns a file given a system native path
func (fsb *fsBundle) get(p string) (Resource, error) {
	native_path := fs_bundle_path(fsb.Root, p)
	path := filepath.ToSlash(native_path)
	full_native := filepath.Join(fsb.Root, native_path)

	if _, err := os.Stat(full_native); err != nil {
		return nil, err
	}

	return &fsResource{path: path, bundle: fsb}, nil
}

// Get returns a File representing the file at the given
// path inside the bundle.
func (fsb *fsBundle) Get(path string) (Resource, error) {
	return fsb.get(filepath.FromSlash(path))
}

// Glob returns a list of Files matching the pattern inside
// the bundle.
func (fsb *fsBundle) Glob(pattern string) (files []Resource, err error) {
	var paths []string
	full_pattern := filepath.Join(fsb.Root, filepath.FromSlash(pattern))
	if paths, err = filepath.Glob(full_pattern); err != nil {
		return
	}

	for _, path := range paths {
		file, _ := fsb.get(path)
		files = append(files, file)
	}
	return
}

// Lists every known resource in the Bundle. This is essentially
// every non-directory file found by calling filepath.Walk() on
// the FSBundle's root directory.
func (fsb *fsBundle) List() ([]Resource, error) {
	var list []Resource
	err := filepath.Walk(fsb.Root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !info.IsDir() {
			r, _ := fsb.get(path)
			list = append(list, r)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}

func fs_bundle_path(root string, path string) string {
	root = filepath.Clean(root)
	path = filepath.Clean(path)

	if !filepath.IsAbs(path) {
		return path
	}

	return path[len(root)+1:]
}

func (fsb *fsBundle) String() string {
	return "FSBundle(" + fsb.Root + ")"
}
