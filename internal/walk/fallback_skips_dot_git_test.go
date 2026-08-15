package walk

import (
	"errors"
	"reflect"
	"testing"
)

// R-GKEC-L5SZ
func TestFallbackSkipsDotGit(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "one.txt", "one")
	writeTestFile(t, root, "nested/two.txt", "two")
	writeTestFile(t, root, ".git/private", "private")
	runGit := func(string, ...string) ([]byte, error) { return nil, errors.New("git unavailable") }

	files, err := (&Walker{Root: root, RunGit: runGit}).Files([]string{"."})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"nested/two.txt", "one.txt"}
	if !reflect.DeepEqual(files, want) {
		t.Fatalf("Files() = %v, want %v", files, want)
	}
}
