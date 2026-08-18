# Classification Rollout and Incident

**Date:** 2026-08-18
**Session focus:** Shipped confidentiality classification for installed
projects end-to-end (schema through enforcement), a scaffold rename it
depended on, real E2E testing against a clone of this repo, and response
to a live security near-miss the testing work surfaced.
**Tasks completed:** 004 (Archive content lifecycle, v0.3.0), 007 (Rename
`.agentic-cms/bin` → `.agentic-cms/scripts`, v0.4.0), 005 (First-class
content classification, v0.5.0)
**Tasks created:** 003 (Register mode for reference files, Not Started),
008 (Global classification scanner, Not Started)
**Tasks referenced:** 002 (Draft content state — already PR Open at
session start), 006 (Shell installer bootstrap — another session's
parallel work)

## Actions Taken

- Resolved task 002's push blocker (GitHub token lacked write access;
  retried after the user refreshed it) and confirmed its merge/release.
- Scoped two new capabilities from user requests: reassessed an initial
  "attach mode" design for reference documents into a leaner **Register
  mode** (task 003) after the user clarified the real scenario (e.g.
  invoices — index-only, no conversion); confirmed **Archive** (task 004)
  as the mirror of the existing Drafts convention.
- Implemented and shipped **task 004** — `status: archived`,
  `docs/<topic>/archive/`, `ac-page archive`, archived items re-filed
  under `wiki/index.md`'s `## Archived` section (index auto-creates the
  section on older files), `content-manage-item` Archive mode. Renamed
  `content-new-item` → `content-manage-item` as part of the same release,
  since the old name no longer described the skill's lifecycle scope.
- Created and heavily refined **task 005** (content classification):
  first pass covered write-time rating only; refined to cover
  update-drift (hooks, pre-commit gate, lint sweep); researched Claude
  Code/Codex/Copilot hook schemas directly against their docs; then, on
  the user's explicit SRP/KISS challenge, redesigned around a single
  `ac-classify` detection engine with every enforcement point as a thin
  caller — added an explicit "no integration point re-implements
  detection" acceptance criterion, and deferred Copilot's hook (no
  documented blocking capability) to a follow-up.
- Diagnosed and fixed a VS Code visibility issue (`files.exclude:
  **/bin/**` was hiding `.agentic-cms/bin/`). The user rejected patching
  the exclude and correctly reframed it as an architecture problem — a
  directory holding committed source colliding with the generic
  build-output convention. Assessed migration complexity (Tier 1 rename
  vs. Tier 2 update-time migration for already-installed projects); user
  chose Tier 1, sequenced **before** completing task 005.
- Implemented and shipped **task 007** — `git mv` plus a literal-string
  sweep of `.agentic-cms/bin` → `.agentic-cms/scripts` across every living
  reference, explicitly preserving historical records (changelog,
  compendium, session history) unchanged.
- Resumed and shipped **task 005**: built `ac-classify` (`check`/`sweep`/
  `hook`), Claude Code + Codex agent hooks, a git pre-commit gate
  (materializes the staged tree via `git checkout-index`, reuses the same
  engine), skill verify-tail wiring across four write-path skills plus the
  importer subagent, and content-lint/content-list integration.
  Reconciled the entire implementation against task 007's rename via
  `git mv` + sweep + `git reset --soft origin/main` (chosen over
  rebase/stash specifically to avoid git's rename-detection ambiguity).
- Created and ran a real E2E test playbook for task 005 — refined per the
  user's suggestion to install into a `git clone` of this actual repo
  (genuine brownfield target) rather than a synthetic `git init` sandbox.
  Executed all 4 steps live: build/test gate, install/verify, a 4-turn
  git-commit trigger sequence against the real pre-commit gate, and the
  agent-hook payload contract. Caught and corrected a flaw in the test's
  own Turn 4 fixture mid-run (reused an already-escalated page, masking
  the effect being tested) rather than accepting a misleading pass.
- Investigated and remediated a real security incident: the user
  committed a file with genuine PII (name, SSN, bank account, medical
  diagnosis) directly to this repo's `main`, unblocked, because this repo
  has never run `agentic-cms init` on itself. Confirmed the commit was
  unpushed, safely removed it via `git reset --hard origin/main` (stashing
  and restoring unrelated pre-existing uncommitted work from another
  session first), and disclosed residual considerations (local git object
  GC lag, conversation-transcript exposure).
- Scoped **task 008** (global classification scanner) directly from that
  incident: a global Claude Code skill installed on first binary run, plus
  an explicit warning-gated `agentic-cms hooks enable-global` git hook via
  `core.hooksPath` with chaining into any project's own per-project gate
  so task 005's richer behavior isn't silently downgraded.

## Problems Solved

- **GitHub token lost write access twice** mid-session (different token
  IDs each time) — resolved both times by retry after the user refreshed
  credentials; not a code issue.
- **`filepath.Join`'s split string arguments defeated a literal-string
  rename sweep** in `scaffold_test.go`, twice (tasks 007 and 005
  independently) — a plain grep for `".agentic-cms/bin"` doesn't match
  `filepath.Join(dir, ".agentic-cms", "bin", ...)`. Caught both times by
  the test suite actually failing, not by assumption.
- **Installer shipped `.agentic-cms/hooks/pre-commit` non-executable** —
  the executable-mode check only covered `.agentic-cms/bin/` (later
  `scripts/`), not the new `hooks/` directory; the installed hook
  silently no-op'd via its own `[ -x ]` guard until fixed.
- **A serious test-isolation leak**: `make smoke-test`'s sandboxes live
  inside this repo, so git's upward repo-discovery found *this repo's
  own* `.git/hooks/pre-commit` and the installer correctly-but-unwantedly
  wired the classification gate into it — a real side effect on the dev
  repo from running tests. Fixed via `GIT_CEILING_DIRECTORIES` scoped to
  the sandbox root; verified the accidentally-written hook was removed
  and stays absent on every subsequent run.
- **Task 007's own rename swept `.agentic-cms/bin` as a literal path but
  missed `bin` as a standalone word** in two tree-diagram lines
  (README.md, CONTENT.md) — already reached the released v0.4.0. Found
  and fixed as an incidental side effect of task 005's reconciliation
  pass over the same files, not a separate cleanup effort.
- **Test playbook's Turn 4 had a fixture flaw**: bypassing new PII into a
  page already rated at the ceiling (C3) can never demonstrate a floor
  violation, since the existing rating already covers anything lower.
  Not a product bug — `ac-classify`'s `clean` semantics (staleness
  advisory, floor-violation blocking) are correct by design; the test
  needed a fresh, lower-rated page to demonstrate the intended behavior,
  corrected live during execution.

## Decisions Made

- **Register mode, not "attach mode"** for reference documents — files
  stay in `raw/`, get a `wiki/sources/` sidecar (tags: [reference]) and
  nothing else; no new folder, no new mechanism.
- **SRP guardrail for task 005**: exactly one detection engine
  (`ac-classify`); every enforcement point is a thin caller varying only
  in response handling (block/warn/auto-fix/report). Made an explicit,
  checkable acceptance criterion, not just a stated intent.
- **Agent hooks scoped to Claude Code + Codex only** — both support
  blocking per verified docs; Copilot's hook deferred (no documented
  blocking capability), with the skill-tail check as its fallback there.
- **`.agentic-cms/bin` → `.agentic-cms/scripts`, Tier 1 only** — rename the
  scaffold source; no update-time migration for already-installed
  projects in this pass (same accepted gap as the earlier
  `content-new-item` → `content-manage-item` rename). Sequenced before
  task 005 completed, at the user's explicit direction, even though it
  meant reopening already-tested work.
- **`git reset --soft`, not rebase, for reconciling 005 against 007**:
  since the rename is a pure textual substitution, applying it
  independently to both branches' overlapping files converges without
  needing git's merge/rename-detection machinery at all — avoided
  entirely rather than risked.
- **Sensitive commit removal: `git reset --hard`, not `revert`** — since
  it was unpushed and at the tip, a revert would have left the PII
  sitting in a past commit forever; a hard reset to `origin/main` removed
  it from history completely.
- **Task 008's global hook: `core.hooksPath`, not `init.templateDir`** —
  the only mechanism that retroactively covers already-existing repos
  (including this one); its "silently overrides local hooks" risk is
  addressed by making activation explicit and warning-gated, plus
  building chaining into per-project gates as a non-optional default.

## Files Modified

Spans three merged PRs (#3, #4, #5) plus this session's task-planning
files. See each task's own PR description and `CHANGELOG.md`'s `## [0.3.0]`
through `## [0.5.0]` sections for the authoritative file lists — not
reproduced here in full given the volume (task 005 alone touched 28 files).
Key new files: `scaffold/tree/.agentic-cms/scripts/ac-classify`,
`scaffold/tree/.agentic-cms/hooks/pre-commit`, `scaffold/githook.go` (+
test), `scaffold/tree/.claude/settings.json`, `scaffold/tree/.codex/hooks.json`,
`.smaqit/user-testing/tests/005_first-class-content-classification.md`,
`.smaqit/user-testing/2026-08-18_test-report.md`,
`.smaqit/tasks/003_register_mode_reference_files.md`,
`.smaqit/tasks/008_global_classification_scanner.md`.

## Next Steps

- **Task 008** (global classification scanner) is specced and ready to
  start — covers the actual gap the security incident exposed.
- **Narrower, immediate alternative to 008** noted but not done: simply
  running `agentic-cms init` on this repo itself would give it task 005's
  project-scoped gate right away, independent of whether 008 ships.
- **Task 003** (Register mode) is specced and ready to start — no longer
  blocked (its stated dependency, task 002, merged).
- **Task 006** (Shell installer bootstrap) is in progress in a separate
  session; its "no global agent/skill directories to seed" design
  decision is now stale given task 008 — flagged in `PLANNING.md` for
  that session to revisit, not resolved here.
- **Tier 2 migration** (already-installed projects reconciling their own
  `.agentic-cms/bin/` to `.agentic-cms/scripts/`) remains an open,
  unscoped gap, same class as the earlier skill-rename gap.
- Two residual, non-blocking items from the security incident: the
  removed commit's raw bytes may persist in this repo's local `.git/`
  object database until normal GC; the sensitive content is also present
  in this session's conversation transcript, outside git's reach entirely.

## Session Metrics

- **Duration:** Single extended session, 2026-08-18
- **Tasks completed:** 3 (004, 007, 005) — releases v0.3.0, v0.4.0, v0.5.0
- **Tasks created:** 2 (003, 008)
- **PRs merged:** 3 (#3, #4, #5)
- **Real bugs found and fixed pre-release:** 4 (non-executable pre-commit
  script; test-isolation git-hook leak into this repo; two
  `filepath.Join`-split rename misses; two standalone-word rename misses
  in README/CONTENT.md)
- **Security incident:** 1 found, contained (unpushed), and remediated
  within the same session
