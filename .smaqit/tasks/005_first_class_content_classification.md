---
status: In Progress
created: "2026-08-16"
mode: Assisted
started: "2026-08-18"
---

# First-class content classification (auto CIA rating)

## Description

Add confidentiality classification as a first-class schema concept: a
`classification: C0 | C1 | C2 | C3` frontmatter field on every content page,
**auto-assigned by the authoring agent at write time** against a rubric
defined in `CONTENT.md`, validated and reported by the toolkit, and enforced
by `content-lint` (a bleed rule for index/log summaries, a ratchet rule for
re-rating).

Classification must survive **updates to existing content**, not just
creation: a page rated C1 can become C2+ through a single appended note or
edit. All detection is one engine — `ac-classify` — that every integration
point calls rather than re-implementing: a deterministic checker detects
stale ratings via a body-hash stamped at rating time and enforces heuristic
*floors* (auto-raise only). It fires at exactly two enforcement moments,
add/edit and pre-commit, via a post-tool-use hook (Claude Code, Codex — both
support blocking) and a git pre-commit gate; a skill-tail check guarantees
coverage where no blocking hook is installed or available, and `content-lint`
gains classification as one more category in its existing audit sweep rather
than a new mechanism.

Surfaced from real usage of an installed project's CMS, where uniform
content handling nearly pushed confidential material (financial figures,
PII, private correspondence) to a shared git remote — the CMS had no
concept of confidentiality at all, so `wiki/index.md` one-liners,
`wiki/log.md` narratives, and cross-referencing entity pages all leaked
sensitive detail from the pages they described. Classification fixes the
root cause: every page knows its own sensitivity, and the bookkeeping
layers are constrained by it.

Scale (standard CIA confidentiality axis): **C0 Public** (no harm if
published), **C1 Internal** (members-only; low harm), **C2 Confidential**
(need-to-know; financials, personal data, private strategy), **C3
Restricted** (severe harm; verbatim private correspondence, legal
instruments, anything credential-adjacent).

## Issue Triage Context

**Mode:** Auto
**Technologies:** Claude Code, OpenAI Codex CLI, Git
**Platforms/Environments:** Linux
**Features/Integrations:** agent lifecycle hooks (PostToolUse), git pre-commit hooks
**Versions/Constraints:** Hook schemas verified against platform docs 2026-08-18; Copilot CLI hook deferred to a follow-up (no documented blocking capability — see Design Decisions); task 004 dependency satisfied (merged as v0.3.0)

## Design Decisions

- **"Automatic" = agent-assigned at write time against a written rubric,
  not a toolkit regex classifier.** The toolkit validates the enum and
  reports distributions. Heuristic *detection* (credential-shaped strings,
  email addresses, currency amounts, person-names-with-figures) lives in a
  new deterministic `ac-classify` script — advisory for *ratings*, but
  authoritative as a **floor**: a pattern hit may auto-raise a page to that
  pattern's minimum level (credential-adjacent → C3; PII/financial → C2),
  never lower it. Floor-raising is safe by the ratchet's own asymmetry —
  over-protection is the recoverable direction, and "only the user may
  lower" already covers false positives. Full C0–C3 rating judgment stays
  with the agent; the tool never carries false authority on the way down.
- **Single engine, thin callers (SRP guardrail; refined 2026-08-18 after
  review).** All detection logic — enum validity, `classified-hash`
  staleness, heuristic floor patterns — lives exactly once, in
  `ac-classify`. Every integration point below is a thin wrapper that
  calls `ac-classify check`/`sweep`/`hook` and varies only in *response
  handling* (block, warn, auto-fix, report); none may re-implement
  detection. This is an explicit acceptance criterion (verified by
  inspection at completion), not just a stated intent — the earlier draft
  of this task described four "enforcement layers" in a way that read as
  four parallel mechanisms; they are four callers of one mechanism.
- **Two enforcement moments, one guaranteed fallback, one reused audit**
  (refined 2026-08-18): write-time rating alone leaves edits to existing
  pages unguarded — skill-mediated updates are instruction-only, and
  direct edits bypass skills entirely. The two moments that matter are
  add/edit and pre-commit:
  - **Add/edit** — a post-tool-use agent hook (Claude Code, Codex — both
    support blocking per verified docs) runs `ac-classify check` on the
    file just written, auto-applies floor bumps, and returns stale
    ratings as blocking feedback so the agent re-rates in the same turn.
  - **Pre-commit** — a git pre-commit gate checks staged content at the
    exfiltration boundary, the layer that would have caught the
    motivating incident regardless of who made the edit or which
    platform (or no platform) they used.
  - **Skill verify tail** is not a third moment competing with the two
    above — it is the guaranteed fallback for exactly the cases a
    blocking hook can't cover: a platform whose hook can't block (see the
    Copilot note below), or an installed project that hasn't picked up
    the hook config yet. Every write-path skill's verify tail gains
    `ac-classify check` alongside its existing `ac-index check && ac-links
    check`.
  - **`content-lint`** is not a new feature either — classification
    becomes one more category in the audit sweep it already runs
    (orphans, stale claims, broken links), via `ac-classify sweep`. Its
    job is catching what the two blocking gates never saw at all:
    pre-existing drift, edits never committed, or commits made with
    `--no-verify`.
- **Agent hooks: Claude Code and Codex only in this task; Copilot
  deferred** (scoped 2026-08-18 after review). All three CLIs fire a
  post-tool-use event with a JSON payload on stdin — verified against
  platform docs — so `ac-classify hook` is written to auto-detect payload
  shape (`tool_name`/`tool_input` on Claude Code and Codex vs.
  `toolName`/`toolArgs` on Copilot) and stays extensible to a third
  config later. But Copilot's docs describe hook *execution* only, no
  documented block/feedback channel — shipping its config now would add a
  full installer/brownfield/smoke-test surface for a platform where the
  hook can enforce nothing beyond what the skill-tail fallback already
  guarantees. Filed as follow-up work (see Notes) rather than built now:
  revisit once Copilot documents blocking, or once real usage on that
  platform shows the fallback insufficient.
- **Git pre-commit gate:** a versioned script at
  `.agentic-cms/hooks/pre-commit` checks the **staged blobs** (not the
  worktree) of `docs/`/`wiki/` markdown in the commit: BLOCK on floor
  violations and on C2+ bleed into staged `wiki/index.md`/`wiki/log.md`
  changes; WARN (non-blocking) on stale ratings and unrated pages.
  `git commit --no-verify` remains the standard escape hatch — document
  it rather than fight it. Prior art: the smaqit ecosystem's own
  pre-commit-validate hook follows this install-script pattern.
- **Staleness anchor is a content hash, not date discipline:** rating a
  page (`ac-page new --classification` / `ac-page classify`) stamps a
  short normalized-body hash (`classified-hash:`) in frontmatter;
  `ac-classify check` recomputes it — mismatch = the page changed since it
  was last rated. Deterministic, and independent of `updated:` hygiene.
- **Hook installation respects the non-destructive installer, per
  surface:** Claude Code (`.claude/settings.json`) and Codex
  (`.codex/hooks.json`) are single JSON files with no managed-block
  mechanism — greenfield installs ship them; brownfield installs with an
  existing file are never overwritten, init reports the snippet for
  manual merge. The git pre-commit hook installs into `.git/hooks/` only
  when no hook exists there; an existing hook gets a managed shell block
  appended (shell supports the CLAUDE.md-style append that JSON cannot),
  and a repo using `core.hooksPath` gets a report instead of a write.
- **New log operation `classify`** for re-ratings and floor bumps, added
  to CONTENT.md's operations list (same pattern as task 004's
  `archive`/`unarchive`).
- **Absent field = C1 (Internal)** — back-compat with every installed
  project and existing page, no migration, same reasoning as
  `status: final` in task 002. All *new* writes set the field explicitly.
- **Ratchet asymmetry:** an agent may raise a page's classification on
  update (a note can make a page more sensitive); only the user may lower
  one. Misclassifying downward is the dangerous direction.
- **Bleed rule:** `wiki/index.md` and `wiki/log.md` entries for C2+ pages
  may use only opaque summaries — no figures, no personal details, no
  quoted content. The bookkeeping layer inherits the classification of
  what it describes unless written to avoid doing so.
- **Live template field** (`classification: {{CLASSIFICATION}}`, default
  C1), not a commented example — per task 002's finding: `init`/`update`
  never overwrite existing template files, so there is no migration risk
  to avoid, and live fields keep the toolkit (not hand-editing)
  authoritative.
- **Multi-instance segregation is explicitly out of scope** — where C2+
  content physically lives (vault repo, two-tier split, history
  remediation) is deployment policy, currently being proven out in a
  downstream installed project. This task makes content *classifiable and
  bleed-safe within one CMS instance*; a follow-up task generalizes
  segregation once the downstream pattern settles.
- **Generic scoping:** no installed-project specifics in this tracker,
  same as the drafts task (002).

## Implementation Steps

1. **Schema** — `scaffold/tree/CONTENT.md`: new "Classification" section
   defining the C0–C3 scale, the auto-rating rubric (what content
   characteristics map to each level), the bleed rule for index/log, and
   the ratchet rule. Add the field to the frontmatter example block.
2. **Templates** — add the live `classification:` field to all five
   templates in `scaffold/tree/.agentic-cms/templates/` (doc, entity,
   concept, source, topic).
3. **Toolkit (mutators)** — `scaffold/tree/.agentic-cms/bin/ac-page`:
   validated `--classification C0..C3` on `new` (default C1) and a
   `classify <path> <level>` subcommand for re-rating; both stamp
   `classified-hash:` (short normalized-body hash) and `classify` touches
   `updated:`. The ratchet rule is enforced by skill instructions and the
   hook's raise-only floor logic, not by the mutator — the tool stays
   mechanical. `scaffold/tree/.agentic-cms/bin/ac-inventory`: report
   per-page classification and a distribution summary.
4. **Toolkit (checker)** — new `scaffold/tree/.agentic-cms/bin/ac-classify`:
   `check <path...>` returns per-page JSON — enum validity, staleness
   (`classified-hash` mismatch), and heuristic floor violations with the
   implied minimum level; `sweep` runs `check` across all of `docs/` and
   `wiki/`. Floor patterns: credential-shaped strings → C3; email/PII and
   currency-with-names → C2. Same JSON-out contract as the rest of
   `bin/`; document in `bin/README.md`.
5. **Agent hooks (Claude Code + Codex)** — `ac-classify hook`: reads the
   platform payload from stdin, auto-detects shape
   (`tool_name`/`tool_input` on both platforms today; the parser stays
   structured so a third shape can be added later without rework),
   extracts edited `docs/`/`wiki/` markdown paths, runs `check`,
   auto-applies floor violations via `ac-page classify` (raise-only) +
   `ac-log append classify`, and emits a blocking response (exit 2 +
   stderr, or `decision: block` JSON) carrying "stale rating — re-rate
   against the rubric". Ship two configs:
   `scaffold/tree/.claude/settings.json` (`PostToolUse`, Edit|Write
   matcher) and `scaffold/tree/.codex/hooks.json` (`PostToolUse`).
   Installer behavior: greenfield ships both; brownfield with an existing
   `.claude/settings.json` or `.codex/hooks.json` is skipped with the
   snippet reported. Extend the smoke test's brownfield scenario to cover
   the skip-and-report path. **Copilot's hook config is explicitly out of
   scope for this task** — see Design Decisions and Notes.
6. **Git pre-commit gate** — versioned script at
   `scaffold/tree/.agentic-cms/hooks/pre-commit`: collects staged
   `docs/`/`wiki/` markdown (`git diff --cached --name-only
   --diff-filter=ACM`), runs `ac-classify check` against the **staged
   blobs** (`git show :<path>`), BLOCKs on floor violations and C2+
   index/log bleed, WARNs on stale/unrated. `agentic-cms init` wires it:
   no existing `.git/hooks/pre-commit` → install a caller; existing hook →
   append a managed shell block (mirroring the CLAUDE.md merge);
   `core.hooksPath` set → report, don't write. Document `--no-verify` as
   the escape hatch.
7. **Write-path skills** — `content-manage-item`, `content-import`,
   `content-add-notes`, `content-research`: an explicit "rate before you
   write" step against the CONTENT.md rubric, passing the rating through
   `ac-page`; `content-add-notes` re-rates on every note (raise-only
   without user approval); every write-path skill's verify tail becomes
   `ac-index check && ac-links check && ac-classify check <touched
   paths>`.
8. **Lint & list** — `content-lint`: classification pass via `ac-classify
   sweep` — invalid values, unrated pages (advisory, back-compat), stale
   ratings, floor violations (report if the hooks were bypassed, e.g.
   out-of-session edits committed with `--no-verify`), and C2+ index/log
   leakage. `content-list`: surface classification distribution from
   `ac-inventory`.
9. **Root docs** — one-line mention in `scaffold/tree/CLAUDE.md` alongside
   the existing hard rules; README skills-table touch if warranted.
10. **Verification** — `make smoke-test` and `go test ./...` unaffected
    (plus the new brownfield skip-and-report checks); manual toolkit pass:
    create pages at each level, confirm enum validation refuses C5,
    confirm default C1, confirm `ac-inventory` reporting; edit a rated
    page and confirm `ac-classify check` reports it stale; plant a
    credential-shaped string in a C1 page and confirm the floor reports
    C3 and the hook path auto-raises; pipe recorded sample payloads from
    both platforms through `ac-classify hook` and confirm path extraction
    and blocking response shape; stage a floor-violating page in a
    sandbox repo and confirm the pre-commit gate blocks, then confirm
    `--no-verify` bypasses and lint later reports it; confirm a rating is
    never lowered by any mechanical path; **code-review every integration
    point (skill tail, both hooks, pre-commit, lint) and confirm none
    contains detection logic of its own — all four call `ac-classify`**.

## Known Issues Triage
**Triaged:** 2026-08-18
**Tools searched:** Claude Code (anthropics/claude-code), Git (git/git)
**Result:** Advisory

### Advisory Issues
- [#86405 PreToolUse/PostToolUse hooks not fired for subagent tool calls](https://github.com/anthropics/claude-code/issues/86405) — `anthropics/claude-code` — opened 2026-08-13 — bug, area:hooks, area:agents, needs-info, platform:macos. Reporter's own description is ambiguous about scope (background/async Agent-tool dispatch vs. inline subagent turns) and the issue is unconfirmed (needs-info); platform label is macOS, not Linux, so it does not confirm this task's platform dimension. Relevant because `content-import`'s brownfield sweep delegates to the `content-importer` subagent — if subagent-mediated edits genuinely bypass PostToolUse, the hook layer would miss them there, though the skill-tail check and pre-commit gate still cover it. Worth a manual check during implementation (Step 5/10): have the importer subagent edit a file and confirm the hook actually fires.
- [#87356 Hooks are tool-call-scoped, not filesystem-scoped](https://github.com/anthropics/claude-code/issues/87356) — `anthropics/claude-code` — opened 2026-08-17 — enhancement, area:hooks. Confirms the current model: a hook matcher fires per tool name (Edit, Write, ...), not per file-path pattern natively. The task's design already accounts for this — the Edit|Write matcher plus in-script path filtering to `docs/**`/`wiki/**` inside `ac-classify hook` — so this is confirmation the chosen approach is necessary, not a blocker.

### Historical (Closed)
- [#87542 PostToolUse hook exit-2/stderr warning not delivered to model in VSCode extension](https://github.com/anthropics/claude-code/issues/87542) — `anthropics/claude-code` — closed 2026-08-18 (same day, fast turnaround). Scoped to the VS Code extension surface specifically; not applicable to a CLI/terminal session, but worth knowing the exit-2/stderr delivery path had a recent bug in one surface.

### Unresolvable Tools
- OpenAI Codex CLI — `github-issues.sh resolve` returned `router-for-me/CLIProxyAPI`, an unrelated repo (not OpenAI's Codex CLI); treated as a mis-resolution and not searched. The research map's verified docs URL (https://learn.chatgpt.com/docs/hooks) remains the authoritative source for its hook schema; no GitHub issue tracker was checked for it.

### Omitted Tools
- (none — three tools identified, one unresolvable as above, two searched)

### Search Warnings
- (none)

## Acceptance Criteria

- [ ] `CONTENT.md` defines the C0–C3 scale, the auto-rating rubric, the index/log bleed rule for C2+, and the ratchet rule (agent may raise, only user may lower)
- [ ] All five templates carry a live `classification:` field; absent field means C1 (back-compat, no migration required for installed projects)
- [ ] `ac-page` validates `--classification` on `new` and supports `classify` for re-rating; both stamp `classified-hash`; `ac-inventory` reports per-page classification and distribution
- [ ] `ac-classify check`/`sweep` detects invalid values, stale ratings (hash mismatch), and heuristic floor violations; floors only ever raise, never lower
- [ ] Post-tool-use hooks ship for Claude Code and Codex (`.claude/settings.json`, `.codex/hooks.json`) through one `ac-classify hook` adapter that auto-detects the payload shape, auto-applies floor bumps, and returns stale ratings as blocking feedback; brownfield installs never overwrite an existing single-file config (snippet reported for manual merge); Copilot's config is explicitly not built in this task
- [ ] Git pre-commit gate checks staged `docs/`/`wiki/` blobs — blocks floor violations and C2+ index/log bleed, warns on stale/unrated; `init` installs it non-destructively (fresh install, managed append to an existing hook, report-only under `core.hooksPath`)
- [ ] Write-path skills (content-manage-item, content-import, content-add-notes, content-research) rate before writing and pass the rating through `ac-page`; content-add-notes re-rates on update; verify tails include `ac-classify check`
- [ ] `content-lint` runs `ac-classify sweep` — invalid values, unrated pages (advisory), stale ratings, floor violations, and C2+ index/log leakage
- [ ] `content-list` surfaces classification distribution
- [ ] **No integration point re-implements detection: code review confirms the skill-tail check, both agent hooks, the pre-commit gate, and the lint sweep each call `ac-classify` rather than containing their own floor-pattern or staleness logic**
- [ ] `make smoke-test` and `go test ./...` remain passing

## Findings

[Populated by smaqit.task-complete. Do not fill in manually before task is complete.]

**Implementation approach:**
- TBD

**Decisions made:**
- TBD

**Blockers encountered:**
- TBD

**Follow-up identified:**
- TBD

## Notes

Kept intentionally generic — this tracker describes the CMS schema/tooling
itself, not any particular installed project's content.

Scope watch: the 2026-08-18 refinement added the `ac-classify` script, two
hook configs, and git-hook wiring in `init` (Go changes). If implementation
balloons, the natural split is schema + toolkit + skills + agent hooks here,
with the git pre-commit gate and its installer wiring as a follow-up task —
the layers are independent by design. Decide at task-start, not now.

Follow-up identified during review (2026-08-18): a **Copilot CLI hook
config** (`.github/hooks/`, `postToolUse`) is explicitly deferred out of
this task — see Design Decisions. `ac-classify hook`'s payload-shape
detection is written to accept a third shape without rework, so the
follow-up task is config-and-installer-only, not an engine change. File it
once Copilot documents a block/feedback channel, or once real multi-agent
usage shows the skill-tail fallback insufficient there.

The downstream project that surfaced this need is applying the same scheme
ad hoc first (its own CONTENT.md section + a local segregation task), the
same pattern the drafts convention (task 002) followed: prove locally,
generalize upstream, reconcile the local section on the next
`agentic-cms update`. The segregation/vault half of that downstream work is
deliberately NOT part of this task — see Design Decisions.

Child tasks inherit their active parent's branch, worktree, and workflow
mode. Only a standalone or parent task owns Git lifecycle cleanup.
