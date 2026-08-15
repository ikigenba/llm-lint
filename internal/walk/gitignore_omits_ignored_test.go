package walk

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
)

func gitRunner(dir string, args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	return cmd.Output()
}

func writeTestFile(t *testing.T, root, name, contents string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func initGitRepo(t *testing.T, root string) {
	t.Helper()
	cmd := exec.Command("git", "init", "-q")
	cmd.Dir = root
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, output)
	}
}

// R-GJ6G-7E2A
func TestGitignoreOmitsIgnored(t *testing.T) {
	root := t.TempDir()
	initGitRepo(t, root)
	writeTestFile(t, root, ".gitignore", "ignored.txt\n")
	writeTestFile(t, root, "tracked.txt", "tracked")
	writeTestFile(t, root, "untracked.txt", "untracked")
	writeTestFile(t, root, "ignored.txt", "ignored")
	cmd := exec.Command("git", "add", "tracked.txt")
	cmd.Dir = root
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git add: %v: %s", err, output)
	}

	files, err := (&Walker{Root: root, RunGit: gitRunner}).Files([]string{"."})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{".gitignore", "tracked.txt", "untracked.txt"}
	if !reflect.DeepEqual(files, want) {
		t.Fatalf("Files() = %v, want %v", files, want)
	}
}
