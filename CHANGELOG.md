# Changelog

## v0.3.0

- Added a `--verbose` option that prints, for every examined file and enabled
  rule, whether its verdict was served from cache or freshly judged and whether
  the rule passed or flagged the file, on stderr and without altering the
  findings written to stdout.
- Made `--stats` report the run's real input-token, output-token, and cost
  totals drawn from the model calls it made; previously these figures always
  printed as zero.
- Removed the silent skip of files larger than the model's context window: an
  oversized file is now sent to the model like any other and, if the call
  cannot complete, surfaces as an ordinary operational failure rather than
  being quietly passed over.

## v0.2.0

- Let a run select any supported model by name through the `model` config key,
  running it against that model's default provider automatically; a bare model
  name must be one of the catalogued models, while naming a provider explicitly
  with `provider` accepts any model as pass-through.
- Added provider and authentication selection via the `provider`, `auth`, and
  `auth_file` config keys, spanning the Anthropic, Google, OpenAI, OpenRouter,
  xAI, and Z.AI providers, each with its own API-key environment variable.
- Allowed a run to authenticate against a provider subscription instead of a
  metered API key, the cheaper path, by setting `auth=sub` and pointing
  `auth_file` at a subscription token file (defaulting to
  `~/.llm-lint/<provider>-auth.json`); subscription auth is offered where the
  provider supports it.
- Extended `llm-lint --help` to list every selectable model grouped by
  provider, each provider's authentication options, and its API-key environment
  variable, so a user can see what they can select without leaving the terminal.
- Kept all model configuration on the repeatable `-c key=value` override
  mechanism with no dedicated per-setting flags.

## v0.1.0

- Introduced `llm-lint`, a CLI that lints a code tree against prompt-defined
  rules by asking a cheap language model for each verdict and reporting
  violations in vet-style `file:line: message (rule)` form, printing nothing on
  a clean run.
- Shipped one built-in rule, `no-sleep-in-tests`, that flags fixed-duration
  sleeps used to wait out a timing or synchronization effect in tests while
  sparing legitimate delays; every available rule is active by default.
- Added project configuration through `.llm-lint.json`: an `enable` allowlist
  to restrict a run to chosen rules, a `disable` list to silence rules by name,
  `exclude` path globs, project-authored rule files, and per-project model
  settings, all overridable per invocation with `-c key=value`.
- Produced text and JSON output, the JSON carrying the exact offending source
  line so tooling can build stable identities for findings.
- Honored inline `llm-lint:ignore` comments so a finding judged acceptable
  disappears from the next run.
- Cached verdicts by file and rule content, so re-running on an unchanged tree
  spends nothing and makes no new model calls; the model is deliberately not
  part of the cache key, so changing models never forces a re-judge.
- Provided the `--list-rules`, `--stats`, `--no-cache`, `--concurrency`,
  `--model`, `--rules`, `--format`, and `-V`/`--version` flags, with a vet-style
  exit code taxonomy (0 clean, 1 findings, 2 usage or config error, 3
  operational failure) so CI can tell a failing lint apart from a broken tool.
- Distributed prebuilt binaries for Linux and macOS, installable with the
  `install.sh` script.
