package config

import (
	"strings"
	"testing"
)

func TestExplicitProviderAllowsPassthruWithoutPricing(t *testing.T) {
	// R-IYAI-Q0P3
	const model = "provider-private-model-12345"
	cfg, err := Load(t.TempDir(), []string{"provider=google", "model=" + model}, func(key string) string {
		if key == "GEMINI_API_KEY" {
			return "secret"
		}
		return ""
	})
	if err != nil {
		t.Fatal(err)
	}
	conversation, err := cfg.NewConversation("system", nil)
	if err != nil {
		t.Fatal(err)
	}
	if conversation.Model != model || conversation.Pricing != nil {
		t.Fatalf("model = %q pricing = %#v, want raw model with nil pricing", conversation.Model, conversation.Pricing)
	}
	if len(cfg.Warnings) != 1 || !strings.Contains(cfg.Warnings[0], "no pricing (cost reports 0), reasoning unchecked") {
		t.Fatalf("Warnings = %#v, want no-pricing warning", cfg.Warnings)
	}
}
