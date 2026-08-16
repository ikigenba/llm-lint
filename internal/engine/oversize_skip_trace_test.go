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
	trace := &Trace{}
	client := &budgetClient{}
	runner := Engine{Client: client, Concurrency: 1, Trace: trace}
	_, stats, err := runner.Run(context.Background(), []rules.Rule{rule}, map[string][]string{rule.ID: {"large.go"}}, func(string) ([]byte, error) {
		return make([]byte, 100), nil
	}, io.Discard)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if client.calls != 1 || stats.Calls != 1 || len(trace.Entries()) != 0 {
		t.Fatalf("calls = %d, stats calls = %d, trace = %#v; want one call and no synthetic trace", client.calls, stats.Calls, trace.Entries())
	}
}
