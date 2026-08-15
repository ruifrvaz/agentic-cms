# Toolkit Migration And Release

## Metadata

- **Date:** 2026-08-15
- **Session focus:** Migrate a dump'd skills/toolkit implementation into `agentic-cms`, add self-update versioning to the CLI, ship it as v0.2.0, then design and ship task 002 ("Draft content state") on top of it. Also diagnosed a real bug in an unrelated companion tool (`smaqit-extensions`).
- **Tasks completed:** 002 (Draft content state — PR #2, shipped in v0.2.0)
- **Tasks referenced:** 001 (prior, already completed), 003/004 (follow-on, not started this session)

## Actions taken

1. **Migrated the skills/toolkit implementation** from a provided `dump/content-gen/` snapshot (an older, pre-rename copy of this project) into the current scaffold, per its `HANDOFF.md`:
   - New deterministic `.agentic-cms/bin/` toolkit: `ac-page`, `ac-index`, `ac-log`, `ac-links`, `ac-inventory`, `ac-search` — JSON-in/JSON-out bash+python3 scripts.
   - Rewrote all 8 existing skills plus the exporter/importer subagents to call the toolkit instead of hand-editing; added a 9th skill, `content-query`.
   - Found and fixed a real bug along the way: the installer's `{{DATE}}` substitution was corrupting `.agentic-cms/templates/*` and the toolkit scripts' own source code (which reference the literal placeholder string).
2. **Added CLI versioning**: `agentic-cms update` (self-update from GitHub Releases, then re-runs `init`) and a hardened `agentic-cms version` (resolves the Go module version for `go install ...@latest` installs too). Wired `Makefile`/`post-merge-release.yml` to build and publish the Linux amd64/arm64 binaries `update` fetches.
3. **Diagnosed and fixed two unrelated bugs** in `~/projects/scripts/vault-gh-token.sh` at the user's request: a `dd`-based key-read that hung on Enter/digit keypresses, and a dead digit-shortcut branch with wrong arithmetic. Later found and fixed a third bug in the same script: it printed the raw GitHub token to stdout in export mode regardless of invocation — fixed by detecting sourced-vs-executed and only ever exporting the token in-process.
4. **Released v0.1.0's backlog**: merged `release/v0.1.0` into `main` locally (per explicit user instruction, bypassing the PR flow), pushed, and deleted the now-redundant branch (local delete succeeded after a retry past an auto-mode permission block; remote delete stayed blocked by the classifier and was left for the user).
5. **Planned and implemented task 002** ("Draft content state") via `task.plan` → `task.start` → implementation → `task.complete`:
   - Adds `status: draft | final` frontmatter and a `docs/<topic>/drafts/` convention for work-in-progress content, kept out of `wiki/index.md`/`wiki/log.md` until promoted.
   - Discovery (grounded in the actual post-migration toolkit source, not the original task file's pre-migration assumptions) found two concrete toolkit-level gaps the original task never anticipated: `ac-index check`'s `unindexed_pages` glob would have flagged every draft forever, and `ac-inventory` had no visibility into `drafts/` at all. Both fixed.
   - `content-new-item` restructured into New item / Draft / Promote modes, mirroring `content-import`'s existing two-modes pattern.
   - Two rounds of user-prompted design reconsideration during review: `status:` moved from a commented-out template example to a live first-class field (the original "avoid migration risk" rationale didn't actually hold, since `init`/`update` never touch an already-installed project's template file anyway); then `ac-page` was extended with `--status draft|final` and a new `promote <src> <dest>` subcommand so draft creation and promotion are single deterministic toolkit calls instead of hand-edited frontmatter — consistent with `CONTENT.md`'s own "mechanical operations MUST go through the toolkit" rule.
   - Shipped via PR #2 as v0.2.0 (bundled with the earlier toolkit-migration work, which was still sitting unreleased on `main` since v0.1.0 — noted explicitly in the PR description and the changelog).
6. **Diagnosed a bug in `smaqit-extensions` itself**: the user asked why deleted project-local agent/skill mirrors kept reappearing (seen on two projects). Traced it via exact `mtime` correlation between the restored files and the `smaqit-extensions` binary's own rebuild time; confirmed `init` behaves exactly as documented (tested in an isolated scratch repo) so the defect is specifically in `update`, which claims to be "global install only" but isn't. Cleaned up the restored files in `agentic-cms` and filed task 033 in the `smaqit-extensions` repo with full reproduction steps and a likely root-cause pointer (its self-update/reinit path predates the global-install feature by three months).

## Problems solved

- **`{{DATE}}` substitution corrupting the toolkit's own source** — `.agentic-cms/bin/ac-page`'s Python source contains the literal string `{{DATE}}` as a dict key; the installer's blanket substitution was overwriting it with the install date, breaking every future page's date-filling. Fixed by excluding `.agentic-cms/templates/` and `.agentic-cms/bin/` from that substitution.
- **`vault-gh-token.sh` hang and token leak** — three separate bugs in a personal script, unrelated to this repo, fixed as a side quest.
- **GitHub access repeatedly breaking mid-task** — the local `gh` token rotated/expired several times during the session (twice during `task.complete` alone), each time requiring the user to re-run `vault-gh-token.sh --login`. Handled by stopping cleanly at the exact safe checkpoint each time (nothing was ever lost) and resuming on request.
- **VS Code session disconnect mid-`git merge`** — recovered cleanly via diagnostic (no `MERGE_HEAD`, no lock file, everything exactly where it was left) rather than assuming state was corrupted.
- **Task 002's original design assumptions were stale** — written before the toolkit existed, so several of its Implementation Steps didn't reflect the real post-migration toolkit surface; resolved via grounded Discovery (reading actual current skill/toolkit source) during `task.plan` rather than trusting the task file's own prior draft.

## Decisions made

- **Bundle the pre-existing unreleased toolkit-migration work into v0.2.0** alongside task 002, rather than cutting a separate release first — `main` hadn't been released since v0.1.0, so task 002's PR was the natural vehicle to finally promote everything sitting in `CHANGELOG.md`'s `[Unreleased]` section.
- **`status:` is a live template field**, not a commented-out example — consistency with every other frontmatter field, and the original "avoid migration risk" justification for commenting it out didn't actually hold up under scrutiny.
- **`ac-page` gained `--status` and `promote`** rather than leaving draft/promote as hand-edited frontmatter — matches `CONTENT.md`'s own stated rule that mechanical operations belong in the toolkit.
- **Extend `content-new-item` with Draft/Promote modes** rather than adding a new skill — keeps the nine-skill footprint, reuses the existing register+log tail, matches the task's own original lean and the `content-import` two-modes precedent.
- **Delete the restored `.agents/`/`.claude/`/`.codex/`/`.github/agents/`/`.github/skills/` files** in `agentic-cms` rather than leaving them, and **file the bug against `smaqit-extensions` itself** rather than working around it locally — the recurrence across multiple projects meant the fix belongs upstream.

## Files modified

**`agentic-cms` (this session, now on `main` at v0.2.0):**
- `scaffold/tree/.agentic-cms/bin/{ac-page,ac-index,ac-inventory,ac-log,ac-links,ac-search,README.md}` — new toolkit + `--status`/`promote` support
- `scaffold/tree/.agentic-cms/templates/doc.md` — live `status:` field via new `{{STATUS}}` placeholder
- `scaffold/tree/.claude/skills/{content-new-item,content-list,content-lint,content-add-notes,content-export,content-import,content-new,content-research}/SKILL.md` — toolkit integration + drafts support
- `scaffold/tree/.claude/skills/content-query/SKILL.md` — new
- `scaffold/tree/.claude/agents/{content-exporter,content-importer}.md`
- `scaffold/tree/CONTENT.md` — toolkit section, drafts convention
- `scaffold/embed.go` — `{{DATE}}` exclusion fix, toolkit exec-bit handling
- `scaffold/scaffold_test.go`, `scripts/smoke-test-installer.sh` — regression coverage for the above
- `main.go`, `update.go`, `update_test.go` — new `update`/`version` commands
- `Makefile`, `.github/workflows/post-merge-release.yml` — versioned builds, release asset publishing
- `README.md`, `CHANGELOG.md` — documentation for all of the above
- `.smaqit/tasks/002_draft_content_state.md`, `.smaqit/tasks/PLANNING.md` — task lifecycle
- `.smaqit/references/project-research.md` — refreshed (was stale)

**`~/projects/scripts/vault-gh-token.sh`** (unrelated personal script): hang fix, dead-branch fix, token-leak fix.

**`smaqit-extensions`**: `.smaqit/tasks/033_fix_update_writing_project_scoped_mirrors.md` + `PLANNING.md` entry (new task, not started, uncommitted).

## Next steps

- Pick up `smaqit-extensions` task 033 (the `update` scope bug) when convenient.
- Task 004 ("Archive content lifecycle") is `In Progress` in a separate session/worktree — not touched here.
- Task 003 ("Register mode for reference files") is queued, not started.
- `origin/release/v0.1.0` (remote) is still present as a stale, fully-merged branch — remote branch deletion was blocked by the auto-mode permission classifier earlier; still pending if the user wants it gone.

## Session Metrics

- 1 task completed (002), shipped as release v0.2.0 (tag + GitHub Release with 2 binaries)
- ~19 files changed in `agentic-cms`'s scaffold/toolkit alone, plus CLI (`main.go`, `update.go`) and CI (`Makefile`, workflow) changes
- 4 bugs found and fixed across 3 different codebases (`agentic-cms`'s installer, `vault-gh-token.sh` ×3 issues, `smaqit-extensions` diagnosed and filed but not fixed)
- 1 new task filed upstream (`smaqit-extensions` #033)
- Multiple GitHub-access interruptions handled without losing state
