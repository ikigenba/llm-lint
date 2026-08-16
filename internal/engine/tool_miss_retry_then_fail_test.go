package engine

import (
	"io"
	"strings"
	"testing"

	"github.com/ikigenba/agentkit"
	"github.com/ikigenba/llm-lint/internal/rules"
)

func TestAgentkitClientRetriesToolMissOnceThenNamesPair(t *testing.T) {
	// R-GU5J-NBQJ
	provider := &scriptedProvider{responses: []agentkit.Message{assistantText("no tool"), assistantText("still no tool")}}
	client := AgentkitClient{NewConversation: func(system string, log io.Writer) (*agentkit.Conversation, error) {
		return &agentkit.Conversation{Provider: provider, Model: "fake", System: system, Log: log}, nil
	}}
	_, _, err := client.Judge(t.Context(), rules.Rule{ID: "rule-x"}, "file.go", []byte("package p"))
	if err == nil || !strings.Contains(err.Error(), "rule-x") || !strings.Contains(err.Error(), "file.go") {
		t.Fatalf("error = %v, want operational detail naming rule and file", err)
	}
	if provider.calls != 2 {
		t.Fatalf("provider calls = %d, want initial call plus one retry", provider.calls)
	}
}
