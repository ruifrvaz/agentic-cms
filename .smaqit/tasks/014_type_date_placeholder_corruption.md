---
status: In Progress
created: "2026-09-15"
mode: Assisted
started: "2026-09-15"
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

- [ ] A fresh `agentic-cms init --type candidate-interview` no longer corrupts `.agentic-cms/TYPE.md`'s placeholder-name prose.
- [ ] Regression test added and passing; `make test` and `make smoke-test` both pass.
- [ ] Released via PR flow as v0.8.1.
- [ ] Verified live against the real released v0.8.1 binary.

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
| `scaffold/types_install.go` | Modify |
| `scaffold/types_test.go` | Modify |
| `CHANGELOG.md` | Modify |

## Notes

Discovered mid-flow while planning a broader live e2e validation (layer 2:
running all four `candidate-interview` skills against a synthetic engagement
to check the resulting tree shape) — this patch fix is a prerequisite so
that validation runs against a correct binary, not a known-buggy one.

Child tasks inherit their active parent's branch, worktree, and workflow mode. Only a standalone or parent task owns Git lifecycle cleanup.
