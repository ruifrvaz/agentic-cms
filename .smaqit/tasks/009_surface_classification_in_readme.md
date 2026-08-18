---
status: PR Open
created: "2026-08-18"
mode: Assisted
started: "2026-08-18"
pr: 7
---

# Surface classification feature in README

## Description

Give the C0-C3 confidentiality classification feature (shipped in task 005,
v0.5.0) more prominence in `README.md`. Today it's a single low-level
Features bullet ("One detector, many thin callers") describing the
enforcement mechanism, not the concept or why it matters. This is a
genuinely distinctive feature for an agentic CMS and deserves headline
visibility: mention the CIA triad explicitly, and give an impactful
explanation of why classification matters specifically for a system where
an agent autonomously synthesizes raw sources into a compounding wiki.

## Issue Triage Context

**Mode:** Skip
**Technologies:** None
**Platforms/Environments:** None
**Features/Integrations:** None
**Versions/Constraints:** None

## Design Decisions

- **Elevate the Features bullet, move implementation detail to a new
  dedicated section** — replace the existing "One detector, many thin
  callers" bullet (which only describes enforcement plumbing) with a
  punchier headline bullet ("Confidentiality-aware by design") that links
  to a new `## Classification` section, mirroring the existing
  "Three-layer content model" → "See [How it works](#how-it-works)"
  pattern already used in this README.
- **New `## Classification` section placed after "How it works"** — both
  are core conceptual sections; keeps the doc flow (concept → concept →
  compatibility → getting started → commands) intact.
- **Explicitly name the CIA triad** (Confidentiality, Integrity,
  Availability) and clarify agentic-cms uses its confidentiality axis —
  matches `CONTENT.md`'s own framing ("the standard CIA confidentiality
  axis") so the README doesn't invent new terminology.
- **"Why it matters" grounded in the mechanism, not a specific incident** —
  explains the risk generically (agent-driven synthesis across `raw/` into
  a compounding `wiki/` can silently propagate sensitive detail into
  summaries/exports) rather than referencing this repo's own internal
  history, since that's not appropriate content for a public README.
- **Keep the C0-C3 table and "One detector, many thin callers" detail**,
  just relocated into the new section rather than deleted — the
  implementation-detail content is still accurate and worth keeping for
  readers who want the mechanics.

## Implementation Steps

1. In `README.md`'s `## Features` section, replace the "One detector, many
   thin callers" bullet with a headline "Confidentiality-aware by design"
   bullet mentioning the CIA triad and linking to `#classification`.
2. Insert a new `## Classification` section after `## How it works` (before
   `## Compatibility`): a C0-C3 table (Public/Internal/Confidential/
   Restricted with examples, matching `CONTENT.md`'s rubric), a "Why it
   matters" paragraph explaining the risk for agent-driven synthesis, and
   the "One detector, many thin callers" enforcement-mechanism paragraph
   (moved from the old Features bullet).
3. Proofread against `scaffold/tree/CONTENT.md`'s Classification section
   for accuracy (ratchet-up-only rule, bleed rule for C2+ index/log
   entries, `ac-classify` as floor-only detector).
4. Run `make smoke-test` and `go test ./...` to confirm the docs-only
   change doesn't affect anything (no code touched, but confirms nothing
   else regressed).

## Known Issues Triage

Triage skipped — explicitly marked `Skip` in task Issue Triage Context
(docs-only change, no third-party tools or platforms implicated).

## Acceptance Criteria

- [x] README's Features section has a headline bullet mentioning
      classification and the CIA triad
- [x] A dedicated `## Classification` section explains the C0-C3 levels
      (with examples) and why enforcement matters for agent-driven content
      synthesis
- [x] Content is accurate against `scaffold/tree/CONTENT.md`'s
      Classification rubric (ratchet-up-only, bleed rule, floor-only
      detector)
- [x] `make smoke-test` and `go test ./...` remain passing

## Findings

**Implementation approach:**
- Replaced the "One detector, many thin callers" Features bullet with a
  headline "Confidentiality-aware by design" bullet naming the CIA triad
  and linking to `#classification`.
- Added a new `## Classification` section after "How it works": a C0-C3
  table (levels, meanings, examples matching `CONTENT.md` verbatim), a
  "Why it matters" paragraph, and the enforcement-mechanics paragraph
  (moved down from the old Features bullet, wording unchanged).

**Decisions made:**
- "Why it matters" is grounded in the general mechanism (agent-driven
  synthesis across `raw/` into a compounding, cross-linked `wiki/` can
  silently propagate sensitive detail into summaries/exports), not this
  repo's own internal incident history — not appropriate content for a
  public README.
- Kept the C0-C3 table and enforcement-mechanics detail rather than
  cutting it for brevity — still accurate and useful for readers who want
  the mechanics, just relocated out of the Features section.

**Blockers encountered:**
- GitHub push access broke mid-task-start (403 on both a plain `git push`
  and an explicit `gh`-credential-helper push, despite `gh api` reporting
  full push/admin rights) — a known recurring pattern where this user's
  managed PAT rotates/expires. Resolved by the user re-running their
  token-refresh script; retried successfully.

**Follow-up identified:**
- None — task scoped to README only, fully delivered.

## Files to Create / Modify

| File | Action |
|------|--------|
| README.md | Modify |

## Notes

Docs-only change, no code touched. Triage set to Skip — no third-party
tools or platforms implicated by a README edit.

Child tasks inherit their active parent's branch, worktree, and workflow
mode. Only a standalone or parent task owns Git lifecycle cleanup.
