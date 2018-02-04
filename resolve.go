package dots

import (
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Resolve accepts a slice of paths with optional "..." placeholder and a slice with paths to be skipped.
// The final result is the set of all files from the selected directories subtracted with
// the files in the skip slice.
func Resolve(includePatterns, skipPatterns []string) ([]string, []error) {
	skip, errs := resolvePatterns(skipPatterns)
	filter := newPathFilter(skip)

	pathSet := map[string]bool{}
	include, includeErrs := resolvePatterns(includePatterns)
	errs = append(errs, includeErrs...)

	var result []string
	for _, i := range include {
		if _, ok := pathSet[i]; !ok && !filter(i) {
			pathSet[i] = true
			result = append(result, i)
		}
	}
	return result, errs
}

func isDir(filename string) bool {
	fi, err := os.Stat(filename)
	return err == nil && fi.IsDir()
}

func isFile(filename string) bool {
	fi, err := os.Stat(filename)
	return err == nil && !fi.IsDir()
}

func readDir(dirname string, recurse bool) ([]string, error) {
	var files []string
	var appendFile = func(file string, info os.FileInfo, err error) error {
		if strings.HasSuffix(file, ".go") {
			files = append(files, file)
		}
		return nil
	}
	if recurse {
		err := filepath.Walk(dirname, appendFile)
		return files, err
	}
	res, err := ioutil.ReadDir(dirname)
	if err == nil {
		for _, f := range res {
			appendFile(filepath.Join(dirname, f.Name()), f, nil)
		}
	}
	return files, err
}

func readNestedPackages(root string, current string, recurse bool) (files []string, errs []error) {
	if strings.HasPrefix(current, root) {
		pkg, e := build.Import(current, ".", 0)
		if e != nil {
			if _, nogo := e.(*build.NoGoError); !nogo {
				errs = append(errs, e)
			}
		} else {
			var pkgFiles []string
			pkgFiles = append(pkgFiles, pkg.GoFiles...)
			pkgFiles = append(pkgFiles, pkg.CgoFiles...)
			pkgFiles = append(pkgFiles, pkg.TestGoFiles...)
			if pkg.Dir != "." {
				for i, f := range pkgFiles {
					pkgFiles[i] = filepath.Join(pkg.Dir, f)
				}
			}
			files = append(files, pkgFiles...)
		}
		if recurse {
			for _, i := range pkg.Imports {
				fls, e := readNestedPackages(root, i, recurse)
				files = append(files, fls...)
				errs = append(errs, e...)
			}
		}
	}
	return files, errs
}

func readPackage(packageName string, recurse bool) ([]string, []error) {
	return readNestedPackages(packageName, packageName, recurse)
}

func resolvePattern(pattern string) ([]string, []error) {
	recurse := false
	if strings.HasSuffix(pattern, "/...") {
		recurse = true
		pattern = strings.Replace(pattern, "/...", "", 1)
	}
	if isDir(pattern) {
		abs, err := filepath.Abs(pattern)
		if err != nil {
			return nil, []error{err}
		}
		res, err := readDir(abs, recurse)
		if err != nil {
			return res, []error{err}
		}
		return res, nil
	}
	if isFile(pattern) {
		return []string{pattern}, nil
	}
	return readPackage(pattern, recurse)
}

func resolvePatterns(patterns []string) ([]string, []error) {
	var paths []string
	var errs []error
	for _, s := range patterns {
		res, err := resolvePattern(s)
		if err != nil {
			errs = append(errs, err...)
		} else {
			paths = append(paths, res...)
		}
	}
	return paths, errs
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
