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

// R-2J16-KM75
func TestRunShortVersionAliasMatchesLongVersionAndAppearsInHelp(t *testing.T) {
	oldVersion := version
	version = "test-version"
	t.Cleanup(func() { version = oldVersion })

	var shortOut, shortErr bytes.Buffer
	shortCode := run([]string{"-V"}, bytes.NewReader(nil), &shortOut, &shortErr, noEnv, t.TempDir())
	var longOut, longErr bytes.Buffer
	longCode := run([]string{"--version"}, bytes.NewReader(nil), &longOut, &longErr, noEnv, t.TempDir())
	var helpOut, helpErr bytes.Buffer
	helpCode := run([]string{"--help"}, bytes.NewReader(nil), &helpOut, &helpErr, noEnv, t.TempDir())

	if shortCode != 0 || shortOut.String() != longOut.String() || shortOut.String() != "test-version\n" || shortErr.Len() != 0 {
		t.Fatalf("run(-V) = code %d, stdout %q, stderr %q; run(--version) = code %d, stdout %q, stderr %q", shortCode, shortOut.String(), shortErr.String(), longCode, longOut.String(), longErr.String())
	}
	if longCode != 0 || longErr.Len() != 0 {
		t.Fatalf("run(--version) = code %d, stderr %q; want 0 and empty stderr", longCode, longErr.String())
	}
	if helpCode != 0 || !bytes.Contains(helpOut.Bytes(), []byte("-V")) || helpErr.Len() != 0 {
		t.Fatalf("run(--help) = code %d, stdout %q, stderr %q; want 0, -V in help, empty stderr", helpCode, helpOut.String(), helpErr.String())
	}
}
