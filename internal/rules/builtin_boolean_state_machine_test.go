package rules

import (
	"reflect"
	"strings"
	"testing"
)

// R-4IQ0-3CA4
func TestBuiltInsContainsBooleanStateMachineWithDefaultCodeGlobs(t *testing.T) {
	var rule Rule
	for _, candidate := range BuiltIns() {
		if candidate.ID == "boolean-state-machine" {
			rule = candidate
			break
		}
	}
	if rule.ID != "boolean-state-machine" || rule.Severity != SeverityError || !rule.BuiltIn {
		t.Fatalf("catalog rule = %#v, want built-in boolean-state-machine error rule", rule)
	}
	if !reflect.DeepEqual(rule.Include, DefaultCodeGlobs) {
		t.Fatalf("boolean-state-machine Include = %v, want DefaultCodeGlobs %v", rule.Include, DefaultCodeGlobs)
	}
	if len(rule.Exclude) != 0 {
		t.Fatalf("boolean-state-machine Exclude = %v, want none", rule.Exclude)
	}
	for _, anchor := range []string{"invalid or meaningless", "held", "sold"} {
		if !strings.Contains(rule.Prompt, anchor) {
			t.Fatalf("boolean-state-machine prompt does not contain %q: %q", anchor, rule.Prompt)
		}
	}
}
