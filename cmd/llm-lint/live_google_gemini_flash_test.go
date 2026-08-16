package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// R-GZ15-6EPB
func TestRunLiveGoogleGeminiFlashReportsSleepPastRace(t *testing.T) {
	if os.Getenv("GEMINI_API_KEY") == "" {
		t.Skip("GEMINI_API_KEY is not set")
	}

	cwd := filepath.Join("..", "..", "testdata", "live-smoke")
	var out, errOut bytes.Buffer
	code := run([]string{"-c", "model=gemini-3.7-flash", "--no-cache", "--stats"}, bytes.NewReader(nil), &out, &errOut, os.Getenv, cwd)
	if code != 1 {
		t.Fatalf("run() = %d, want 1; stdout %q; stderr %q", code, out.String(), errOut.String())
	}
	for _, want := range []string{"sleep_past_race_test.go:", "(no-sleep-in-tests)"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout %q does not contain %q; stderr %q", out.String(), want, errOut.String())
		}
	}
	stats := errOut.String()
	if !strings.Contains(stats, "llm-lint: 1 rules, 1 files, 1 pairs, 1 calls") || strings.Contains(stats, ", 0 in /") {
		t.Fatalf("stderr %q does not contain stats with non-zero input tokens", stats)
	}
}
