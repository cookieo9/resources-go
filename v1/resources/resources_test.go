package resources

import (
	. "testing"
)

func TestSequence(t *T) {
	zpath := "foo.zip"
	pattern := "*_test.go"

	t.Log("Opening", zpath)
	zb, err := OpenZip(zpath)
	if err != nil {
		t.Fatal(err)
	}
	defer zb.Close()

	cp, err := OpenCurrentPackage()
	if err != nil {
		t.Fatalf("OpenCurrentPackage(): %v", err)
	}

	bundle := BundleSequence{cp, zb}
	bundle = append(bundle, DefaultBundle)
	t.Log("Bundle:", bundle)

	t.Log("List():", bundle.List())

	files, err := bundle.Glob(pattern)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Glob(%s): %+v", pattern, files)

	f, err := bundle.Get(files[0].Path())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Get(%s): %+v", files[0].Path(), f)
}
