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
	if !errors.Is(err, ErrAuth) || !strings.Contains(err.Error(), "GEMINI_API_KEY") {
		t.Fatalf("error = %v, want ErrAuth naming GEMINI_API_KEY", err)
	}
}
