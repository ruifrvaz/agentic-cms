---
status: In Progress
created: "2026-08-18"
mode: Assisted
started: "2026-08-18"
---

# Rename .agentic-cms/bin to .agentic-cms/scripts

## Description

The installed scaffold's toolkit directory (`ac-page`, `ac-index`, `ac-log`,
`ac-links`, `ac-inventory`, `ac-search`) lives at `.agentic-cms/bin/` in
every installed project. This is a naming mistake: `bin/` is committed,
version-controlled **source** — bash+python3 scripts meant to be read,
audited, and modified — not compiled build output. Naming it `bin/`
collides with the generic IDE/tooling convention that a directory named
`bin` holds compiled artifacts safe to hide from the file explorer (VS
Code's `files.exclude` on `**/bin/**` is exactly this convention, applied
by the smaqit worktree tooling's generated workspace file). That collision
surfaced directly: the user could not find `ac-classify` in VS Code because
the workspace file's default exclude was hiding the entire directory.

This task renames the scaffold source's toolkit directory from
`.agentic-cms/bin/` to `.agentic-cms/scripts/` everywhere it is referenced,
so fresh installs never hit this collision again.

**Scope is Tier 1 only** (explicit user decision 2026-08-18): rename the
scaffold source and every reference to it. **Not in scope**: migration
logic for already-installed projects (their `.agentic-cms/bin/` stays as
committed content on their side; `agentic-cms update`'s non-destructive
installer would only ever *add* a new `.agentic-cms/scripts/` alongside the
old one, never move/rewrite existing files — same gap task 004 already left
for the `content-new-item` → `content-manage-item` skill rename, tracked
generically under the README's "scaffold diffing" roadmap item). A future
task may build real update-time migration if that gap becomes a real
problem for the downstream project(s) mentioned in other task notes.

**Sequencing note**: this task was deliberately started *before* task 005
(First-class content classification) completes, on explicit user
instruction — task 005 is fully implemented in its worktree
(`.agentic-cms/bin/ac-classify` plus ~150 other `.agentic-cms/bin`
references across its own diff) but not yet committed/PR'd. After this
task's PR merges to `main`, task 005's branch must be brought up to date
with the renamed baseline and its own new content (which currently also
says `bin/`) updated to say `scripts/` before task 005 completes. That
reconciliation is 005's responsibility at its next session, not this task's
— this task's scope ends at a clean rename of what exists on `main` today.

## Issue Triage Context

**Mode:** Skip
**Technologies:** None
**Platforms/Environments:** None
**Features/Integrations:** None
**Versions/Constraints:** Scoped to what exists on `main` at task start (6 scripts + README under `.agentic-cms/bin/`); does not touch task 005's uncommitted worktree content

## Design Decisions

- **Tier 1 only** (explicit user decision): rename the scaffold source;
  no update-time migration for already-installed projects in this task.
- **New name: `.agentic-cms/scripts/`** (explicit user instruction —
  "scripts should be under a scripts folder").
- **`.agentic-cms/hooks/` is unaffected** — that directory doesn't exist on
  `main` yet (it's task-005-only, uncommitted); when it lands it already
  uses a real semantic name (mirroring git's own `.git/hooks/`), not a
  build-output-colliding one, so it needs no equivalent rename.
- **Historical records are never rewritten**: `.smaqit/compendium.md`,
  `.smaqit/history/*`, `CHANGELOG.md`'s past version sections, and
  completed task files (002, 004) describe what was true at the time and
  keep saying `.agentic-cms/bin/` — rewriting them would falsify the audit
  trail. Only *living* documentation (scaffold content, root README,
  Go source, the smoke test) gets renamed.
- **Mechanical, single find-and-replace**: every occurrence of the literal
  string `.agentic-cms/bin` becomes `.agentic-cms/scripts` — no other
  content changes. This is a textual substitution with no structural
  interaction with any other in-flight work, which is exactly why it's
  safe to land as an independent task right before task 005 completes
  rather than folded into it.
- **`git mv`, not copy-then-delete**, for the directory itself, so history
  is preserved and the rename is visible as a rename in `git log --follow`.

## Implementation Steps

1. `git mv scaffold/tree/.agentic-cms/bin scaffold/tree/.agentic-cms/scripts`
   (moves all 7 files: `README.md`, `ac-page`, `ac-index`, `ac-log`,
   `ac-links`, `ac-inventory`, `ac-search` — preserves the executable bit).
2. Sweep every remaining reference to the literal string
   `.agentic-cms/bin` and replace with `.agentic-cms/scripts` in:
   - `scaffold/tree/CONTENT.md` (directory map, "The toolkit" section)
   - `scaffold/tree/.agentic-cms/scripts/README.md` (composition example)
   - Every `scaffold/tree/.claude/skills/*/SKILL.md` that references
     `ac-*` commands (content-new, content-manage-item, content-import,
     content-research, content-add-notes, content-list, content-lint,
     content-export, content-query)
   - `scaffold/tree/.claude/agents/*.md` (content-importer, content-exporter,
     content-researcher)
   - `scaffold/embed.go`: the executable-mode prefix check
     (`strings.HasPrefix(rel, filepath.FromSlash(".agentic-cms/bin/"))`)
   - `scaffold/scaffold_test.go`: `wantFiles` list entries and the
     executable-bit assertion's path
   - `scripts/smoke-test-installer.sh`: the toolkit functional check
     section's `.agentic-cms/bin/ac-*` invocations
   - Root `README.md`: the "What gets installed" tree diagram and the
     "Deterministic core, judgment at the edges" design-notes bullet
   Use `grep -rl '\.agentic-cms/bin' <paths above>` to enumerate exactly,
   then verify zero remaining hits afterward (excluding the
   never-touched historical files listed in Design Decisions).
3. Run `go vet ./...`, `go test ./...`, and `make smoke-test` — all must
   pass unchanged in *behavior*, only paths differ. Confirm
   `scaffold_test.go`'s greenfield test finds `.agentic-cms/scripts/ac-index`
   installed and executable, and that no test still references the old
   `.agentic-cms/bin/*` path.
4. Grep the final tree for any remaining literal `.agentic-cms/bin`
   occurrence outside the explicitly-excluded historical files, and
   confirm zero.

## Known Issues Triage
**Triaged:** 2026-08-18
**Tools searched:** none — triage skipped
**Result:** Clear

Triage skipped — explicitly marked `Mode: Skip` in Issue Triage Context.
Task scope is a pure textual rename (path string substitution) across the
project's own scaffold source and Go code; no third-party product,
library, or service dependency is introduced or at risk.

## Acceptance Criteria

- [ ] `scaffold/tree/.agentic-cms/bin/` no longer exists; `scaffold/tree/.agentic-cms/scripts/` contains all 7 files (`README.md` + 6 `ac-*` scripts), moved via `git mv`
- [ ] Zero remaining `.agentic-cms/bin` references in any *living* file (CONTENT.md, all skills, all agents, scaffold/embed.go, scaffold/scaffold_test.go, scripts/smoke-test-installer.sh, root README.md)
- [ ] Historical files (`.smaqit/compendium.md`, `.smaqit/history/*`, `CHANGELOG.md`, completed task files 002/004) are untouched — still say `bin/`, accurately describing what shipped at the time
- [ ] `scaffold/embed.go`'s executable-mode logic marks `.agentic-cms/scripts/*` files executable (not `.agentic-cms/bin/*`)
- [ ] `go vet ./...`, `go test ./...`, and `make smoke-test` all pass against the renamed paths
- [ ] A fresh `agentic-cms init` into an empty directory installs `.agentic-cms/scripts/ac-index` (executable) and no `.agentic-cms/bin/` at all

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
| scaffold/tree/.agentic-cms/bin/ → scaffold/tree/.agentic-cms/scripts/ | Rename (git mv) |
| scaffold/tree/CONTENT.md | Modify |
| scaffold/tree/.claude/skills/*/SKILL.md | Modify (all 9) |
| scaffold/tree/.claude/agents/*.md | Modify (all 3) |
| scaffold/embed.go | Modify |
| scaffold/scaffold_test.go | Modify |
| scripts/smoke-test-installer.sh | Modify |
| README.md | Modify |

## Notes

Task 005 (First-class content classification) is fully implemented in its
own worktree but not yet committed, sitting at ~150 `.agentic-cms/bin`
references of its own (new `ac-classify` script, two agent-hook configs,
a git pre-commit gate, and skill/CONTENT.md edits). This task deliberately
does not touch that worktree — after this task's PR merges, task 005 needs
its own reconciliation pass (pull the renamed baseline, rename its own new
`.agentic-cms/bin/ac-classify` → `.agentic-cms/scripts/ac-classify` and
every reference to it) before it can complete. Since this rename is a pure
textual substitution with no structural changes, applying the identical
substitution independently to task 005's content commutes cleanly with
this task's changes to shared files — expect a low-conflict reconciliation,
not a rewrite.

Child tasks inherit their active parent's branch, worktree, and workflow
mode. Only a standalone or parent task owns Git lifecycle cleanup.
