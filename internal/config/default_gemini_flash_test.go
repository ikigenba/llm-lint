package config

import (
	"testing"

	"github.com/ikigenba/agentkit"
)

func TestDefaultModelResolvesToGoogle(t *testing.T) {
	// R-G6ZG-DONC
	cfg, err := Load(t.TempDir(), nil, func(string) string { return "secret" })
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Model.ModelID != "gemini-3.7-flash" || cfg.Model.Provider != string(agentkit.ProviderGoogle) {
		t.Fatalf("model = %q provider = %q", cfg.Model.ModelID, cfg.Model.Provider)
	}
}
