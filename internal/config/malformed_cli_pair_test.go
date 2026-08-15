package config

import (
	"errors"
	"testing"
)

func TestLoadRejectsMalformedCLIOverrides(t *testing.T) {
	// R-G9F9-584Q
	for _, pair := range []string{"model", "=value", "mystery=value"} {
		t.Run(pair, func(t *testing.T) {
			_, err := Load(t.TempDir(), []string{pair}, func(string) string { return "" })
			if !errors.Is(err, ErrConfig) {
				t.Fatalf("error = %v, want ErrConfig", err)
			}
		})
	}
}
