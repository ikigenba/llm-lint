package config

import (
	"errors"
	"strings"
	"testing"
)

func TestNewConversationMissingProviderEnvironment(t *testing.T) {
	// R-G87C-RGE1
	cfg, err := Load(t.TempDir(), nil, func(string) string { return "" })
	if err != nil {
		t.Fatal(err)
	}
	_, err = cfg.NewConversation("system", nil)
	if !errors.Is(err, ErrAuth) || !strings.Contains(err.Error(), "GOOGLE_API_KEY") {
		t.Fatalf("error = %v, want ErrAuth naming GOOGLE_API_KEY", err)
	}
}
