package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ikigenba/agentkit"
	"github.com/ikigenba/agentkit/anthropic"
	"github.com/ikigenba/agentkit/catalog"
	"github.com/ikigenba/agentkit/google"
	"github.com/ikigenba/agentkit/openai"
	"github.com/ikigenba/agentkit/openrouter"
	"github.com/ikigenba/agentkit/xai"
	"github.com/ikigenba/agentkit/zai"
)

const defaultModel = "gemini-3.7-flash"

type Config struct {
	Root    string
	Enable  []string
	Disable []string
	Rules   []string
	Exclude []string
	Model   Model

	getenv func(string) string
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

type fileConfig struct {
	Enable  []string        `json:"enable"`
	Disable json.RawMessage `json:"disable"`
	Rules   []string        `json:"rules"`
	Exclude []string        `json:"exclude"`
	Model   json.RawMessage `json:"model"`
}

func Load(cwd string, cliPairs []string, getenv func(string) string) (*Config, error) {
	c := &Config{Root: cwd, getenv: getenv}
	name, err := findFile(cwd)
	if err != nil {
		return nil, err
	}
	if name != "" {
		if err := loadFile(name, c); err != nil {
			return nil, err
		}
		c.Root = filepath.Dir(name)
	}
	for _, pair := range cliPairs {
		key, value, ok := strings.Cut(pair, "=")
		if !ok || key == "" {
			return nil, fmt.Errorf("%w: malformed override %q", ErrConfig, pair)
		}
		if err := setModelValue(&c.Model, key, value, value == "default"); err != nil {
			return nil, fmt.Errorf("%w: override %q: %v", ErrConfig, pair, err)
		}
	}
	if c.Model.ModelID == "" {
		c.Model.ModelID = defaultModel
	}
	resolution := catalog.Resolve(agentkit.ProviderID(c.Model.Provider), c.Model.ModelID)
	if resolution.Coverage == catalog.Unrouted {
		return nil, fmt.Errorf("%w: unknown model %q", ErrConfig, c.Model.ModelID)
	}
	c.Model.Provider = string(resolution.Provider)
	return c, nil
}

func findFile(cwd string) (string, error) {
	dir, err := filepath.Abs(cwd)
	if err != nil {
		return "", fmt.Errorf("%w: resolve %q: %v", ErrConfig, cwd, err)
	}
	for {
		name := filepath.Join(dir, ".llm-lint.json")
		_, err := os.Stat(name)
		if err == nil {
			return name, nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("%w: inspect %s: %v", ErrConfig, name, err)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", nil
		}
		dir = parent
	}
}

func loadFile(name string, c *Config) error {
	data, err := os.ReadFile(name)
	if err != nil {
		return fmt.Errorf("%w: read %s: %v", ErrConfig, name, err)
	}
	var raw fileConfig
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&raw); err != nil {
		return fmt.Errorf("%w: %s: %v", ErrConfig, name, err)
	}
	if err := ensureEOF(decoder); err != nil {
		return fmt.Errorf("%w: %s: %v", ErrConfig, name, err)
	}
	c.Enable, c.Rules, c.Exclude = raw.Enable, raw.Rules, raw.Exclude
	if len(raw.Disable) != 0 {
		if string(raw.Disable) == "null" || json.Unmarshal(raw.Disable, &c.Disable) != nil {
			return fmt.Errorf("%w: %s: disable must be an array of strings", ErrConfig, name)
		}
	}
	if len(raw.Model) != 0 && string(raw.Model) != "null" {
		if err := decodeModel(raw.Model, &c.Model); err != nil {
			return fmt.Errorf("%w: %s: %v", ErrConfig, name, err)
		}
	}
	return nil
}

func ensureEOF(decoder *json.Decoder) error {
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return errors.New("multiple JSON values")
		}
		return err
	}
	return nil
}

func decodeModel(data []byte, model *Model) error {
	var values struct {
		Provider         *string  `json:"provider"`
		Model            *string  `json:"model"`
		Temperature      *float64 `json:"temperature"`
		TopP             *float64 `json:"top_p"`
		MaxTokens        *int     `json:"max_tokens"`
		Effort           *string  `json:"effort"`
		ThinkingBudget   *int     `json:"thinking_budget"`
		ThinkingLevel    *string  `json:"thinking_level"`
		Thinking         *bool    `json:"thinking"`
		MaxAttempts      *int     `json:"max_attempts"`
		BaseDelay        *string  `json:"base_delay"`
		MaxDelay         *string  `json:"max_delay"`
		MaxElapsed       *string  `json:"max_elapsed"`
		IgnoreRetryAfter *bool    `json:"ignore_retry_after"`
		BaseURL          *string  `json:"base_url"`
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&values); err != nil {
		return fmt.Errorf("model: %v", err)
	}
	ordered := []modelValue{
		modelPointer("provider", values.Provider), modelPointer("model", values.Model),
		modelPointer("temperature", values.Temperature), modelPointer("top_p", values.TopP),
		modelPointer("max_tokens", values.MaxTokens), modelPointer("effort", values.Effort),
		modelPointer("thinking_budget", values.ThinkingBudget), modelPointer("thinking_level", values.ThinkingLevel),
		modelPointer("thinking", values.Thinking), modelPointer("max_attempts", values.MaxAttempts),
		modelPointer("base_delay", values.BaseDelay), modelPointer("max_delay", values.MaxDelay),
		modelPointer("max_elapsed", values.MaxElapsed), modelPointer("ignore_retry_after", values.IgnoreRetryAfter),
		modelPointer("base_url", values.BaseURL),
	}
	for _, item := range ordered {
		if !item.present {
			continue
		}
		if err := applyModelValue(model, item.key, item.value); err != nil {
			return err
		}
	}
	return nil
}

type modelValue struct {
	key     string
	value   any
	present bool
}

func modelPointer[T any](key string, value *T) modelValue {
	if value == nil {
		return modelValue{key: key}
	}
	return modelValue{key: key, value: *value, present: true}
}

func setModelValue(model *Model, key, text string, reset bool) error {
	if reset {
		return resetModelValue(model, key)
	}
	var value any = text
	var err error
	switch key {
	case "temperature", "top_p":
		value, err = strconv.ParseFloat(text, 64)
	case "max_tokens", "thinking_budget", "max_attempts":
		value, err = strconv.Atoi(text)
	case "thinking", "ignore_retry_after":
		value, err = strconv.ParseBool(text)
	}
	if err != nil {
		return fmt.Errorf("bad %s value %q", key, text)
	}
	return applyModelValue(model, key, value)
}

func applyModelValue(model *Model, key string, value any) error {
	switch key {
	case "provider":
		if value.(string) == "" {
			return errors.New("provider cannot be empty")
		}
		model.Provider = value.(string)
	case "model":
		if value.(string) == "" {
			return errors.New("model cannot be empty")
		}
		model.ModelID = value.(string)
	case "temperature":
		n := value.(float64)
		model.Temperature = &n
	case "top_p":
		n := value.(float64)
		model.TopP = &n
	case "max_tokens":
		if value.(int) < 0 {
			return errors.New("max_tokens cannot be negative")
		}
		model.MaxTokens = value.(int)
	case "effort", "thinking_level":
		model.Reasoning = agentkit.Level(value.(string))
	case "thinking_budget":
		if value.(int) < 0 {
			return errors.New("thinking_budget cannot be negative")
		}
		model.Reasoning = agentkit.Budget(value.(int))
	case "thinking":
		if value.(bool) {
			model.Reasoning = agentkit.EnableReasoning()
		} else {
			model.Reasoning = agentkit.DisableReasoning()
		}
	case "max_attempts":
		if value.(int) < 0 {
			return errors.New("max_attempts cannot be negative")
		}
		model.Retry.MaxAttempts = value.(int)
	case "base_delay":
		d, err := time.ParseDuration(value.(string))
		if err != nil {
			return fmt.Errorf("bad base_delay value %q", value)
		}
		model.Retry.BaseDelay = d
	case "max_delay":
		d, err := time.ParseDuration(value.(string))
		if err != nil {
			return fmt.Errorf("bad max_delay value %q", value)
		}
		model.Retry.MaxDelay = d
	case "max_elapsed":
		d, err := time.ParseDuration(value.(string))
		if err != nil {
			return fmt.Errorf("bad max_elapsed value %q", value)
		}
		model.Retry.MaxElapsed = d
	case "ignore_retry_after":
		model.Retry.IgnoreRetryAfter = value.(bool)
	case "base_url":
		model.BaseURL = value.(string)
	default:
		return fmt.Errorf("unknown model key %q", key)
	}
	return nil
}

func resetModelValue(model *Model, key string) error {
	switch key {
	case "provider":
		model.Provider = ""
	case "model":
		model.ModelID = ""
	case "temperature":
		model.Temperature = nil
	case "top_p":
		model.TopP = nil
	case "max_tokens":
		model.MaxTokens = 0
	case "effort", "thinking_budget", "thinking_level", "thinking":
		model.Reasoning = agentkit.ReasoningValue{}
	case "max_attempts":
		model.Retry.MaxAttempts = 0
	case "base_delay":
		model.Retry.BaseDelay = 0
	case "max_delay":
		model.Retry.MaxDelay = 0
	case "max_elapsed":
		model.Retry.MaxElapsed = 0
	case "ignore_retry_after":
		model.Retry.IgnoreRetryAfter = false
	case "base_url":
		model.BaseURL = ""
	default:
		return fmt.Errorf("unknown model key %q", key)
	}
	return nil
}

func (c *Config) NewConversation(system string, log io.Writer) (*agentkit.Conversation, error) {
	resolution := catalog.Resolve(agentkit.ProviderID(c.Model.Provider), c.Model.ModelID)
	if resolution.Coverage == catalog.Unrouted {
		return nil, fmt.Errorf("%w: unknown model %q", ErrConfig, c.Model.ModelID)
	}
	envName, provider, err := c.provider(resolution.Provider)
	if err != nil {
		return nil, err
	}
	key := ""
	if c.getenv != nil {
		key = c.getenv(envName)
	}
	if key == "" {
		return nil, fmt.Errorf("%w: missing %s", ErrAuth, envName)
	}
	provider = constructProvider(resolution.Provider, key, c.Model.BaseURL)
	return &agentkit.Conversation{
		Provider: provider,
		Model:    resolution.WireModel,
		Pricing:  resolution.Offering.Pricing,
		System:   system,
		Log:      log,
		Gen:      agentkit.GenSettings{Temperature: c.Model.Temperature, TopP: c.Model.TopP, MaxTokens: c.Model.MaxTokens, Reasoning: c.Model.Reasoning},
		Retry:    c.Model.Retry,
	}, nil
}

func (c *Config) provider(id agentkit.ProviderID) (string, agentkit.Provider, error) {
	switch id {
	case agentkit.ProviderGoogle:
		return "GOOGLE_API_KEY", nil, nil
	case agentkit.ProviderOpenAI:
		return "OPENAI_API_KEY", nil, nil
	case agentkit.ProviderAnthropic:
		return "ANTHROPIC_API_KEY", nil, nil
	case agentkit.ProviderOpenRouter:
		return "OPENROUTER_API_KEY", nil, nil
	case agentkit.ProviderXAI:
		return "XAI_API_KEY", nil, nil
	case agentkit.ProviderZAI:
		return "ZAI_API_KEY", nil, nil
	default:
		return "", nil, fmt.Errorf("%w: unsupported provider %q", ErrConfig, id)
	}
}

func constructProvider(id agentkit.ProviderID, key, baseURL string) agentkit.Provider {
	switch id {
	case agentkit.ProviderGoogle:
		if baseURL != "" {
			return google.New(google.APIKey(key), google.WithBaseURL(baseURL))
		}
		return google.New(google.APIKey(key))
	case agentkit.ProviderOpenAI:
		if baseURL != "" {
			return openai.New(openai.APIKey(key), openai.WithBaseURL(baseURL))
		}
		return openai.New(openai.APIKey(key))
	case agentkit.ProviderAnthropic:
		if baseURL != "" {
			return anthropic.New(anthropic.APIKey(key), anthropic.WithBaseURL(baseURL))
		}
		return anthropic.New(anthropic.APIKey(key))
	case agentkit.ProviderOpenRouter:
		if baseURL != "" {
			return openrouter.New(openrouter.APIKey(key), openrouter.WithBaseURL(baseURL))
		}
		return openrouter.New(openrouter.APIKey(key))
	case agentkit.ProviderXAI:
		if baseURL != "" {
			return xai.New(xai.APIKey(key), xai.WithBaseURL(baseURL))
		}
		return xai.New(xai.APIKey(key))
	case agentkit.ProviderZAI:
		if baseURL != "" {
			return zai.New(zai.APIKey(key), zai.WithBaseURL(baseURL))
		}
		return zai.New(zai.APIKey(key))
	default:
		return nil
	}
}
