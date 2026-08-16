package config

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestLoadDisableArray(t *testing.T) {
	// R-2K92-YDXU
	t.Run("preserves file order", func(t *testing.T) {
		dir := t.TempDir()
		writeDisableConfig(t, dir, `{"disable":["rule-z","rule-a","rule-m"]}`)

		cfg, err := Load(dir, nil, func(string) string { return "" })
		if err != nil {
			t.Fatal(err)
		}
		want := []string{"rule-z", "rule-a", "rule-m"}
		if !slices.Equal(cfg.Disable, want) {
			t.Fatalf("Disable = %#v, want %#v", cfg.Disable, want)
		}
	})

	for _, tt := range []struct {
		name    string
		content string
	}{
		{name: "non-array", content: `{"disable":"rule-a"}`},
		{name: "non-string element", content: `{"disable":["rule-a",42]}`},
		{name: "null", content: `{"disable":null}`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			writeDisableConfig(t, dir, tt.content)

			_, err := Load(dir, nil, func(string) string { return "" })
			if !errors.Is(err, ErrConfig) || !strings.Contains(err.Error(), "disable") {
				t.Fatalf("error = %v, want ErrConfig naming disable", err)
			}
		})
	}
}

func writeDisableConfig(t *testing.T, dir, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, ".llm-lint.json"), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}
