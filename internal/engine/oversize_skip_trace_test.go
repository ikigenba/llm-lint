package engine

import (
	"context"
	"io"
	"testing"

	"github.com/ikigenba/llm-lint/internal/rules"
)

func TestLargeFileDoesNotCreateSyntheticSkipTrace(t *testing.T) {
	// R-H40E-2VXC
	rule := rules.Rule{ID: "large-rule", Prompt: "inspect this file"}
	client := &budgetClient{}
	runner := Engine{Client: client, Concurrency: 1}
	_, stats, err := runner.Run(context.Background(), []rules.Rule{rule}, map[string][]string{rule.ID: {"large.go"}}, func(string) ([]byte, error) {
		return make([]byte, 100), nil
	}, io.Discard)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if client.calls != 1 || stats.Calls != 1 {
		t.Fatalf("calls = %d, stats calls = %d; want one ordinary call with no synthetic skip path", client.calls, stats.Calls)
	}
}
