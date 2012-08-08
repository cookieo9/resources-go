package resources

import (
	"archive/zip"
	"io"
	"os"
)

type zipResource struct {
	*zip.File
}

func (zf *zipResource) Stat() (os.FileInfo, error) { return zf.FileInfo(), nil }
func (zf *zipResource) String() string             { return "ZipFile(" + zf.Name + ")" }
func (zf *zipResource) Path() string               { return zf.Name }

type zipBundleCloser struct {
	file *os.File
	*zipBundle
}

func (zbc *zipBundleCloser) Close() error {
	return zbc.file.Close()
}

// OpenZip opens a file on disk as a zip archive,
// and returns a Bundle of its contents.
func OpenZip(path string) (BundleCloser, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	finfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	bndl, err := OpenZipReader(file, finfo.Size(), path)
	if err != nil {
		return nil, err
	}
	return &zipBundleCloser{file: file, zipBundle: bndl.(*zipBundle)}, nil
}

type zipBundle struct {
	name string
	Bundle
}

func (zrb *zipBundle) String() string {
	return "ZipBundle(" + zrb.name + ")"
}

// OpenZipReader opens a bundle from an open io.ReaderAt. The name parameter
// is simply a symbolic name for debug purposes.
func OpenZipReader(rda io.ReaderAt, size int64, name string) (Bundle, error) {
	var rdr *zip.Reader

	rdr, err := zip.NewReader(rda, size)
	if err != nil {
		rdr2, err2 := zipExeReader(rda, size)
		if err2 != nil {
			return nil, &bundleError{"zip", name, err}
		}
		rdr = rdr2
	}

	var files []Resource
	for _, file := range rdr.File {
		if !file.FileInfo().IsDir() {
			files = append(files, &zipResource{file})
		}
	}
	return &zipBundle{Bundle: OpenList(files), name: name}, nil
}
