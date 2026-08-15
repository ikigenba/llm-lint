package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

// R-FYG5-PAGH
func TestRunOperationalFailurePrintsErrorAndReturnsThree(t *testing.T) {
	bindClient(t, &fakeClient{err: errors.New("provider unavailable after retries")})
	var out, errOut bytes.Buffer
	code := run(nil, bytes.NewReader(nil), &out, &errOut, noEnv, lintTree(t))
	if code != 3 || !strings.Contains(errOut.String(), "provider unavailable after retries") {
		t.Fatalf("run() = code %d, stderr %q; want 3 and operational failure", code, errOut.String())
	}
}
