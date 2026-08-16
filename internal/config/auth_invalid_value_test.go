package config

import (
	"errors"
	"strings"
	"testing"
)

func TestLoadRejectsInvalidAuthValue(t *testing.T) {
	// R-IZIF-3SFS
	_, err := Load(t.TempDir(), []string{"auth=magic"}, func(string) string { return "" })
	if !errors.Is(err, ErrConfig) || !strings.Contains(err.Error(), `"magic"`) {
		t.Fatalf("error = %v, want ErrConfig naming bad auth value", err)
	}
}
