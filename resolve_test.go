package dots

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveNoArgs(t *testing.T) {
	result, err := Resolve([]string{}, []string{})
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 0 {
		t.Errorf("Matched different number of files: got=%v, want=0", len(result))
	}
}

func TestResolve(t *testing.T) {
	result, err := Resolve([]string{"fixtures/dummy/..."}, []string{"fixtures/dummy/foo", "fixtures/dummy/UNKNOWN"})

	files := []string{
		filepath.FromSlash("fixtures/dummy/bar/bar1.go"),
		filepath.FromSlash("fixtures/dummy/bar/bar2.go"),
		filepath.FromSlash("fixtures/dummy/baz/baz1.go"),
		filepath.FromSlash("fixtures/dummy/baz/baz2.go"),
		filepath.FromSlash("fixtures/dummy/baz/baz3.go"),
	}

	if err != nil {
		t.Fatal(err)
	}
	if len(result) != len(files) {
		t.Fatalf("Matched different number of files: got=%v, want=%v", len(result), len(files))
	}
	for _, r := range result {
		matched := false
		for _, e := range files {
			matched = matched || strings.HasSuffix(r, e)
		}
		if !matched {
			t.Errorf("Not supposed to match: %v", r)
		}
	}
}

func TestResolveTests(t *testing.T) {
	result, err := Resolve([]string{"fixtures/withtests/..."}, []string{})

	files := []string{
		filepath.FromSlash("fixtures/withtests/gofile.go"),
		filepath.FromSlash("fixtures/withtests/testgo_test.go"),
		filepath.FromSlash("fixtures/withtests/xtestgo_test.go"),
	}

	if err != nil {
		t.Fatal(err)
	}
	if len(result) != len(files) {
		t.Fatalf("Matched different number of files: got=%v, want=%v", len(result), len(files))
	}
	for _, r := range result {
		matched := false
		for _, e := range files {
			matched = matched || strings.HasSuffix(r, e)
		}
		if !matched {
			t.Errorf("Not supposed to match: %v", r)
		}
	}
}

func TestPackageResolve(t *testing.T) {
	result, err := Resolve([]string{"github.com/mgechev/dots"}, []string{"resolve_test.go"})

	files := []string{
		"resolve.go",
	}

	if err != nil {
		t.Fatal(err)
	}
	if len(result) != len(files) {
		t.Fatalf("Matched different number of files: got=%v, want=%v", len(result), len(files))
	}
	for _, r := range result {
		matched := false
		for _, e := range files {
			matched = matched || strings.HasSuffix(r, e)
		}
		if !matched {
			t.Errorf("Not supposed to match: %v", r)
		}
	}
}

func TestSkipWildcard(t *testing.T) {
	result, err := Resolve([]string{"fixtures/dummy/"}, []string{"fixtures/dummy/..."})

	files := map[string]bool{}

	if err != nil {
		t.Fatal(err)
	}
	if len(result) != len(files) {
		t.Fatalf("Matched different number of files: got=%v, want=%v", len(result), len(files))
	}
}

func TestPackageWildcard(t *testing.T) {
	result, err := Resolve([]string{"github.com/mgechev/dots/fixtures/pkg/foo/...", "github.com/mgechev/dots/fixtures/pkg/baz"}, []string{})
	files := []string{
		"baz1.go",
		"baz2.go",
	}

	if err != nil {
		t.Fatal(err)
	}
	if len(result) != len(files) {
		t.Fatalf("Matched different number of files: got=%v, want=%v", len(result), len(files))
	}
	for _, r := range result {
		matched := false
		for _, e := range files {
			matched = matched || strings.HasSuffix(r, e)
		}
		if !matched {
			t.Errorf("Not supposed to match: %v", r)
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
		t.Fatal(err)
	}
	if len(result) != len(files) {
		t.Fatalf("Matched different number of files: got=%v, want=%v", len(result), len(files))
	}
	for _, r := range result {
		matched := false
		for _, e := range files {
			matched = matched || strings.HasSuffix(r, e)
		}
		if !matched {
			t.Errorf("Not supposed to match: %v", r)
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
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("Matched different number of files: got=%v, want=%v", len(result), 1)
	}
	for _, pkg := range result {
		if len(pkg) == 0 {
			t.Error("Empty package")
		}
	}
}
