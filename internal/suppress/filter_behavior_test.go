package suppress

import (
	"errors"
	"strings"
	"testing"

	"github.com/ikigenba/llm-lint/internal/engine"
)

func TestFilterWrapsReadError(t *testing.T) {
	readErr := errors.New("permission denied")
	_, err := Filter([]engine.Finding{{Rule: "rule", File: "unreadable.go", Line: 1}}, func(string) ([]byte, error) {
		return nil, readErr
	})
	if !errors.Is(err, readErr) {
		t.Fatalf("Filter() error = %v, want wrapped read error", err)
	}
	if !strings.Contains(err.Error(), "suppress: read \"unreadable.go\"") {
		t.Fatalf("Filter() error = %q, want package and path context", err)
	}
}

func TestEmptyListItemIsIgnoredAsText(t *testing.T) {
	findings := []engine.Finding{{Rule: "first-rule", File: "example.go", Line: 1}}
	got, err := Filter(findings, readerWith("work() // llm-lint:ignore first-rule,,second-rule\n"))
	if err != nil {
		t.Fatalf("Filter() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("Filter() = %#v, want malformed marker ignored", got)
	}
}
