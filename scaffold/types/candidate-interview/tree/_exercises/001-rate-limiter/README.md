# Coding drill 001 — per-key rate limiter

Simulates a HackerRank-style coding-round format: "complete an unfinished
codebase by fixing unit tests," in Go. Stdlib-only Go, no third-party deps —
a per-key request rate limiter, the kind of thing a real backend service
plausibly needs in front of a public API.

## Task

Two files: `ratelimiter.go` (the codebase, with gaps) and
`ratelimiter_test.go` (the spec — do not edit it). Make every test pass.

There are three things to find and fix, each marked in the source:

1. `Limiter.Allow` — unimplemented (`TODO`). Token-bucket algorithm: refill
   based on elapsed time, cap at burst, consume one token if available.
2. `MultiLimiter.EvictIdle` — unimplemented (`TODO`).
3. `MultiLimiter.Allow` — implemented, compiles, passes under light load.
   Marked `BUG`. Read the comment above it for a pointed hint. `go test
   -race` will **not** catch this one — it's not a data race, it's a logic
   bug. You have to reason about it.

## Rules (matching the real round)

- Stdlib only. No new imports beyond what's already there.
- Don't change the test file — it's your spec, exactly like a real
  HackerRank harness.
- **Set a 40-minute timer before you start** (a guess at this stage's
  per-section split of the 90 minutes — the guide doesn't say — but a
  reasonable one to rehearse against). Stop when it goes off, whether or
  not you're green, and note where you were.

## Running it

```bash
cd <this folder>
go vet ./...
go test ./... -v          # see what's failing and why
go test ./... -race       # for the concurrency test specifically
```

All 6 tests passing (including 3 clean `-race` runs in a row) = done.

## When you're done (or time's up)

Come back and tell me: which of the 3 you fixed, in what order, where you
got stuck, and how much time each took. I'll debrief it like a rehearsal
session — including the part of this bug a real panel would ask about out
loud once your tests are green ("why does `-race` miss this?").
