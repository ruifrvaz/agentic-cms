# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Deterministic `.agentic-cms/bin/` toolkit (`ac-page`, `ac-index`, `ac-log`, `ac-links`, `ac-inventory`, `ac-search`) that skills call for mechanical wiki/frontmatter/index/log operations instead of hand-editing.
- `content-query` skill for cited, directed questions over the compiled wiki.
- `agentic-cms update`: self-updates the binary from the latest GitHub release, then re-runs `init` in the current directory to pick up new scaffold files.
- `agentic-cms version` now resolves the Go module version for `go install .../agentic-cms@latest` installs, in addition to the ldflags-injected version used by `make build`.
- **Draft content state** (pending v0.2.0 · PR #2) — adds an optional `status: draft | final` frontmatter field and a `docs/<topic>/drafts/` convention for work-in-progress content, kept out of `wiki/index.md`/`wiki/log.md` until promoted; `ac-page` gained `--status` and a `promote` subcommand for the toolkit-driven create/promote flow.

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

[Unreleased]: https://github.com/ruifrvaz/agentic-cms/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/ruifrvaz/agentic-cms/releases/tag/v0.1.0
