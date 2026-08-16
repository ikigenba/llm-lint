package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestRunVerboseJSONKeepsOutputStreamsSeparate(t *testing.T) {
	// R-H2SH-P46N
	bindClient(t, &perPairClient{})
	root := lintTree(t)
	var out, errOut bytes.Buffer
	code := run([]string{"--verbose", "--format", "json"}, bytes.NewReader(nil), &out, &errOut, noEnv, root)
	if code != 1 {
		t.Fatalf("run() code = %d, stderr = %q", code, errOut.String())
	}
	var findings []map[string]any
	if err := json.Unmarshal(out.Bytes(), &findings); err != nil || len(findings) != 1 {
		t.Fatalf("stdout = %q, JSON error = %v; want one-finding array", out.String(), err)
	}
	wantTrace := "sample_test.go: no-sleep-in-tests miss fail\n"
	if errOut.String() != wantTrace {
		t.Fatalf("stderr = %q, want %q", errOut.String(), wantTrace)
	}
	if strings.Contains(out.String(), " miss fail") || strings.Contains(errOut.String(), "\"rule\"") {
		t.Fatalf("streams intermingled: stdout = %q, stderr = %q", out.String(), errOut.String())
	}
}
