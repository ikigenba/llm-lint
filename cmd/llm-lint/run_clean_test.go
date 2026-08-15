package main

import (
	"bytes"
	"testing"
)

// R-FR4R-EO0B
func TestRunCleanTreeWritesNothingAndReturnsZero(t *testing.T) {
	client := &fakeClient{}
	bindClient(t, client)
	var out, errOut bytes.Buffer
	code := run(nil, bytes.NewReader(nil), &out, &errOut, noEnv, lintTree(t))
	if code != 0 || out.Len() != 0 || client.calls == 0 {
		t.Fatalf("run() = code %d, stdout %q, calls %d; want 0, empty, positive", code, out.String(), client.calls)
	}
}
