# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **Delta-scoped git pre-commit classification gate** (pending v0.7.0 · PR #9) — blocks only on the staged commit's files (plus bleed into a staged `wiki/index.md`/`wiki/log.md`, regardless of source page); pre-existing tree-wide drift is summarized as one non-blocking warning instead of blocking every commit until the whole legacy backlog is rated.
- **Classification floor acknowledgment** (pending v0.7.0 · PR #9) — `ac-page classify <path> <level> --ack-floor` records a user-reviewed floor false positive, bound to the page's content hash; any edit invalidates the ack. Acking is a user-only decision, documented as such in every write-path skill.
- **CONTENT.md schema reconciliation** (pending v0.7.0 · PR #9) — `init`/`update` over an install with a customized `CONTENT.md` missing upstream schema sections now emit a report and write `.agentic-cms/CONTENT.upstream.md` for manual merge; the user's file is never edited.

### Changed

- **Currency classification floor widened** (pending v0.7.0 · PR #9) — now also catches ISO-coded, symbol-suffixed, and word-form currency amounts, keeping the floor's zero-false-negative recall contract; false positives are resolved via the new ack mechanism instead of narrower detection.
- **`.agentic-cms/VERSION` auto-stamped** (pending v0.7.0 · PR #9) — now set to the installing binary's version on every `init`/`update` run instead of a static literal that had drifted to `0.1.0`.

## [0.6.2] - 2026-08-19

### Fixed

- **Scaffolding logic files now refresh on init/update** — `scaffold.Install()` previously skipped every existing file uniformly, so `.claude/skills/`, `.claude/agents/`, `.agentic-cms/templates/`, `.agentic-cms/scripts/`, `.agentic-cms/hooks/`, and `.codex/` never actually refreshed on re-init despite `update.go`'s documented intent. Those framework-owned paths are now always overwritten with the embedded version, reported under a new `updated` result bucket; user-content paths keep their skip-if-exists behavior.

## [0.6.1] - 2026-08-18

### Changed

- **Classification surfaced in README** — promotes the C0-C3 confidentiality classification feature to a headline Features bullet naming the CIA triad, plus a new dedicated Classification section explaining the levels and why enforcement matters for agent-driven content synthesis.

## [0.6.0] - 2026-08-18

### Added

- **Shell installer bootstrap** — adds `install.sh`, a `curl | bash` bootstrap that downloads the latest release binary with no Go toolchain required (`AGENTIC_CMS_VERSION` pins a version); `agentic-cms update` now resolves the project root from any subdirectory instead of only the cwd; README restructured to mirror `ruifrvaz/smaqit`'s section style.

## [0.5.0] - 2026-08-18

### Added

- **First-class content classification** — adds a `classification: C0 | C1 | C2 | C3` frontmatter field, rated by the agent at write time against a rubric in `CONTENT.md`. All detection (enum validity, `classified-hash` staleness, heuristic floor patterns) lives in one new `ac-classify` engine; every enforcement point — Claude Code/Codex agent hooks, a git pre-commit gate, write-path skill verify tails, and `content-lint`'s sweep — is a thin caller of it. `ac-page` gained `--classification`/`classify`; `ac-inventory` gained a distribution tally.

## [0.4.0] - 2026-08-18

### Changed

- **`.agentic-cms/bin/` renamed to `.agentic-cms/scripts/`** — the toolkit directory held committed source (bash+python3 scripts meant to be read and audited), not compiled build output, colliding with the generic IDE convention that `bin/` is safe to hide as build artifacts. Fresh installs only; already-installed projects keep a stale `.agentic-cms/bin/` until reconciled, same gap as the `content-new-item` → `content-manage-item` rename.

## [0.3.0] - 2026-08-16

### Added

- **Archive content state** — adds `status: archived` and a `docs/<topic>/archive/` convention mirroring drafts, for retired content that stays retrievable; `ac-page` gained an `archive` subcommand (sibling of `promote`), archived items stay in `wiki/index.md` under a new `## Archived` section (`ac-index` auto-creates it on older indexes), `ac-inventory` reports `archived` per topic, and `content-list`/`content-lint` treat archived pages as retired rather than drifted.

### Changed

- **`content-new-item` renamed to `content-manage-item`** — the skill now owns the full item lifecycle (create, draft, promote, archive); upgraded installs keep the stale `content-new-item/` skill directory until `agentic-cms update` learns scaffold diffing/renames.

## [0.2.0] - 2026-08-15

### Added

- Deterministic `.agentic-cms/bin/` toolkit (`ac-page`, `ac-index`, `ac-log`, `ac-links`, `ac-inventory`, `ac-search`) that skills call for mechanical wiki/frontmatter/index/log operations instead of hand-editing.
- `content-query` skill for cited, directed questions over the compiled wiki.
- `agentic-cms update`: self-updates the binary from the latest GitHub release, then re-runs `init` in the current directory to pick up new scaffold files.
- `agentic-cms version` now resolves the Go module version for `go install .../agentic-cms@latest` installs, in addition to the ldflags-injected version used by `make build`.
- **Draft content state** — adds an optional `status: draft | final` frontmatter field and a `docs/<topic>/drafts/` convention for work-in-progress content, kept out of `wiki/index.md`/`wiki/log.md` until promoted; `ac-page` gained `--status` and a `promote` subcommand for the toolkit-driven create/promote flow.

### Fixed

- Installer's `{{DATE}}` substitution no longer corrupts `.agentic-cms/templates/*` (which must keep the placeholder live for `ac-page new` to fill in at page-creation time) or `.agentic-cms/bin/*` (whose source code references the literal placeholder string).

## [0.1.0] - 2026-08-01

### Added

- Initial agentic Markdown CMS scaffold and Go CLI installer. (4c6224a)
- Installer smoke-test coverage for validating generated project scaffolds. (e0d5307)
- Multi-root VS Code workspace configuration for task worktrees. (48f2b9e)
- Post-merge release automation that tags releases and creates GitHub Releases. (18d9647)

### Changed

- Build the CLI binary into the installer directory for distribution. (e6b7a26)

[Unreleased]: https://github.com/ruifrvaz/agentic-cms/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/ruifrvaz/agentic-cms/releases/tag/v0.2.0
[0.1.0]: https://github.com/ruifrvaz/agentic-cms/releases/tag/v0.1.0
