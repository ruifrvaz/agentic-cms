---
status: PR Open
created: "2026-09-15"
mode: Assisted
started: "2026-09-15"
pr: 12
---

# Fix stale AGENTS.md references in candidate-interview type skill files

## Description

Found via a live end-to-end run of all four `candidate-interview` skills
(`interview-setup` → `interview-round` ×2 → `interview-rehearsal` →
`interview-refresher`) against a synthetic engagement on the real released
`v0.8.1` binary — the first time these skills were ever actually executed
rather than just read.

Three of the four type skill files reference `AGENTS.md` where they clearly
mean `CONTENT.md`:

```
.claude/skills/interview-setup/SKILL.md:24
  "Record company and role in AGENTS.md's Domain Context placeholders."

.claude/skills/interview-refresher/SKILL.md:26
  "...see AGENTS.md's "Company stack in the wiki" section..."

.claude/skills/interview-rehearsal/SKILL.md:70
  "...it runs against a real exercises/NNN-*/ module (per AGENTS.md's convention)."

.claude/skills/interview-rehearsal/SKILL.md:78
  "A per-attempt debrief page also belongs in the CMS per AGENTS.md..."
```

`"Company stack in the wiki"` and the `exercises/NNN-*/` convention are real
subsection titles inside `CONTENT.md`'s own `## Type: candidate-interview`
section (task 013). `AGENTS.md` is smaqit's own file — `TYPE.md`'s own "not
part of the deployed state" list explicitly excludes it from what a typed
install ships. These are near-certainly a leftover mix-up from when the
fixture repo (`magnificah-interview-kb`) was itself a smaqit project with
its own `AGENTS.md`; the reference carried through into the shipped type
content unnoticed because nobody had run these skills for real before this
validation.

`interview-setup`'s occurrence is a slightly different case: "Domain Context
placeholders" isn't a `CONTENT.md` concept at all (it's specific to *this*
repo's own smaqit-authored `AGENTS.md`), so redirecting it to `CONTENT.md`
doesn't make sense either — the instruction should just be dropped, since
`docs/company/company-overview.md` and `docs/interview-prep/role-research.md`
already capture the same company/role facts a few steps later in the same
skill's own flow.

`interview-round/SKILL.md` has no such reference — confirmed via grep across
all four type skill files.

## Issue Triage Context

**Mode:** Auto
**Technologies:** None
**Platforms/Environments:** None
**Features/Integrations:** `candidate-interview` type skills (`interview-setup`, `interview-refresher`, `interview-rehearsal`)
**Versions/Constraints:** Regression present since v0.8.0 (task 013); targets v0.8.2 (PATCH — content/documentation correction only, no behavior change)

## Design Decisions

- **Three straightforward redirects**: `interview-refresher` line 26 and `interview-rehearsal` lines 70/78 change `AGENTS.md` → `CONTENT.md`, wording otherwise unchanged (the section titles they reference already exist verbatim in `CONTENT.md`'s Type section).
- **Drop, don't redirect, `interview-setup` line 24**: remove the "Record company and role in `AGENTS.md`'s Domain Context placeholders" sentence entirely rather than pointing it at `CONTENT.md` (which has no matching concept). The engagement facts are already captured properly via `docs/company/company-overview.md` and `docs/interview-prep/role-research.md`, created later in the same skill's own steps — this sentence is redundant as well as wrong.
- **Source of truth stays the fixture-derived files in `scaffold/types/candidate-interview/tree/`** — these are edited directly (not regenerated from the fixture repo, which is out of scope per the prior session's "the fixture is disposable, not golden" ruling).
- **Scope**: text-only correction to installed skill instructions. No Go code changes, no new tests needed beyond confirming the fixed text installs correctly (a plain grep-based check is sufficient — this isn't logic the Go installer touches).

## Implementation Steps

1. Edit `scaffold/types/candidate-interview/tree/.claude/skills/interview-refresher/SKILL.md` line 26: `AGENTS.md` → `CONTENT.md`.
2. Edit `scaffold/types/candidate-interview/tree/.claude/skills/interview-rehearsal/SKILL.md` lines 70 and 78: `AGENTS.md` → `CONTENT.md`.
3. Edit `scaffold/types/candidate-interview/tree/.claude/skills/interview-setup/SKILL.md` line 24: remove the "Record company and role in `AGENTS.md`'s Domain Context placeholders" sentence.
4. Confirm no other `AGENTS.md` references remain in any of the four type skill files (`grep -rn "AGENTS.md" scaffold/types/candidate-interview/tree/.claude/skills/`).
5. `make test` and `make smoke-test` both green (unaffected by this change, but confirm no regression).
6. Release through the normal PR flow, expected v0.8.2 (PATCH).
7. Re-verify live: download the newly-released binary, fresh `init --type candidate-interview`, confirm the installed skill files read correctly.

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
text-only correction inside this repo's own type-content files, with no
external dependency involved.

## Acceptance Criteria

- [x] All four `AGENTS.md` mentions across the type's skill files are corrected (3 redirected to `CONTENT.md`, 1 removed).
- [x] No remaining `AGENTS.md` reference anywhere under `scaffold/types/candidate-interview/tree/.claude/skills/`. Verified via grep on both the source tree and a locally-built install.
- [x] `make test` and `make smoke-test` both pass.
- [ ] Released via PR flow as v0.8.2.
- [ ] Verified live against the real released v0.8.2 binary. *(Post-merge step, same pattern as tasks 013/014 — done once v0.8.2 is actually released.)*

## Findings

**Implementation approach:**
- Three straight text substitutions (`AGENTS.md` → `CONTENT.md`) in
  `interview-refresher` and `interview-rehearsal`, plus removal of the one
  sentence in `interview-setup` that had no `CONTENT.md` equivalent to
  redirect to.
- Verified the fix two ways: a source-tree grep confirming zero remaining
  `AGENTS.md` mentions under the type's skills directory, and a live
  install from a freshly built dev binary confirming the installed files
  read correctly — including noticing the fixed lines now read consistently
  with other, already-correct `CONTENT.md` references already present
  elsewhere in the same files (`interview-refresher` line 11,
  `interview-rehearsal` line 115), which is a good sign this was a genuine
  isolated slip rather than a systematic pattern needing broader rework.

**Decisions made:**
- Dropped rather than redirected `interview-setup`'s "Domain Context
  placeholders" instruction, per the task's own Design Decisions — the
  facts it asked to record are already captured properly a few steps later
  via `docs/company/company-overview.md` and `role-research.md`, so there
  was nothing meaningful to redirect it to.

**Blockers encountered:**
- None.

**Follow-up identified:**
- Acceptance criterion 5 (verify against the real released v0.8.2 binary)
  is a post-merge step, same pattern as tasks 013/014.
- No further live skill-run findings expected from this specific defect
  class — the broader Layer 2 e2e validation (task 013/014's follow-on)
  otherwise passed cleanly; this was its only finding.

## Files to Create / Modify

| File | Action |
|------|--------|
| `scaffold/types/candidate-interview/tree/.claude/skills/interview-setup/SKILL.md` | Modify |
| `scaffold/types/candidate-interview/tree/.claude/skills/interview-refresher/SKILL.md` | Modify |
| `scaffold/types/candidate-interview/tree/.claude/skills/interview-rehearsal/SKILL.md` | Modify |
| `CHANGELOG.md` | Modify |

## Notes

Found as a side effect of a live "Layer 2" validation exercise (running all
four type skills against a synthetic fictional engagement — company
Northwind Analytics, candidate Jordan Ellis — in a fresh v0.8.1 sandbox),
requested after v0.8.1 shipped to prove the type mechanism can regenerate a
coherent engagement from scratch rather than treating the original
`magnificah-interview-kb` fixture repo as a golden copy to maintain. That
validation otherwise passed cleanly: `ac-index check`/`ac-links
check`/`ac-classify sweep` all clean, classification ratings matched the
type's stated defaults, and the rehearsal → refresher → exam mechanics
(including the "still-open gap carries forward" behavior) worked correctly
end to end. This was the only defect the run surfaced.

Child tasks inherit their active parent's branch, worktree, and workflow mode. Only a standalone or parent task owns Git lifecycle cleanup.
