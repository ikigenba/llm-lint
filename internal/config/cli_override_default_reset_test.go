package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCLIOverrideDefaultResetsFileValue(t *testing.T) {
	// R-G5RJ-ZWWN
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".llm-lint.json"), []byte(`{"model":{"temperature":0.2}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(dir, []string{"temperature=default"}, func(string) string { return "" })
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Model.Temperature != nil {
		t.Fatalf("Temperature = %v, want unset", *cfg.Model.Temperature)
	}
}
