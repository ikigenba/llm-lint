package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCLIOverrideReplacesFileModelKey(t *testing.T) {
	// R-G3BR-8DF9
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".llm-lint.json"), []byte(`{"model":{"temperature":0.2}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(dir, []string{"temperature=0.8"}, func(string) string { return "" })
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Model.Temperature == nil || *cfg.Model.Temperature != 0.8 {
		t.Fatalf("Temperature = %v, want 0.8", cfg.Model.Temperature)
	}
}
