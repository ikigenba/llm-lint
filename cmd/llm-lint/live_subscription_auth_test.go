package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ikigenba/llm-lint/internal/config"
)

// R-J81P-S6MN
func TestRunLiveSubscriptionReportsSleepPastRace(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skipf("cannot locate the default subscription token: %v", err)
	}
	tokenFile := filepath.Join(home, ".llm-lint", "openai-auth.json")
	if _, err := os.Stat(tokenFile); err != nil {
		t.Skipf("default OpenAI subscription token is unavailable: %v", err)
	}

	cwd := filepath.Join("..", "..", "testdata", "live-smoke")
	pairs := []string{"provider=openai", "model=gpt-5.6-luna", "auth=sub"}
	if _, err := config.Load(cwd, pairs, os.Getenv); err != nil {
		t.Skipf("default OpenAI subscription token does not load: %v", err)
	}

	var out, errOut bytes.Buffer
	code := run([]string{"-c", pairs[0], "-c", pairs[1], "-c", pairs[2], "--no-cache"}, bytes.NewReader(nil), &out, &errOut, os.Getenv, cwd)
	if code != 1 {
		t.Fatalf("run() = %d, want 1; stdout %q; stderr %q", code, out.String(), errOut.String())
	}
	for _, want := range []string{"sleep_past_race_test.go:", "(no-sleep-in-tests)"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout %q does not contain %q; stderr %q", out.String(), want, errOut.String())
		}
	}
}
