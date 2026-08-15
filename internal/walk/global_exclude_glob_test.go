package walk

import (
	"errors"
	"reflect"
	"testing"
)

// R-GO21-QH12
func TestGlobalExcludeGlob(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "keep.go", "package keep")
	writeTestFile(t, root, "generated/deep/drop.go", "package drop")
	runGit := func(string, ...string) ([]byte, error) { return nil, errors.New("not a repo") }

	files, err := (&Walker{Root: root, Exclude: []string{"generated/**"}, RunGit: runGit}).Files([]string{"."})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"keep.go"}; !reflect.DeepEqual(files, want) {
		t.Fatalf("Files() = %v, want %v", files, want)
	}
}
