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
- **Providers**: `openai.New(cred, opts...)`, `google.New(cred, opts...)` etc.
  — one constructor per provider package, taking a typed `Credential`. The
  key credential is `openai.APIKey(key)` (and the analogue per package).
  autotune's `internal/config/config.go` is the reference for building a
  provider from a catalog resolution plus `getenv`.
- **Subscription credential**: `openai` and `xai` additionally expose
  `Subscription(ts TokenSource) Credential` and a sibling package
  `openai/subscription` (and `xai/subscription`) whose `Load(path) (*Store,
  error)` reads a raw OAuth token-endpoint response file; the returned `*Store`
  satisfies `TokenSource`. `Load` requires the file to be JSON with a non-empty
  `access_token`, and derives a ChatGPT account id by base64url-decoding the
  JWT payload of `id_token` (falling back to `access_token`) and reading the
  `https://api.openai.com/auth` → `chatgpt_account_id` claim; a missing token
  or account claim is an error. The **signature is never verified** (only the
  payload is base64-decoded), so a test can construct a synthetic-but-loadable
  token file — agentkit's own `openai/subscription/subscription_test.go`
  (`TestLoadDerivesAccountFromRawTokenResponse`) is the reference. `Load`
  performs no discovery and reads no ambient credentials: the caller owns path
  selection and the initial login. `agentkit` has **no** subscription support
  for the other four providers.
- **Provider env keys and auth methods** (agent-repl's catalog wrapper is the
  reference for the exact strings): `anthropic`/`ANTHROPIC_API_KEY` (key only),
  `google`/`GEMINI_API_KEY` (key only), `openai`/`OPENAI_API_KEY` (**sub**,
  key — first is the default), `openrouter`/`OPENROUTER_API_KEY` (key only),
  `x-ai`/`XAI_API_KEY` (**sub**, key), `z-ai`/`ZAI_API_KEY` (key only).
  agent-repl's default auth file is `~/.agentrepl/<provider>-auth.json` via
  `config.DefaultAuthFile(dir, providerID)`.
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
  `base_url`, and the auth keys `auth` (`key`|`sub`) and `auth_file`. The
  literal value `default` resets a key to unset.
- **Model/provider/auth selection** (agent-repl `internal/config/config.go`,
  `internal/repl/help.go`): a bare `model=` derives its provider from the
  catalog's first offering for that model (`catalog.Resolve` with an empty
  provider); setting `provider=` pins the provider explicitly so derivation is
  skipped. `auth`, when unset, takes the provider's **first** listed method
  (so openai/x-ai default to `sub`, the rest to `key`); setting `auth=sub` on a
  key-only provider is rejected. `WriteHelp` renders a `providers:` block (each
  provider's `auth=key (ENVVAR)` and, where supported, `auth=sub
  (auth_file=…)`) followed by each provider's models with reasoning clauses
  (`*` marks the default), then a footer describing bare-`model=` vs.
  explicit-provider selection. agent-repl warns and *passes through* an
  uncatalogued model when a provider is set explicitly (no pricing, reasoning
  unchecked).
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
