---
status: Completed
created: "2026-08-15"
mode: Assisted
started: "2026-08-15"
completed: "2026-08-16"
---

# Archive content lifecycle

## Description

Add an "archived" terminal state to the content lifecycle for first-class
items that are no longer active but must remain retrievable. Task 002 gave
the schema a pre-first-class state (`status: draft`, `docs/<topic>/drafts/`);
this task adds the post-first-class mirror: `status: archived` and
`docs/<topic>/archive/`, completing the lifecycle
draft → final → archived.

Unlike a draft — which was never in the wiki — an archived page *was*
first-class: it has an index entry, log history, and possibly inbound links.
Archiving therefore isn't invisibility but re-filing: the item moves to the
topic's `archive/` folder, its index entry moves to an Archived section
(archive ≠ delete — `wiki/index.md` is the retrieval system, and
`content-query` should still find archived material and know it's archived),
and the transition is logged.

## Issue Triage Context

**Mode:** Auto
**Technologies:** None
**Platforms/Environments:** None
**Features/Integrations:** None
**Versions/Constraints:** None

## Design Decisions

- **Status enum extends:** `status: draft | final | archived`. Absent still
  means `final` — no migration for installed projects, same principle as 002.
- **Location:** `docs/<topic>/archive/<item>.md` — per-topic, mirroring
  `drafts/`, so un-archiving is a small move. `raw/` never archives
  (immutability already covers it); `exports/` never archives (derived
  artifacts are regenerable).
- **Toolkit: `ac-page archive <src> <dest>`**, sibling of the `promote`
  subcommand 002 added: move + set `status: archived` + touch `updated:`,
  refusing on destination conflict.
- **Archived items stay indexed, re-filed under an "archived" index
  section** (`ac-index remove` + `ac-index add archived ...`) — the key
  difference from drafts. Because they remain indexed, `ac-index check`
  needs no new exclusion glob for `archive/` (verify: `check`'s
  `docs/**/*.md` glob finds them, and the index still lists them → clean).
- **Un-archive = `ac-page promote`** from `archive/` back up: mechanically
  identical (move + `status: final` + touch), then re-file the index entry
  back to its normal section. No new toolkit command.
- **New log operation `archive`** (and `unarchive`), added to CONTENT.md's
  operations list; `ac-log append` takes free-form names, no script change.
- **Inbound-link cleanup is a forcing function:** after the move,
  `ac-links check` flags now-broken links from active pages — desirable;
  the archive flow ends by fixing or removing those references.
- **`content-lint` orphan check exempts the archived index section:**
  archived pages intentionally lose inbound links, so "no inbound links" is
  their expected steady state, not a defect.
- **`ac-inventory` reports `archived`** per topic (count + files) alongside
  the `drafts` field from 002; the `topics` item glob (`docs/<topic>/*.md`)
  is non-recursive so archived items never inflate active item counts.
- **Skill surface: an Archive mode on `content-new-item`**, mirroring
  Promote in reverse — keeps the nine-skill footprint, same modes pattern
  as 002.
- **Skill renamed `content-new-item` → `content-manage-item`** (user
  decision during implementation): with Draft/Promote/Archive modes the old
  name became misleading — the skill owns create + the item lifecycle, not
  just creation. All cross-references updated (root README, CONTENT.md,
  content-new, content-list, content-add-notes, scaffold_test.go). Upgraded
  installs keep a stale `content-new-item/` skill dir until `update` learns
  scaffold diffing/renames — recorded in the root README roadmap.

## Implementation Steps

1. Update `scaffold/tree/CONTENT.md`: extend the `status` enum, document the
   `docs/<topic>/archive/` convention, the archive/un-archive flows, the
   archived index section, and add `archive`/`unarchive` to the log
   operations list.
2. Add `archive <src> <dest>` to `scaffold/tree/.agentic-cms/bin/ac-page`
   (move + `status: archived` + touch `updated:`, refuse on existing
   destination), update its usage header and
   `scaffold/tree/.agentic-cms/bin/README.md`.
3. Verify `scaffold/tree/.agentic-cms/bin/ac-index` handles an `archived`
   category/section in `index.md` for `add`/`remove`/`list`/`check`; adjust
   only if the section handling isn't already generic.
4. Update `scaffold/tree/.agentic-cms/bin/ac-inventory`: add an `archived`
   field (per-topic count + files from `docs/*/archive/*.md`), mirroring the
   `drafts` field; document it in `bin/README.md`.
5. Add an **Archive** mode to
   `scaffold/tree/.claude/skills/content-new-item/SKILL.md`: `ac-page
   archive`, re-file the index entry under archived, `ac-log append archive`,
   then run `ac-links check` and fix flagged inbound references; document
   un-archive as Promote from `archive/` + index re-filing.
6. Update `scaffold/tree/.claude/skills/content-list/SKILL.md`: report
   archived items in their own section from `ac-inventory`'s `archived`
   field.
7. Update `scaffold/tree/.claude/skills/content-lint/SKILL.md`: exempt the
   archived index section from the orphan check, with a one-line rationale.
8. Run `make smoke-test` and `go test ./...`. Manually exercise in a
   sandbox: create an item, archive it, confirm the index entry moved to the
   archived section, the log records it, `ac-index check` and `ac-links
   check` return clean after reference cleanup, `ac-inventory` reports it;
   un-archive and confirm the reverse.

## Known Issues Triage
**Triaged:** 2026-08-15
**Tools searched:** none — no third-party tools identified
**Result:** Clear

Task scope is entirely internal: schema documentation (`CONTENT.md`), skill
instructions, and stdlib bash/python3 edits to the project's own
`.agentic-cms/bin/` toolkit. No third-party product, library, or service
dependency is introduced or at risk. GitHub issue search is not applicable.

### Unresolvable Tools
- (none)

### Omitted Tools
- (none)

### Search Warnings
- (none)

## Acceptance Criteria

- [x] `CONTENT.md` documents `status: archived`, `docs/<topic>/archive/`,
      archive/un-archive flows, and the new log operations
- [x] `ac-page archive` moves a page, sets `status: archived`, touches
      `updated:`, and refuses on destination conflict
- [x] Archived items remain in `wiki/index.md` under an archived section;
      `ac-index check` stays clean with no new exclusion glob
- [x] `ac-inventory` reports archived pages per topic
- [x] `content-manage-item` (renamed from `content-new-item`) has an Archive
      mode including index re-filing, log entry, and inbound-link cleanup;
      un-archive documented via Promote; all cross-references use the new name
- [x] `content-list` reports archived items separately
- [x] `content-lint` exempts archived pages from orphan detection
- [x] `make smoke-test` and `go test ./...` remain passing

## Findings

**Implementation approach:**
- Mirrored task 002's drafts pattern throughout: `docs/<topic>/archive/` as
  the structural mirror of `drafts/`, `archived` inventory field shaped like
  `drafts`, Archive mode added to the same lifecycle skill
- `ac-page archive` implemented as a sibling of `promote` on one shared
  move-with-status code path (move, set status, touch updated, refuse on
  destination conflict)
- Verified end-to-end in a sandbox install: archive, re-file, both checks
  clean, inventory reporting, un-archive, header auto-create, and
  status-field insertion all exercised against a real `init` scaffold

**Decisions made:**
- Archived items stay in `wiki/index.md`, re-filed under `## Archived` —
  unlike drafts — because the index is the retrieval system and archive is
  not delete; no new `ac-index check` exclusion needed as a result
- `ac-index add` auto-creates a known section header missing from an older
  `index.md` (init never rewrites existing files, so upgraded installs would
  otherwise hard-fail on first archive); fresh installs ship the section
- Archiving a page whose frontmatter lacks `status:` inserts the field
  explicitly — absent status means final, so an archived page must carry it
- Skill renamed `content-new-item` → `content-manage-item` (user decision):
  it owns create + the draft/promote/archive lifecycle, not just creation
- `content-lint` also skips archived pages in the stale-claims check, not
  only the orphan check — retired content is not expected to stay current

**Blockers encountered:**
- Task files were created in the template's legacy bold-header format, which
  the lifecycle resolver rejects — converted 003/004 to YAML frontmatter
- Metadata push to origin/main failed with a persistent 403 (fine-grained
  PAT lacked write access); resolved by the user updating the token

**Follow-up identified:**
- `agentic-cms update` cannot clean up renamed scaffold paths: upgraded
  installs keep a stale `content-new-item/` skill dir alongside
  `content-manage-item/` until scaffold diffing lands (recorded in roadmap)
- Task files 001/002 remain in legacy format and produce resolver warnings;
  converting them would silence the noise
- Whole-topic archiving deliberately out of scope; revisit if per-item
  archiving proves insufficient

## Files to Create / Modify

| File | Action |
|------|--------|
| scaffold/tree/CONTENT.md | Modify |
| scaffold/tree/.agentic-cms/bin/ac-page | Modify |
| scaffold/tree/.agentic-cms/bin/ac-index | Modify (only if archived section isn't already generic) |
| scaffold/tree/.agentic-cms/bin/ac-inventory | Modify |
| scaffold/tree/.agentic-cms/bin/README.md | Modify |
| scaffold/tree/.claude/skills/content-manage-item/SKILL.md | Rename from content-new-item + Modify |
| scaffold/tree/.claude/skills/content-list/SKILL.md | Modify |
| scaffold/tree/.claude/skills/content-lint/SKILL.md | Modify |

## Notes

Directly stacks on task 002 (draft content state, PR #2): extends the
`status` enum, the `ac-page promote` command family, the `ac-inventory`
`drafts` pattern, and `content-new-item`'s modes structure. Must start after
002 merges. Independent of task 003 (register mode), though both edit
`CONTENT.md` and `content-list` — whichever merges second rebases trivially.

Whole-topic archiving (retiring an entire `docs/<topic>/`) is out of scope;
revisit if per-item archiving proves insufficient.

Child tasks inherit their active parent's branch, worktree, and workflow
mode. Only a standalone or parent task owns Git lifecycle cleanup.
