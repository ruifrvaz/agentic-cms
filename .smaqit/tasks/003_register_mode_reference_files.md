---
status: Not Started
created: "2026-08-15"
---

# Register mode for reference files

## Description

Real usage of an installed CMS surfaced two distinct file scenarios that the
schema currently treats as one: (1) files kept in the repo mostly for archival
and reference purposes with rare consultation (e.g. an invoice PDF — store it,
index it, not much else), and (2) files dropped in the repo to be processed and
converted to markdown (the existing `content-import` pipeline). Today only
scenario 2 exists as a workflow, so reference files either go through a full
conversion they don't need or sit in `raw_uningested` nagging forever.

The "drop and sit until processed" flow already exists structurally: `raw/` is
the drop point (its README says "Drop source material here"), and
`ac-inventory`'s `raw_uningested` field — raw files not referenced by any
docs/wiki page — is the pending queue, surfaced by `content-list`. This task
adds the missing lightweight exit from that queue: a **Register** mode on
`content-import` that catalogs a file (wiki/sources sidecar + index entry +
log entry) without converting it to markdown.

## Issue Triage Context

**Mode:** Auto
**Technologies:** None
**Platforms/Environments:** None
**Features/Integrations:** None
**Versions/Constraints:** None

## Design Decisions

- **Drop point is `raw/`, unchanged:** both scenarios share the same landing
  zone (user-curated, immutable once landed). Subfolder organization such as
  `raw/invoices/` already works — the `raw_uningested` scan is recursive.
- **Register = catalog entry only:** create the `wiki/sources/<slug>.md`
  sidecar via `ac-page new source ... --raw-path raw/<path>`, then
  `ac-index add sources` + `ac-log append register`. No docs/ page, no
  conversion, no entity/concept work.
- **The sidecar clears the queue by construction:** it references the raw
  path, so the file drops out of `raw_uningested` automatically — the sidecar
  IS the "processed" marker. Zero changes to `ac-inventory`.
- **The sidecar is the findability layer:** `ac-search`/`content-query` work
  over markdown; a binary PDF is invisible to them. The sidecar's frontmatter
  and one-liner ("Invoice, Acme, June 2026, €1,200") are the greppable trace
  that "rare consultation" needs.
- **`wiki/sources/` semantics broaden:** from "one summary page per ingested
  raw source" to "one page per registered raw file" — reference-only or fully
  ingested. Same template, same category. Register-only sidecars carry
  `tags: [reference]` so lint/list can distinguish them if ever needed.
- **New log operation `register`:** added to CONTENT.md's operations list so
  the journal distinguishes registrations from full imports (`ac-log append`
  takes free-form operation names — no script change).
- **Upgrade path is free:** a registered file can be fully ingested later —
  the sidecar/index/log half already exists, so only the conversion half runs.
- **No toolkit script changes:** the design deliberately reuses `ac-page new
  source`, `ac-index add`, `ac-log append`, and the existing `raw_uningested`
  logic unchanged.

## Implementation Steps

1. Update `scaffold/tree/CONTENT.md`: document the two file scenarios
   (register vs. ingest) in the Operations section, broaden the
   `wiki/sources/` definition to "one page per registered raw file", add
   `register` to the log operations list, and note `raw/` subfolder
   organization as supported.
2. Restructure `scaffold/tree/.claude/skills/content-import/SKILL.md` into
   three modes: **Register** (new — steps: secure raw copy if not already in
   `raw/`, create sources sidecar with `tags: [reference]`, index, log,
   verify with `ac-index check`), **Ingest** (existing full pipeline,
   unchanged), **Brownfield sweep** (existing — its file→topic mapping
   proposal now also proposes register-vs-ingest treatment per file).
3. Update `scaffold/tree/.claude/agents/content-importer.md` so sweep
   processing applies the approved per-file treatment (register or ingest).
4. Update `scaffold/tree/raw/README.md` with one line on the triage flow:
   dropped files are either registered (reference) or imported (converted).
5. Update `scaffold/tree/.claude/skills/content-list/SKILL.md`: the
   `raw_uningested` line now suggests triage via `content-import` (register
   or ingest), not just import.
6. Run `make smoke-test` and `go test ./...` — no new scaffold files, so the
   installer diff should be unaffected. Manually exercise: drop a dummy file
   into a sandbox project's `raw/`, confirm it appears in `raw_uningested`,
   register it, confirm it leaves the queue and `ac-index check` stays clean.

## Known Issues Triage

[Populated by smaqit.task-start via smaqit.utils.triage-issues. Do not edit manually.]

## Acceptance Criteria

- [ ] `CONTENT.md` documents both file scenarios, the broadened
      `wiki/sources/` semantics, and the `register` log operation
- [ ] `content-import` has a Register mode: sources sidecar + index + log,
      no markdown conversion, no docs/ page
- [ ] A registered file drops out of `raw_uningested` with no `ac-inventory`
      change (verified manually in a sandbox)
- [ ] Register-only sidecars carry `tags: [reference]`
- [ ] Brownfield sweep proposes register-vs-ingest treatment per file
- [ ] No `.agentic-cms/bin/` script is modified
- [ ] `make smoke-test` and `go test ./...` remain passing

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
| scaffold/tree/CONTENT.md | Modify |
| scaffold/tree/.claude/skills/content-import/SKILL.md | Modify |
| scaffold/tree/.claude/agents/content-importer.md | Modify |
| scaffold/tree/raw/README.md | Modify |
| scaffold/tree/.claude/skills/content-list/SKILL.md | Modify |

## Notes

Stacks on task 002 (draft content state, PR #2): touches the same
`CONTENT.md` and `content-list` files and assumes 002's conventions are
merged. Start after 002 completes. Independent of task 004 (archive), though
both edit `CONTENT.md` and `content-list` — whichever merges second rebases
trivially.

Child tasks inherit their active parent's branch, worktree, and workflow
mode. Only a standalone or parent task owns Git lifecycle cleanup.
