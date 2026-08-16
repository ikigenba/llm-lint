package main

import (
	"bytes"
	"fmt"
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

	cwd := filepath.Join("..", "..", "testdata", "live-smoke")
	candidates := []struct {
		provider string
		model    string
	}{
		{provider: "openai", model: "gpt-5.6-luna"},
		{provider: "x-ai", model: "grok-4.3"},
	}
	var pairs []string
	var unavailable []string
	for _, candidate := range candidates {
		tokenFile := filepath.Join(home, ".llm-lint", candidate.provider+"-auth.json")
		if _, err := os.Stat(tokenFile); err != nil {
			unavailable = append(unavailable, fmt.Sprintf("%s: %v", candidate.provider, err))
			continue
		}
		candidatePairs := []string{"provider=" + candidate.provider, "model=" + candidate.model, "auth=sub"}
		if _, err := config.Load(cwd, candidatePairs, os.Getenv); err != nil {
			unavailable = append(unavailable, fmt.Sprintf("%s: %v", candidate.provider, err))
			continue
		}
		pairs = candidatePairs
		break
	}
	if pairs == nil {
		t.Skipf("no default subscription token is available and loadable: %s", strings.Join(unavailable, "; "))
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
