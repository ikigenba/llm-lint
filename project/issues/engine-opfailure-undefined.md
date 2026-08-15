# engine.OpFailure named in D01 but not defined in D05

Observed while briefing Phase 01: `project/design/D01.md` names
`engine.OpFailure` as the operational-failure error surface that `run` maps
to exit 3, but `project/design/D05.md` never declares a type, constructor, or
sentinel of that name. D05 only describes operational failures as errors
returned from `Engine.Run` after retries.

Out of scope for Phase 01 gather: this is a design-doc inconsistency, not
something the brief can settle. Phase 01 still has to expose an
`engine.OpFailure` surface because D01's error-mapping text and R-FYG5-PAGH
require it.
