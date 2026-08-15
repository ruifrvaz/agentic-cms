# Draft content state

**Status:** Not Started
**Mode:** Assisted
**Created:** 2026-08-14

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
  (index entry, log entry). Open question for the implementer: whether this
  is a mode/flag on `content-new-item` or warrants its own short skill —
  lean toward extending `content-new-item` first to keep the skill surface
  small, per the project's existing nine-skill footprint.

## Implementation Steps

1. Update `scaffold/tree/CONTENT.md`: document the `status` field, the
   `docs/<topic>/drafts/` convention, and the promotion flow, alongside the
   existing layer/type/frontmatter definitions.
2. Update `scaffold/tree/.agentic-cms/templates/doc.md` to show the optional
   `status:` field (commented example, not a default value, so untouched
   installs don't silently gain the field).
3. Update `scaffold/tree/.claude/skills/content-new-item/SKILL.md` to support
   starting an item as a draft, and to support promoting an existing draft to
   first-class content.
4. Update `scaffold/tree/.claude/skills/content-lint/SKILL.md` so draft items
   are excluded from orphan-page and staleness checks.
5. Update `scaffold/tree/.claude/skills/content-list/SKILL.md` to report
   drafts in their own section, separate from first-class content.
6. Consider a one-line mention in the root `README.md`'s "How it works" or
   skills table if the addition is significant enough to surface there.
7. Run `make smoke-test` — confirm the installer diff against
   `scaffold/tree/` and existing greenfield/idempotency/brownfield checks
   still pass unaffected (no new files ship in the base scaffold; drafts are
   created per-project on demand).

## Known Issues Triage

[Populated by smaqit.task-start via smaqit.utils.triage-issues. Do not edit manually.]

## Acceptance Criteria

- [ ] `CONTENT.md` documents `status: draft | final` (default `final`) and
      the `docs/<topic>/drafts/` convention, including the promotion flow
- [ ] `doc.md` template shows the optional `status:` field
- [ ] `content-new-item` supports creating a draft (no wiki wiring) and
      promoting an existing draft to first-class content
- [ ] `content-lint` does not flag drafts as orphans or stale content
- [ ] `content-list` reports drafts separately from first-class content
- [ ] `make smoke-test` and `go test ./...` remain unaffected and passing

## Notes

Kept intentionally generic — this tracker describes the CMS schema/tooling
itself, not any particular installed project's content.

Child tasks inherit their active parent's branch, worktree, and workflow
mode. Only a standalone or parent task owns Git lifecycle cleanup.
