# User Testing Report

**Date:** 2026-08-18
**Repository:** ruifrvaz/agentic-cms
**Branch:** main
**Commit:** cde49f939f288f791cbe6d4c200badd8b9c88612
**OS/Arch:** Linux/x86_64
**Duration:** ~10 minutes

## Scope
- Test file: `.smaqit/user-testing/tests/005_first-class-content-classification.md`
- Commands executed:
   - `make build`, `make test`
   - `git clone --quiet "$(pwd)" "$D"` then `agentic-cms init "$D"`
   - `ac-page new`, `ac-index add`, `git commit` (Turn 1)
   - Planted a credential string, `git commit` (Turn 2, expected block)
   - `ac-page classify docs/ai/notes.md C3`, `git commit` (Turn 3)
   - `ac-page new docs/ai/expenses.md --classification C1`, planted PII,
     `git commit --no-verify`, `ac-classify sweep` (Turn 4, corrected mid-run)
   - `ac-classify hook` piped two simulated Claude Code/Codex payloads

## Checklist
- [x] Test command discovered and confirmed (`Makefile`: `build`, `test`)
- [x] Dependencies installed (none beyond Go/git/python3, already present)
- [x] Test suite executed (all 4 playbook steps, 17 checkpoints)
- [x] Results captured (pass/fail + key errors)
- [x] Evidence collected (full command/output pairs for every checkpoint)

## Execution Log (Timestamped)
- Step 1 — `make build` exit 0; `make test` exit 0 (go vet + go test, both packages ok)
- Step 2 — cloned this repo (commit cde49f9) into an isolated `mktemp -d`
  target; ran `agentic-cms init` against it. Every predicted line matched
  exactly: 39 created / 1 merged (`CLAUDE.md`) / 0 skipped, including the
  toolkit, both hook configs, and the git pre-commit hook, all fresh-created
  since this repo's own root carries none of them
- Step 2 verification — remote origin confirmed as the local clone path
  (not GitHub); both `ac-classify` and `.git/hooks/pre-commit` executable;
  hook marker present; `CLAUDE.md` correctly preserved `@AGENTS.md` with
  the managed block appended below it; this repo's own working tree and
  `.git/hooks/` confirmed untouched before and after
- Step 3 Turn 1 — clean C1 page committed cleanly, exit 0
- Step 3 Turn 2 — planted a credential-shaped string, committed: **blocked**,
  exit 1, correct implied level (C3) in the message; `git log --oneline`
  confirmed the blocked commit never landed
- Step 3 Turn 3 — `ac-page classify ... C3` then re-commit: exit 0, unblocked
- Step 3 Turn 4 (first attempt) — bypassed with `--no-verify` on the
  already-C3-rated `notes.md`; sweep reported `clean: true`. **This
  contradicted the playbook's stated expectation** — investigated live
  rather than accepted at face value
- Step 3 Turn 4 (corrected) — re-ran against a fresh `docs/ai/expenses.md`
  rated C1, planted the same PII, bypassed with `--no-verify`: exit 0;
  sweep now correctly reported `clean: false`, `floor_violation: true`,
  `implied_level: "C2"` — the intended behavior, on the right fixture
- Step 4 — `ac-classify hook` on the floor-violating `expenses.md`:
  auto-raised C1→C2, exit 0, no silent skip; on a `Bash`/non-content
  payload: `{"ok":true,"skip":true}`, exit 0
- Cleanup — temp clone removed; confirmed this repo's `.git/hooks/pre-commit`
  still absent (no leak from the test run)

## Results
- Overall: **PASS**
- Summary:
   - 17/17 playbook checkpoints passed after one live correction to the
     playbook itself (not the product)
   - Zero product defects found — every enforcement point (pre-commit
     gate, `ac-classify sweep`, hook payload contract, brownfield install
     paths) behaved exactly as the implementation and its documented
     design specify
   - The one discrepancy was in the test's own Turn 4 fixture choice
     (reusing an already-C3 page masked the effect being tested), not in
     `ac-classify`'s `clean` semantics, which are intentionally
     staleness-is-advisory / floor-violation-is-blocking — confirmed
     correct by re-deriving the case properly

## Pain Points
- Blockers:
   - None
- Issues:
   - None in the shipped code. The playbook's original Turn 4 fixture
     (bypassing new PII into a page already rated at the ceiling, C3)
     could not have demonstrated "sweep catches what slipped through" —
     any page already at C3 absorbs a lower-floor addition by design.
     This was a test-authoring gap, caught only by actually running it
- UX Friction:
   - None observed during execution — every documented command matched
     its actual behavior once the Turn 4 fixture was corrected
- Performance:
   - Full 4-step playbook (build, install, 4 commit turns, 2 hook checks)
     completed in well under a minute of actual command time; no live
     service or network dependency to wait on

## Recommendations
- None for the product — task 005's implementation matches its design
  exactly, including the subtler staleness-vs-floor-violation distinction
  that the corrected Turn 4 now demonstrates correctly
- Playbook fixed in place (`.smaqit/user-testing/tests/005_first-class-content-classification.md`):
  Turn 4 now uses a fresh, lower-rated page instead of reusing the
  already-escalated one from Turns 1–3, with a note explaining why that
  matters — future re-runs of this playbook will get the correct result
  on the first pass
- Worth carrying forward as a general playbook-authoring principle: when
  a test step depends on a *prior* step's mutated state (a page's current
  rating, here), explicitly check whether that state could already
  satisfy the condition being tested before asserting the expected outcome
