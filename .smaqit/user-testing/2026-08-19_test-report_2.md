# User Testing Report

**Date:** 2026-08-19
**Repository:** ruifrvaz/agentic-cms
**Branch:** main
**Commit:** 4cfd922
**OS/Arch:** Linux/x86_64
**Duration:** ~25 minutes

## Scope

- Test file: `.smaqit/user-testing/tests/011_classification-adoption-gaps.md`
- Commands executed:
  - `make build`, `make test`, `make smoke-test`
  - `curl -fsSL https://raw.githubusercontent.com/ruifrvaz/agentic-cms/main/install.sh | bash` (real published `v0.7.0` release, not a local build)
  - `agentic-cms version`, `agentic-cms init` (fresh install, then re-init over a customized `CONTENT.md`)
  - `.agentic-cms/scripts/ac-classify check`, `.agentic-cms/scripts/ac-page classify ... --ack-floor`
  - Five real `git commit` scenarios against the installed pre-commit gate
  - Currency-floor fixture checks for four widened shapes

## Checklist

- [x] Test command discovered and confirmed (`Makefile`: `build`, `test`, `smoke-test` targets)
- [x] Dependencies installed (Go toolchain, git, python3, bash — all present)
- [x] Test suite executed (`make test`, `make smoke-test`, plus the live playbook scenarios)
- [x] Results captured (pass/fail + key errors, all command output tee'd to `/tmp/011-*.txt`)
- [x] Evidence collected per test file — all 19 playbook checkboxes verified against real command output

## Execution Log (Timestamped)

- Step 1 — `make build` and `make test`: both exit 0, all packages `ok`, including the five new test cases added by task 011
- Step 1 — `make smoke-test`: all checks pass, including the two new sections (floor ack, CONTENT.md reconciliation)
- Step 2 — installed the real `v0.7.0` release via `install.sh`; `agentic-cms version` reports `v0.7.0`; fresh `init` into a `mktemp -d` target; `.agentic-cms/VERSION` reads exactly `v0.7.0`, no `{{VERSION}}` placeholder leak
- Step 3, scenario 1 — committed a clean staged file against a tree with 2 pre-existing floor violations: exit 0, exactly one summarized backlog warning (`4 pre-existing classification issue(s)...`), not per-file noise
- Step 3, scenario 2 — staged the actual floor-violating file: exit 1, block message names the file and offers both the raise command and the `--ack-floor` remedy
- Step 3, scenario 3 — staged a `wiki/index.md` edit leaking a C2 page's figures: exit 1, blocked citing `wiki/index.md:33` even though the leaking source page (`docs/finance/captable.md`) itself was not staged
- Step 4 — `ac-classify check` on the pricing page reported `floor_violation: true` pre-ack; `ac-page classify ... --ack-floor` stamped `classification-ack`; the acked commit (scenario 4) succeeded after a test-harness fix (see Pain Points); a further body edit (scenario 5) re-blocked the commit, confirming the ack does not survive a content change
- Step 4 — confirmed via `grep` that all five write-path skills reference `ack-floor`, and `content-manage-item`'s text explicitly frames it as a user decision
- Step 5 — all four widened currency shapes (`100 EUR`, `250000 USD`, `9,50€`, `1,000 dollars`) reported `floor_violation: true`
- Step 6 — renamed the installed `## Classification` heading to `## My Local Rules`, re-ran `init`: reconciliation report printed naming `Classification` as missing, `.agentic-cms/CONTENT.upstream.md` was written containing the real `## Classification` section, and `grep -c` confirmed the user's `CONTENT.md` still had `## My Local Rules` and zero occurrences of `## Classification` — the file was provably never touched

## Results

- Overall: **PASS**
- Summary:
  - 19/19 playbook checkboxes verified against real command output
  - All 8 of task 011's acceptance criteria exercised live against the actual published `v0.7.0` binary, not just unit tests or the task-worktree verification done during implementation
  - Zero product defects found

## Pain Points

- Blockers:
  - None
- Issues:
  - None found in the product itself
- UX Friction:
  - One test-harness-only mistake during execution (not a product bug): after scenario 3's blocked bleed-check commit, `git checkout -- wiki/index.md` was used to restore the working tree, but that command restores from the **index**, not `HEAD` — since the leaked content was already staged, it silently remained staged. The next commit attempt (scenario 4, the ack) then failed with the *old* scenario 3 bleed error, momentarily looking like an ack-mechanism bug. Diagnosed via `git status --short` and fixed with `git reset -q wiki/index.md` before re-restoring. Worth noting for future playbook authors: always pair a blocked/aborted scenario's cleanup with `git reset` (unstage) before `git checkout --` (restore working tree), not `checkout --` alone.
- Performance:
  - No concerns; the entire pre-commit gate (materializing the staged tree via `git checkout-index` and running the full `ac-classify sweep`) completed well under a second on every scenario, even with the backlog present.

## Recommendations

- Fold the "unstage before restore" cleanup pattern above into a shared helper or explicit callout in this project's own playbook-writing convention, since it's specific to testing a gate that blocks commits by design and will recur in future pre-commit-gate playbooks.
- No product-code recommendations — all five of task 011's gaps behave exactly as specified against the real released artifact.
