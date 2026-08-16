package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunPrintsEachConfigWarningOnceBeforeLinting(t *testing.T) {
	client := &fakeClient{}
	bindClient(t, client)
	getenv := func(key string) string {
		if key == "GEMINI_API_KEY" {
			return "test-key"
		}
		return ""
	}

	var out, errOut bytes.Buffer
	code := run([]string{"-c", "provider=google", "-c", "model=private-google-model"}, bytes.NewReader(nil), &out, &errOut, getenv, lintTree(t))
	if code != 0 {
		t.Fatalf("run() = %d, want 0; stdout %q; stderr %q", code, out.String(), errOut.String())
	}
	if got := strings.Count(errOut.String(), "has no pricing"); got != 1 {
		t.Fatalf("stderr contains pricing warning %d times, want once: %q", got, errOut.String())
	}
}
