package config

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"
)

func TestSubscriptionAuthMissingTokenFileNamesExactPath(t *testing.T) {
	// R-J5LX-0N59
	path := filepath.Join(t.TempDir(), "missing-token.json")
	cfg, err := Load(t.TempDir(), []string{"provider=openai", "model=private-openai-model", "auth=sub", "auth_file=" + path}, func(string) string { return "" })
	if err != nil {
		t.Fatal(err)
	}
	_, err = cfg.NewConversation("system", nil)
	if !errors.Is(err, ErrAuth) || !strings.Contains(err.Error(), path) {
		t.Fatalf("error = %v, want ErrAuth naming %q", err, path)
	}
}
