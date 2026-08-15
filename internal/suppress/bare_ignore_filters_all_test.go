package suppress

import (
	"testing"

	"github.com/ikigenba/llm-lint/internal/engine"
)

func TestBareIgnoreFiltersAll(t *testing.T) {
	// R-HA08-MCDK
	findings := []engine.Finding{
		{Rule: "first-rule", File: "example.go", Line: 1},
		{Rule: "second-rule", File: "example.go", Line: 1},
	}
	got, err := Filter(findings, readerWith("work() // llm-lint:ignore\n"))
	if err != nil {
		t.Fatalf("Filter() error = %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("Filter() = %#v, want no findings", got)
	}
}
