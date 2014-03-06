package resources

import (
	"path/filepath"
)

func ExampleBundleSequence() {
	var (
		searchPath                         BundleSequence
		zipBundle, exeBundle, pluginBundle Bundle
	)

	zipBundle, _ = OpenZip("foo.zip")
	if exe, err := ExecutablePath(); err == nil {
		exeBundle, _ = OpenZip(exe)
		exeDir := filepath.Dir(exe)
		pluginBundle = OpenFS(filepath.Join(exeDir, "plugins"))
	}

	searchPath = BundleSequence{
		pluginBundle,
		zipBundle,
		exeBundle,
	}

	if _, err := searchPath.Glob("*.png"); err != nil {
		panic("No PNG Files!")
	}
}
