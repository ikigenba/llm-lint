package report

import (
	"bytes"
	"encoding/json"
	"reflect"
	"sort"
	"testing"

	"github.com/ikigenba/llm-lint/internal/engine"
	"github.com/ikigenba/llm-lint/internal/rules"
)

// R-HEVU-5FCC
func TestJSONHasExactShapeEvidenceAndEmptyArray(t *testing.T) {
	finding := engine.Finding{Rule: "clarity", Severity: rules.SeverityWarning, File: "sample.go", Line: 4, Evidence: "\treturn value", Explanation: "explain the value"}
	var out bytes.Buffer
	if err := JSON(&out, []engine.Finding{finding}); err != nil {
		t.Fatal(err)
	}
	var got []map[string]any
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatalf("JSON emitted invalid JSON %q: %v", out.String(), err)
	}
	wantKeys := []string{"evidence", "explanation", "file", "line", "rule", "severity"}
	keys := make([]string, 0, len(got[0]))
	for key := range got[0] {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	if !reflect.DeepEqual(keys, wantKeys) || got[0]["evidence"] != finding.Evidence {
		t.Fatalf("JSON object = %#v; want keys %v and evidence %q", got[0], wantKeys, finding.Evidence)
	}

	out.Reset()
	if err := JSON(&out, nil); err != nil {
		t.Fatal(err)
	}
	if out.String() != "[]\n" {
		t.Fatalf("JSON(nil) = %q, want empty array", out.String())
	}
}
