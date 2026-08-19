---
status: PR Open
created: "2026-08-19"
started: "2026-08-19"
mode: Assisted
pr: 9
---

# Classification adoption gaps (first real-world rollout)

## Description

Repair the product gaps surfaced by the classification feature's (task 005,
v0.5.0) first real-world adoption in an installed brownfield project — a
repo with ~50 pre-existing pages, ~30 of them tripping heuristic floors on
day one. The gates *work* (the pre-commit hook was verified live: blocked
with exit 1 on staged C2-shaped content), but four gaps make honest
adoption either painful or impossible without `--no-verify`:

1. **The pre-commit gate blocks on the whole staged tree, not the staged
   delta.** `hooks/pre-commit` materializes the entire index
   (`git checkout-index -a`) and sweeps ALL of docs/wiki — so committing
   one unrelated file is blocked by every pre-existing floor violation
   anywhere in the project. On a brownfield install, that means *no
   commit of anything* until the entire legacy backlog is rated. The
   design language ("checks staged blobs") is technically satisfied but
   the practical effect punishes adoption and trains users toward
   `--no-verify`.
2. **No floor-override mechanism — false positives are permanent blocks.**
   A floor violation's only offered remedy is "raise it." If the user
   judges the pattern a false positive (see gap 3) and deliberately rates
   the page C1, `floor_violation` fires forever (rated C1, implied C2).
   Worse: obeying the gate means over-classifying to C2 — and under any
   placement policy that means *removing the page from the repo*, an
   absurd outcome for a public vendor-pricing table. The ratchet's own
   escape ("only the user may lower") exists as doctrine but has no
   mechanical representation the engine honors.
3. **Currency-amount floor is over-broad.** A bare currency figure
   (`$100/month`, `€9/user`) implies C2 — flagging ~20 vendor-pricing
   decision docs in the adopting project, all public list prices. Real
   confidential-finance pages (cap tables, budgets) were correctly
   flagged, but the noise ratio (~2/3 false positives) undermines trust
   in the floor mechanism.
4. **Schema evolution never reaches installed projects' CONTENT.md.** The
   non-destructive installer (correctly) never overwrites a customized
   `CONTENT.md` — but there is no reconciliation surface at all: after
   `update`, the adopting project had the toolkit, hooks, and templates
   live while `CONTENT.md` contained zero mention of classification. The
   agent-facing half of the feature (rubric, placement rules, bleed rule,
   ratchet) silently fails to arrive for every pre-existing install.

Plus one hygiene item: **`scaffold/tree/.agentic-cms/VERSION` reads
`0.1.0`** while the product is at v0.6.1 — installed projects cannot tell
which scaffold generation they run, which is exactly the signal a
reconciliation surface (gap 4) would key on.

## Issue Triage Context

**Mode:** Auto
**Technologies:** Git, Claude Code, OpenAI Codex CLI
**Platforms/Environments:** Linux
**Features/Integrations:** ac-classify engine, git pre-commit gate, `agentic-cms update`
**Versions/Constraints:** Baseline v0.6.1; must not weaken the gate for *newly staged* violations — all changes concern pre-existing drift, overrides, and delivery

## Design Decisions

- **Gate scope: block on the staged delta, warn on the rest.** The
  pre-commit gate should BLOCK only for violations in files actually
  changed in the commit (`git diff --cached --name-only`), and demote
  tree-wide pre-existing drift to a summarized WARN (count + `ac-classify
  sweep` pointer, not 30 lines of noise). `content-lint` remains the
  audit that owns the backlog. This preserves the exfiltration-boundary
  guarantee for new/edited content while making brownfield adoption
  incremental instead of all-or-nothing.
- **Floor acknowledgment: explicit, user-only, hash-anchored.** Add an
  ack mechanism (e.g. `ac-page classify <path> C1 --ack-floor`, stamping
  `classification-ack:` frontmatter bound to the same normalized-body
  hash as `classified-hash:`) recording "user reviewed the floor hit and
  keeps this rating." `ac-classify` treats an acked page's floor
  violation as a non-blocking note; any content change invalidates the
  ack together with the hash, forcing re-review. Skills must instruct
  agents that acking is a user decision to relay, never an agent
  decision — the ratchet's asymmetry stays intact.
- **Currency floor stays maximally broad — recall is the floor's
  contract (REVISED 2026-08-19):** the floor must have a 0% false
  negative rate over the pattern classes it enumerates, so detection is
  never narrowed to buy precision. Bare currency amounts keep implying
  the C2 floor; all false-positive relief routes through the hash-bound
  ack mechanism (previous bullet), never through weaker patterns. This
  step instead WIDENS recall: the current currency regex only matches
  symbol-prefixed figures (`$100`), missing `USD 100`, `100 EUR`, and
  word-form amounts (`1,000 dollars`) — extend it to ISO-code,
  symbol-suffixed, and word forms. The widening propagates to the
  index/log bleed scan automatically (it reuses `floor_level()`).
  Scope of the guarantee: regexes give 0% FN only for enumerated
  shapes; semantically confidential content with no detectable shape
  remains the agent rubric's responsibility. Watch item: ack fatigue —
  the delta-scoped gate and per-page, logged, hash-invalidated acks are
  the mitigations.
- **Schema reconciliation surface for CONTENT.md:** `update` (and `init`
  on existing installs) should compare the installed CONTENT.md against
  the shipped one and, when upstream schema sections are missing, emit a
  reconciliation report naming the absent sections and write the current
  upstream copy to `.agentic-cms/CONTENT.upstream.md` for manual merge.
  Explicitly NOT auto-editing the user's CONTENT.md — it is
  user-co-evolved by design; precedent for the sidecar-plus-report shape
  over managed blocks: CONTENT.md has no managed-marker convention,
  unlike CLAUDE.md's `agentic-cms:begin/end` block.
- **VERSION stamps the release:** `scaffold/tree/.agentic-cms/VERSION`
  carries the real release version and the release process bumps it;
  `update` uses it for the reconciliation report ("installed scaffold
  vX, shipping vY").
- **Out of scope:** the global scanner (task 008), Copilot hooks (005's
  deferred follow-up), and any adopting project's own backlog rating —
  that's per-project work the delta-scoped gate makes incremental.

## Implementation Steps

1. Rework `scaffold/tree/.agentic-cms/hooks/pre-commit`: sweep the staged
   snapshot as today, but partition results into staged-delta files
   (BLOCK on floor/invalid/bleed) vs. rest-of-tree (single summarized
   WARN line with counts). Keep bleed checks blocking for staged
   index/log changes regardless of which page leaked.
2. Add the floor-ack mechanism to `ac-classify` (honor
   `classification-ack:` bound to the body hash) and `ac-page`
   (`--ack-floor` on `classify`, user-decision-only semantics documented
   in the skill instructions and CONTENT.md's ratchet section).
3. Widen the currency floor pattern in `ac-classify` (ISO-code,
   symbol-suffixed, and word-form amounts) without narrowing any
   existing pattern; add engine test fixtures: public-pricing page
   (trips the floor, passes only when acked), cap-table page (trips the
   floor), and the widened currency forms (all caught).
4. Implement the CONTENT.md reconciliation report + `CONTENT.upstream.md`
   sidecar in `update` (and `init` over an existing install), keyed on
   missing schema section headings; cover with installer tests.
5. Stamp `scaffold/tree/.agentic-cms/VERSION` with the real release
   version and wire the bump into the release checklist; `update`'s
   reconciliation report cites installed-vs-shipped versions.
6. Update CONTENT.md (shipped) + affected skill texts for the ack
   semantics and the delta-scoped gate description; `make smoke-test`
   and `go test ./...` green; manually re-run the brownfield scenario
   (one staged clean file + dirty tree → commit passes with summary
   warning; staged floor violation → blocked; acked false positive →
   passes; ack invalidated by edit → blocked again).

## Known Issues Triage
**Triaged:** 2026-08-19
**Tools searched:** Git (`git/git`), Claude Code (`anthropics/claude-code`), OpenAI Codex CLI (`openai/codex`)
**Result:** Advisory

### Advisory Issues
- [#77341 PostToolUse hooks do not fire in daemon/background-job sessions (PreToolUse + SessionStart do)](https://github.com/anthropics/claude-code/issues/77341) — `anthropics/claude-code` — opened 2026-07-14 — bug, has repro, platform:linux, area:hooks. Partial dimension match: affects hook *dispatch* in daemon contexts, not the ac-classify callee this task changes; the skill verify-tail fallback is the designed mitigation for exactly this class of gap.
- [#84439 PostToolUse hooks in settings.json never register; documented updatedToolOutput example silently fails for structured-output tools](https://github.com/anthropics/claude-code/issues/84439) — `anthropics/claude-code` — opened 2026-08-06 — no labels. Title scope appears limited to structured-output tools (this project's settings.json PostToolUse hook was verified firing live in session 002); ambiguous, so kept Advisory.
- [#73586 PostToolUse hooks never fire for MCP tool calls — native tools fine](https://github.com/anthropics/claude-code/issues/73586) — `anthropics/claude-code` — opened 2026-07-02 — bug, has repro, area:mcp, area:hooks, platform:wsl. Not our path: content writes go through native file tools.
- [#27133 Project-level .codex/hooks.json is silently ignored when Codex runs inside a git worktree](https://github.com/openai/codex/issues/27133) — `openai/codex` — opened 2026-06-09 — bug, CLI, hooks, config. Directly concerns this project's Codex enforcement layer — and this task is implemented in a git worktree, so Codex-side hook verification must not be trusted from the worktree; verify via pre-commit gate + verify tails instead.
- [#18067 Codex CLI hooks fail silently on Linux/Windows when editing large files](https://github.com/openai/codex/issues/18067) — `openai/codex` — opened 2026-04-16 — bug, windows-os, hooks.

### Historical (Closed)
- [#87542 PostToolUse hook exit-2/stderr warning not delivered to model in VSCode extension](https://github.com/anthropics/claude-code/issues/87542) — `anthropics/claude-code` — closed 2026-08-18. Relevant precedent: our hook's block path uses exit 2 + stderr; fixed upstream.
- [#83877 Docs: hooks reference doesn't document the format of tool_input.file_path (absolute)](https://github.com/anthropics/claude-code/issues/83877) — `anthropics/claude-code` — closed 2026-08-05. ac-classify's hook path already normalizes absolute→relative.
- [#26452 codex exec still does not dispatch hooks with valid hooks.json shape](https://github.com/openai/codex/issues/26452) — `openai/codex` — closed 2026-08-06.

### Unresolvable Tools
- None — but the resolver helper misresolved "OpenAI Codex CLI" twice (`router-for-me/CLIProxyAPI`, `steipete/CodexBar`); `openai/codex` was taken from agent knowledge instead.

### Omitted Tools
- GitHub REST API (Releases) — documentation anchor for `agentic-cms update`, not an issue-searchable dependency
- Python 3 (re module) — standard library; project-layer baseline, no meaningful issue-search dimension for this task

### Search Warnings
- None — all six searches (open+closed × 3 repos) returned complete results; `git/git` returned zero matches both ways.

## Acceptance Criteria

- [x] Committing a clean staged change in a project with pre-existing tree drift succeeds, emitting one summarized warning (not per-file noise); staged-delta violations still block
- [x] Bleed into staged wiki/index.md or wiki/log.md still blocks regardless of source page
- [x] `ac-page classify <path> <level> --ack-floor` stamps a body-hash-bound ack; `ac-classify` reports the floor hit as non-blocking for an acked, unchanged page; any body change invalidates the ack
- [x] Skill and CONTENT.md text state acking is a user decision only (ratchet preserved)
- [x] Floor recall is never narrowed: bare list prices still trip the C2 floor and pass only via an explicit user ack; widened currency forms (`USD 100`, `100 EUR`, word amounts) are caught (fixture-tested)
- [x] `update` over an install with a customized CONTENT.md missing upstream schema sections emits a reconciliation report and writes `.agentic-cms/CONTENT.upstream.md`; never edits the user's CONTENT.md
- [x] `scaffold/tree/.agentic-cms/VERSION` matches the release version; release checklist includes the bump
- [x] `make smoke-test` and `go test ./...` pass

## Findings

**Implementation approach:**
- Engine changes (`ac-classify`, `ac-page`) implement gaps 2 and 3 (ack mechanism, widened currency recall); `hooks/pre-commit` implements gap 1 (delta-scoped blocking, tree-wide backlog summarized); a new `scaffold/reconcile.go` implements gap 4 (CONTENT.md section-heading diff, upstream sidecar, never edits the user's file); `scaffold/embed.go` + `main.go` implement the VERSION stamp (gap 5) by making `.agentic-cms/VERSION` a `{{VERSION}}` placeholder Install() always overwrites with the installing binary's ldflags version.
- Added three new test files (`engine_test.go`, `reconcile_test.go`, plus VERSION coverage in `scaffold_test.go`) exercising the real installed bash/python scripts and a real git repo end-to-end — the engine had no fixture tests before this task.
- Documentation pass across CONTENT.md (promoted Classification from a Conventions bullet to its own `##` section covering the ack and delta-gate), the scripts README contract, five write-path/lint skills (ack is user-only, never agent-initiated), and the root README (also fixed stale "files are never overwritten" claims left over from task 010/v0.6.2).
- Verified live end-to-end with the built binary against a real git sandbox: clean-staged-file-with-dirty-tree (passes, one summary warning), staged floor violation (blocks), user ack (passes), ack invalidated by edit (blocks again) — all four brownfield scenarios from Implementation Step 6 confirmed against real git commits, not just unit tests.

**Decisions made:**
- Currency floor was **widened**, not narrowed, per an explicit mid-task user ruling: the floor's contract is 0% false negatives over its enumerated pattern classes (symbol-prefixed/suffixed, ISO-coded, word-form); all false-positive relief routes through the user-only ack, never a weaker regex. This reverses the task file's original gap-3 framing (session-start commentary proposed context-narrowing; the user corrected that before implementation began) — the Design Decisions section above already carries the revision marker.
- VERSION stamping is automated (`{{VERSION}}` placeholder, always overwritten from the binary's own version at every `init`/`update` run) rather than a manually-bumped literal plus a release-checklist step. No release-checklist document exists in this repo to add a step to, and automatic stamping satisfies the AC's intent — VERSION always matches the release — more strongly than a manual step that could be forgotten (the original hygiene bug this gap reports). Flagged to the user for review; not walked back since.
- `ac-page classify` without `--ack-floor` withdraws any standing ack (not just leaves it stale) — a plain re-rate is a fresh rating decision that supersedes an old ack, matching the ratchet's spirit that classification decisions are current-state, not additive.

**Blockers encountered:**
- None. One self-caught test bug during implementation: an early version of `TestPreCommitDeltaScope`'s "staged floor violation blocks" scenario re-`git add`ed an already-committed unchanged file, producing no staged delta — the gate correctly ignored it, which looked like a false failure until the test itself was fixed to actually edit the file before staging.
- One latent bug fixed opportunistically while rewriting the hook: the old pre-commit's early-exit check only looked for staged `docs/`/`wiki/(entities|concepts|sources)/` paths, so a commit touching only `wiki/index.md` or `wiki/log.md` skipped the gate entirely — silently exempting the exact bleed-into-index/log case gap 1's own acceptance criterion requires blocking. Included in this task's `hooks/pre-commit` rewrite and covered by `TestPreCommitDeltaScope`'s scenario 4.

**Follow-up identified:**
- Tier 2 migration (already-installed projects reconciling `.agentic-cms/bin/` → `.agentic-cms/scripts/`, from task 007) remains open and unrelated to this task.
- Task 008 (global classification scanner) inherits the ack mechanism and widened floor for free via the shared `ac-classify` engine, as anticipated in this task's Notes.
- The adopting brownfield project's own legacy backlog (the ~30 floor hits, ~20 of them pricing false positives) still needs its owner to walk through and ack/re-rate — that was always explicitly out of scope for this task (per Design Decisions' "Out of scope" bullet).

## Notes

Kept intentionally generic — evidence comes from the first brownfield
adoption of v0.5.0/v0.6.x classification in an installed project (~50
pages, ~30 floor hits, ~2/3 of them public-pricing false positives; gate
verified blocking correctly on genuinely staged C2-shaped content). The
adopting project's own backlog rating and schema reconciliation are its
local work, not this task's.

Interacts with task 008 (global scanner) only in that both consume
`ac-classify` — the ack mechanism and floor refinement land in the shared
engine, so 008 inherits them for free.

Child tasks inherit their active parent's branch, worktree, and workflow
mode. Only a standalone or parent task owns Git lifecycle cleanup.
