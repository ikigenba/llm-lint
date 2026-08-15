package main

import (
	"bytes"
	"strings"
	"testing"
)

// R-FUSG-JZ8E
func TestRunUnknownFlagPrintsUsageToStderrAndReturnsTwo(t *testing.T) {
	var out, errOut bytes.Buffer
	code := run([]string{"--not-a-flag"}, bytes.NewReader(nil), &out, &errOut, noEnv, t.TempDir())
	if code != 2 || !strings.Contains(errOut.String(), "Usage: llm-lint") {
		t.Fatalf("run() = code %d, stderr %q; want 2 and usage", code, errOut.String())
	}
}
