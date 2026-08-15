package suppress

import (
	"reflect"
	"testing"

	"github.com/ikigenba/llm-lint/internal/engine"
)

func TestCommentOnlyMarkerFiltersOnlyNextLine(t *testing.T) {
	// R-HB85-0449
	findings := []engine.Finding{
		{Rule: "target-rule", File: "example.go", Line: 1},
		{Rule: "target-rule", File: "example.go", Line: 2},
		{Rule: "target-rule", File: "example.go", Line: 3},
	}
	want := []engine.Finding{findings[0], findings[2]}
	got, err := Filter(findings, readerWith("// llm-lint:ignore target-rule\nwork()\nmoreWork()\n"))
	if err != nil {
		t.Fatalf("Filter() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Filter() = %#v, want %#v", got, want)
	}
}
