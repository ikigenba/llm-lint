# Changelog

## v0.7.0

- Updated the bundled agentkit client to v0.22.0, which adds four models to the
  advisory catalog. Three are served through OpenRouter —
  `nemotron-3.5-lightning` (NVIDIA), `qwen3.8-max`, and `qwen3.8-27b` (Qwen) —
  and `claude-opus-5` is served natively by Anthropic with an OpenRouter
  alternative. Each can now be selected with `-c model=<name>` (or a `model`
  entry in `.llm-lint.json`) and appears in the `--help` model listings.

## v0.6.0

- Redesigned the `--help` screen to follow agent-repl's layout: a `usage:`
  line, a `flags:` block listing each option with a short description, a
  `defaults:` section showing the built-in default provider, model, and auth, an
  aligned `providers:` block, and per-provider model listings whose reasoning
  choices are shown in braces with the default marked by `*`, with the sections
  separated by blank lines.
- Changed the response to an unknown flag to print the flag error followed by a
  single one-line usage summary, instead of printing the full usage block
  twice.

## v0.5.0

- Added a `boolean-state-machine` built-in rule that flags groups of two or
  more boolean fields, variables, or columns that together encode one concept
  with more than two states, and recommends replacing them with a single enum
  or status value.
- Gave every rule a default set of code file globs, so a rule that omits its
  `include` now examines common source files (Go, Python, JavaScript,
  TypeScript, Java, Ruby, Rust, C, C++, C#) instead of having to list them
  explicitly; an `include` present but empty is still rejected.
- Reformatted each `--verbose` trace line to lead with a pass or fail circle,
  then the bracketed rule id, then the file path, so the outcome and the rule
  read first at a glance.

## v0.4.0

- Changed `--verbose` to stream its per-file, per-rule audit trace to stderr
  live as the run proceeds: each line now appears the moment that pair is
  judged, in the order verdicts complete, so progress is visible while the run
  is still going and a run that aborts still shows the work it finished. It
  previously buffered every line and printed them sorted only after the run
  completed.

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
