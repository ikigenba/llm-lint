package suppress

import (
	"testing"

	"github.com/ikigenba/llm-lint/internal/engine"
)

func TestSameLineNamedRuleFilters(t *testing.T) {
	// R-H8SC-8KMV
	finding := engine.Finding{Rule: "no-sleep-in-tests", File: "example.go", Line: 1}
	got, err := Filter([]engine.Finding{finding}, readerWith("time.Sleep(time.Second) // llm-lint:ignore no-sleep-in-tests\n"))
	if err != nil {
		t.Fatalf("Filter() error = %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("Filter() = %#v, want no findings", got)
	}
}
