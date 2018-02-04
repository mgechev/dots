package dots

import (
	"os"
	"path/filepath"
	"strings"
)

// Resolve accepts a slice of paths with optional "..." placeholder and a slice with paths to be skipped.
// The final result is the set of all files from the selected directories subtracted with
// the files in the skip slice.
func Resolve(paths, skip []string) ([]string, error) {
	if len(paths) == 0 {
		return []string{"."}, nil
	}

	skipPath := newPathFilter(skip)
	dirs := newStringSet()
	for _, path := range paths {
		if strings.HasSuffix(path, "/...") {
			root := filepath.Dir(path)
			err := filepath.Walk(root, func(p string, i os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				skip := skipPath(p)
				switch {
				case i.IsDir() && skip:
					return filepath.SkipDir
				case !i.IsDir() && !skip && strings.HasSuffix(p, ".go"):
					dirs.add(filepath.Clean(filepath.Dir(p)))
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
		} else {
			dirs.add(filepath.Clean(path))
		}
	}
	out := make([]string, 0, dirs.size())
	for _, d := range dirs.asSlice() {
		out = append(out, relativePackagePath(d))
	}
	return loadFiles(out, skipPath)
}

func newPathFilter(skip []string) func(string) bool {
	filter := map[string]bool{}
	for _, name := range skip {
		filter[name] = true
	}

	return func(path string) bool {
		base := filepath.Base(path)
		if filter[base] || filter[path] {
			return true
		}
		return base != "." && base != ".." && strings.ContainsAny(base[0:1], "_.")
	}
}

// loadFiles from a list of directories
func loadFiles(paths []string, skipPath func(string) bool) ([]string, error) {
	filePaths := []string{}
	for _, dir := range paths {
		paths, err := filepath.Glob(filepath.Join(dir, "*.go"))
		if err != nil {
			return nil, err
		}
		for _, f := range paths {
			if !skipPath(f) {
				filePaths = append(filePaths, f)
			}
		}
	}
	return filePaths, nil
}

type stringSet struct {
	items map[string]struct{}
}

func newStringSet(items ...string) *stringSet {
	setItems := make(map[string]struct{}, len(items))
	for _, item := range items {
		setItems[item] = struct{}{}
	}
	return &stringSet{items: setItems}
}

func (s *stringSet) add(item string) {
	s.items[item] = struct{}{}
}

func (s *stringSet) asSlice() []string {
	items := []string{}
	for item := range s.items {
		items = append(items, item)
	}
	return items
}

func (s *stringSet) size() int {
	return len(s.items)
}

func relativePackagePath(dir string) string {
	if filepath.IsAbs(dir) || strings.HasPrefix(dir, ".") {
		return dir
	}
	return "./" + dir
}
