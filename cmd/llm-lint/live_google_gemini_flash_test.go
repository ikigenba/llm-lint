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
	if os.Getenv("GOOGLE_API_KEY") == "" {
		t.Skip("GOOGLE_API_KEY is not set")
	}

	cwd := filepath.Join("..", "..", "testdata", "live-smoke")
	var out, errOut bytes.Buffer
	code := run([]string{"-c", "model=gemini-3.7-flash", "--no-cache"}, bytes.NewReader(nil), &out, &errOut, os.Getenv, cwd)
	if code != 1 {
		t.Fatalf("run() = %d, want 1; stdout %q; stderr %q", code, out.String(), errOut.String())
	}
	for _, want := range []string{"sleep_past_race_test.go:", "(no-sleep-in-tests)"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout %q does not contain %q; stderr %q", out.String(), want, errOut.String())
		}
	}
}
