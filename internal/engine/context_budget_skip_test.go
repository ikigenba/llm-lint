package engine

import (
	"bytes"
	"context"
	"testing"

	"github.com/ikigenba/llm-lint/internal/rules"
)

type budgetClient struct{ calls int }

func (c *budgetClient) Judge(context.Context, rules.Rule, string, []byte) ([]Finding, Usage, error) {
	c.calls++
	return nil, Usage{}, nil
}

func TestRunSendsLargeFileToClientWithoutPrecheck(t *testing.T) {
	// R-GWLC-EV7X
	client := &budgetClient{}
	var warning bytes.Buffer
	rule := rules.Rule{ID: "no-secrets", Prompt: "find secrets"}
	runner := Engine{Client: client, Concurrency: 1}
	findings, stats, err := runner.Run(context.Background(), []rules.Rule{rule}, map[string][]string{rule.ID: {"large.go"}}, func(string) ([]byte, error) {
		return bytes.Repeat([]byte("x"), 100), nil
	}, &warning)
	if err != nil || len(findings) != 0 {
		t.Fatalf("findings = %v, error = %v; want clean success", findings, err)
	}
	if client.calls != 1 || stats.Calls != 1 {
		t.Fatalf("judge calls = %d, stats calls = %d; want one", client.calls, stats.Calls)
	}
	if warning.Len() != 0 {
		t.Fatalf("warning = %q, want none", warning.String())
	}
}
