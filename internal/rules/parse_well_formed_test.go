package rules

import (
	"reflect"
	"testing"
)

// R-GAN5-IZVF
func TestParseWellFormedRulePreservesAllFieldsAndBody(t *testing.T) {
	prompt := "Check the code exactly.\n\nKeep this final newline.\n"
	src := testRuleSource("find risky calls", "warning", `["**/*.go", "src/*.ts"]`, `["vendor/**"]`, prompt)

	got, err := Parse("path/risky-calls.md", src)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	want := Rule{
		ID:          "risky-calls",
		Description: "find risky calls",
		Severity:    SeverityWarning,
		Include:     []string{"**/*.go", "src/*.ts"},
		Exclude:     []string{"vendor/**"},
		Prompt:      prompt,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Parse() = %#v, want %#v", got, want)
	}
}
