package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadRejectsMalformedJSONAndUnknownKeys(t *testing.T) {
	// R-G23U-ULOK
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{name: "malformed", content: `{"enable":`, want: ".llm-lint.json"},
		{name: "top level key", content: `{"mystery":true}`, want: "mystery"},
		{name: "model key", content: `{"model":{"mystery":true}}`, want: "mystery"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.WriteFile(filepath.Join(dir, ".llm-lint.json"), []byte(tt.content), 0o600); err != nil {
				t.Fatal(err)
			}
			_, err := Load(dir, nil, func(string) string { return "" })
			if !errors.Is(err, ErrConfig) || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %v, want ErrConfig containing %q", err, tt.want)
			}
		})
	}
}
