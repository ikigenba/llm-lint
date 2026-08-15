package walk

import (
	"errors"
	"reflect"
	"testing"
)

// R-GLM8-YXJO
func TestNULByteExcludesBinary(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "binary.dat", "text\x00more")
	writeTestFile(t, root, "plain.txt", "plain")
	runGit := func(string, ...string) ([]byte, error) { return nil, errors.New("not a repo") }

	files, err := (&Walker{Root: root, RunGit: runGit}).Files([]string{"."})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"plain.txt"}; !reflect.DeepEqual(files, want) {
		t.Fatalf("Files() = %v, want %v", files, want)
	}
}
