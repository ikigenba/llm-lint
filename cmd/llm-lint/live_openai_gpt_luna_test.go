package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// R-H091-K6G0
func TestRunLiveOpenAIGPTLunaReportsSleepPastRace(t *testing.T) {
	if os.Getenv("OPENAI_API_KEY") == "" {
		t.Skip("OPENAI_API_KEY is not set")
	}

	cwd := filepath.Join("..", "..", "testdata", "live-smoke")
	var out, errOut bytes.Buffer
	code := run([]string{"-c", "model=gpt-5.6-luna", "--no-cache"}, bytes.NewReader(nil), &out, &errOut, os.Getenv, cwd)
	if code != 1 {
		t.Fatalf("run() = %d, want 1; stdout %q; stderr %q", code, out.String(), errOut.String())
	}
	for _, want := range []string{"sleep_past_race_test.go:", "(no-sleep-in-tests)"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout %q does not contain %q; stderr %q", out.String(), want, errOut.String())
		}
	}
}
