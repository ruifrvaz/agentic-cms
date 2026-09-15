---
status: In Progress
created: "2026-09-15"
mode: Assisted
started: "2026-09-15"
---

# Second content type: recruiter-interview

## Description

Task 013 built the content-type mechanism and its first type,
`candidate-interview`. This task adds the second type, **`recruiter-interview`**
— the hiring side of the same domain — for the explicit purpose the original
handoff named: proving the mechanism generalizes. Its own `TYPE.md` states
the test directly: *"once the installer handles the first type, installing
this one should require nothing but its `scaffold/types/recruiter-interview/`
payload — if it needs more, the mechanism is not general yet."*

Fixture: `~/projects/magnificah/recruitment/recruiter-interview-kb`, single
commit `ed2b1bc` ("Initial commit: deployed-state fixture for the
recruiter-interview type"), no tags. Unlike `candidate-interview`, it has
**no lived end-state reference** — its own `TYPE.md` says so explicitly and
defines acceptance as a synthetic dry-run instead (see Acceptance Criteria).

Verified live against the actual fixture repo this session (not just recalled
from the prior handoff):

- **Base layer**: byte-identical to `scaffold/tree/` except
  `.agentic-cms/VERSION` (same as `candidate-interview`'s fixture).
- **Templates**: 34 total in `.agentic-cms/templates/`, 29 type-specific
  (5 base: `doc`, `entity`, `concept`, `source`, `topic`).
- **Skills**: `recruiter-setup`, `recruiter-round`, `recruiter-candidate`,
  `recruiter-decision`, alongside the unchanged 9 base `content-*` skills.
- **Exercises**: `exercises/001-rate-limiter/` — the same rate-limiter
  module as `candidate-interview`'s fixture, but as a take-home: adds a
  candidate-facing framing plus an interviewer-only `GRADING.md` (no
  `reference/` solution in this fixture).
- **`CONTENT.md` diff**: directory map (`exercises/`, `templates/`
  description, `TYPE.md`, `skills/` description — same 4 lines as
  `candidate-interview`), a filenames exception (candidate-slug-prefixed
  items + dated decision records — **no `type:` frontmatter enum
  addition**, since this type introduces no new page type, every page
  stays a base `doc`/`entity`/`concept`/`source`), a new
  `## Type: recruiter-interview` section, an `## Operations` trailing
  paragraph, a `## Greenfield vs brownfield` Greenfield-bullet sentence,
  **and** a Brownfield-bullet sentence (`candidate-interview`'s fragment
  has no Brownfield edit — this type's does: "only the organization's own
  documents are import candidates — never candidate documents"). **Seven
  edit points, not six** — the anchor-fragment mechanism must not assume a
  fixed count.
- **`CLAUDE.md` diff**: same shape as `candidate-interview` — a second
  managed block, `<!-- agentic-cms:type:recruiter-interview:begin/end -->`,
  after the leading `@AGENTS.md` line (smaqit-owned, excluded) and the base
  block.
- **`TYPE.md`**: prose + code-fence metadata block (not yet YAML
  frontmatter — same starting shape `candidate-interview`'s fixture had
  before task 013 redesigned it), documents ownership, installer gaps
  (identical five gaps — already closed by task 013's mechanism), what's
  excluded, and its own **Acceptance** section (the synthetic dry-run test,
  see below).

## Issue Triage Context

**Mode:** Auto
**Technologies:** None
**Platforms/Environments:** Linux
**Features/Integrations:** `scaffold.InstallType`, `scaffold.ComposeTypeContentMD`, `scaffold.InstallTypeClaudeMD`, `scaffold.applyTypeFragment` (all from task 013 — expected to need zero changes; that expectation is the test)
**Versions/Constraints:** Fixture built against `agentic-cms v0.7.0`; this task targets the current release (v0.8.2 as of filing) plus whatever MINOR bump adding a type warrants

## Design Decisions

- **Reuse the task-013 mechanism completely unmodified, if possible.** No
  new Go code is the expected, desired outcome — `InstallType`,
  `ComposeTypeContentMD`, `InstallTypeClaudeMD`, and `applyTypeFragment`
  are already fully parameterized by type name and impose no fixed count
  or shape on a fragment's anchor edits. If implementing this type turns
  out to need a Go change, that is itself a finding to report clearly
  (per the fixture's own "if it needs more, the mechanism is not general
  yet" framing) — not something to route around silently.
- **Manifest format**: YAML frontmatter on `TYPE.md`, matching
  `candidate-interview`'s task-013 redesign (name/version/base/templates/
  skills/paths) — for consistency, not because this task re-derives the
  decision. Keep the prose below it (What the type owns, Acceptance)
  adapted for a shipped install the same way task 013 adapted
  `candidate-interview`'s (drop "Installer gaps" — closed; drop the
  fixture-authoring self-check recipe — not meaningful in a real install;
  generalize "Acceptance" away from this specific fixture's synthetic
  dry-run wording into a description of the shape a real requisition
  should produce).
- **Known embedding gotchas from task 013, reused, not rediscovered**:
  `exercises/001-rate-limiter/go.mod` must be authored as `go.mod.embedded`
  (a literal `go.mod` makes `go:embed` treat the subtree as a separate
  module and silently exclude it, `all:` prefix notwithstanding), and the
  `exercises` directory must be authored as `_exercises` (so this repo's
  own `go build/vet/test ./...` — which does not understand `go:embed`'s
  `all:` prefix — ignores the leading-underscore path). `unmangleEmbeddedPath`
  already handles both renames generically by directory/filename pattern,
  not type name, so it should need no change either.
- **No shared kernel with `candidate-interview`**, despite real overlap
  (company context, hiring-process shape, a Rounds-table README,
  culture/values, `round-plan`↔`interview-guide`,
  `round-debrief`↔`scorecard`, `entity-panelist`↔`entity-interviewer`,
  `offer-outcome`↔`candidate-offer`, the `exercises/` take-home pattern) —
  carried over from task 013's already-settled decision, not re-opened
  here.
- **Live validation mirrors task 013/014/015's own methodology**: after
  release, actually install the type from the real binary and run a
  synthetic dry-run through all four skills (a fictional requisition,
  fictional candidates) — the exact test this type's own `TYPE.md`
  Acceptance section already prescribes — rather than treating a clean
  `go build`/`smoke-test` as sufficient proof. Tasks 014 and 015 each
  found a real, otherwise-invisible bug this way; assume this type has at
  least one similar latent issue until a live run says otherwise.

## Implementation Steps

1. **Import the overlay** into
   `scaffold/types/recruiter-interview/tree/`: the 29 type templates, the
   four `recruiter-*` skills, `exercises/001-rate-limiter/` (renamed
   `go.mod` → `go.mod.embedded`, directory renamed `exercises` →
   `_exercises`, same as task 013's `candidate-interview` import).
2. **Author `scaffold/types/recruiter-interview/tree/.agentic-cms/TYPE.md`**
   with YAML frontmatter (name `recruiter-interview`, version `0.1.0`,
   `base: "{{VERSION}}"`, `templates:`/`skills:`/`paths:` lists) plus
   adapted prose (see Design Decisions).
3. **Extract `scaffold/types/recruiter-interview/content.fragment.md`** —
   seven `<!-- ac-type-anchor: after|before|replace "<exact base line>" -->`
   blocks (directory map ×4, filenames exception, Type section, Operations
   paragraph, Greenfield sentence, Brownfield sentence — count them
   precisely against the real `diff -U1` output, don't assume task 013's
   six carries over). Verify byte-for-byte against the fixture with the
   same throwaway-script technique task 013 used (`apply_fragment.py`
   equivalent) before trusting the Go path.
4. **Extract `scaffold/types/recruiter-interview/claude.fragment.md`** —
   the single managed block, verified the same way.
5. **No Go changes expected** — run `agentic-cms init --type
   recruiter-interview` in a sandbox using the *existing* `InstallType`/
   `ComposeTypeContentMD`/`InstallTypeClaudeMD`/`ReconcileContentMD` code
   unmodified, confirm it works. If it doesn't, diagnose why before writing
   any new Go — the gap is the finding.
6. **Tests beside the code**: extend `scaffold/types_test.go` (or add
   `scaffold/types_recruiter_test.go`) mirroring `candidate-interview`'s
   coverage — overlay install, `go.mod`/`_exercises` unmangling,
   `{{VERSION}}` stamped, `TYPE.md`'s own `{{DATE}}`-mention safety (the
   task-014 regression class — confirm it can't recur for a second type,
   since the exclusion in `InstallType` is by relative path, not type
   name, so it already covers this automatically; write the test to prove
   it rather than assume it), `ComposeTypeContentMD` produces the exact
   fixture `CONTENT.md`, reconciliation detects a missing
   `## Type: recruiter-interview` section.
7. **Extend `scripts/smoke-test-installer.sh`** with a second typed-install
   section mirroring the `candidate-interview` one (template smoke over
   all 34 templates, idempotent re-run, implicit-type re-init,
   reconciliation), or generalize the existing section to loop over both
   types if that's cleaner — judgment call, note the reasoning either way.
8. **Manual fixture diff**: `diff -rq` a freshly `init --type
   recruiter-interview`'d sandbox against the live fixture repo, same
   exclusion list task 013 used (`.git .smaqit .github AGENTS.md
   *.code-workspace .gitignore`) plus the same expected deviations
   (`VERSION`, `TYPE.md` redesign, `CLAUDE.md`'s leading `@AGENTS.md` line,
   and possibly the same `wiki/index.md`/`log.md` `{{DATE}}`-vs-real-date
   fixture-authoring artifact task 013 already flagged and decided not to
   chase).
9. **Docs**: README.md's Content types section gains a short mention of
   the second type (or stays generic enough not to need one — judgment
   call); `CHANGELOG.md` entry.
10. **Release** through the normal PR flow.
11. **Live validation** (post-release, mirroring tasks 014/015's
    methodology): download the real released binary, fresh `init --type
    recruiter-interview`, then actually run all four skills
    (`recruiter-setup` → `recruiter-round` → `recruiter-candidate` →
    `recruiter-decision`) through a synthetic dry-run — a fictional
    requisition, a fictional company, 2-3 fictional candidates through
    intake → scoring → a decision — independently re-verifying
    `ac-index check`/`ac-links check`/`ac-classify sweep`, and confirming
    every roster/index/log line refers to candidates by slug only (this
    type's own stated acceptance bar). Treat any subagent self-report the
    same way tasks 013-015 did: verify directly, don't trust blindly.
    File and fix whatever this surfaces before calling the type done.

## Known Issues Triage
**Triaged:** 2026-09-15
**Tools searched:** None
**Result:** Clear

### Blocking Issues
- None

### Advisory Issues
- None

### Historical (Closed)
- None

### Unresolvable Tools
- None

### Omitted Tools
- None

### Search Warnings
- None

Notes: No third-party tools identified — triage not applicable. This is a
data-import task reusing this repo's own existing Go mechanism (task 013);
no external dependency is implicated.

## Acceptance Criteria

- [ ] `agentic-cms init --type recruiter-interview` on an empty directory
      produces a tree matching `scaffold/tree` ∪
      `scaffold/types/recruiter-interview/tree` (post-substitution), with a
      composed `CONTENT.md` containing all seven type edits and the base
      `CONTENT.md`/`candidate-interview`-typed installs completely
      unaffected.
- [ ] Implemented with **zero changes to existing Go logic** — or, if not,
      the specific gap is documented as a finding, not silently patched
      around.
- [ ] `make test` and `make smoke-test` both pass, including new coverage
      for `recruiter-interview` mirroring `candidate-interview`'s.
- [ ] A freshly typed sandbox install `diff -rq`s clean against the real
      fixture repo (except `VERSION` and the same already-understood
      deviation classes from task 013), run once manually against the
      released binary.
- [ ] Released via PR flow.
- [ ] Live synthetic-dry-run validation (Implementation Step 11) completed
      against the real released binary, independently re-verified (not
      trusted from a subagent's self-report alone), with any findings
      fixed and released before this task closes.

## Findings

[Populated by smaqit.task-complete. Do not fill in manually before task is complete.]

**Implementation approach:**
- TBD

**Decisions made:**
- TBD

**Blockers encountered:**
- TBD

**Follow-up identified:**
- TBD

## Files to Create / Modify

| File | Action |
|------|--------|
| `scaffold/types/recruiter-interview/tree/.agentic-cms/templates/*.md` (29 files) | Create |
| `scaffold/types/recruiter-interview/tree/.agentic-cms/TYPE.md` | Create |
| `scaffold/types/recruiter-interview/tree/.claude/skills/{recruiter-setup,recruiter-round,recruiter-candidate,recruiter-decision}/SKILL.md` | Create |
| `scaffold/types/recruiter-interview/tree/exercises/001-rate-limiter/*` (incl. `GRADING.md`) | Create |
| `scaffold/types/recruiter-interview/content.fragment.md` | Create |
| `scaffold/types/recruiter-interview/claude.fragment.md` | Create |
| `scaffold/types_test.go` or new `scaffold/types_recruiter_test.go` | Create/Modify |
| `scripts/smoke-test-installer.sh` | Modify |
| `README.md` | Modify (maybe) |
| `CHANGELOG.md` | Modify |
| `main.go`, `scaffold/*.go` | Modify only if Step 5 finds a real gap — not expected |

## Notes

- This task is the direct test of task 013's own design claim. A clean
  result (steps 1-4 + zero Go changes) is the strongest possible evidence
  the content-type mechanism is genuinely general; a Go change requirement
  is itself the most valuable possible finding from this task, not a
  failure to route around quietly.
- Do not build the `interview-process` kernel the overlap between the two
  types gestures at (round-plan/interview-guide, round-debrief/scorecard,
  etc.) — carried over from task 013's explicit decision to keep types
  independent. If overlap maintenance becomes genuinely painful across a
  third type, that's a future decision to raise with the user, not to
  make unilaterally here.
- The fixture's own `TYPE.md` Acceptance section is written assuming no
  lived reference exists (unlike `candidate-interview`'s
  `end-state-reference` tag) — Implementation Step 11's synthetic dry-run
  is this type's *only* end-to-end acceptance signal, so treat it as load
  -bearing, not optional polish.

Child tasks inherit their active parent's branch, worktree, and workflow mode. Only a standalone or parent task owns Git lifecycle cleanup.
