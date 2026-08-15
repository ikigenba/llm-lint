# llm-lint — Research

Non-contractual ground truth gathered so design does not re-derive it. All
facts below were read from the named sources on 2026-08-15.

## agentkit (github.com/ikigenba/agentkit, v0.20.0)

The suite's multi-provider LLM client. The parts llm-lint uses:

- **Conversation** (`orchestration.go`) is the unit of work: a struct with
  exported fields `Provider`, `Model`, `Pricing`, `System`, `Log io.Writer`,
  `Gen GenSettings`, `Retry RetryPolicy`, `Tools []Tool`, `History`,
  `MaxToolIterations`. `Send(ctx, userText) *Stream` runs one turn including
  the internal tool-call loop. Conversations retain history and must not be
  shared between independent calls — autotune documents this on its
  `ProviderFactory` seam (`autotune/internal/config/provider_factory.go`):
  build a fresh Conversation per independent call.
- **Stream** exposes `Events()`, `Err()`, `Usage()`, `Cost()`, `Warnings()`.
  `Cost` is computed from the catalog `Pricing` when set.
- **Structured output**: agentkit has **no forced tool choice** in its public
  surface (the only trace is `WarnToolChoiceForced`, emitted when a provider
  downgrades to auto). The suite idiom for structured results is a typed tool:
  `agentkit.NewTool[In](name, description, fn)` derives a JSON schema from the
  `In` struct via reflection and decodes the model's call into it. Forcing is
  done by instruction (system prompt demands the tool call); the caller must
  handle the model answering in prose without calling the tool.
- **Providers**: `openai.New(openai.APIKey(key), opts...)`,
  `google.New(google.APIKey(key), opts...)` etc. — one constructor per
  provider package, credential as a typed value. autotune's
  `internal/config/config.go` is the reference for building a provider from a
  catalog resolution plus `getenv`.
- **Catalog** (`agentkit/catalog`): `Lookup(model) (Entry, bool)`,
  `Resolve(provider, model)`, `ListCurated(provider)`, plus reasoning
  validation via `Check`. Entries carry context window, pricing, and the
  provider's env key. Both target models are curated:
  `gemini-3.7-flash` (Google, 1,048,576-token context) and `gpt-5.6-luna`
  (OpenAI, 400,000-token context) — `catalog/data.go:122,178`.
- **Raw exchange log**: setting `Conversation.Log` to an open file writes the
  complete raw exchange as JSONL. agent-repl writes
  `~/.agentkit/<session-id>.jsonl`, unbuffered, write-only ("forensic output —
  never read back to resume"); its `internal/session/session.go` is the
  reference (`Open(dir)` creates the dir and the timestamped file).
- **Retry**: `agentkit.RetryPolicy` on the Conversation; agent-repl exposes
  its knobs as the flat config keys `max_attempts`, `base_delay`, `max_delay`,
  `max_elapsed`, `ignore_retry_after`.

## Sibling conventions (ralph, agent-repl, autotune)

- **Shallow main**: agent-repl's `cmd/agentrepl/main.go` is the exact idiom —
  `func main() { os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr, …)) }`
  with every exit path flowing through `run`'s returned code, and a
  `var version = "dev"` overridden via `-ldflags "-X main.version=..."`.
- **Flat model config keys** (agent-repl README): `provider`, `model`,
  `system`, `temperature`, `top_p`, `max_tokens`, the native reasoning keys
  (`effort`, `thinking_budget`, `thinking_level`, `thinking`), the retry keys,
  `base_url` (zai only). The literal value `default` resets a key to unset.
- **Live-test gating**: agentkit's own integration tests skip when the
  provider key is absent (`deferred_tools_integration_test.go:19-21`:
  `os.Getenv("ANTHROPIC_API_KEY")` → `t.Skip`). This keeps `go test ./...`
  green on machines without credentials while proving the live path where keys
  exist.
- **Frontmatter**: ralph hand-parses a strict `key: value` frontmatter subset
  (`ralph/internal/ralph/frontmatter.go`) rather than depending on a YAML
  library. The suite tolerates small third-party deps where warranted (ralph
  pulls chroma; autotune pulls a jsonschema stack) but leans stdlib-first.
- **Release**: ralph/agent-repl release via a `v*` tag-push GitHub Actions
  workflow (`.github/workflows/release.yml`: checkout with full depth,
  setup-go 1.26, `goreleaser/goreleaser-action@v6` with `release --clean`) and
  a goreleaser v2 config (CGO off, linux/darwin × amd64/arm64, tar.gz
  `{{ .Binary }}_{{ .Os }}_{{ .Arch }}`, checksums, github release under
  `ikigenba/<tool>`), plus a POSIX `install.sh` that resolves OS/arch, fetches
  the latest (or `$<TOOL>_VERSION`) release asset, and installs into
  `${BINDIR:-${PREFIX:-$HOME/.local}/bin}`. ralph keeps a `release_test.go`
  asserting these files stay consistent.
- **Toolchain**: go 1.26 across the suite; go1.26.5 on this machine.

## Options evaluated and not chosen

- **Forced JSON / native structured-output modes**: not exposed by agentkit
  v0.20.0; would mean bypassing the suite client. Rejected in favor of the
  typed-tool idiom above.
- **YAML library for rule frontmatter**: the needed subset (two scalars, two
  string lists) does not justify the dependency given the ralph precedent;
  design specifies a strict hand-parsed subset instead.
- **A gitignore-matching library**: exact `.gitignore` semantics (nested
  files, negations) are notoriously fiddly to reimplement. When the linted
  tree is a git work tree, `git ls-files --cached --others
  --exclude-standard` enumerates exactly the non-ignored files using git's own
  matcher at zero dependency cost; a plain directory walk is the non-repo
  fallback. Chosen over both reimplementation and a library dependency.
- **Glob matching**: stdlib `path.Match` has no `**`. The standard Go answer
  is `github.com/bmatcuk/doublestar/v4`; one small, widely-used dependency,
  accepted.
