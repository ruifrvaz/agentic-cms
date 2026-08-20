---
status: Not Started
created: "2026-08-20"
---

# Classifier scope is hardcoded to docs/+wiki/, missing everything else

## Description

`ac-classify`'s file discovery is hardcoded to exactly two trees:
`docs/**/*.md` and `wiki/{entities,concepts,sources}/*.md` (plus
`wiki/index.md`/`wiki/log.md` for the bleed check). Every other path in a
project — `.claude/` (agent config and, in one observed case, an
accidentally-tracked memory-store directory), `.smaqit/history/` (session
history files, which quote and summarize raw content freely), `.smaqit/tasks/`
(task files, which routinely restate decision-doc figures for context),
`raw/` (immutable sources — PDFs, exports — which can themselves be C2/C3),
and any project-specific convention outside `docs/`/`wiki/` — is entirely
invisible to the engine. No floor, no bleed check, no sweep coverage, no
enforcement moment ever looks at these paths.

Surfaced during a real segregation task in an installed project: a
pre-squash manual grep (not the classifier) caught a session-history file
under `.smaqit/history/` carrying real capital figures and an investor's
identity — content the classifier's own sweep had reported zero problems
on, because it was never in scope to check. The same repo also had
`.claude/projects/` (a coding agent's own cross-session memory store)
accidentally git-tracked since the project's first week; non-sensitive in
that instance, but structurally the same blind spot — a path outside
`docs/`/`wiki/` that no classifier moment ever examines, discovered only
by a human-directed manual sweep during an unrelated cleanup, not by the
tool this task's users are meant to be able to trust for exactly that job.

## Issue Triage Context

**Mode:** Auto
**Technologies:** None
**Platforms/Environments:** None
**Features/Integrations:** ac-classify engine (all four callers inherit this fix)
**Versions/Constraints:** Baseline v0.7.0 (task 011); must not regress the existing docs/wiki behavior or its delta-scoped gate semantics

## Design Decisions

- **Scope should be configurable per project, not a second hardcoded
  pair of trees.** Projects vary in what non-docs/wiki paths carry real
  content (`.smaqit/history/`, `raw/`, a project-specific convention) —
  a fixed list just relocates the same blind-spot problem. Proposed:
  `CONTENT.md` (or a small `.agentic-cms/classify-scope` file) declares
  additional path globs the engine should include, defaulting to none
  (today's behavior preserved for every existing install) unless a
  project opts in.
- **`raw/` needs a distinct policy, not blanket inclusion.** It holds
  immutable sources, often binary (PDFs) or pre-classification-scheme
  legacy material — scanning it for markdown-oriented heuristics doesn't
  translate directly. At minimum, flag `raw/` as a known gap in
  documentation even if full detection isn't built this task; a
  same-filename-pattern or manifest-based classification (a sidecar
  noting each raw file's rating) is one plausible shape, not decided
  here.
- **`.claude/`, `.codex/`, and similar tool-config directories should
  never be treated as content to classify** — the actual fix for the
  memory-store incident isn't "classify it," it's "don't track it in
  the first place." Scope: add a gitignore/installer-level check
  (`agentic-cms init`/`update` warns, or a `content-lint`-style check
  flags, if paths matching known agent-tool memory/state directories are
  git-tracked) rather than pulling them into the classification engine.
  Out of scope for this task if it turns out to need its own design
  pass — file as a follow-up rather than block this task on it.
- **`.smaqit/tasks/` and `.smaqit/history/` are the clearest, highest-value
  addition** — they're markdown, they're project-authored (not tool
  internals), and the observed incident was exactly this path. Prioritize
  bringing these into scope (via the opt-in mechanism above) over the
  harder `raw/` and tool-config cases.

## Implementation Steps

1. Design the scope-declaration mechanism (CONTENT.md section vs. a
   dedicated file) and its default (empty — no behavior change for
   existing installs without an explicit opt-in).
2. Extend `ac-classify`'s page-discovery glob to read the declared scope
   in addition to the existing hardcoded `docs/`/`wiki/` trees; extend the
   bleed check's target list the same way if a project declares
   additional index/log-like files.
3. Document the opt-in in the shipped `CONTENT.md`, with `.smaqit/tasks/`
   and `.smaqit/history/` as the recommended first addition for smaqit-
   convention projects specifically (not a default — this repo and most
   agentic-cms installs don't use `.smaqit/`).
4. Add a `content-lint`/installer-level check (separate from ac-classify)
   that flags known agent-tool state/memory directories
   (`.claude/projects/`, similar Codex/Copilot equivalents if
   identifiable) if git-tracked, with guidance to gitignore rather than
   classify them.
5. Test fixtures: a project with an opted-in `.smaqit/tasks/` scope
   catches a planted figure there; a project without the opt-in behaves
   identically to today (regression guard); the tool-state check flags a
   planted `.claude/projects/` directory.
6. `make smoke-test` and `go test ./...` pass; manually verify against a
   project using the opt-in.

## Known Issues Triage

[Populated by smaqit.task-start via smaqit.utils.triage-issues. Do not edit manually.]

## Acceptance Criteria

- [ ] A project can declare additional path globs for `ac-classify` to include, defaulting to no change from today's docs/+wiki/-only behavior
- [ ] `.smaqit/tasks/` and `.smaqit/history/` are documented as the recommended opt-in scope for smaqit-convention projects
- [ ] A separate check flags git-tracked agent-tool state/memory directories (starting with `.claude/projects/`) and recommends gitignoring them, without pulling them into classification
- [ ] `raw/`'s gap is at minimum documented explicitly (full detection may be a follow-up, not required here)
- [ ] Regression: a project with no scope declaration behaves identically to pre-task behavior
- [ ] `make smoke-test` and `go test ./...` pass

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

Kept intentionally generic — the observed incident was in an installed
project using the `smaqit` task-planning convention (`.smaqit/tasks/`,
`.smaqit/history/`), but the underlying gap (classifier scope hardcoded
to exactly two trees) applies to any project with sensitive content
outside `docs/`/`wiki/`, smaqit or not.
