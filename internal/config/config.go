package config

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/ikigenba/agentkit"
)

type Config struct {
	Root    string
	Enable  []string
	Rules   []string
	Exclude []string
	Model   Model
}

type Model struct {
	Provider    string
	ModelID     string
	BaseURL     string
	Temperature *float64
	TopP        *float64
	MaxTokens   int
	Reasoning   agentkit.ReasoningValue
	Retry       agentkit.RetryPolicy
}

var (
	ErrConfig = errors.New("config: invalid configuration")
	ErrAuth   = errors.New("config: authentication failed")
)

func Load(cwd string, cliPairs []string, getenv func(string) string) (*Config, error) {
	c := &Config{Root: cwd, Model: Model{ModelID: "default"}}
	for _, pair := range cliPairs {
		key, value, ok := strings.Cut(pair, "=")
		if !ok || key == "" {
			return nil, fmt.Errorf("%w: malformed override %q", ErrConfig, pair)
		}
		switch key {
		case "model":
			if value == "" {
				return nil, fmt.Errorf("%w: model cannot be empty", ErrConfig)
			}
			c.Model.ModelID = value
		case "enable":
			if value != "" {
				c.Enable = strings.Split(value, ",")
			}
		case "max_tokens":
			n, err := strconv.Atoi(value)
			if err != nil || n < 0 {
				return nil, fmt.Errorf("%w: bad max_tokens %q", ErrConfig, value)
			}
			c.Model.MaxTokens = n
		default:
			return nil, fmt.Errorf("%w: unknown key %q", ErrConfig, key)
		}
	}
	_ = getenv
	return c, nil
}

func (c *Config) NewConversation(system string, log io.Writer) (*agentkit.Conversation, error) {
	return nil, fmt.Errorf("%w: conversation construction is not implemented", ErrAuth)
}
