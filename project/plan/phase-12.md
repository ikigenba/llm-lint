# Phase 12 — rules: default-on selection with `disable` subtraction

*Realizes design Decision 3 (Rule model and catalog).*

`rules.Select` takes its new signature `Select(all []Rule, enable, disable
[]string) ([]Rule, error)` and computes the active set: an empty `enable`
selects every rule in `all` (in `all` order); a non-empty `enable` selects only
the named rules (in `enable` order); then any id present in `disable` is removed
from the result. An `enable` or `disable` entry naming no known rule is an
`ErrRule` naming the id. The signature change ripples to the existing
`Select` call sites and tests, which this phase updates to pass the new
argument.

**Done when:**

- R-2LGZ-C5OJ — `Select` with empty `enable` and empty `disable` returns every
  rule in `all`, in `all` order.
- R-2MOV-PXF8 — a rule id in `disable` is absent from `Select`'s result, both
  from the all-on default and from an explicit `enable` allowlist.
- R-2NWS-3P5X — a `disable` entry naming no known rule returns `ErrRule` naming
  the id.
- R-GD2Y-AJCT and R-GGQN-FUKW remain green under the new signature (non-empty
  `enable` yields exactly the named rule; unknown `enable` id errors).
- The suite is green.
