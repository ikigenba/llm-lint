package main

import (
	"bytes"
	"testing"
)

// R-HEVU-5FCC
func TestRunJSONCleanTreeWritesEmptyArray(t *testing.T) {
	bindClient(t, &fakeClient{})
	var out, errOut bytes.Buffer
	code := run([]string{"--format", "json"}, bytes.NewReader(nil), &out, &errOut, noEnv, lintTree(t))
	if code != 0 || out.String() != "[]\n" {
		t.Fatalf("run() = code %d, stdout %q, stderr %q; want 0 and empty JSON array", code, out.String(), errOut.String())
	}
}
