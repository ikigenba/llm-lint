# Changelog

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
