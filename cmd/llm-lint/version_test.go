package main

import (
	"bytes"
	"testing"
)

// R-FW0C-XQZ3
func TestRunVersionPrintsVersionAndReturnsZero(t *testing.T) {
	var out, errOut bytes.Buffer
	code := run([]string{"--version"}, bytes.NewReader(nil), &out, &errOut, noEnv, t.TempDir())
	if code != 0 || out.String() != "dev\n" {
		t.Fatalf("run() = code %d, stdout %q; want 0, %q", code, out.String(), "dev\\n")
	}
}
