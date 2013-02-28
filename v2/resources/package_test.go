package resources

import (
	. "testing"
)

func TestPackageBundle(t *T) {
	cp, err := OpenCurrentPackage()
	if err != nil {
		t.Fatalf("OpenCurrentPacakge(): %v", err)
	}

	t.Log("CurrentPackage:", cp)
	list, err := cp.(Lister).List()
	if err != nil {
		t.Fatal("cp.List():", err)
	}
	t.Log("cp.List():", list)
	if fs, err := cp.(Searcher).Glob("*.go"); err != nil {
		t.Fatal("Glob(*.go):", err)
	} else {
		t.Log("Glob(*.go):", fs)
	}
}
