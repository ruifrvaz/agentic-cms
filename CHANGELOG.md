# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **Shell installer bootstrap** (pending v0.6.0 · PR #6) — adds `install.sh`, a `curl | bash` bootstrap that downloads the latest release binary with no Go toolchain required (`AGENTIC_CMS_VERSION` pins a version); `agentic-cms update` now resolves the project root from any subdirectory instead of only the cwd; README restructured to mirror `ruifrvaz/smaqit`'s section style.

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
