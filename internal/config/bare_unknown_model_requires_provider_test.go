package config

import (
	"errors"
	"strings"
	"testing"
)

func TestBareUnknownModelRequiresExplicitProvider(t *testing.T) {
	// R-IX2M-C8YE
	_, err := Load(t.TempDir(), []string{"model=not-in-the-catalog-12345"}, func(string) string { return "" })
	if !errors.Is(err, ErrConfig) || !strings.Contains(err.Error(), "set provider explicitly") {
		t.Fatalf("error = %v, want ErrConfig instructing explicit provider", err)
	}
}
