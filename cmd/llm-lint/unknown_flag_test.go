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
