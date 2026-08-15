package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFindsAncestor(t *testing.T) {
	// R-FZO2-3276
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, ".llm-lint.json"), []byte(`{"enable":["local-rule"]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(root, "a", "b")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(nested, nil, func(string) string { return "" })
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Root != root {
		t.Fatalf("Root = %q, want %q", cfg.Root, root)
	}
	if len(cfg.Enable) != 1 || cfg.Enable[0] != "local-rule" {
		t.Fatalf("Enable = %#v", cfg.Enable)
	}
}
