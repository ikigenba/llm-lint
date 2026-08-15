package suppress

import (
	"reflect"
	"testing"

	"github.com/ikigenba/llm-lint/internal/engine"
)

func TestDifferentRuleIDLeavesFindingIntact(t *testing.T) {
	// R-HCG1-DVUY
	finding := engine.Finding{Rule: "wanted-rule", File: "example.go", Line: 1}
	got, err := Filter([]engine.Finding{finding}, readerWith("work() // llm-lint:ignore other-rule\n"))
	if err != nil {
		t.Fatalf("Filter() error = %v", err)
	}
	if !reflect.DeepEqual(got, []engine.Finding{finding}) {
		t.Fatalf("Filter() = %#v, want original finding", got)
	}
}

func readerWith(contents string) func(string) ([]byte, error) {
	return func(string) ([]byte, error) {
		return []byte(contents), nil
	}
}
