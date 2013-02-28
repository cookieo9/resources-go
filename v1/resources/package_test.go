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
	if list, err := cp.List(); err != nil {
		t.Error("cp.List():", err)
	} else {
		t.Log("cp.List():", list)
	}
	if fs, err := cp.Glob("*.go"); err != nil {
		t.Fatal("Glob(*.go):", err)
	} else {
		t.Log("Glob(*.go):", fs)
	}
}
