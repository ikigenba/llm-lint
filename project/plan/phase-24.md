# Phase 24 — Default code globs and the boolean-state-machine built-in rule

*Realizes design Decision 3 (rule model and catalog — the optional-include
slice), Decision 12 (built-in rule: boolean-state-machine), and Decision 9's
count-free restatement.*

Three coupled pieces in `internal/rules`:

- The `DefaultCodeGlobs` exported list and the optional-`include` frontmatter
  semantics (D03): a rule file omitting the key parses with the default
  globs; a present-but-empty `include` array remains `ErrRule`. The existing
  malformed-frontmatter test (R-GBV1-WRM4) updates to the new error set.
- The new embedded rule file
  `internal/rules/builtin/boolean-state-machine.md` per D12's frame
  (description, severity error, no `include` key, no excludes, the
  invalid-combination prompt with the ticket held/sold anchor and one
  flagged plus one spared example).
- The catalog test currently asserting exactly one built-in updates to D09's
  count-free restatement of R-HIJJ-AQKF (no-sleep-in-tests present with its
  globs), alongside D12's new catalog assertion.

**Done when:**

- R-4HI3-PKJF — a rule file that omits the `include` key parses successfully
  with `Include` equal to `DefaultCodeGlobs`, and `DefaultCodeGlobs` is
  exactly the eleven code globs D03 lists, in that order — tagged by a
  genuine test.
- R-GBV1-WRM4 — malformed frontmatter (unterminated fence, unknown key,
  invalid severity, a present-but-empty `include` array) returns `ErrRule`
  naming the file — its existing test updated to this behavior.
- R-4IQ0-3CA4 — the embedded catalog contains `boolean-state-machine` with
  severity error, `Include` equal to `rules.DefaultCodeGlobs`, no excludes,
  and the invalid-combination/held-sold prompt anchors; it appears in
  `--list-rules` output — tagged by a genuine test.
- R-HIJJ-AQKF — the embedded catalog contains `no-sleep-in-tests` with
  severity error, its test-file globs, and no excludes; it appears in
  `--list-rules` output — its existing test updated to drop the catalog
  count.
- The suite is green (per design CONVENTIONS.md).
