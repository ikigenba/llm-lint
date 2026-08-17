package main

import (
	"bytes"
	"strings"
	"testing"
)

// R-4CGH-YUDC
func TestRunHelpRendersReasoningChoicesRangesDynamicAndNone(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := run([]string{"--help"}, bytes.NewReader(nil), &out, &errOut, noEnv, t.TempDir()); code != 0 || errOut.Len() != 0 {
		t.Fatalf("run(--help) = code %d, stderr %q; want 0 and empty stderr", code, errOut.String())
	}
	help := out.String()
	wants := map[string]string{
		"gemini-3.7-flash":     "thinking_level={low|*medium|high}",
		"gemini-2.5-flash":     "thinking_budget={off|0–24576|*dynamic}",
		"gemini-2.5-pro":       "thinking_budget={128–32768|*dynamic}",
		"gemini-embedding-001": "—",
	}
	for model, clause := range wants {
		var row string
		for _, line := range strings.Split(help, "\n") {
			fields := strings.Fields(line)
			if len(fields) == 2 && fields[0] == model {
				row = line
				break
			}
		}
		if fields := strings.Fields(row); len(fields) != 2 || fields[0] != model || fields[1] != clause {
			t.Errorf("row for %s = %q; want aligned model and clause %q", model, row, clause)
		}
	}
	for model, clause := range wants {
		if model == "gemini-embedding-001" {
			continue
		}
		if count := strings.Count(clause, "*"); count != 1 {
			t.Errorf("%s clause %q marks %d defaults; want exactly one", model, clause, count)
		}
	}
}
