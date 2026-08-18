# Scaffold Overwrite Semantics

**Date:** 2026-08-19
**Session focus:** Diagnose and fix a real-world installer bug reported by the user, ship it as a release, then verify it end-to-end against the actual published artifact.
**Tasks completed:** 010 (Always overwrite scaffolding logic files on init/update)

## Actions taken

- Investigated a user report from a separate project (`magnificah`): after `agentic-cms update` to v0.6.0 then `agentic-cms init`, every pre-existing skill/agent/template file was reported `skipped (exists)`, with no way to tell it was stale relative to the newly installed CLI version. Also addressed a secondary claim — that git showed no untracked files for the new scaffold — which turned out to be incorrect on inspection; `git status` did show them correctly.
- Read `scaffold/embed.go` and `update.go` to confirm the root cause: `Install()` treated every existing target file identically (skip-if-exists), including files `update.go`'s own comments described as being refreshed on re-init (skills, agents, templates). Confirmed via `scaffold_test.go` that this blanket behavior was even locked in by a test using `wiki/index.md` as its (correct, for that file) "never overwrite" example.
- Created task 010, scoped narrowly per user direction: no drift/checksum detection, just always overwrite a fixed bucket of framework-owned paths.
- Implemented the fix in `scaffold/embed.go`: added `isFrameworkOwned()` path-prefix classification, split the existing-file branch so framework-owned paths (`.claude/skills/`, `.claude/agents/`, `.agentic-cms/templates/`, `.agentic-cms/scripts/`, `.agentic-cms/hooks/`, `.codex/`) are always overwritten and reported in a new `Result.Updated` bucket, while user-content paths keep skip-if-exists. Updated `main.go`'s summary line and added `TestInstallOverwritesFrameworkFiles`.
- Ran `make test` and `make smoke-test` (all green), plus a manual repro of the exact user scenario before opening the PR.
- Completed task 010 through the full PR-gated release flow: Phase 1 (commit, `release-analysis`/`release-approval`, PR #8 "Prepare release v0.6.2", pending CHANGELOG entry, promoted on the branch), user merged and released, Phase 2 (verified merge via `gh pr view`, pulled `main`, cleaned up worktree/branch). Released as v0.6.2.
- Created an E2E test playbook (`.smaqit/user-testing/tests/010_...md`) modeled on task 005's brownfield-clone convention, then executed it live — deviating from the script at the user's explicit request to install the actual GitHub release via `install.sh` rather than a local build, so the test validated the real shipped artifact.
- Verified the fix by content, not just printed labels: hand-edited `content-lint/SKILL.md`, re-ran `init`, confirmed via `grep -c` that the edit was actually reset (not just relabeled `updated`), while `wiki/index.md` and `CLAUDE.md`'s merge-block state survived untouched. Wrote the formal PASS report.

## Problems solved

- Root cause: a single uniform skip-if-exists rule in `scaffold.Install()` silently defeated the documented promise that `init`/`update` refreshes skills, agents, and templates — once a scaffolding file existed on disk, it was frozen at whatever version was first installed, with releases never reaching it again.
- Fix scope was deliberately kept minimal (static path-prefix bucketing, no checksums/diffing) per explicit user direction, avoiding over-engineering a problem that only needed "always vs never overwrite," not drift detection.

## Decisions made

- **Two static buckets, not per-file config or content hashing:** `.claude/skills/`, `.claude/agents/`, `.agentic-cms/templates/`, `.agentic-cms/scripts/`, `.agentic-cms/hooks/`, `.codex/` are framework-owned (always overwrite); everything else stays skip-if-exists. `.claude/settings.json` deliberately stays outside the framework-owned bucket — it's user-configurable, not scaffolding logic.
- **New `Result.Updated` bucket instead of reusing `Skipped`/`Created`:** an overwritten file is neither a fresh create nor a real skip; conflating either would misrepresent `agentic-cms init`'s output.
- **Test execution used the real released binary, not a local build:** at the user's request, validated `v0.6.2` as fetched via the public `install.sh` bootstrap, closing the loop against the artifact a real user would actually get.

## Files modified

| File | Change |
|------|--------|
| `scaffold/embed.go` | Added `frameworkOwnedPrefixes`/`isFrameworkOwned()`; reworked `Install()`'s existing-file branch to always-overwrite framework-owned paths into a new `Result.Updated` bucket |
| `scaffold/scaffold_test.go` | Added `TestInstallOverwritesFrameworkFiles` |
| `main.go` | Extended the `Done: ...` summary line with the `updated` count |
| `CHANGELOG.md` | Added the v0.6.2 entry (via pending-entry → promote flow) |
| `.smaqit/tasks/010_overwrite_scaffolding_logic_files_on_init.md` | Created, taken through full lifecycle to Completed |
| `.smaqit/tasks/PLANNING.md` | Task 010 added, moved through In Progress → PR Open → Completed |
| `.smaqit/references/project-research.md` | Added Task 010's research block (Go/embed.FS) |
| `.smaqit/user-testing/tests/010_always-overwrite-scaffolding-logic-files-on-init-update.md` | Created E2E playbook, executed, checkboxes filled in with real results |
| `.smaqit/user-testing/2026-08-19_test-report.md` | Created — PASS, 12/12 |
| `agentic-cms.code-workspace` | Rebuilt after worktree cleanup |

## Next steps

- Task 003 (Register mode for reference files) and task 008 (global classification scanner) remain Not Started — both unrelated to this session's work.
- No follow-up identified for task 010 itself; the fix, release, and live verification are all closed out.

## Session Metrics

- **Duration:** ~1 session (task creation through release, test, and report)
- **Tasks completed:** 1 (010)
- **PRs merged:** 1 (#8)
- **Release shipped:** v0.6.2
- **Files created:** 4 (task file, playbook, test report, this history entry)
- **Files modified:** 6 (embed.go, scaffold_test.go, main.go, CHANGELOG.md, PLANNING.md, workspace file)
- **Test result:** PASS, 12/12 checks, verified against the real released artifact
