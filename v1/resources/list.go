package resources

import (
	"os"
	"path"
)

// A Lister represents a object with a collection
// of Resources. It can be used by the OpenLister
// function to create a fully featured Bundle,
// and is exported so that external code can
// easily generate their own Bundles.
type Lister interface {
	List() []Resource
}

type listerBundle struct {
	Lister
}

func (lb *listerBundle) Get(resourcePath string) (Resource, error) {
	resourcePath = path.Clean(resourcePath)

	for _, item := range lb.List() {
		if item.Path() == resourcePath {
			return item, nil
		}
	}

	return nil, &bundleError{"get", resourcePath, os.ErrNotExist}
}

func (lb *listerBundle) Glob(pattern string) ([]Resource, error) {
	var resources []Resource

	for _, item := range lb.List() {
		ok, err := path.Match(pattern, item.Path())
		if err != nil {
			return nil, &bundleError{"glob", pattern, err}
		}
		if ok {
			resources = append(resources, item)
		}
	}

	return resources, nil
}

type sliceLister []Resource

func (sl sliceLister) List() []Resource {
	return ([]Resource)(sl)
}

// OpenList takes a list of resources, and
// returns a bundle that represents them.
func OpenList(list []Resource) Bundle {
	return &listerBundle{sliceLister(list)}
}

// OpenLister converts a Lister into a Bundle. If the
// Lister is already a bundle it is returned un-altered.
func OpenLister(lister Lister) Bundle {
	if bundle, ok := lister.(Bundle); ok {
		return bundle
	}
	return &listerBundle{
		Lister: lister,
	}
}
