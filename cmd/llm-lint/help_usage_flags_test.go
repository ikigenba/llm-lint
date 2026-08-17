package main

import (
	"bytes"
	"strings"
	"testing"
)

// R-4A0P-7AVY
func TestRunHelpPrintsUsageAndAlignedFlagCatalog(t *testing.T) {
	var out, errOut bytes.Buffer
	code := run([]string{"--help"}, bytes.NewReader(nil), &out, &errOut, noEnv, t.TempDir())
	if code != 0 || errOut.Len() != 0 {
		t.Fatalf("run(--help) = code %d, stderr %q; want 0 and empty stderr", code, errOut.String())
	}

	lines := strings.Split(out.String(), "\n")
	if len(lines) < 3 || lines[0] != "usage: llm-lint [flags] [path ...]" || lines[1] != "flags:" {
		t.Fatalf("help starts %q; want exact usage line followed by flags block", strings.Join(lines[:min(3, len(lines))], "\n"))
	}
	flagLines := lines[2:12]
	for _, want := range []string{"-c,", "--rules", "--format", "--concurrency", "--no-cache", "--stats", "--verbose", "--list-rules", "-V, --version", "--help"} {
		found := false
		for _, line := range flagLines {
			if strings.Contains(strings.Join(strings.Fields(line), " "), want) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("flags block does not contain %q:\n%s", want, strings.Join(flagLines, "\n"))
		}
	}
	descriptionColumn := strings.Index(flagLines[0], "override")
	if descriptionColumn < 0 {
		t.Fatalf("first flag row has no description: %q", flagLines[0])
	}
	for _, line := range flagLines[1:] {
		if len(line) <= descriptionColumn || line[descriptionColumn] == ' ' {
			t.Errorf("flag description is not aligned at column %d: %q", descriptionColumn, line)
		}
	}
}
