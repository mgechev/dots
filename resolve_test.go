package dots

import (
	"log"
	"testing"
)

func TestResolve(t *testing.T) {
	result, err := Resolve([]string{"fixtures/dummy/..."}, []string{"fixtures/dummy/foo"})

	files := map[string]bool{
		"fixtures/dummy/bar/bar1.go": true,
		"fixtures/dummy/bar/bar2.go": true,
		"fixtures/dummy/baz/baz1.go": true,
		"fixtures/dummy/baz/baz2.go": true,
		"fixtures/dummy/baz/baz3.go": true,
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
		if _, ok := files[r]; !ok {
			t.Error("Not supposed to match: " + r)
		}
	}
}

func TestPackageResolve(t *testing.T) {
	result, err := Resolve([]string{"github.com/mgechev/dots"}, []string{"resolve_test.go"})

	files := map[string]bool{
		"resolve.go": true,
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
		if _, ok := files[r]; !ok {
			t.Error("Not supposed to match" + r)
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

	for _, r := range result {
		if _, ok := files[r]; !ok {
			t.Error("Not supposed to match" + r)
		}
	}
}

func TestPackageWildcard(t *testing.T) {
	result, err := Resolve([]string{"github.com/mgechev/dots/fixtures/pkg/foo/...", "github.com/mgechev/dots/fixtures/pkg/baz"}, []string{})

	files := map[string]bool{
		"baz1.go": true,
		"baz2.go": true,
		"foo1.go": true,
		"foo2.go": true,
		"bar1.go": true,
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
		if _, ok := files[r]; !ok {
			t.Error("Not supposed to match" + r)
		}
	}
}

func TestPackageWildcardWithSkip(t *testing.T) {
	result, err := Resolve([]string{"github.com/mgechev/dots/fixtures/pkg/baz"}, []string{"github.com/mgechev/dots/fixtures/pkg/foo/..."})

	files := map[string]bool{
		"baz1.go": true,
		"baz2.go": true,
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
		if _, ok := files[r]; !ok {
			t.Error("Not supposed to match" + r)
		}
	}
}
