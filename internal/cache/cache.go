package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"sync/atomic"
	"time"

	"github.com/ikigenba/llm-lint/internal/engine"
	"github.com/ikigenba/llm-lint/internal/rules"
)

type Store struct {
	Dir string
	Now func() time.Time
}

type CachingClient struct {
	Store  *Store
	Next   engine.Client
	NoRead bool
	Hits   *atomic.Int64
}

func (c *CachingClient) Judge(ctx context.Context, r rules.Rule, file string, content []byte) ([]engine.Finding, error) {
	return c.Next.Judge(ctx, r, file, content)
}

func Key(model string, rulePromptAndMeta, fileContent []byte) string {
	h := sha256.New()
	h.Write([]byte(model))
	h.Write(rulePromptAndMeta)
	h.Write(fileContent)
	return hex.EncodeToString(h.Sum(nil))
}
