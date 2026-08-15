package engine

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/ikigenba/llm-lint/internal/rules"
)

type budgetClient struct{ calls int }

func (c *budgetClient) ContextWindow() int64 { return 10 }
func (c *budgetClient) Judge(context.Context, rules.Rule, string, []byte) ([]Finding, error) {
	c.calls++
	return nil, nil
}

func TestRunSkipsFileOverContextBudget(t *testing.T) {
	// R-GWLC-EV7X
	client := &budgetClient{}
	var warning bytes.Buffer
	rule := rules.Rule{ID: "no-secrets", Prompt: "find secrets"}
	engine := Engine{Client: client, Concurrency: 1}
	findings, stats, err := engine.Run(context.Background(), []rules.Rule{rule}, map[string][]string{rule.ID: {"large.go"}}, func(string) ([]byte, error) {
		return bytes.Repeat([]byte("x"), 100), nil
	}, &warning)
	if err != nil || len(findings) != 0 {
		t.Fatalf("findings = %v, error = %v; want clean success", findings, err)
	}
	if client.calls != 0 || stats.Calls != 0 {
		t.Fatalf("judge calls = %d, stats calls = %d; want zero", client.calls, stats.Calls)
	}
	if text := warning.String(); !strings.Contains(text, "large.go") || !strings.Contains(text, rule.ID) {
		t.Fatalf("warning %q does not name file and rule", text)
	}
}
