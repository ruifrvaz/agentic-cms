# Task Planning

## Active

| ID | Title | Status | Notes |
|----|-------|--------|-------|
| 003 | Register mode for reference files | Not Started | Lightweight `content-import` Register mode: wiki/sources sidecar + index + log for archival/reference files (e.g. invoices), no markdown conversion; start after 002 merges |
| 011 | Classification adoption gaps (first real-world rollout) | PR Open (#9) | Delta-scoped pre-commit gate; user-only hash-bound floor ack; currency floor widened (never narrowed — 0% FN ruling) with false positives resolved via ack; CONTENT.md schema reconciliation report + CONTENT.upstream.md sidecar; VERSION now auto-stamped from the binary. Targets v0.7.0. Lands in shared ac-classify engine — 008 inherits |
| 008 | Global classification scanner | Not Started | Global Claude Code skill (`~/.claude/skills/agentic-cms-classify/`, installs on first binary run) + optional warning-gated `agentic-cms hooks enable-global` git pre-commit hook via `core.hooksPath`, chaining into per-project `.agentic-cms/hooks/pre-commit` when present. Motivated by a real incident: this repo (never `init`-ed on itself) had zero coverage when a direct git commit landed real PII, caught only by manual review. Depends on task 005 (Completed, v0.5.0). NOTE: will supersede task 006's "no global agent/skill directories to seed" design decision — revisit that bullet |

## Completed

| ID | Title | Completed | Notes |
|----|-------|-----------|-------|
| 001 | Installer smoke test script | 2026-07-31 | Bash end-to-end smoke test of `agentic-cms init` via `make smoke-test` |
| 002 | Draft content state | 2026-08-15 | PR #2; `status: draft` frontmatter + `docs/<topic>/drafts/` convention, `ac-page --status`/`promote` toolkit support; shipped as part of v0.2.0 |
| 004 | Archive content lifecycle | 2026-08-16 | PR #3; `status: archived` + `docs/<topic>/archive/` mirroring drafts, `ac-page archive`, archived items stay indexed under `## Archived`; renamed content-new-item → content-manage-item; released as v0.3.0 |
| 007 | Rename .agentic-cms/bin to .agentic-cms/scripts | 2026-08-18 | PR #4; `bin/` collided with the generic IDE build-output-exclude convention but held committed source; Tier 1 rename (scaffold source + all living references) via `git mv` + literal-string sweep; released as v0.4.0. Task 005 must reconcile its own ~150 bin/ references against this baseline before completing |
| 005 | First-class content classification (auto CIA rating) | 2026-08-18 | PR #5; single `ac-classify` engine (`classified-hash` staleness + heuristic floors, raise-only) with four thin callers — Claude Code/Codex agent hooks, a git pre-commit gate on staged blobs, skill verify tails as the guaranteed fallback, content-lint's sweep; `ac-page --classification`/`classify`, `ac-inventory` distribution tally; reconciled against task 007's bin→scripts rename; released as v0.5.0. Segregation/vault out of scope |
| 006 | Shell installer bootstrap (no Go dependency) | 2026-08-18 | PR #6; `curl \| bash` `install.sh` (no Go toolchain required, `AGENTIC_CMS_VERSION` pins a release) modeled selectively on `smaqit-extensions/install.sh`; ancestor-dir resolution added to `agentic-cms update`; README restructured to mirror `ruifrvaz/smaqit`'s section style; released as v0.6.0. Also fixed: repo was accidentally private, breaking every install path — user made it public, both `install.sh` and `update` re-verified live |
| 009 | Surface classification feature in README | 2026-08-19 | PR #7; headline Features bullet naming the CIA triad + dedicated Classification section (C0-C3 table, "why it matters" for agent-driven synthesis); docs-only, released as v0.6.1. A follow-on direct commit to `main` ("cleaned up readme") further reorganized the README post-merge — user's own edit, outside this task's scope |
| 010 | Always overwrite scaffolding logic files on init/update | 2026-08-19 | PR #8; `scaffold.Install()` now always overwrites framework-owned scaffolding paths (`.claude/skills/`, `.claude/agents/`, `.agentic-cms/templates/`, `.agentic-cms/scripts/`, `.agentic-cms/hooks/`, `.codex/`) on every install run instead of skipping them once they exist, reported under a new `Result.Updated` bucket; user-content paths unchanged. No drift/checksum detection, per scope. Released as v0.6.2 |

## Abandoned

| ID | Title | Date | Reason |
|----|-------|------|--------|
