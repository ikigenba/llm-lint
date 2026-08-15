package report

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/ikigenba/llm-lint/internal/engine"
)

// R-HHBM-WYTQ
func TestFormatsSortFindingsByFileLineAndRule(t *testing.T) {
	findings := []engine.Finding{
		{Rule: "z-rule", File: "b.go", Line: 1, Explanation: "fourth"},
		{Rule: "z-rule", File: "a.go", Line: 2, Explanation: "third"},
		{Rule: "z-rule", File: "a.go", Line: 1, Explanation: "second"},
		{Rule: "a-rule", File: "a.go", Line: 1, Explanation: "first"},
	}
	var textOut bytes.Buffer
	if err := Text(&textOut, "/root", "/root", findings); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(textOut.String()), "\n")
	wantText := []string{"a.go:1: first (a-rule)", "a.go:1: second (z-rule)", "a.go:2: third (z-rule)", "b.go:1: fourth (z-rule)"}
	if !reflect.DeepEqual(lines, wantText) {
		t.Fatalf("Text order = %q, want %q", lines, wantText)
	}

	var jsonOut bytes.Buffer
	if err := JSON(&jsonOut, findings); err != nil {
		t.Fatal(err)
	}
	var objects []struct {
		Rule string `json:"rule"`
		File string `json:"file"`
		Line int    `json:"line"`
	}
	if err := json.Unmarshal(jsonOut.Bytes(), &objects); err != nil {
		t.Fatal(err)
	}
	gotJSON := make([]string, 0, len(objects))
	for _, object := range objects {
		gotJSON = append(gotJSON, fmt.Sprintf("%s:%s:%d", object.File, object.Rule, object.Line))
	}
	wantJSON := []string{"a.go:a-rule:1", "a.go:z-rule:1", "a.go:z-rule:2", "b.go:z-rule:1"}
	if !reflect.DeepEqual(gotJSON, wantJSON) {
		t.Fatalf("JSON order = %q, want %q", gotJSON, wantJSON)
	}
}
