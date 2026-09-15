---
status: PR Open
created: "2026-09-15"
mode: Assisted
started: "2026-09-15"
pr: 11
---

# Fix {{DATE}} substitution corrupting TYPE.md's own placeholder documentation

## Description

Found via a live end-to-end validation of the released `v0.8.0` binary
(downloaded via the public `install.sh`, not a local build), run immediately
after task 013 shipped: `agentic-cms init --type candidate-interview` writes
a corrupted `.agentic-cms/TYPE.md`.

`TYPE.md`'s own prose documents the page-template placeholders `ac-page new`
substitutes: "...the placeholders `ac-page new` substitutes (`{{TITLE}}`,
`{{TOPIC}}`, `{{RAW_PATH}}`, `{{STATUS}}`, `{{CLASSIFICATION}}`,
`{{CLASSIFIED_HASH}}`, `{{DATE}}`)." — a literal list of placeholder *names*
for a human/agent reader, not live placeholders meant to resolve.

`scaffold.InstallType` (`scaffold/types_install.go`) substitutes `{{DATE}}`
everywhere in the type tree except two prefixes: `.agentic-cms/templates/`
and `exercises/`. `.agentic-cms/TYPE.md` matches neither exclusion, so its
own `{{DATE}}` mention gets replaced with the literal install date, turning
the sentence into "...`{{CLASSIFIED_HASH}}`, 2026-09-15)." — nonsensical
documentation, permanently baked into every fresh typed install.

The task 013 smoke test never caught this because its only relevant
assertion was "no leftover `{{ }}`" — the *opposite* failure mode. A
substitution that "succeeds" against the wrong text passes that check
silently.

## Issue Triage Context

**Mode:** Auto
**Technologies:** Go 1.21
**Platforms/Environments:** Linux
**Features/Integrations:** `scaffold.InstallType`, `.agentic-cms/TYPE.md`
**Versions/Constraints:** Fixes a regression in v0.8.0 (task 013); targets v0.8.1 (PATCH — bug fix only, no new capability)

## Design Decisions

- **Fix:** treat `.agentic-cms/TYPE.md` like the existing `.agentic-cms/templates/`/`exercises/` exclusions in `InstallType`'s `{{DATE}}` substitution — exclude it entirely. `TYPE.md` has no field of its own that legitimately needs today's date resolved (unlike `.agentic-cms/VERSION`, which has its own explicit, separate `{{VERSION}}` substitution path that is unaffected and stays correct).
- **Regression test:** assert that installing the type leaves the literal string `{{DATE}}` present in `TYPE.md`'s prose (the opposite of the smoke test's existing "no leftover `{{`" check elsewhere) — this is the exact inversion that let the bug through originally, so the test must check for presence, not absence.
- **Scope:** this is the only known instance of the pattern (no other type-tree file outside `templates/`/`exercises/` mentions a placeholder name in prose), so no broader exclusion mechanism is being built — just closing this one specific gap. If a future type file needs the same treatment, extend the same exclusion list.

## Implementation Steps

1. In `scaffold/types_install.go`'s `InstallType`, add `.agentic-cms/TYPE.md` (or the general `typeManifestFile` constant already in scope) to the `{{DATE}}` substitution exclusion alongside `.agentic-cms/templates/` and `exercises/`.
2. Add a regression test in `scaffold/types_test.go` asserting the installed `TYPE.md` still contains the literal substring `{{DATE}}` after `InstallType` runs.
3. `make test` and `make smoke-test` both green.
4. Release through the normal PR flow (`smaqit.release-analysis` → PR → merge), expected v0.8.1 (PATCH).
5. Re-verify live: download the newly-released binary, fresh `init --type candidate-interview`, confirm `TYPE.md`'s prose is intact.

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

Notes: No third-party tools/libraries are implicated — the bug is a plain
string-substitution ordering mistake entirely within this repo's own
`scaffold/types_install.go` (unlike task 013, this doesn't touch `go:embed`
semantics or any other third-party behavior). "Go 1.21" is the only listed
technology and is the language itself, not a specific feature at play here,
so a GitHub search against `golang/go` would produce the same kind of
generic keyword noise task 013's triage already found irrelevant — not
run again for the same reason.

## Acceptance Criteria

- [x] A fresh `agentic-cms init --type candidate-interview` no longer corrupts `.agentic-cms/TYPE.md`'s placeholder-name prose. Verified via a local dev build with the fix applied.
- [x] Regression test added and passing; `make test` and `make smoke-test` both pass. Test verified to actually catch the bug (confirmed it fails against the pre-fix code, passes with the fix).
- [ ] Released via PR flow as v0.8.1.
- [ ] Verified live against the real released v0.8.1 binary. *(Post-merge step, same pattern as task 013's criterion 8 — done once v0.8.1 is actually released.)*

## Findings

**Implementation approach:**
- Added `.agentic-cms/TYPE.md` (via the existing `typeManifestFile` constant)
  as a third exclusion in `InstallType`'s `{{DATE}}` substitution check,
  alongside the pre-existing `.agentic-cms/templates/` and `exercises/`
  exclusions — same mechanism, one more path.
- Added the regression test right next to `TestInstallTypeOverlay`'s
  existing `{{VERSION}}`-is-stamped assertion, since it's the natural
  contrasting case: `{{VERSION}}` must resolve, `{{DATE}}` in this file
  must not.

**Decisions made:**
- Verified the regression test actually catches the bug before trusting it:
  stashed the fix, confirmed the test fails against the pre-fix code
  (`TYPE.md's placeholder-name documentation lost its literal {{DATE}}
  mention`), restored the fix, confirmed it passes. A test that was never
  seen to fail is not verified.
- Verified the fix live: built a local dev binary, fresh `init --type
  candidate-interview`, confirmed `TYPE.md`'s prose reads intact with
  `{{DATE}}` still literally present.
- Scoped narrowly per the task's own Design Decisions — just excluded
  `TYPE.md` from this one substitution, did not build a general
  "documentation file" exclusion mechanism for a single known instance.

**Blockers encountered:**
- None.

**Follow-up identified:**
- Acceptance criterion 4 (verify against the real released v0.8.1 binary)
  is a post-merge step, same pattern as task 013's criterion 8 — done once
  v0.8.1 is actually released.
- Layer 2 of the broader e2e validation (live skill run against a
  synthetic engagement) was paused when this bug surfaced during Layer 1;
  resumes once v0.8.1 is released and re-verified.

## Files to Create / Modify

| File | Action |
|------|--------|
| `scaffold/types_install.go` | Modify |
| `scaffold/types_test.go` | Modify |
| `CHANGELOG.md` | Modify |

## Notes

Discovered mid-flow while planning a broader live e2e validation (layer 2:
running all four `candidate-interview` skills against a synthetic engagement
to check the resulting tree shape) — this patch fix is a prerequisite so
that validation runs against a correct binary, not a known-buggy one.

Child tasks inherit their active parent's branch, worktree, and workflow mode. Only a standalone or parent task owns Git lifecycle cleanup.
