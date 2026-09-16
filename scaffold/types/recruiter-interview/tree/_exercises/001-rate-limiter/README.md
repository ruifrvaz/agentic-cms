# Take-home 001 — per-key rate limiter

You are given a small, unfinished Go module: a per-key request rate
limiter of the kind a backend service would put in front of a public API.
Your task is to complete it so that the provided test suite passes.

## Task

Two files: `ratelimiter.go` (the codebase, with gaps) and
`ratelimiter_test.go` (the specification — do not edit it). Make every
test pass.

There are three things to find and fix, each marked in the source:

1. `Limiter.Allow` — unimplemented (`TODO`). Token-bucket algorithm: refill
   based on elapsed time, cap at burst, consume one token if available.
2. `MultiLimiter.EvictIdle` — unimplemented (`TODO`).
3. `MultiLimiter.Allow` — implemented, compiles, passes under light load.
   Marked `BUG`. Read the comment above it. `go test -race` will **not**
   catch this one — it is not a data race, it is a logic bug.

## Rules

- Standard library only. No new imports beyond what is already there.
- Do not change the test file — it is your specification.
- Time box: **40 minutes**. Stop when it is up whether or not you are
  green, and note where you were.
- Keep a short log as you go — which item you started with, what each
  failing test told you, what you changed and why. Send it back with your
  code; we read it as carefully as the diff.

## Running it

```bash
cd <this folder>
go vet ./...
go test ./... -v          # see what is failing and why
go test ./... -race       # for the concurrency test specifically
```

All 6 tests passing (including 3 clean `-race` runs in a row) = done.

## Submitting

Reply with the modified `ratelimiter.go` and your log. Do not include the
test file.
