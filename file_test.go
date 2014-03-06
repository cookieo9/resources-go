package resources

import (
	. "testing"
)

func tryFSGlob(t *T, b Bundle, pattern string) {
	rsrcs, err := b.(Searcher).Glob(pattern)
	if err != nil {
		t.Fatalf("FSBundle.(Searcher).Glob(%q): %v", pattern, err)
	}
	t.Logf("FSBundle.(Searcher).Glob(%q): %v", pattern, rsrcs)
}

func TestFS(t *T) {
	t.Log("Opening CWD")
	b := OpenFS(".")

	tryFSGlob(t, b, "*")
	tryFSGlob(t, b, "*/*")
}
