---
status: Completed
created: "2026-08-14"
mode: Assisted
started: "2026-08-15"
completed: "2026-08-15"
---

# Draft content state

## Description

Add a lower-ceremony "draft" state to the content schema for work-in-progress
material that isn't ready to become first-class content yet.

Today the schema only has two homes for anything an agent writes: `raw/`
(immutable — never touched again after landing) and `docs/` (organized,
wiki-wired, meant to already be clean and readable). There's no place for a
brainstorm or half-formed note that will be iterated on repeatedly before
it's ready to graduate — `raw/`'s immutability rule makes it the wrong fit
for anything actively edited, and filing straight into `docs/` immediately
pulls it into `wiki/index.md`/`wiki/log.md` bookkeeping before it's settled.
This gap surfaced from real usage of an installed project's CMS, where the
user wanted to capture an evolving idea before committing it as a proper
content item.

## Design Decisions

- **New optional frontmatter field:** `status: draft | final`, defaulting to
  `final` when absent — so every already-installed project and existing page
  stays valid with no migration, consistent with the installer's own
  non-destructive/idempotent principle.
- **Location:** drafts live at `docs/<topic>/drafts/<item>.md`, using the
  existing `doc.md` template with `status: draft` set — near their eventual
  home, so "promoting" one is a small move rather than a cross-tree jump.
- **Not wired into the wiki while draft:** `wiki/index.md` and `wiki/log.md`
  are skipped for draft items until promotion. `content-lint` must not flag
  an unlinked draft as an orphan page; `content-list` should instead surface
  drafts in their own section so they stay visible without being treated as
  finished content.
- **Promotion = a small, explicit step:** flip `status: draft` → remove the
  field (or set `final`), move the file from `docs/<topic>/drafts/` up to
  `docs/<topic>/`, then run the normal `content-new-item` wiki-wiring
  (index entry, log entry).
- **Resolved: extend `content-new-item`, not a new skill.** Add "Draft" and
  "Promote" modes alongside its existing flow, mirroring the "two modes"
  pattern `content-import` already uses (single-source vs. brownfield
  sweep). Keeps the nine-skill footprint; the promote path reuses
  `content-new-item`'s existing register+log tail rather than duplicating
  it in a new skill.
- **Toolkit-level findings (from post-migration Discovery — the toolkit
  didn't exist yet when this task was first drafted):**
  - `.agentic-cms/bin/ac-index`'s `check` command globs `docs/**/*.md`
    recursively for `unindexed_pages`, which would flag every draft
    forever and break the `"clean": true` gate every write-skill depends
    on. Needs a one-line exclusion for `docs/*/drafts/*.md`.
  - `.agentic-cms/bin/ac-inventory` has no visibility into `drafts/` at
    all — it must start reporting draft pages (path + topic) for
    `content-list` to have anything to show.
  - No change needed in `content-lint`: its orphan check only walks
    `ac-index list`, which never contains drafts by construction — once
    `ac-index` is fixed, drafts are structurally unreachable by that
    check. `ac-links check` also needs no change — checking links *from*
    drafts is desirable, not a false positive.
  - Promotion refuses on a destination conflict (mirrors `ac-page new`'s
    existing refuse-overwrite behavior) rather than attempting a merge.

## Implementation Steps

1. Update `scaffold/tree/CONTENT.md`: document the `status` field, the
   `docs/<topic>/drafts/` convention, and the promotion flow, alongside the
   existing layer/type/frontmatter definitions.
2. Update `scaffold/tree/.agentic-cms/templates/doc.md` to show the optional
   `status:` field (commented example, not a default value, so untouched
   installs don't silently gain the field).
3. Fix `scaffold/tree/.agentic-cms/bin/ac-index` (`check` command, the
   `unindexed_pages` glob around line 100): exclude `docs/*/drafts/*.md` so
   drafts never trip drift detection. Do this before step 5 — it depends on
   it.
4. Fix `scaffold/tree/.agentic-cms/bin/ac-inventory`: add a `drafts` field
   to its JSON output (path + topic per draft page). Parallel with step 3;
   step 6 depends on it.
5. Restructure `scaffold/tree/.claude/skills/content-new-item/SKILL.md` into
   a "Modes" section: **New item** (existing flow, unchanged), **Draft**
   (same through page creation, `status: draft` set, but skip the
   index/log register step), **Promote** (move the file out of `drafts/`,
   clear `status:`, then run the normal register+log tail against the new
   path). Depends on steps 3–4.
6. Update `scaffold/tree/.claude/skills/content-list/SKILL.md` to report
   drafts in their own section, sourced from `ac-inventory`'s new `drafts`
   field. Depends on step 4.
7. Add one clarifying sentence to
   `scaffold/tree/.claude/skills/content-lint/SKILL.md` noting drafts are
   excluded from orphan detection by construction (no functional change —
   see Design Decisions) so this isn't mistaken for a gap later.
8. Consider a one-line mention in the root `README.md`'s "How it works" or
   skills table if the addition is significant enough to surface there.
9. Run `make smoke-test` and `go test ./...` — confirm the installer diff
   against `scaffold/tree/` and existing greenfield/idempotency/brownfield
   checks still pass unaffected (no new files ship in the base scaffold;
   drafts are created per-project on demand). Manually exercise the
   toolkit too: create a draft via `ac-page new`, confirm `ac-index check`
   stays clean, confirm `ac-inventory` lists it, promote it, confirm
   `ac-index check` is still clean post-promotion.

## Known Issues Triage
**Triaged:** 2026-08-15
**Tools searched:** none — no third-party tools identified
**Result:** Clear

Task scope is entirely internal (schema docs, skill instructions, and stdlib
bash/python3 edits to the project's own toolkit) — no third-party
product, library, or service dependency is introduced or at risk. GitHub
issue search is not applicable.

### Unresolvable Tools
- (none)

### Omitted Tools
- (none)

### Search Warnings
- (none)

## Acceptance Criteria

- [x] `CONTENT.md` documents `status: draft | final` (default `final`) and
      the `docs/<topic>/drafts/` convention, including the promotion flow
- [x] `doc.md` template shows the optional `status:` field
- [x] `ac-index check` excludes `docs/*/drafts/*.md` from `unindexed_pages`
- [x] `ac-inventory` reports draft pages (path + topic)
- [x] `content-new-item` supports creating a draft (no wiki wiring) and
      promoting an existing draft to first-class content
- [x] `content-lint` does not flag drafts as orphans or stale content
      (verify this holds once `ac-index` is fixed — expected by
      construction, not a separate content-lint code change)
- [x] `content-list` reports drafts separately from first-class content
- [x] `make smoke-test` and `go test ./...` remain unaffected and passing

## Findings

**Implementation approach:**
- `status:` shipped as a live, first-class frontmatter field (`status: final`
  default in `doc.md`, substituted via a new `{{STATUS}}` placeholder) rather
  than the originally-planned commented-out example — see Decisions.
- `ac-page` gained `--status draft|final` on `new` (validated, defaults to
  `final`) and a new `promote <src> <dest>` subcommand (refuse-on-conflict,
  moves the file, sets `status: final`, touches `updated:`) so both draft
  creation and promotion are one deterministic toolkit call each, with no
  hand-editing of frontmatter.
- `ac-index check`'s `unindexed_pages` glob excludes `docs/*/drafts/*.md`;
  `ac-inventory` gained a per-topic `drafts` field. `content-lint` needed no
  functional change — its orphan check only walks `ac-index list`, which
  never contains drafts by construction.
- `content-new-item` restructured into New item / Draft / Promote modes,
  mirroring `content-import`'s existing two-modes pattern.

**Decisions made:**
- Extend `content-new-item` rather than add a new skill (matches the task's
  original lean; confirmed during planning).
- `status:` is a live template field, not a commented example as originally
  planned — the commented-example rationale ("untouched installs don't
  silently gain the field") doesn't actually hold: `init`/`update` never
  overwrite an existing template file regardless, so there was no real
  migration risk to avoid. Live keeps it consistent with every other
  frontmatter field and matches `CONTENT.md`'s own rule that mechanical
  operations belong in the toolkit, not hand-edited — which is also why
  `ac-page` gained `--status`/`promote` instead of the originally-planned
  manual frontmatter edits during Draft/Promote. Both changes were raised by
  the user during review and implemented after confirmation.
- Promotion refuses on a destination conflict (mirrors `ac-page new`'s
  existing refuse-overwrite behavior) rather than attempting a merge.
- `ac-links check` needs no change — checking links *from* drafts is
  desirable, not a false positive.

**Blockers encountered:**
- None.

**Follow-up identified:**
- None beyond what's already in `CONTENT.md`/README's existing roadmap notes.

## Notes

Kept intentionally generic — this tracker describes the CMS schema/tooling
itself, not any particular installed project's content.

Child tasks inherit their active parent's branch, worktree, and workflow
mode. Only a standalone or parent task owns Git lifecycle cleanup.
