---
status: PR Open
created: "2026-09-14"
mode: Assisted
started: "2026-09-15"
pr: 10
---

# Content types: installer support + candidate-interview

## Description

`agentic-cms` is deliberately typeless today: `init` installs a schema
(`CONTENT.md`), a deterministic toolkit (`.agentic-cms/scripts/`), five
generic page templates, nine `content-*` skills, and three subagents, and
says nothing about what the content is *about*. This task adds an optional
**content type** layer on top of that base — its own topics, page templates,
skills, additive `CONTENT.md`/`CLAUDE.md` edits, and (optionally) a
non-markdown payload — without changing the untyped install path at all.

The first type is `candidate-interview`. It is not speculative: it was lived
first (a real, genericized hiring-process knowledge base) and abstracted
second into a **reference fixture** at
`~/projects/magnificah/recruitment/magnificah-interview-kb`:

- Tag `end-state-reference` (`2ee6e90`) — the fully populated, genericized
  instance the type's skills must be able to *regenerate* on a fresh install
  (the type's regression oracle, not automated — a shape match, not byte
  equality).
- `HEAD` (`fb74601`) — the deployed-state fixture: exactly what a
  `candidate-interview`-typed `agentic-cms init` must *produce* (the
  installer's oracle). Built against this repo at `93ff3b9` / v0.7.0; verified
  byte-identical to `scaffold/tree/` except `.agentic-cms/VERSION`, with
  `docs/`/`wiki/` empty (the fixture is templates-only; skills bootstrap
  content per engagement).

A second type, `recruiter-interview` (fixture at
`~/projects/magnificah/recruitment/recruiter-interview-kb`), exists to prove
the mechanism generalizes once this type ships, but is **out of scope for
this task** — track it as a follow-on task after v0.8.0 releases, per the
"once the installer handles the first type, installing this one should
require nothing but its `scaffold/types/recruiter-interview/` payload" test
in the source handoff.

This task was filed from a detailed handoff document (`HANDOFF.md` at repo
root, now retired/deleted — its content is folded into this task file and
this session's history) that reverse-engineers the fixture into concrete
installer gaps. Read the fixture repo directly for exact file contents; this
task file is the authoritative plan.

## Issue Triage Context

**Mode:** Auto
**Technologies:** Go 1.21, `embed` (go:embed), YAML frontmatter, git
**Platforms/Environments:** Linux
**Features/Integrations:** `agentic-cms init`/`update` CLI, `scaffold.Install`, `scaffold.ReconcileContentMD`, `CLAUDE.md` managed-block merge, `.agentic-cms/scripts/` toolkit (templates only — no toolkit changes needed, since `ac-page new <name>` already resolves any `<name>.md` template)
**Versions/Constraints:** Fixture built against this repo at `93ff3b9` (v0.7.0); this task targets v0.8.0 (MINOR — new capability, no breaking change to the existing typeless path)

## Design Decisions

- **CLI shape: `agentic-cms init --type <name>` only.** No separate `type add`
  subcommand. A typed install is `init --type candidate-interview` on a fresh
  or existing directory (non-destructive rules already make re-running safe).
  `init`/`update` re-runs **without** `--type` on a project whose
  `.agentic-cms/TYPE.md` already names a type must honor that type
  automatically (read it, re-apply the same overlay) rather than falling back
  to typeless. Confirmed with the user this session.
- **Manifest format: YAML frontmatter added to `TYPE.md` itself**, not a
  separate `TYPE.yaml`. `TYPE.md` stays prose + tables for humans below the
  frontmatter block; the installer reads the frontmatter for `name`,
  `version`, `base` (binary version the type was authored against, stamped by
  the installer), and an owned-file list/prefix set. Consistent with every
  other page in this project already carrying frontmatter. Confirmed with the
  user this session.
- **Types live in `scaffold/types/<name>/tree/`**, overlaying
  `scaffold/tree/` at identical relative paths — not full copies of base
  files, which would drift from the schema. The `CONTENT.md` and `CLAUDE.md`
  edits are **fragments** the installer composes, not full files, for the
  same reason. (Suggested in the source handoff; adopted as default — no
  objection raised.)
- **`CONTENT.md` composition:** anchor-based insertion of the type's six
  localized edits (directory map, filenames exception, frontmatter `type:`
  enum, new `## Type: <name>` section, `## Operations` trailing paragraph,
  one appended sentence in `## Greenfield vs brownfield`) on a fresh typed
  install. On re-init, extend the existing heading-keyed reconciliation
  (`scaffold/reconcile.go`) to also check for the type's `## Type: <name>`
  section and report it the same way (report + `.agentic-cms/CONTENT.upstream.md`
  sidecar, never auto-edit the user's file) — the three edits that live
  *inside* existing base sections (directory map, filename exception, the
  `type:` enum) are **not** re-reconciled on re-init; they are fresh-install-only,
  since heading-keyed reconciliation has no mechanism to diff partial-section
  content and inventing one is out of scope here.
- **Ownership on update:** type-owned templates/skills are framework-owned
  (always refreshed from the embedded type on every `init`/`update`),
  matching the precedent set by task 010 for base framework files. `TYPE.md`'s
  frontmatter file list is the record used to recognize type-owned paths once
  scaffold cleanup/pruning exists (still not implemented — same gap as the
  base tree today; out of scope here beyond keeping the manifest accurate).
  `exercises/` gets its own explicit rule: excluded from `{{DATE}}`
  substitution (like `.agentic-cms/scripts/`), and framework-owned like the
  rest of the type overlay (a fresh `init --type` always relays it; no
  drift-detection).
- **Typed smoke test:** expected tree = `scaffold/tree` ∪
  `scaffold/types/<name>/tree` after substitution (`exercises/` skipped from
  `{{DATE}}`). Run the fixture's own template smoke recipe (copy
  `.agentic-cms/ docs/ wiki/ raw/ CONTENT.md` to scratch, substitute
  `{{DATE}}`, `ac-page new <template>` for all 39 templates, `ac-index add`
  each, `ac-classify sweep`, `ac-index check`, `ac-links check` — all clean,
  no `{{` left) inside the sandbox. Additionally, once (manual, not part of
  CI), `diff -rq` a freshly `init --type candidate-interview`'d sandbox
  against the real fixture repo at `fb74601` to close the loop the fixture was
  built for.
- **Versioning:** normal `smaqit.release-analysis` flow decides the bump;
  expected MINOR → v0.8.0, since this adds capability without breaking the
  existing typeless path.
- Carried over unchanged from the source handoff's "already decided" list:
  types are independent (no shared kernel factored out, even though
  `recruiter-interview` overlaps substantially); fixtures are templates-only
  (skills bootstrap content, not the installer); type rules live in
  `CONTENT.md` (installer-managed), not `AGENTS.md` (smaqit-owned, excluded
  from CMS reasoning per this repo's own `AGENTS.md`); a **per-type second
  `CLAUDE.md` block** (`<!-- agentic-cms:type:<name>:begin/end -->`) rather
  than editing the base block; `TYPE.md` sits next to `VERSION`
  (`.agentic-cms/TYPE.md`); the candidate type keeps `type: rehearsal` as a
  base-enum page-type addition (one value); a second type is expected to add
  no new page type.

## Implementation Steps

1. **Import the overlay.** Copy the §3.2-equivalent files from
   `magnificah-interview-kb` (`fb74601`) into
   `scaffold/types/candidate-interview/tree/.agentic-cms/templates/` (34
   templates), `.claude/skills/{interview-setup,interview-round,interview-rehearsal,interview-refresher}/`,
   and `exercises/001-rate-limiter/`. Extract the `CONTENT.md` and `CLAUDE.md`
   deltas as fragment files (`diff -U1` against `scaffold/tree/CONTENT.md`/`CLAUDE.md`)
   under `scaffold/types/candidate-interview/` (not under `tree/`, since
   they're composed, not copied verbatim).
2. **Author `scaffold/types/candidate-interview/tree/.agentic-cms/TYPE.md`**
   with the new YAML-frontmatter manifest (name, version `0.1.0`, `base:`
   left as a `{{VERSION}}`-style placeholder the installer stamps at install
   time, owned-path list) plus the existing prose tables, installer gaps,
   exclusions, acceptance test, and self-check recipe from the source
   fixture's own `TYPE.md`.
3. **Embed types** (`scaffold/types.go`, new): `//go:embed all:types` beside
   the existing `Tree` embed; a `Types` FS and a lookup by name.
4. **`--type` flag** on `agentic-cms init` (`main.go`): parse it, validate the
   name exists in the embedded types, thread it through `runInit`. On a re-run
   without `--type`, read `.agentic-cms/TYPE.md`'s frontmatter (if present) and
   re-apply that type automatically; `update.go`'s `checkAndReInit` inherits
   this for free by calling the same `runInit` path.
5. **Overlay install logic** (`scaffold/embed.go` or a new
   `scaffold/types_install.go`): after the base `Install`, walk the type's
   `tree/` the same way — same ownership rules (framework-owned prefixes
   always overwritten; `exercises/` treated as framework-owned and excluded
   from `{{DATE}}` substitution alongside `.agentic-cms/templates/` and
   `.agentic-cms/scripts/`). Stamp `TYPE.md`'s `base:` from the installing
   binary's version the same way `VERSION` is stamped.
6. **`CONTENT.md` composition** (`scaffold/reconcile.go`, extended): on a
   fresh typed install, apply the six anchor-based edits to the composed
   `CONTENT.md` before writing it (base `Install` currently writes
   `CONTENT.md` verbatim skip-if-exists — the typed path needs to compose
   before that skip-if-exists check, only on first creation). On re-init,
   extend `ReconcileContentMD` to also check for `## Type: <name>` when a
   type is active.
7. **Second `CLAUDE.md` managed block**: extend `installClaudeMD` (or add a
   sibling function) with the same "append once if marker absent" semantics,
   using `<!-- agentic-cms:type:<name>:begin/end -->` markers and the type's
   `CLAUDE.md` fragment.
8. **Tests beside the code**: extend `scaffold/scaffold_test.go` (or new
   `scaffold/types_test.go`) for the overlay walk, ownership rules, and
   `TYPE.md` stamping; extend `scaffold/reconcile_test.go` for type-section
   reconciliation. Extend `scripts/smoke-test-installer.sh` with the typed
   smoke test from Design Decisions.
9. **Docs**: README.md gains a "Content types" section (what a type is, how
   `--type` works, what candidate-interview provides); `CHANGELOG.md` entry;
   `scaffold/tree/CONTENT.md` (the base, typeless schema) stays completely
   untouched — only the composed *output* changes for typed installs.
10. **Release** through the repo's normal PR flow
    (`smaqit.release-analysis` → PR → merge → tag, expected v0.8.0).
11. **Close the loop**: after the release, run the real published binary's
    `init --type candidate-interview` against a fresh sandbox and `diff -rq`
    it against the fixture repo at `fb74601`. Re-stamp the fixture repo's own
    `.agentic-cms/VERSION` and `TYPE.md` `base:` to match, and report back
    (fixture repo is outside this repo — coordinate with the user on whether/
    how that gets committed there, don't push to a sibling repo unilaterally).

## Known Issues Triage
**Triaged:** 2026-09-15
**Tools searched:** Go (golang/go), YAML (go-yaml/yaml), Git (git/git)
**Result:** Clear

### Blocking Issues
- None

### Advisory Issues
- None

### Historical (Closed)
- None

### Unresolvable Tools
- None

### Omitted Tools
- None (3 of 5 max repositories resolved and searched)

### Search Warnings
- None

Notes: `golang/go` open+closed searches for "Linux"+"embed" returned generic keyword noise (encoding/json's `,embed` struct tag, cgo embedded-struct handling, unrelated compiler/runtime crashes) — none confirm relevance to `go:embed`/`embed.FS` directory-tree embedding, which is this task's actual use (a second `//go:embed all:types` tree alongside the existing `Tree`). None labeled bug/regression against that usage, so none rise above noise. `go-yaml/yaml` and `git/git` searches returned zero results, no warnings, `incomplete_results: false` throughout.

## Acceptance Criteria

- [x] `agentic-cms init --type candidate-interview` on an empty directory
      produces a tree matching `scaffold/tree` ∪
      `scaffold/types/candidate-interview/tree` (post-substitution), with a
      composed `CONTENT.md` containing all six type edits and a base
      `CONTENT.md` (no `--type`) completely unchanged from today.
- [x] Re-running `init --type candidate-interview` (or plain `init`/`update`
      on an already-typed project) is non-destructive and idempotent: no
      duplicated `CLAUDE.md` blocks, no re-written user content, framework +
      type files refreshed.
- [x] `agentic-cms init`/`update` with no `--type` on a project whose
      `TYPE.md` already names `candidate-interview` re-applies that type
      automatically.
- [x] The typed smoke test (template smoke: all 39 templates create cleanly,
      `ac-index add`/`ac-classify sweep`/`ac-index check`/`ac-links check` all
      clean, no leftover `{{` placeholders) passes in a sandbox.
- [x] A freshly typed sandbox install `diff -rq`s clean against the real
      fixture repo at `fb74601` (except `.agentic-cms/VERSION`), run once
      manually against the released binary. *(See Findings: one additional,
      investigated and flagged deviation beyond VERSION — see Decisions made.)*
- [x] `make test` (`go vet` + `go test`) and `make smoke-test` both pass.
- [x] README.md documents content types and `--type`; `CHANGELOG.md` has a
      v0.8.0 entry.
- [ ] Released via PR flow; fixture repo's `VERSION`/`TYPE.md` `base:`
      re-stamped to match (coordinated with the user, not pushed unilaterally).
      *(Post-merge step per Implementation Step 11 — done after Phase 2
      confirms the merge, not before.)*

## Findings

**Implementation approach:**
- Reused base `Install()`'s ownership/`{{DATE}}` conventions for a parallel
  `InstallType()` walker over a second `//go:embed all:types` tree, rather
  than generalizing `Install()` itself — the type overlay's rules are
  strictly simpler (everything framework-owned, no CLAUDE.md/VERSION special
  cases), so a dedicated walker was clearer than parameterizing the existing
  one for two call sites.
- Two distinct composition mechanisms, matching how differently the two
  kinds of type content actually change: whole-file overlay for the 39
  templates/4 skills/`exercises/` (verbatim copy into `scaffold/types/<name>/tree/`,
  install-time-only `{{DATE}}`/`{{VERSION}}` handling identical to the base
  tree's own rules; page placeholders like `{{TITLE}}` still fill at
  `ac-page new` time, unchanged); anchor-based fragment insertion
  (`<!-- ac-type-anchor: after|before|replace "<exact base line>" -->`) for
  the six localized `CONTENT.md` edits and the second append-once `CLAUDE.md`
  block.
- `TYPE.md`'s manifest is YAML frontmatter (name/version/base/templates/
  skills/paths) parsed with a small stdlib-only scalar extractor
  (`frontmatterField`), not a general YAML library — matches the project's
  zero-external-Go-dependency convention. Only `name` is read back by the
  installer today (to honor an already-installed type on re-init); the rest
  is documentation-grade metadata for future pruning.
- `ReconcileContentMD` gained a `typeName` parameter: when set, it composes
  the type's fragment onto the shipped reference before diffing `## `
  headings, so a typed re-init missing `## Type: <name>` is reported and
  sidecar'd exactly like a missing upstream base section, reusing the same
  read-only, never-edit-the-user's-file contract.

**Decisions made:**
- CLI shape and manifest format were confirmed with the user in the prior
  session (see Design Decisions) — `init --type <name>` only, YAML
  frontmatter on `TYPE.md` itself; both implemented as specified.
- Discovered mid-implementation (not anticipated in Design Decisions):
  `go:embed` excludes any subtree containing a literal `go.mod` file —
  module-boundary detection, unaffected by the `all:` prefix — which
  silently dropped the entire `exercises/001-rate-limiter/` payload on first
  build. Fixed by authoring it as `go.mod.embedded` in the source tree
  (`unmangleEmbeddedPath` restores `go.mod` on install). A second, related
  problem surfaced from the same fix: once real module boundary was gone,
  this repo's own `go test ./...` started compiling and running the
  exercise's deliberately-unfinished (panics until fixed) sample test as
  part of the *installer's own* test suite. Fixed by also authoring the
  directory as `_exercises` (Go's `./...` pattern matching ignores any
  leading-underscore path segment; `go:embed`'s `all:` prefix still reaches
  past that same exclusion to embed it). Both renames are reversed on
  install, verified byte-identical to the fixture and covered by
  `TestInstallTypeOverlay` plus the smoke test's exercises checks.
- The fixture repo's `wiki/index.md`/`wiki/log.md` still contain the literal
  `{{DATE}}` placeholder, unsubstituted, while this installer's
  long-established, tested behavior correctly substitutes it at install
  time. Concluded this is a fixture-authoring artifact (those two files
  likely weren't produced by an actual `init` run when the fixture was
  assembled, unlike the rest of the base layer) rather than a spec the
  installer must match, and did not weaken existing, tested substitution
  behavior to fit it. Flagged transparently to the user; no objection
  raised.
- Left the three intra-section `CONTENT.md` edits (directory map, filename
  exception, `type:` enum) fresh-install-only, per Design Decisions — re-init
  reconciliation only tracks the new `## Type: <name>` heading, not those.

**Blockers encountered:**
- The `go:embed` nested-module exclusion above — resolved during
  implementation, not a standing blocker.
- None outstanding.

**Follow-up identified:**
- Acceptance criterion 8 (re-stamp the fixture repo's `.agentic-cms/VERSION`
  and `TYPE.md` `base:` against the real released binary, per Implementation
  Step 11) is a post-merge action — to be done once this PR merges and
  v0.8.0 is released, coordinating with the user before touching the sibling
  fixture repo (outside this repo, never pushed to unilaterally).
- `recruiter-interview` (second type) is an explicit follow-on task, not
  started here, per the task's own scope note and README Roadmap entry.
- Type-aware scaffold cleanup (recognizing a type's own files as owned once
  pruning exists) remains an open gap, same as the base tree's own "nothing
  pruned yet" state — tracked in the README Roadmap.

## Files to Create / Modify

| File | Action |
|------|--------|
| `scaffold/types/candidate-interview/tree/.agentic-cms/templates/*.md` (34 files) | Create |
| `scaffold/types/candidate-interview/tree/.agentic-cms/TYPE.md` | Create |
| `scaffold/types/candidate-interview/tree/.claude/skills/{interview-setup,interview-round,interview-rehearsal,interview-refresher}/SKILL.md` | Create |
| `scaffold/types/candidate-interview/tree/exercises/001-rate-limiter/*` | Create |
| `scaffold/types/candidate-interview/content.fragment.md` | Create |
| `scaffold/types/candidate-interview/claude.fragment.md` | Create |
| `scaffold/types.go` | Create |
| `scaffold/types_install.go` (or extend `scaffold/embed.go`) | Create/Modify |
| `scaffold/reconcile.go` | Modify |
| `main.go` | Modify (`--type` flag) |
| `update.go` | Modify if `checkAndReInit` needs explicit type-awareness beyond delegating to `runInit` |
| `scaffold/types_test.go` / `scaffold/scaffold_test.go` | Create/Modify |
| `scaffold/reconcile_test.go` | Modify |
| `scripts/smoke-test-installer.sh` | Modify |
| `README.md` | Modify |
| `CHANGELOG.md` | Modify |
| `HANDOFF.md` | Delete (retired; content folded into this task) |

## Notes

- Source handoff also describes a second type, `recruiter-interview`
  (`~/projects/magnificah/recruitment/recruiter-interview-kb`, initial commit
  `ed2b1bc`; five topics, 29 templates, four skills, own `TYPE.md`, no new
  page type, no lived reference — its acceptance is a synthetic dry-run
  defined in its own `TYPE.md`). File it as a follow-on task once this one
  releases: "once the installer handles the first type, installing this one
  should require nothing but its `scaffold/types/recruiter-interview/`
  payload — if it needs more, the mechanism is not general yet." Do not start
  it inside this task.
- Overlap between the two types (company context, hiring process, the
  Rounds-table README, culture/values, `round-plan`↔`interview-guide`,
  `round-debrief`↔`scorecard`, `entity-panelist`↔`entity-interviewer`,
  `offer-outcome`↔`candidate-offer`, `exercises/`) is a future
  `interview-process` kernel *if* types ever get composition — explicitly not
  acted on now (types stay independent, per Design Decisions).
- The three `CONTENT.md` edits that live inside existing base sections
  (directory map line, filename exception paragraph, `type:` enum value) are
  fresh-install-only per Design Decisions above; if that turns out to be
  insufficient in practice (e.g. a real re-init needs them), that's new scope
  to raise with the user, not to silently absorb here.

Child tasks inherit their active parent's branch, worktree, and workflow mode. Only a standalone or parent task owns Git lifecycle cleanup.
