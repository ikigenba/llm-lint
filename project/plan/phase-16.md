# Phase 16 — agent-repl model/provider/auth selection in `internal/config`

*Realizes design Decision 2 (Configuration resolution).*

Bring `internal/config` to the D02 surface so model, provider, and
authentication are selected exactly the way agent-repl selects them, with the
one deliberate difference that the default stays `gemini-3.7-flash` on `google`
with `auth=key`.

Observable end state:

- The `Model` struct carries `ProviderExplicit`, `Auth`, and `AuthFile`; the
  flat model keys accepted by `.llm-lint.json` and `-c` include `auth`
  (`key`|`sub`) and `auth_file`, both resettable with the literal `default`.
- A provider table (`Providers()`, `ProviderInfo`, `ProviderSpec`) is the
  authority for the six providers, their env keys — google's is
  **`GEMINI_API_KEY`** — and their auth methods, with `openai` and `x-ai`
  listing `sub` first (their default) and the other four `key` only.
- Setting `provider` marks it explicit and skips catalog derivation; a bare
  `model` derives its provider from the catalog and, when the model is not
  catalogued, is rejected for lack of a routable provider. An explicit provider
  with a non-catalogued model is a pass-through: the conversation is built with
  the raw model and no pricing, and a no-pricing notice is appended to
  `Config.Warnings`.
- `NewConversation` builds the provider per the resolved auth mode:
  `auth=key` reads the env key (missing → `ErrAuth` naming it); `auth=sub`
  loads the token file via `openai/subscription.Load` or
  `xai/subscription.Load` and builds a `Subscription` credential, defaulting
  the path to `~/.llm-lint/<provider>-auth.json` (`getenv("HOME")`) unless
  `auth_file` overrides it (missing/unreadable/malformed → `ErrAuth` naming the
  path). `auth=sub` on a key-only provider and an out-of-range `auth` value are
  `ErrConfig`.
- The existing `R-G87C-RGE1` test is updated so the missing-key `ErrAuth`
  message names `GEMINI_API_KEY` for the default google provider.

**Done when** the following are covered by genuine `*_test.go` tests in
`internal/config` and the suite is green (`go build ./...`, `go vet ./...`,
`go test ./...` exit 0; `gofmt -l .` excluding `project/` prints nothing):

- R-IVUP-YH7P — explicit `provider=Q` overrides the catalog default provider P
  for a model, and `ProviderExplicit` is set.
- R-IX2M-C8YE — a bare `model=` unknown to the catalog with no provider returns
  `ErrConfig` telling the caller to set a provider explicitly.
- R-IYAI-Q0P3 — explicit provider + non-curated model builds a conversation
  with the raw model as `Model` and nil `Pricing`, appending a no-pricing
  warning to `Config.Warnings` (no error).
- R-IZIF-3SFS — an `auth` value other than `key`/`sub` returns `ErrConfig`
  naming the bad value.
- R-J0QB-HK6H — `auth=sub` on a key-only provider (e.g. google) returns
  `ErrConfig` naming the provider as key-only.
- R-J364-93NV — with `auth` unset, an openai/x-ai model resolves to `sub` and a
  key-only-provider model resolves to `key`.
- R-J4E0-MVEK — `auth=sub` with a present, loadable (synthetic) token file
  builds a subscription-backed conversation, honoring the default
  `~/.llm-lint/<provider>-auth.json` path and the `auth_file` override.
- R-J5LX-0N59 — `auth=sub` with a missing/unreadable token file returns
  `ErrAuth` naming the exact path.
- R-G87C-RGE1 — its updated test asserts the missing-key `ErrAuth` for the
  default provider names `GEMINI_API_KEY`.
