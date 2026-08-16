package config

import (
	"errors"
	"strings"
	"testing"
)

func TestLoadRejectsSubscriptionAuthForKeyOnlyProvider(t *testing.T) {
	// R-J0QB-HK6H
	_, err := Load(t.TempDir(), []string{"provider=google", "auth=sub"}, func(string) string { return "" })
	if !errors.Is(err, ErrConfig) || !strings.Contains(err.Error(), "google") || !strings.Contains(err.Error(), "key-only") {
		t.Fatalf("error = %v, want ErrConfig naming google as key-only", err)
	}
}
