package config

import "testing"

func TestDefaultAuthFollowsProviderTable(t *testing.T) {
	// R-J364-93NV
	tests := []struct {
		name     string
		pairs    []string
		wantAuth string
	}{
		{name: "openai", pairs: []string{"provider=openai", "model=private-openai-model"}, wantAuth: "sub"},
		{name: "google", pairs: nil, wantAuth: "key"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := Load(t.TempDir(), tt.pairs, func(string) string { return "" })
			if err != nil {
				t.Fatal(err)
			}
			if cfg.Model.Auth != tt.wantAuth {
				t.Fatalf("Auth = %q, want %q", cfg.Model.Auth, tt.wantAuth)
			}
		})
	}
}
