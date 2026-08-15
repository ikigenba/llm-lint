package rules

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// R-GEAU-OB3I
func TestLoadReadsMarkdownFromDirectoryAndExplicitFile(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "catalog")
	if err := os.Mkdir(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	writeTestRule(t, filepath.Join(dir, "from-dir.md"), "directory rule")
	if err := os.WriteFile(filepath.Join(dir, "ignored.txt"), []byte("not a rule"), 0o600); err != nil {
		t.Fatal(err)
	}
	writeTestRule(t, filepath.Join(root, "explicit.md"), "explicit rule")

	got, err := Load(root, []string{"catalog", "explicit.md"})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	ids := make([]string, len(got))
	for i, rule := range got {
		ids[i] = rule.ID
	}
	if !reflect.DeepEqual(ids, []string{"from-dir", "explicit"}) {
		t.Fatalf("Load() ids = %v, want directory and explicit rules", ids)
	}
	if got[0].Description != "directory rule" || got[1].Description != "explicit rule" || got[0].BuiltIn || got[1].BuiltIn {
		t.Fatalf("Load() rules = %#v, want parsed local rules", got)
	}
}

func writeTestRule(t *testing.T, path, description string) {
	t.Helper()
	src := testRuleSource(description, "error", `["**/*"]`, `[]`, "judge this\n")
	if err := os.WriteFile(path, src, 0o600); err != nil {
		t.Fatal(err)
	}
}
