# Classification Adoption Gaps Release

**Date:** 2026-08-19
**Session focus:** Repair the four gaps surfaced by classification's first real brownfield rollout (a repo with ~50 pages, ~30 pre-existing floor hits, ~2/3 of them false positives), release v0.7.0, and verify the fix live against the real published binary.
**Tasks completed:** 011 (Classification adoption gaps — first real-world rollout)

## Actions taken

- Reviewed task 011 (already specced from a prior session) and, before implementation, revised its currency-floor design decision at the user's explicit direction: the floor's contract is 0% false negatives over its enumerated pattern classes, so gap 3 flipped from *narrowing* the currency regex (context-requiring) to *widening* it (ISO-coded, symbol-suffixed, word-form amounts) — all false-positive relief routes through a new ack mechanism instead, never through weaker detection.
- Ran `smaqit.task-start 011`: created the task branch/worktree, refreshed the task-relevant research-map block, ran issue triage (Advisory only — notably openai/codex#27133 on `.codex/hooks.json` being ignored inside git worktrees, and anthropics/claude-code#77341 on PostToolUse not firing in daemon sessions).
- Implemented all five gaps in the task worktree:
  - **Delta-scoped pre-commit gate** (`hooks/pre-commit`): blocks only violations in files staged in the commit, plus bleed into a staged `wiki/index.md`/`wiki/log.md` regardless of source page; pre-existing tree-wide drift collapses into one summary warning. Fixed a latent bug in the same rewrite — the old early-exit check didn't cover commits touching only the index/log files, silently exempting them from the gate.
  - **Floor acknowledgment**: `ac-page classify <path> <level> --ack-floor` stamps a `classification-ack:` field bound to the same body hash as `classified-hash:`; `ac-classify` honors it as non-blocking while the hash matches; any edit invalidates it.
  - **Currency floor widened**: added symbol-suffixed, ISO-coded, and word-form amount patterns to `C2_PATTERNS`, keeping every existing pattern intact.
  - **CONTENT.md schema reconciliation** (`scaffold/reconcile.go`, new): `init`/`update` diff installed vs. shipped `## ` heading text; on a mismatch, print a report and write `.agentic-cms/CONTENT.upstream.md` — the user's `CONTENT.md` is never programmatically edited, by design.
  - **VERSION auto-stamp**: `.agentic-cms/VERSION` became a `{{VERSION}}` placeholder that `Install()` always overwrites with the installing binary's ldflags version, replacing a static literal that had drifted to `0.1.0` at product v0.6.2.
- Added the classification engine's first-ever fixture tests (`engine_test.go`, `reconcile_test.go`, plus VERSION coverage in `scaffold_test.go`), extended `make smoke-test` with ack and reconciliation scenarios, and updated CONTENT.md (promoted Classification to its own `##` section), the scripts README contract, and five write-path/lint skills to document the ack as user-only.
- Took task 011 through the full PR-gated release flow: `smaqit.release-analysis` (Task mode) computed v0.7.0 (MINOR — new capabilities, no breaking changes) against boundary v0.6.2; PR #9 ("Prepare release v0.7.0") opened, pending CHANGELOG entry pushed to `main` then promoted on the branch; user merged; Phase 2 confirmed the merge via `gh pr view`, pulled `main`, cleaned up the worktree and branch.
- Created a new E2E test playbook (`.smaqit/user-testing/tests/011_classification-adoption-gaps.md`) and executed it live against the **real published v0.7.0 binary** (installed via the public `install.sh` bootstrap, not a local build) in a fresh `mktemp` sandbox: all 8 acceptance criteria exercised through real `git commit` scenarios and real `ac-classify`/`ac-page` invocations. Result: PASS, 19/19 checks. One self-inflicted test-harness mistake (restoring a blocked commit's staged file with `git checkout --` alone, which doesn't unstage) was diagnosed and fixed mid-run; no product defect found.
- Answered follow-up conceptual questions: how CONTENT.md reconciliation works (read-only diff + report + sidecar, never edits the user's file — clarified this is *expected* even when nothing visibly changes in `CONTENT.md` itself), diagnosed a real-world confusion on the `magnificah` project where `agentic-cms update` reported no CONTENT.md changes (walked through the two actual reconciliation signals — console report and `.agentic-cms/CONTENT.upstream.md` — rather than the file itself, and provided diagnostic commands), and explained the ack mechanism's purpose, user-only restriction, and concrete intervention triggers.

## Problems solved

- The classification gate's original design punished brownfield adoption: any single pre-existing floor violation anywhere in the tree blocked every commit, training users toward `--no-verify`. Now only the staged delta blocks; the backlog is `content-lint`'s job.
- False positives previously had no legitimate way to persist — the ratchet's "only the user may lower" escape existed as doctrine with no mechanical representation. The ack now makes it real and auditable (hash-bound, always reversible by an edit).
- `CONTENT.md`'s classification rubric silently never reached pre-existing installs after task 005 shipped it — `update` had no signal at all that the installed schema was behind. The reconciliation report + sidecar closes that gap without violating the project's non-destructive-install guarantee.
- `.agentic-cms/VERSION` had drifted to a stale `0.1.0` literal at product v0.6.2, with no mechanism to keep it current. It's now automatically correct on every `init`/`update`, permanently.

## Decisions made

- **Floor recall over precision, no exceptions**: reversed the task's original gap-3 design (context-narrowing) mid-session per explicit user direction — the floor must have 0% false negatives over its enumerated shapes; all noise reduction happens via the ack, never via weaker patterns.
- **VERSION automation over a manual checklist step**: the acceptance criterion asked for "release checklist includes the bump," but no release-checklist document exists in this repo — automatic stamping (always-overwrite from the binary's own version) satisfies the underlying intent more robustly than a step someone could forget, which is the exact bug this gap exists to fix. Flagged to the user; not contested.
- **CONTENT.md stays read-only, always**: no managed-block convention (unlike `CLAUDE.md`'s `<!-- agentic-cms:begin/end -->`) — sidecar-plus-report is the permanent shape for this reconciliation, not a stepping stone toward auto-merging.
- **A plain `ac-page classify` (no `--ack-floor`) withdraws any standing ack** — a fresh rating decision supersedes an old ack rather than leaving it stale.

## Files modified

| File | Change |
|------|--------|
| `scaffold/tree/.agentic-cms/scripts/ac-classify` | Widened `C2_PATTERNS` (ISO/symbol-suffixed/word-form currency); ack-awareness in `check_one()` |
| `scaffold/tree/.agentic-cms/scripts/ac-page` | `classify` gained `--ack-floor`; plain re-rate withdraws a standing ack |
| `scaffold/tree/.agentic-cms/hooks/pre-commit` | Rewritten for delta-scoped blocking + tree-wide backlog summary; fixed index/log-only-commit exemption bug |
| `scaffold/reconcile.go` | New — `ReconcileContentMD`, `InstalledVersion`, sidecar writer |
| `scaffold/embed.go` | `Install()` takes a version param; VERSION always restamped |
| `main.go`, `update.go` | Threaded version through `runInit`/`checkAndReInit`; wired `ReconcileContentMD` into `init` |
| `scaffold/tree/CONTENT.md` | Classification promoted to its own `##` section; ack + delta-gate documented |
| `scaffold/tree/.agentic-cms/scripts/README.md` | `ac-classify`/`ac-page` contract updated for ack and widened recall |
| `scaffold/tree/.claude/skills/{content-manage-item,content-import,content-add-notes,content-research,content-lint}/SKILL.md` | Ack documented as user-only in every write-path/lint skill |
| `scaffold/engine_test.go`, `scaffold/reconcile_test.go` | New — first fixture tests for the classification engine |
| `scaffold/scaffold_test.go` | VERSION-stamping coverage; updated `Install()` call sites |
| `scripts/smoke-test-installer.sh` | Extended with ack and CONTENT.md-reconciliation scenarios |
| `README.md` | Classification section expanded (ack, delta gate); fixed stale "never overwritten" claims |
| `CHANGELOG.md` | v0.7.0 entry |
| `.smaqit/tasks/011_classification_adoption_gaps.md` | Design decisions revised; triage, Findings, criteria, status through full lifecycle |
| `.smaqit/tasks/PLANNING.md` | Task 011 moved through In Progress → PR Open → Completed |
| `.smaqit/references/project-research.md` | Task 011 research block |
| `.smaqit/user-testing/tests/011_classification-adoption-gaps.md` | Created, executed, all 19 checks verified |
| `.smaqit/user-testing/2026-08-19_test-report_2.md` | Created — PASS, 19/19, verified against the real released artifact |
| `agentic-cms.code-workspace` | Rebuilt after worktree cleanup |

## Next steps

- Task 008 (global classification scanner) is unblocked and now inherits the ack mechanism and widened floor for free via the shared `ac-classify` engine — good next candidate.
- Task 003 (Register mode for reference files) remains Not Started, unrelated to this session.
- The `magnificah` project's `agentic-cms update` run is mid-diagnosis as of session end — user was asked to check `.agentic-cms/VERSION`, `CONTENT.md`'s heading list, and `.agentic-cms/CONTENT.upstream.md`'s presence there; unresolved pending their reply.
- Tier 2 migration (`.agentic-cms/bin/` → `.agentic-cms/scripts/` in already-installed projects, from task 007) remains an open, unscoped gap, unchanged from prior sessions.
- The playbook and test report from this session's live verification (`.smaqit/user-testing/tests/011_classification-adoption-gaps.md`, `.smaqit/user-testing/2026-08-19_test-report_2.md`) were offered for commit to the user but not yet committed as of session end.

## Session Metrics

- **Duration:** ~1 session (design revision through release, live E2E verification, and follow-up Q&A)
- **Tasks completed:** 1 (011)
- **PRs merged:** 1 (#9)
- **Release shipped:** v0.7.0
- **Files created:** 7 (2 Go source, 2 Go test, 1 playbook, 1 test report, this history entry)
- **Files modified:** ~14 (engine scripts, hooks, skills, docs, task/planning state)
- **Test result:** PASS, 19/19 checks, verified against the real released artifact
