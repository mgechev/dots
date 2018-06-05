package dots

import (
	"log"
	"strings"
	"testing"
)

func TestResolveNoArgs(t *testing.T) {
	result, err := Resolve([]string{}, []string{})

	files := []string{}

	if err != nil {
		t.Error("Got errors")
	}

	if len(result) != len(files) {
		t.Error("Matched different number of files")
	}

	if err != nil {
		log.Fatal(err)
	}
}

func TestResolve(t *testing.T) {
	result, err := Resolve([]string{"fixtures/dummy/..."}, []string{"fixtures/dummy/foo"})

	files := []string{
		"fixtures/dummy/bar/bar1.go",
		"fixtures/dummy/bar/bar2.go",
		"fixtures/dummy/baz/baz1.go",
		"fixtures/dummy/baz/baz2.go",
		"fixtures/dummy/baz/baz3.go",
	}

	if err != nil {
		t.Error("Got errors")
	}

	if len(result) != len(files) {
		t.Error("Matched different number of files")
	}

	if err != nil {
		log.Fatal(err)
	}

	for _, r := range result {
		matched := false
		for _, e := range files {
			matched = matched || strings.HasSuffix(r, e)
		}
		if !matched {
			t.Error("Not supposed to match: " + r)
		}
	}
}

func TestPackageResolve(t *testing.T) {
	result, err := Resolve([]string{"github.com/mgechev/dots"}, []string{"resolve_test.go"})

	files := []string{
		"resolve.go",
	}

	if err != nil {
		t.Error("Got errors")
	}

	if len(result) != len(files) {
		t.Error("Matched different number of files")
	}

	if err != nil {
		log.Fatal(err)
	}

	for _, r := range result {
		matched := false
		for _, e := range files {
			matched = matched || strings.HasSuffix(r, e)
		}
		if !matched {
			t.Error("Not supposed to match: " + r)
		}
	}
}

func TestSkipWildcard(t *testing.T) {
	result, err := Resolve([]string{"fixtures/dummy/"}, []string{"fixtures/dummy/..."})

	files := map[string]bool{}

	if err != nil {
		t.Error("Got errors")
	}

	if len(result) != len(files) {
		t.Error("Matched different number of files")
	}

	if err != nil {
		log.Fatal(err)
	}
}

func TestPackageWildcard(t *testing.T) {
	result, err := Resolve([]string{"github.com/mgechev/dots/fixtures/pkg/foo/...", "github.com/mgechev/dots/fixtures/pkg/baz"}, []string{})
	files := []string{
		"baz1.go",
		"baz2.go",
		"foo1.go",
		"foo2.go",
		"bar1.go",
	}

	if err != nil {
		t.Error("Got errors")
	}

	if len(result) != len(files) {
		t.Error("Matched different number of files")
	}

	if err != nil {
		log.Fatal(err)
	}

	for _, r := range result {
		matched := false
		for _, e := range files {
			matched = matched || strings.HasSuffix(r, e)
		}
		if !matched {
			t.Error("Not supposed to match: " + r)
		}
	}
}

func TestPackageWildcardWithSkip(t *testing.T) {
	result, err := Resolve([]string{"github.com/mgechev/dots/fixtures/pkg/baz"}, []string{"github.com/mgechev/dots/fixtures/pkg/foo/..."})

	files := []string{
		"baz1.go",
		"baz2.go",
	}

	if err != nil {
		t.Error("Got errors")
	}

	if len(result) != len(files) {
		t.Error("Matched different number of files")
	}

	if err != nil {
		log.Fatal(err)
	}

	for _, r := range result {
		matched := false
		for _, e := range files {
			matched = matched || strings.HasSuffix(r, e)
		}
		if !matched {
			t.Error("Not supposed to match: " + r)
		}
	}
}

func TestComplainForMissingDirectories(t *testing.T) {
	_, err := Resolve([]string{"./fixturess"}, []string{})

	if err == nil {
		t.Error("Should get an error")
	}
}

func TestComplainForMissingPackages(t *testing.T) {
	_, err := Resolve([]string{"github.com/mgechev/bazbaz"}, []string{})

	if err == nil {
		t.Error("Should get an error")
	}
}

func TestResolvePackages(t *testing.T) {
	result, err := ResolvePackages([]string{"github.com/mgechev/dots/fixtures/pkg/foo/...", "github.com/mgechev/dots/fixtures/pkg/baz"}, []string{})

	if err != nil {
		t.Error("Got errors")
	}

	if len(result) != 3 {
		t.Error("Matched different number of files")
	}

	for _, pkg := range result {
		if len(pkg) == 0 {
			t.Error("Empty package")
		}
	}

}
