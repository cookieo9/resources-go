package resources

import (
	. "testing"
)

func test_zip(t *T, zb Bundle) {
	pattern := "*_test.go"
	t.Log("Testing:", zb)
	t.Log("List():", zb.List())

	files, err := zb.Glob(pattern)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Glob(%s): %+v", pattern, files)
}

const ZipTestFile = "foo.zip"

func TestZip(t *T) {
	t.Log("Opening", ZipTestFile)
	zb, err := OpenZip(ZipTestFile)
	if err != nil {
		t.Fatal(err)
	}
	defer zb.Close()

	test_zip(t, zb)
}
