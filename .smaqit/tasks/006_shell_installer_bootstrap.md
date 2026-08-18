---
status: PR Open
created: "2026-08-18"
mode: Assisted
started: "2026-08-18"
pr: 6
---

# Shell installer bootstrap (no Go dependency)

## Description

Add a `curl | bash`-style `install.sh` at the repo root that downloads the
prebuilt platform binary from GitHub Releases, mirroring
`smaqit-extensions/install.sh` selectively — no global-install or
multi-agent scaffolding, no `--scope`/`uninstall` surface. Also adopt
ancestor-directory resolution into `agentic-cms update` so it works
correctly from a project subdirectory.

This is gradual, selective adoption of `smaqit-extensions`' distribution
mechanism — not a port of its CLI. `agentic-cms init`/`update`/`version`
already exist and are not being rebuilt; the gap being closed is
specifically that today's only documented install paths (`go install
...@latest`, `make build && sudo make install`) require a Go toolchain,
unlike `smaqit-extensions`' `curl | bash` bootstrap.

## Issue Triage Context

**Mode:** Auto
**Technologies:** Bash, Go
**Platforms/Environments:** Linux (macOS/Windows detection included for forward-compat, not yet released)
**Features/Integrations:** GitHub Releases API
**Versions/Constraints:** None

## Design Decisions

- **`install.sh` lives at repo root** (not `scripts/`), matching
  `smaqit-extensions` precedent and giving a short raw-GitHub curl URL.
- **No `--install-global` step** — agentic-cms's project-scoped `init`
  already covers that need; there are no global agent/skill directories to
  seed.
- **No `uninstall` command, no `--scope` flag, no multi-agent
  (claude/codex/copilot) global directory split** — explicitly out of
  scope, per "gradual, not full adoption."
- **Ancestor-dir resolution is simplified from smaqit-extensions' version**
  (drops git-root-takes-precedence) since agentic-cms has only one marker
  directory (`.agentic-cms/`) and no existing git-root behavior to
  preserve.
- **No new CI test for `install.sh` itself** — matches `smaqit-extensions`
  precedent (its `install.sh` is also untested) and avoids mocking GitHub
  Releases in CI; covered by manual verification instead.
- **macOS/Windows platform detection included in `install.sh` for
  forward-compat** (matches `update.go`'s existing `runtime.GOOS`
  branching) even though only Linux binaries are currently published —
  fails cleanly with "no asset found" until `build-all`/the release
  workflow grows those targets.

## Implementation Steps

1. Write `install.sh` at repo root — bash script with `set -e`, functions
   `detect_platform` (linux/darwin/windows × amd64/arm64, matching
   `update.go`'s `runtime.GOOS`/`GOARCH` combinations), `get_latest_version`
   (GitHub API `releases/latest`, overridable via an `AGENTIC_CMS_VERSION`
   env var — `latest` default, `vX.Y.Z` pin), `download_binary` (asset name
   `agentic-cms_<os>_<arch>`, matching `update.go`'s existing convention and
   `.github/workflows/post-merge-release.yml`'s build step — no
   server-side change needed), `install_binary` (→
   `~/.local/bin/agentic-cms`, `chmod +x`), `verify_installation`
   (`--version`), `check_path` (PATH warning). Model directly on
   `smaqit-extensions/install.sh`.
2. Add `resolveDefaultProjectDir`/ancestor-walk helper to `update.go`
   (nearest ancestor containing `.agentic-cms/`, else cwd — no git-root
   special case), wire into `runUpdate` in place of the hardcoded `"."`,
   threading the resolved dir through `checkAndReInit`/`reinitWithBinary`.
3. Add test coverage for the new resolution function in `update_test.go`.
4. Update `README.md` Install section: `curl | bash` one-liner promoted to
   primary path; `go install`/`make build && make install` demoted to
   secondary/contributor paths.
5. Manual verification: run `install.sh` end-to-end against the real
   latest GitHub release into a scratch `HOME`; run `agentic-cms update`
   from a nested subdirectory of an installed project and confirm it
   re-inits the project root (not the subdirectory).
6. Run `make smoke-test` and `go test ./...`.

## Known Issues Triage
**Triaged:** 2026-08-18
**Tools searched:** Bash, Go, GitHub REST API
**Result:** Clear

### Blocking Issues
(none)

### Advisory Issues
(none)

### Historical (Closed)
(none)

### Unresolvable Tools
(none)

### Omitted Tools
(none)

### Search Warnings
- Repository resolution for `Bash` and `Go` is generic, not task-specific: Bash has no canonical GitHub issue tracker (upstream is GNU/Savannah), so the helper fell back to `dylanaraps/pure-bash-bible`; `Go` resolved to the language's own `golang/go` tracker. Both open/closed searches on `golang/go` returned only broad language/roadmap items (e.g. platform end-of-support policy, `x/tools` build issues) — none carry a `bug`/`regression` label confirming a match on both the `Linux` platform and `GitHub Releases API` feature dimensions relevant to this task. `dylanaraps/pure-bash-bible` and `github/rest-api-description` returned no results at all. Treated as no relevant findings (Clear), not evidence of the absence of an issue.

## Acceptance Criteria

- [x] `install.sh` downloads the correct `agentic-cms_<os>_<arch>` asset
      from the latest GitHub release, installs it to `~/.local/bin`, and
      verifies via `--version`
- [x] `AGENTIC_CMS_VERSION` env var can pin a specific release (mirroring
      `SMAQIT_EXT_VERSION`)
- [x] README's Install section leads with the `curl | bash` one-liner;
      `go install`/`make install` remain as secondary paths
- [x] `agentic-cms update` resolves the project root by walking up to the
      nearest ancestor containing `.agentic-cms/` when run from a
      subdirectory, instead of only checking the cwd
- [x] No global install/uninstall, `--scope` flag, or multi-agent
      (claude/codex/copilot) global directory logic is introduced
- [x] `make smoke-test` and `go test ./...` remain passing

## Findings

**Implementation approach:**
- Wrote `install.sh` at repo root, mirroring `smaqit-extensions/install.sh`'s
  functions (`detect_platform`, `get_version`, `download_binary`,
  `install_binary`, `verify_installation`, `check_path`) but dropping its
  `--install-global` step, since agentic-cms has no global agent/skill
  directories to seed.
- Added `resolveDefaultProjectDir`/`findAncestorWithEntry` to `update.go`,
  wired into `runUpdate` in place of the hardcoded `"."`; simplified from
  smaqit-extensions' version by dropping its git-root-precedence branch,
  since agentic-cms has only one marker directory (`.agentic-cms/`).
- Added two unit tests (`update_test.go`) covering the ancestor-found and
  no-ancestor-fallback cases.
- Restructured `README.md` to mirror `ruifrvaz/smaqit`'s README section
  skeleton and style (Features / Compatibility / Getting Started / Commands
  with a Reinstallation-and-Updates subsection / Documentation), per a
  follow-up user request mid-task; folded the prior standalone "Design
  notes" section into the new Features bullets to match smaqit's pattern,
  and fixed a link that would have pointed at a nonexistent root-level
  `CONTENT.md` (it only exists at `scaffold/tree/CONTENT.md` in this repo).

**Decisions made:**
- `install.sh` lives at repo root, not `scripts/`, for a short raw-GitHub
  curl URL — matches `smaqit-extensions` precedent.
- No `--install-global`, `--scope`, `uninstall`, or multi-agent
  (claude/codex/copilot) global directory scaffolding — explicitly out of
  scope per the approved plan's "gradual, not full adoption" framing.
- No new CI test for `install.sh` itself (matches `smaqit-extensions`
  precedent); covered by shellcheck plus manual end-to-end verification
  instead.
- The private-repo blocker discovered mid-task was resolved by making
  `ruifrvaz/agentic-cms` public (user-confirmed this was accidental, not
  deliberate) rather than by adding token-based auth to the installers —
  the repo's existing `go install .../agentic-cms@latest` path already
  implicitly required public visibility, so this wasn't a new constraint
  introduced by this task, just the first thing to actually exercise it.

**Blockers encountered:**
- `ruifrvaz/agentic-cms` was private, which broke unauthenticated
  `curl`/`go install` access entirely (GitHub returns 404 for both the
  Releases API and asset downloads without auth on a private repo).
  Diagnosed via an authenticated `gh api`/`gh release download` check
  against the real `v0.3.0` release to confirm `install.sh`'s URL
  construction and download/verify logic were correct despite the
  network failures, then the user made the repo public and both
  `install.sh` and `update`'s ancestor-dir resolution were re-verified
  fully unauthenticated end-to-end against the live public repo (v0.5.0).
- Forwarding a `gh`-issued Bearer token through `browser_download_url`'s
  redirect (during private-repo diagnosis) returned 404 rather than a
  clean auth error — GitHub's release-asset redirect doesn't accept a
  forwarded API token; `gh release download` (or the assets API with
  `Accept: application/octet-stream`) is the correct authenticated path,
  not the plain `browser_download_url`. Not relevant to the shipped code
  (which is intentionally unauthenticated), but worth remembering if
  token-based private-repo support is ever added later.

**Follow-up identified:**
- None — no new tasks required. The `SMAQIT_EXT_VERSION`-style env var,
  ancestor-dir resolution, and README restructure are all complete as
  scoped.

## Files to Create / Modify

| File | Action |
|------|--------|
| install.sh | Create |
| update.go | Modify |
| update_test.go | Modify |
| README.md | Modify |

## Notes

Modeled on `smaqit-extensions/install.sh` and
`smaqit-extensions/installer/main.go`'s `resolveDefaultProjectDir`/
`findAncestorWithEntry` (both on this machine at
`~/projects/smaqit-extensions`), adopted selectively per explicit user
instruction: this is gradual adoption based on need, not a hard copy of
the `smaqit-extensions` CLI. Do not add `--install-global`, `--scope`,
`uninstall`, or multi-agent (claude/codex/copilot) global directory
scaffolding — those remain out of scope.

**Repo visibility (discovered during implementation):** `ruifrvaz/agentic-cms`
was private, which silently broke every documented install path — not just
`install.sh`, but also the pre-existing `go install .../agentic-cms@latest`
(the public Go module proxy can't fetch a private module without the
installing machine having `GOPRIVATE` + its own git credentials configured).
User confirmed this was never deliberate and made the repo public
(2026-08-18). `install.sh` and `agentic-cms update`'s ancestor-dir
resolution were both re-verified fully unauthenticated end-to-end against
the live public repo (v0.5.0) after the flip — see Findings.

Child tasks inherit their active parent's branch, worktree, and workflow
mode. Only a standalone or parent task owns Git lifecycle cleanup.
