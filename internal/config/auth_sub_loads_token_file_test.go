package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ikigenba/agentkit"
)

func TestSubscriptionAuthLoadsDefaultAndOverrideTokenFiles(t *testing.T) {
	// R-J4E0-MVEK
	home := t.TempDir()
	defaultPath := filepath.Join(home, ".llm-lint", "openai-auth.json")
	overridePath := filepath.Join(t.TempDir(), "custom-auth.json")
	writeOpenAISubscriptionToken(t, defaultPath)
	writeOpenAISubscriptionToken(t, overridePath)

	for _, tt := range []struct {
		name  string
		pairs []string
	}{
		{name: "default path", pairs: nil},
		{name: "auth_file override", pairs: []string{"auth_file=" + overridePath}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			pairs := append([]string{"provider=openai", "model=private-openai-model", "auth=sub"}, tt.pairs...)
			cfg, err := Load(t.TempDir(), pairs, func(key string) string {
				if key == "HOME" {
					return home
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
			identity := conversation.Provider.Identity()
			if identity.Provider != agentkit.ProviderOpenAI || identity.Auth != agentkit.AuthSubscription {
				t.Fatalf("identity = %s, want openai subscription", identity)
			}
		})
	}
}

func writeOpenAISubscriptionToken(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"https://api.openai.com/auth":{"chatgpt_account_id":"acct"}}`))
	data := []byte(fmt.Sprintf(`{"access_token":"h.%s.s"}`, payload))
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
}
