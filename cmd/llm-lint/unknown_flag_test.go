package main

import (
	"bytes"
	"strings"
	"testing"
)

// R-FUSG-JZ8E
func TestRunUnknownFlagPrintsUsageToStderrAndReturnsTwo(t *testing.T) {
	for _, args := range [][]string{{"--not-a-flag"}, {"--model", "gemini-3.7-flash"}} {
		var out, errOut bytes.Buffer
		code := run(args, bytes.NewReader(nil), &out, &errOut, noEnv, t.TempDir())
		if code != 2 || !strings.Contains(errOut.String(), "Usage: llm-lint") {
			t.Errorf("run(%q) = code %d, stderr %q; want 2 and usage", args, code, errOut.String())
		}
	}
}

// R-WBQM-PELO
func TestRunUnknownFlagPrintsUsageExactlyOnceAfterFlagError(t *testing.T) {
	var out, errOut bytes.Buffer
	code := run([]string{"--not-a-flag"}, bytes.NewReader(nil), &out, &errOut, noEnv, t.TempDir())

	const flagError = "flag provided but not defined: -not-a-flag"
	const usageLine = "Usage: llm-lint"
	stderr := errOut.String()
	if code != 2 {
		t.Errorf("run() code = %d; want 2", code)
	}
	if count := strings.Count(stderr, usageLine); count != 1 {
		t.Errorf("stderr contains %q %d times; want exactly once; stderr %q", usageLine, count, stderr)
	}
	errorAt, usageAt := strings.Index(stderr, flagError), strings.Index(stderr, usageLine)
	if errorAt < 0 || usageAt < 0 || errorAt >= usageAt {
		t.Errorf("stderr = %q; want flag error before usage", stderr)
	}
}
