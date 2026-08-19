# Classification adoption gaps (first real-world rollout) — E2E Test Playbook

**Test ID:** 011
**Title:** Classification adoption gaps (first real-world rollout)
**Date:** 2026-08-19
**Tester:** User Testing Agent
**Task:** 011

## Objectives

Validate all five of task 011's repaired gaps against the **real published
v0.7.0 release binary**, not just the Go unit tests' synthetic `t.TempDir()`
fixtures or the task-worktree verification already run during
implementation: (1) the pre-commit gate blocks only the staged delta and
summarizes pre-existing drift, (2) a user floor-ack silences a false
positive and any edit invalidates it, (3) the currency floor is widened
(never narrowed) and still catches every enumerated shape, (4) `init`
reconciles a customized `CONTENT.md` against the shipped schema without
ever editing the user's file, and (5) `.agentic-cms/VERSION` is
auto-stamped with the real release version. Mirrors this project's
established convention of testing the actual installed artifact — see
`.smaqit/user-testing/tests/010_always-overwrite-scaffolding-logic-files-on-init-update.md`.

## Prerequisites

- Go toolchain, `git`, `python3`, and `bash` on `PATH`
- This repository's primary checkout in a normal, clonable state
- Network access to fetch the real `v0.7.0` release via `install.sh`

## Test Steps

### Step 1 — Build & Unit Test Gate

- [x] `make build` exits 0
- [x] `make test` exits 0 — runs `go vet ./...` and `go test ./...`, zero
      failures, including the new `TestCurrencyFloorRecall`,
      `TestFloorAckLifecycle`, `TestPreCommitDeltaScope`,
      `TestInstallStampsVersion`, and the two `TestReconcileContentMD*` cases
- [x] `make smoke-test` exits 0 — the extended installer smoke test,
      including its new ack and CONTENT.md-reconciliation sections

### Step 2 — Install the real release & verify VERSION stamp (AC7)

```sh
curl -fsSL https://raw.githubusercontent.com/ruifrvaz/agentic-cms/main/install.sh | bash
BIN="$HOME/.local/bin/agentic-cms"
"$BIN" version
D=$(mktemp -d)
git -C "$D" init -q
git -C "$D" config user.email test@test.local
git -C "$D" config user.name test
"$BIN" init "$D" | tee /tmp/011-init-out.txt
cat "$D/.agentic-cms/VERSION"
```

- [x] `"$BIN" version` reports `v0.7.0`
- [x] `$D/.agentic-cms/VERSION` reads exactly `v0.7.0` (or the resolved dev
      string if the binary was built locally instead — record which)
- [x] Init output shows no leftover `{{VERSION}}` placeholder anywhere

### Step 3 — Delta-scoped pre-commit gate (AC1, AC2)

```sh
cd "$D"
mkdir -p docs/vendors docs/finance docs/notes
printf -- '---\ntitle: Vendor pricing\nclassification: C1\n---\nAcme: $100/month.\n' > docs/vendors/pricing.md
printf -- '---\ntitle: Cap table\nclassification: C1\n---\nAlice invested $2,500,000 at seed.\n' > docs/finance/captable.md
git add -A && git commit -qm "legacy backlog" --no-verify

printf -- '---\ntitle: Standup notes\nclassification: C1\n---\nNo sensitive shapes here.\n' > docs/notes/standup.md
git add docs/notes/standup.md
git commit -m "add standup notes" 2>&1 | tee /tmp/011-scenario1.txt
```

- [x] Scenario 1 commit **succeeds** (exit 0)
- [x] Output shows exactly **one** summarized backlog warning line (not
      one line per pre-existing violation)

```sh
echo 'New tier: $500/month.' >> docs/vendors/pricing.md
git add docs/vendors/pricing.md
git commit -m "update pricing" 2>&1 | tee /tmp/011-scenario2.txt; echo "exit=$?"
```

- [x] Scenario 2 commit is **blocked** (non-zero exit), naming
      `docs/vendors/pricing.md` and mentioning both the raise option and
      the `--ack-floor` remedy

```sh
git reset -q
ENTRY=$(cat wiki/index.md)
printf '%s\n- [Budget](docs/finance/captable.md) — total is $4,000,000\n' "$ENTRY" > wiki/index.md
git add wiki/index.md
git commit -m "index update" 2>&1 | tee /tmp/011-scenario3.txt; echo "exit=$?"
git checkout -- wiki/index.md
```

- [x] Scenario 3 commit is **blocked** — bleed into a staged
      `wiki/index.md` is caught even though `docs/finance/captable.md`
      itself is not staged in this commit

### Step 4 — Floor acknowledgment lifecycle (AC3, AC4)

```sh
.agentic-cms/scripts/ac-classify check docs/vendors/pricing.md
.agentic-cms/scripts/ac-page classify docs/vendors/pricing.md C1 --ack-floor | tee /tmp/011-ack-out.txt
git add docs/vendors/pricing.md
git commit -m "acked pricing" 2>&1 | tee /tmp/011-scenario4.txt; echo "exit=$?"
echo 'Enterprise: $900/month.' >> docs/vendors/pricing.md
git add docs/vendors/pricing.md
git commit -m "edit after ack" 2>&1 | tee /tmp/011-scenario5.txt; echo "exit=$?"
```

- [x] `ac-classify check` before acking reports `"floor_violation": true`
- [x] `ac-page classify ... --ack-floor` output includes
      `"classification-ack"`
- [x] Scenario 4 commit (acked, unchanged body) **succeeds**
- [x] Scenario 5 commit (body edited after ack) is **blocked again** — the
      ack does not survive a content change

- [x] Skill text confirms acking is documented as user-only, never an
      agent decision:
  ```sh
  grep -rl "ack-floor" "$D/.claude/skills/"
  grep -A2 "ack-floor" "$D/.claude/skills/content-manage-item/SKILL.md" | grep -i "user"
  ```

### Step 5 — Currency floor recall (AC5)

```sh
for CASE in '100 EUR total' '250000 USD budget' '9,50€ per seat' '1,000 dollars paid'; do
  printf -- '---\ntitle: fixture\nclassification: C1\n---\n%s\n' "$CASE" > docs/notes/fixture.md
  echo "== $CASE =="
  .agentic-cms/scripts/ac-classify check docs/notes/fixture.md | grep -o '"floor_violation": [a-z]*'
done
```

- [x] All four widened currency shapes (`USD`/`EUR` ISO-prefixed,
      symbol-suffixed, word-form) report `"floor_violation": true`

### Step 6 — CONTENT.md schema reconciliation (AC6)

```sh
sed -i 's/^## Classification$/## My Local Rules/' "$D/CONTENT.md"
"$BIN" init "$D" 2>&1 | tee /tmp/011-reconcile-out.txt
cat "$D/.agentic-cms/CONTENT.upstream.md" | grep -c "## Classification"
grep -c "## My Local Rules" "$D/CONTENT.md"
grep -c "## Classification" "$D/CONTENT.md"
```

- [x] Init output includes a "CONTENT.md reconciliation" report naming
      `Classification` as a missing section
- [x] `.agentic-cms/CONTENT.upstream.md` exists and contains the full
      upstream `## Classification` section
- [x] The user's `CONTENT.md` still contains `## My Local Rules` and does
      **not** contain `## Classification` — confirms the user's file was
      never edited

## Pass/Fail Criteria

**PASS** — All checkboxes checked. Steps 1–2 commands exit 0. Every
pre-commit scenario in Steps 3–4 produces the documented block/pass
outcome, verified by actual git exit codes, not just printed text. Step 5's
four widened currency forms all report `floor_violation: true`. Step 6's
reconciliation report appears and the user's `CONTENT.md` is provably
unedited (`grep -c` on both markers).

**FAIL** — Any checkbox unchecked; any scenario in Steps 3–5 produces the
wrong block/pass outcome; the reconciliation report is missing, or the
user's `CONTENT.md` was modified, or the sidecar is missing/incomplete.

## Evidence to Capture

- Full output of `make test` and `make smoke-test`
- `"$BIN" version` output and `$D/.agentic-cms/VERSION` contents
- All five pre-commit scenario transcripts (`/tmp/011-scenario{1..5}.txt`)
  with their exit codes
- `ac-classify check` / `ac-page classify --ack-floor` JSON output
- Step 5's four `floor_violation` results
- Step 6's reconciliation report text and the diff (or lack thereof) in
  the user's `CONTENT.md`
