package walk

import (
	"reflect"
	"testing"

	"github.com/ikigenba/llm-lint/internal/rules"
)

// R-GP9Y-48RR
func TestCandidatesDoublestarFilter(t *testing.T) {
	files := []string{"main.go", "internal/a.go", "internal/deep/b.go", "internal/generated/c.go", "README.md"}
	rule := rules.Rule{Include: []string{"**/*.go"}, Exclude: []string{"**/generated/**"}}

	got := Candidates(files, rule)
	want := []string{"main.go", "internal/a.go", "internal/deep/b.go"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Candidates() = %v, want %v", got, want)
	}
}
