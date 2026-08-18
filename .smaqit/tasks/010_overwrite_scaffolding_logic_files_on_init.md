---
status: In Progress
created: "2026-08-18"
mode: Assisted
started: "2026-08-19"
---

# Always overwrite scaffolding logic files on init/update

## Description

`scaffold.Install()` in [scaffold/embed.go](../../scaffold/embed.go) applies one uniform rule to every file in the embedded tree: if a file already exists at the target path, it is skipped and reported as `skipped (exists)`. That rule is correct for user-owned content and data files (`wiki/index.md`, `wiki/log.md`, `raw/`, `CONTENT.md`) — those must never be touched by a re-init. But it is applied identically to versioned scaffolding logic files that the project itself owns and ships updates to: `.claude/skills/**`, `.claude/agents/**`, `.agentic-cms/templates/**`, `.agentic-cms/scripts/**`, `.agentic-cms/hooks/**`, and `.codex/**`.

This contradicts the documented intent in [update.go](../../update.go)'s `checkAndReInit`/`reinitWithBinary` comments, which state that re-running init "deploy[s] updated skills, agents, and templates." In practice, once a scaffolding file exists on disk in an installed project, `agentic-cms init`/`agentic-cms update` will never refresh it again — new releases silently fail to reach already-initialized projects for anything under those paths.

Discovered when a user ran `agentic-cms update` to v0.6.0 in a separate project, then `agentic-cms init`, and found every pre-existing skill/agent/template file reported as `skipped (exists)` with no way to know it was stale relative to the newly installed CLI version.

Fix scope is intentionally narrow: no drift detection, no checksums, no diffing. Just reclassify scaffold paths into two buckets and always overwrite the framework-owned bucket unconditionally on every install run.

## Issue Triage Context

**Mode:** Auto
**Technologies:** Go, embed.FS
**Platforms/Environments:** None
**Features/Integrations:** agentic-cms init, agentic-cms update, scaffold install
**Versions/Constraints:** None

## Design Decisions

- **No drift/checksum detection:** explicitly out of scope per user request — always overwrite the framework-owned bucket, regardless of whether the on-disk copy was locally modified.
- **Two static path buckets, not per-file config:** classify by path prefix in `scaffold.Install()`, matching the existing `isScriptDir`-style prefix checks already used in the same function for executable-mode handling.
- **Framework-owned bucket (always overwrite):** `.claude/skills/**`, `.claude/agents/**`, `.agentic-cms/templates/**`, `.agentic-cms/scripts/**`, `.agentic-cms/hooks/**`, `.codex/**`.
- **User-content bucket (skip if exists, unchanged):** `wiki/**`, `raw/**`, `CONTENT.md`, `docs/**`, `.agentic-cms/VERSION`. `CLAUDE.md` keeps its existing special-case merge-block logic untouched.
- **Reporting:** overwritten framework files must not be reported as `skipped` (misleading — they were in fact refreshed). Add a distinct result bucket (e.g. `Updated`) alongside `Created`/`Skipped`/`Merged`.

## Implementation Steps

1. In [scaffold/embed.go](../../scaffold/embed.go), add a helper (e.g. `isFrameworkOwned(rel string) bool`) that matches the framework-owned path prefixes listed above.
2. In the `fs.WalkDir` callback (around line 79), change the existing-file branch: if `isFrameworkOwned(rel)`, always write the file (executable-mode logic already present stays as-is) and record it in a new `Result.Updated` slice instead of `Skipped`; otherwise keep current skip-if-exists behavior.
3. Add `Updated []string` to the `Result` struct and print it in `Result.Print()` (e.g. `"  updated  %s\n"`), placed after `Created` and before `Skipped` in output ordering to match existing created/merged/skipped ordering.
4. Update [scaffold/scaffold_test.go](../../scaffold/scaffold_test.go): `TestInstallIdempotent` (currently asserts `wiki/index.md` is never overwritten) stays as-is and must keep passing unchanged. Add a new test (e.g. `TestInstallOverwritesFrameworkFiles`) that: installs into a temp dir, mutates a framework-owned file (e.g. `.claude/skills/content-lint/SKILL.md`), re-runs `Install()`, and asserts the file content was reset to the embedded version and reported under `Result.Updated`, not `Result.Skipped`.
5. Check `main.go`'s CLI output path for `Result.Print()` usage doesn't need separate wiring beyond the struct/print change.
6. Run `make test` (`go vet ./...` && `go test ./...`).
7. Manually verify against a scratch project: init, modify a skill file, re-run init, confirm the skill file is restored to the shipped version and reported as `updated`, while `wiki/index.md`/`raw/` content is untouched.

## Known Issues Triage
**Triaged:** 2026-08-19
**Tools searched:** Go
**Result:** Clear

### Blocking Issues
- None

### Advisory Issues
- None

### Historical (Closed)
- None — closed `golang/go` results returned (Seek/ReadAt/vendor-embed quirks) but none relate to this task's actual scope, which is agentic-cms's own file-overwrite policy during scaffold install, not a Go `embed.FS` defect.

### Unresolvable Tools
- None

### Omitted Tools
- None

### Search Warnings
- None

## Acceptance Criteria

- [ ] `scaffold.Install()` always overwrites files under `.claude/skills/`, `.claude/agents/`, `.agentic-cms/templates/`, `.agentic-cms/scripts/`, `.agentic-cms/hooks/`, `.codex/`, even when a locally-modified copy already exists at the target path.
- [ ] User-content paths (`wiki/`, `raw/`, `CONTENT.md`, `docs/`, `.agentic-cms/VERSION`) and `CLAUDE.md`'s merge-block behavior remain unchanged — never overwritten.
- [ ] `Result` reports overwritten framework files distinctly (not as `Skipped`), and `agentic-cms init` CLI output reflects this.
- [ ] `TestInstallIdempotent` still passes unmodified.
- [ ] New test confirms a locally-modified framework-owned file is reset to the embedded version on re-install.
- [ ] `make test` passes.

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
| scaffold/embed.go | Modify — add framework-owned path classification, always-overwrite branch, `Updated` result bucket |
| scaffold/scaffold_test.go | Modify — add overwrite test for framework-owned files |

## Notes

Related to task 006 (shell installer bootstrap) and task 007 (bin→scripts rename), both of which touch the same install/update plumbing. No relation to task 008 (global classification scanner) beyond both living under `.agentic-cms/`.

Child tasks inherit their active parent's branch, worktree, and workflow mode. Only a standalone or parent task owns Git lifecycle cleanup.
