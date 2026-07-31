# Installer smoke test script

**Status:** Completed
**Created:** 2026-07-31
**Started:** 2026-07-31
**Completed:** 2026-07-31
**Mode:** Assisted

## Description

Add a bash end-to-end smoke test (`scripts/smoke-test-installer.sh`) that runs the real
compiled `agentic-cms` binary against temp-dir sandboxes to verify install behavior:
greenfield install, idempotent re-run, brownfield `CLAUDE.md` merge, and invalid-target
handling. This complements the existing `scaffold/scaffold_test.go` Go unit tests, which
call `scaffold.Install()` directly and never exercise the actual CLI binary. Wired up via
a new `smoke-test` target in the existing root Makefile (no new/nested Makefile).

Pattern is adapted from the sibling project `~/projects/smaqit-extensions`, which drives
its installer binary end-to-end via `scripts/smoke-test-installer.sh` + `make smoke-test`.
Unlike that project, agentic-cms has no separate `installer/` Go module (its install logic
already lives at repo root in `scaffold/` + `main.go`) and no `uninstall` command, so this
task is scoped to `init` only.

## Design Decisions

- **Location:** `scripts/smoke-test-installer.sh`, not a new `installer/` folder — agentic-cms's
  install logic already lives at repo root (`scaffold/` + `main.go`), unlike smaqit-extensions'
  separate `installer/` Go module.
- **CI:** out of scope for this task — local-only via `make smoke-test`; CI wiring is a follow-up.
- **Uninstall scope:** no `uninstall` command exists in agentic-cms, so the smoke test covers only
  `init` paths (greenfield, idempotent re-run, brownfield merge, invalid-dir error) — no
  install/uninstall round trip like smaqit-extensions.
- **Diff source-of-truth:** `scaffold/tree/` directly — no separate staging-generation step, since
  `embed.FS` reads straight from the committed tree.
- **Makefile:** use the existing root Makefile — do not create a new/nested Makefile.

## Implementation Steps

1. Write `scripts/smoke-test-installer.sh`: `mktemp -d` sandbox, `trap cleanup EXIT`
   (path-prefix-guarded `rm -rf`, `KEEP_SMOKE_DIR=1` override), inline assertion helpers
   (`assert_contains`, `assert_exists`, `assert_tree_matches`).
2. Greenfield case: run `"$binary" init "$sandbox"`; diff installed tree against `scaffold/tree/`
   excluding `CLAUDE.md`; assert `{{DATE}}` substitution in `wiki/log.md`; assert other
   placeholders (e.g. `{{TITLE}}`) survive untouched in `.agentic-cms/templates/`.
3. Idempotency case: mutate an installed file, re-run `init` on the same sandbox, assert 0 new
   files created and the mutation is preserved.
4. Brownfield `CLAUDE.md` case: pre-seed a sandbox with an existing `CLAUDE.md`, run `init`,
   assert `markerBegin`/`markerEnd` present and original content preserved; re-run and assert no
   duplicate marker block.
5. Invalid-target case: run `init` against a nonexistent directory, assert non-zero exit + stderr
   message; add `--version`/`--help` sanity checks.
6. Add `smoke-test` target to Makefile: depends on `build` → `go test ./...` →
   `./$(BINARY) --version`/`--help` sanity → `bash scripts/smoke-test-installer.sh ./$(BINARY)`.
7. Add a short "Development" section to `README.md` documenting `make smoke-test`.

## Known Issues Triage

[Populated by smaqit.task-start via smaqit.utils.triage-issues. Do not edit manually.]

## Acceptance Criteria

- [x] `make smoke-test` builds the binary, runs `go test ./...`, sanity-checks
      `--version`/`--help`, and runs the bash smoke test script successfully end-to-end
- [x] Smoke test asserts a greenfield `init` produces a tree matching `scaffold/tree/`
      (excluding `CLAUDE.md` and `wiki/log.md`, asserted separately for marker-block merge and
      `{{DATE}}` substitution respectively)
- [x] Smoke test asserts an idempotent re-run creates 0 new files and preserves manual
      mutations to installed files
- [x] Smoke test asserts brownfield `CLAUDE.md` merge appends the managed block without
      duplicating markers on repeated runs
- [x] Smoke test asserts `init` against an invalid/nonexistent directory exits non-zero with a
      stderr message
- [x] Existing `go vet ./...` and `go test ./...` (`scaffold_test.go`) remain unaffected and
      passing

## Findings

**Implementation approach:**
- Built `scripts/smoke-test-installer.sh` around a project-local sandbox (`.smoke-test/`,
  gitignored, cleaned on exit unless `KEEP_SMOKE_DIR=1`) and five checks: greenfield tree
  diff, idempotent re-run, brownfield `CLAUDE.md` merge (twice, to catch duplication), invalid
  target directory, `--version`/`--help` sanity.
- Wired into the existing root `Makefile` as `smoke-test: build test`, then `--version`/`--help`
  sanity, then the script — no nested Makefile.
- Added a short Development section to `README.md` documenting `make smoke-test`.

**Decisions made:**
- Refined the planned "exclude CLAUDE.md + wiki/log.md" diff approach: `{{DATE}}` is actually
  substituted across 8 files (not just `wiki/log.md`, per `grep -rl '{{DATE}}' scaffold/tree/`),
  so the script instead copies `scaffold/tree/` to an `expected/` dir, resolves `{{DATE}}` to
  today's date everywhere via `sed`, and diffs that against the installed output with no
  exclusions — more robust and matches `Install()`'s actual substitution behavior exactly.
- Per user follow-up: moved the sandbox from system `/tmp` into a project-local `.smoke-test/`
  directory (added to `.gitignore`) so `KEEP_SMOKE_DIR=1` leaves an inspectable, discoverable
  artifact inside the repo instead of an opaque `/tmp` path.
- Assertion helpers count lines by their `Result.Print()` indent prefix (e.g. `"  created  "`)
  rather than substring-matching words like "created"/"skipped", since the `Done: N created, M
  merged, K skipped.` summary line always contains those words regardless of actual counts.

**Blockers encountered:**
- None.

**Follow-up identified:**
- CI wiring (`.github/workflows/`) was explicitly out of scope for this task; a follow-up task
  can wire `make smoke-test` into CI, mirroring `smaqit-extensions`' `test-integration.yml`.

## Files to Create / Modify

| File | Action |
|------|--------|
| `scripts/smoke-test-installer.sh` | Create |
| `Makefile` | Modify (add `smoke-test` target) |
| `README.md` | Modify (add Development section) |

## Notes

Pattern adapted from `~/projects/smaqit-extensions` (`scripts/smoke-test-installer.sh` +
`installer/Makefile` `smoke-test` target + `.github/workflows/test-integration.yml`), scoped
down since agentic-cms has no separate installer module and no `uninstall` command.

Child tasks inherit their active parent's branch, worktree, and workflow mode. Only a standalone
or parent task owns Git lifecycle cleanup.
