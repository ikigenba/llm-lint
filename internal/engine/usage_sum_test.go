package engine

import (
	"context"
	"io"
	"testing"

	"github.com/ikigenba/llm-lint/internal/rules"
)

func TestRunSumsUsageAcrossJudgedPairs(t *testing.T) {
	// R-ETNZ-VY1Q
	perCall := Usage{Input: 17, Output: 5, CostUSD: 0.0125}
	client := judgeFunc(func(context.Context, rules.Rule, string, []byte) ([]Finding, Usage, error) {
		return nil, perCall, nil
	})
	ruleSet := []rules.Rule{{ID: "alpha"}, {ID: "beta"}}
	files := map[string][]string{
		"alpha": {"one.go", "two.go"},
		"beta":  {"one.go", "two.go"},
	}

	_, stats, err := (&Engine{Client: client, Concurrency: 2}).Run(
		context.Background(), ruleSet, files,
		func(string) ([]byte, error) { return []byte("source"), nil }, io.Discard,
	)
	if err != nil {
		t.Fatal(err)
	}
	const pairs = 4
	if stats.InputTokens != pairs*perCall.Input || stats.OutputTokens != pairs*perCall.Output || stats.CostUSD != pairs*perCall.CostUSD {
		t.Fatalf("usage = (%d in, %d out, $%.4f), want (%d in, %d out, $%.4f)",
			stats.InputTokens, stats.OutputTokens, stats.CostUSD,
			pairs*perCall.Input, pairs*perCall.Output, pairs*perCall.CostUSD)
	}
}
