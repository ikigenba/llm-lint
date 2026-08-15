package walk

import (
	"reflect"
	"testing"
)

// R-GQHU-I0IG
func TestDefaultDotAndExplicitIgnored(t *testing.T) {
	root := t.TempDir()
	initGitRepo(t, root)
	writeTestFile(t, root, ".gitignore", "ignored.txt\n")
	writeTestFile(t, root, "visible.txt", "visible")
	writeTestFile(t, root, "ignored.txt", "ignored")
	walker := &Walker{Root: root, RunGit: gitRunner}

	files, err := walker.Files(nil)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{".gitignore", "visible.txt"}; !reflect.DeepEqual(files, want) {
		t.Fatalf("Files(nil) = %v, want %v", files, want)
	}
	files, err = walker.Files([]string{"ignored.txt"})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"ignored.txt"}; !reflect.DeepEqual(files, want) {
		t.Fatalf("Files(ignored) = %v, want %v", files, want)
	}
}
