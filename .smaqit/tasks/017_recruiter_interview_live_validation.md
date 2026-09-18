---
status: In Progress
created: "2026-09-18"
mode: Assisted
started: "2026-09-18"
---

# Live e2e validation: recruiter-interview candidate workflow

## Description

Task 016 shipped `recruiter-interview` (v0.9.0) with zero changes to the
installer mechanism, verified mechanically (`make test`, `make smoke-test`,
a manual `diff -rq` against the fixture using a local dev build). What it
did **not** get was the live skill-run validation tasks 013/014/015
established as this project's actual acceptance bar for a content type:
actually running the installed skills against a synthetic engagement and
checking the result, independently re-verified rather than trusted from a
single pass. That methodology found real, otherwise-invisible bugs twice
in a row (task 014's `{{DATE}}` corruption, task 015's `AGENTS.md`
mix-up) — both in `candidate-interview`, both missed by mechanical checks
alone. `recruiter-interview` has not had this validation yet.

This task runs it, with the user's explicit focus on the candidate
workflow specifically: intake, scoring, and the slug-discipline rule that
is this type's sharpest design property (a candidate's real name must
never appear outside their own C2 pages — see the "why not X" design
review that closed task 016).

## Issue Triage Context

**Mode:** Auto
**Technologies:** None
**Platforms/Environments:** Linux
**Features/Integrations:** `recruiter-setup`, `recruiter-round`, `recruiter-candidate`, `recruiter-decision` skills; `.agentic-cms/scripts/` toolkit (`ac-page`, `ac-index`, `ac-links`, `ac-classify`)
**Versions/Constraints:** Validates the real released v0.9.0 binary (task 016)

## Design Decisions

- **Two layers, same shape as `candidate-interview`'s validation**:
  Layer 1 (deterministic, done directly, no agent) re-confirms the real
  `v0.9.0` release installs cleanly; Layer 2 (a subagent, driving an
  actual synthetic requisition) is the real test.
- **Fresh fictional data, not reused from `candidate-interview`'s own
  validation run** — a new company/role, so the two validations stay
  independent and neither leaks assumptions into the other.
- **Candidate workflow is the focus**, per explicit user direction: intake
  of 3 fictional candidates, scoring through both rounds, a reference
  check, and a decision — not just a shallow "does setup work" pass.
- **Independent re-verification is mandatory**, not optional — re-run
  `ac-index check`/`ac-links check`/`ac-classify sweep` directly, and
  specifically grep the roster/index/log for any candidate's real
  (fictional) name leaking outside their two permitted C2 pages. Trusting
  a subagent's self-report alone is exactly what let task 014's and 015's
  bugs previously slip past a less careful pass.
- **A defect found here is not fixed inline** — it becomes its own
  follow-up task (018+), mirroring exactly how 014 and 015 emerged from
  013's validation. This task's job is to find and report, not
  necessarily to fix.
- **No code changes are expected from this task itself** — if the
  validation is clean, this task closes with no PR (Rule 4's "owner with
  no committed code changes" path); if it finds something, a new task is
  filed and this task still closes once that's tracked, same as how task
  016 closed before its own live validation ran.

## Implementation Steps

1. **Layer 1**: download the real `v0.9.0` binary via `install.sh`
   (`AGENTIC_CMS_VERSION=v0.9.0`), fresh `init --type recruiter-interview`
   into a temp sandbox, confirm `TYPE.md` intact (no `{{DATE}}`
   corruption), no `AGENTS.md` references anywhere in the installed type
   skills, and the 34-template smoke (create/index/link/classify) clean.
2. **Layer 2**: spawn a subagent with a fully self-contained prompt (fresh
   fictional company/role, not reused from prior validations) to actually
   read and follow each `recruiter-*` SKILL.md in order:
   `recruiter-setup` → `recruiter-round` (both rounds) →
   `recruiter-candidate` (intake 3 candidates, score both rounds, one
   reference check) → `recruiter-decision` (compare, decide, offer,
   archive). Have it report every file created/modified, the final Rounds
   table and candidate roster, every classification rating with a reason,
   and final `ac-index check`/`ac-links check`/`ac-classify sweep` output.
3. **Independently re-verify**: re-run the three health checks directly
   against the subagent's sandbox; spot-read at least one scorecard and
   the decision record for real substance (not placeholder-shaped
   content); grep the roster, `wiki/index.md`, and `wiki/log.md` for any
   candidate's real name — must find none outside their own
   `wiki/entities/cNNN.md` and `docs/candidates/cNNN-profile.md`.
4. **If clean**: write Findings, close this task with no PR.
5. **If a defect is found**: file it as a new task (018+) with the same
   rigor as 014/015 (confirm the regression class, write a test that
   fails without the fix and passes with it, fix, release), then close
   this task noting the follow-up.

## Known Issues Triage
**Triaged:** 2026-09-18
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
validation exercise against this repo's own released binary and skills;
no external dependency is implicated.

## Acceptance Criteria

- [ ] Layer 1 (real `v0.9.0` binary, fresh typed install) confirmed clean.
- [ ] Layer 2 live skill run completed through all four `recruiter-*`
      skills, covering setup, both rounds, 3-candidate intake and scoring,
      a reference check, and a decision.
- [ ] Independently re-verified — `ac-index check`/`ac-links
      check`/`ac-classify sweep` re-run directly (not taken from the
      subagent's report), at least two pages spot-read for real content,
      and a roster/index/log grep confirms zero candidate-name leakage
      outside the two permitted C2 pages per candidate.
- [ ] Findings recorded; any defect found is filed as its own follow-up
      task rather than fixed silently inline.

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
| (none expected — validation only; a follow-up task's own files would list here if a defect requires a fix) | — |

## Notes

Directly requested by the user immediately after task 016 closed:
"the goal is to install locally and run a subagent to test out creating
new candidates." This task exists to track that request as first-class
work, matching how tasks 014/015 were tracked, rather than letting a
live-validation pass happen informally in conversation with no durable
record.

Child tasks inherit their active parent's branch, worktree, and workflow mode. Only a standalone or parent task owns Git lifecycle cleanup.
