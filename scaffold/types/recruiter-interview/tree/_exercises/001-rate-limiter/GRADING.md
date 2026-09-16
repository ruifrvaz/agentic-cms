# Grading notes — take-home 001 (interviewer-only; never in the bundle)

The hand-over bundle is `README.md`, `ratelimiter.go`, `ratelimiter_test.go`,
`go.mod`. This file and any `reference/` solution stay in this repository.

## What the three items test

1. **`Limiter.Allow` (unimplemented)** — can the candidate implement a
   token bucket from a spec given only as tests: elapsed-time refill,
   burst cap, single-token consume, injectable clock respected. Correctness
   and edge cases: zero elapsed time, refill overshoot capped at burst,
   first call on a full bucket.
2. **`MultiLimiter.EvictIdle` (unimplemented)** — reads the surrounding
   code to infer the contract (per-key buckets, idle threshold), respects
   the existing locking discipline, doesn't leak keys.
3. **`MultiLimiter.Allow` (`BUG`)** — the discriminating item. It compiles
   and passes light-load tests; the failure is a logic bug that `-race`
   cannot see. Strong candidates triage the failing test by root cause
   rather than patching symptoms, and can explain *why* the race detector
   misses it.

## Process signals (from the candidate's log)

- Read the tests first and state what each asserts before touching code.
- Triaged failures by root cause, not test by test.
- Narrated hypothesis → evidence → fix → re-run.
- Minimal fixes; no wholesale refactor; no edits to the test file.
- Flagged anything in the tests that looked wrong instead of silently
  working around it.
- Clock awareness: stated where they were at the time box.

## Levels

- **Strong**: all three fixed within the box, item 3 explained correctly,
  log shows root-cause triage.
- **Adequate**: items 1-2 fixed, item 3 identified or partially fixed, log
  coherent.
- **Weak**: item 1 only, or fixes that special-case the tests, or a
  modified test file.

## Rubric mapping

Map each signal to the requisition's competency framework in the
`take-home-assignment` page under `docs/interview-design/` — this file is
the exercise-specific half; the competency weights live there.

## The interview follow-up

Once a submission is green, the live follow-up is: "why does `-race` miss
item 3?" — the answer distinguishes candidates who understood the bug from
those who found a passing patch.
