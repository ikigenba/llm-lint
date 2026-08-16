package config

import (
	"testing"

	"github.com/ikigenba/agentkit"
)

func TestExplicitProviderOverridesCatalogDerivation(t *testing.T) {
	// R-IVUP-YH7P
	cfg, err := Load(t.TempDir(), []string{"model=gemini-3.7-flash", "provider=openai"}, func(string) string { return "" })
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Model.Provider != string(agentkit.ProviderOpenAI) || !cfg.Model.ProviderExplicit {
		t.Fatalf("provider = %q explicit = %t, want openai explicit", cfg.Model.Provider, cfg.Model.ProviderExplicit)
	}
}
