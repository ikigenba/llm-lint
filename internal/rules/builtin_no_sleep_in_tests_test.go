package rules

import (
	"reflect"
	"strings"
	"testing"
)

// R-HIJJ-AQKF
func TestBuiltInsContainsNoSleepInTestsCatalogRule(t *testing.T) {
	var rule Rule
	for _, candidate := range BuiltIns() {
		if candidate.ID == "no-sleep-in-tests" {
			rule = candidate
			break
		}
	}
	wantGlobs := []string{
		"**/*_test.go",
		"**/*_test.py",
		"**/test_*.py",
		"**/*.test.js",
		"**/*.test.ts",
		"**/*.spec.js",
		"**/*.spec.ts",
	}
	if rule.ID != "no-sleep-in-tests" || rule.Severity != SeverityError || !rule.BuiltIn || !reflect.DeepEqual(rule.Include, wantGlobs) {
		t.Fatalf("catalog rule = %#v, want no-sleep-in-tests error rule with globs %v", rule, wantGlobs)
	}
	if len(rule.Exclude) != 0 || !strings.Contains(rule.Prompt, "time.Sleep") || !strings.Contains(rule.Prompt, "deadline") {
		t.Fatalf("built-in exclusions/prompt = %v, %q; want no exclusions and anchored guidance", rule.Exclude, rule.Prompt)
	}
}
