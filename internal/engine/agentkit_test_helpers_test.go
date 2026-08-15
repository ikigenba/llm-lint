package engine

import (
	"context"
	"sync"

	"github.com/ikigenba/agentkit"
)

type scriptedProvider struct {
	mu        sync.Mutex
	responses []agentkit.Message
	calls     int
}

func (p *scriptedProvider) RoundTrip(_ context.Context, _ *agentkit.Request) *agentkit.RoundTrip {
	p.mu.Lock()
	defer p.mu.Unlock()
	index := p.calls
	p.calls++
	if index >= len(p.responses) {
		index = len(p.responses) - 1
	}
	finish := agentkit.FinishStop
	if len(p.responses[index].Blocks) > 0 {
		if _, ok := p.responses[index].Blocks[0].(agentkit.ToolUseBlock); ok {
			finish = agentkit.FinishToolUse
		}
	}
	return agentkit.NewRoundTrip(p.responses[index], finish, agentkit.Usage{}, nil, nil, 0, false)
}

func (p *scriptedProvider) Identity() agentkit.Identity { return agentkit.Identity{} }

func assistantText(text string) agentkit.Message {
	return agentkit.Message{Role: agentkit.RoleAssistant, Blocks: []agentkit.Block{agentkit.TextBlock{Text: text}}}
}
