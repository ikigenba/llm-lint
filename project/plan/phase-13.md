# Phase 13 — cmd: wire default-on, `disable`, and the `-V` alias

*Realizes design Decision 1 (CLI surface and shallow composition root).
Depends on Phase 11 and Phase 12.*

`run` passes both `cfg.Enable` and `cfg.Disable` into `rules.Select`, so the
default-on and opt-out behavior reaches the assembled program. A `-V` flag is
added as a short alias of `--version` with identical output and exit 0, and it
appears in the `--help` usage string. The empty-set warning now fires on an
empty **effective** rule set rather than an empty `enable` list: the retired
`R-GHYJ-TMBL` tag and its test are removed, and `--list-rules` and the
finding/clean `run` tests are updated to the new default (with no config every
rule is active).

**Done when:**

- R-2J16-KM75 — `-V` prints exactly the version variable to stdout and returns
  0, identically to `--version`, and `-V` appears in the `--help` usage string.
- R-2P4O-HGWM — `run` with an empty effective rule set (every active rule
  disabled, or none available) prints `llm-lint: no rules enabled` to stderr
  and returns 0.
- The retired `R-GHYJ-TMBL` tag no longer appears in any test file
  (`grep -rl 'R-GHYJ-TMBL' --include='*_test.go' .` finds nothing).
- R-FX89-BIPS, R-FR4R-EO0B, R-FSCN-SFR0, R-FTKK-67HP remain green under the
  new default-on behavior.
- The suite is green.
