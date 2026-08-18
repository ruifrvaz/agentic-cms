# Installer Bootstrap And Classification Docs

## Metadata

- **Date:** 2026-08-19 (session spanned 2026-08-18 → 2026-08-19)
- **Session focus:** Refine agentic-cms's installer distribution to match `smaqit-extensions`' pattern selectively (a `curl | bash` shell bootstrap, no Go toolchain required), then give the existing C0-C3 confidentiality classification feature headline visibility in the README.
- **Tasks completed:** 006 (Shell installer bootstrap, no Go dependency — PR #6, v0.6.0), 009 (Surface classification feature in README — PR #7, v0.6.1)
- **Tasks referenced:** 005/007/008/010 (parallel/prior work by other sessions, not touched directly here except where their merged content required reconciliation)

## Actions taken

1. **Planned and implemented task 006** via `task.plan` → `task.create` → `task.start` → implementation → `task.complete`:
   - Compared `agentic-cms`'s existing `init`/`update`/`version` CLI against `smaqit-extensions`' `install.sh` + `installer/main.go` (both on this machine) to scope a *selective* adoption — the user was explicit this was gradual adoption, not a port of `smaqit-extensions`' full CLI surface (no `--install-global`, `--scope`, `uninstall`, or multi-agent global directory scaffolding).
   - Added `install.sh` at repo root (platform/arch detection, `AGENTIC_CMS_VERSION` pinning, download/install/verify/PATH-check), modeled directly on `smaqit-extensions/install.sh`.
   - Added `resolveDefaultProjectDir`/`findAncestorWithEntry` to `update.go`, wired into `runUpdate`, so `agentic-cms update` resolves the project root from any subdirectory instead of only the cwd — simplified from `smaqit-extensions`' version (no git-root-precedence branch, since agentic-cms has only one marker directory).
   - Mid-task, on user request, restructured `README.md` to mirror `ruifrvaz/smaqit`'s section skeleton and style (Features / Compatibility / Getting Started / Commands with a Reinstallation-and-Updates subsection / Documentation).
2. **Diagnosed and fixed a real deployment blocker discovered during verification**: `ruifrvaz/agentic-cms` was private, which silently broke every documented install path (`install.sh` and the pre-existing `go install .../agentic-cms@latest`, since the public Go proxy can't fetch a private module without auth). Assessed three alternatives (make public / token-based auth / GitHub Packages) and recommended making the repo public as the only option that didn't compromise the "no server, frictionless install" design goal. User confirmed it was accidental and made the repo public; `install.sh` and `update`'s ancestor-dir resolution were then re-verified fully unauthenticated end-to-end against the live public repo.
3. **Completed task 006**: hit a real rebase conflict against `origin/main` (the branch predated tasks 005/007 merging) in `README.md` — reconciled by carrying task 005's new classification bullet forward into the restructured README rather than dropping it. Release-analysis computed MINOR → v0.6.0. PR #6 opened, pending changelog entry promoted, merged, cleaned up.
4. **Planned and implemented task 009** on direct user feedback ("the classifier feature needs to surface a bit more in the readme... should mention CIA and have some impactful explanation on why it matters"): promoted the C0-C3 classification feature from a low-level implementation-detail bullet to a headline Features bullet naming the CIA triad, plus a new dedicated `## Classification` section (C0-C3 table, a "why it matters" paragraph grounded in the risk mechanism of agent-driven content synthesis, and the enforcement-mechanics detail moved down from the old bullet). Deliberately did not reference this repo's own internal incident history in the public README.
5. **Completed task 009**: clean rebase this time, release-analysis computed PATCH (docs-only) → v0.6.1. PR #7 opened, pending changelog entry promoted, merged, cleaned up.

## Problems solved

- **`ruifrvaz/agentic-cms` was private** — broke every install path; fixed by making the repo public (user-executed after a `gh` permission scope gap prevented me from doing it directly) and re-verifying live.
- **README rebase conflict (task 006)** — task 006's branch deleted the `## Design notes` section (folding it into a new `## Features` section) while tasks 005/007, already merged, had added new content to that same section. Resolved by merging both: kept the restructure, carried the new content forward.
- **GitHub push access broke twice mid-session** (403 on `git push` despite `gh api` reporting full push/admin rights) — a known recurring pattern where this user's managed PAT rotates/expires. Both times resolved by the user refreshing their token; retried successfully immediately after.
- **A same-instruction fine-grained PAT lacked repo-admin scope** for changing visibility via `gh repo edit --visibility public` (403, distinct from the push-access issue) — user made the change manually via GitHub Settings instead.

## Decisions made

- **Selective adoption, not a port** — `install.sh` deliberately omits `smaqit-extensions`' `--install-global`/`--scope`/`uninstall`/multi-agent surface; agentic-cms's `init` is already project-scoped and has no global directories to seed.
- **No new CI test for `install.sh`** — matches `smaqit-extensions` precedent (its own `install.sh` is untested); covered by `shellcheck` plus manual end-to-end verification instead.
- **Fix the private-repo blocker by flipping visibility, not by adding token auth** — the existing `go install` path already implicitly required public visibility, so this wasn't a new constraint introduced by task 006, just the first thing to actually exercise it. Token-based private-repo support would add real adoption friction and was explicitly assessed and rejected for this use case.
- **README restructure mirrors `ruifrvaz/smaqit`'s section skeleton and style** (Features/Compatibility/Getting Started/Commands/Documentation) while keeping agentic-cms-specific sections (How it works, What gets installed, Repository layout, Development, Roadmap, License) that `smaqit`'s own README doesn't need.
- **Classification "why it matters" is grounded in the general risk mechanism**, not this repo's own internal incident history (referenced in task 008's notes) — not appropriate content for a public README.

## Files modified

- `install.sh` (new) — shell installer bootstrap
- `update.go`, `update_test.go` — ancestor-directory resolution for `update`
- `README.md` — restructured (task 006) then given a dedicated Classification section (task 009); also further reorganized by a direct user commit after task 009 merged (outside this session's task scope)
- `CHANGELOG.md` — v0.6.0 and v0.6.1 entries
- `.smaqit/tasks/006_shell_installer_bootstrap.md`, `.smaqit/tasks/009_surface_classification_in_readme.md`, `.smaqit/tasks/PLANNING.md` — task lifecycle
- `.smaqit/references/project-research.md` — task 006's research block added
- `agentic-cms.code-workspace` — worktree add/remove cycles for both tasks

## Next steps

- Task 010 ("Always overwrite scaffolding logic files on init/update") is `PR Open (#8)` in a separate/parallel session — not touched here, nothing pending from this session's side.
- Task 008 ("Global classification scanner") is `Not Started`; its design will supersede task 006's "no global agent/skill directories to seed" rationale for `install.sh` — worth revisiting `install.sh` when task 008 starts.
- Task 003 ("Register mode for reference files") remains `Not Started`, unblocked.
- The direct "cleaned up readme" commit made on `main` right after task 009 merged (by the user, outside any task) further reorganized the README — worth a quick look next session to confirm nothing from task 009's content was inadvertently lost in that pass.

## Session Metrics

- 2 tasks completed (006, 007→009 numbering: 006 and 009), 2 releases shipped (v0.6.0, v0.6.1), 2 PRs merged (#6, #7)
- 1 real production blocker found and fixed (private repo silently breaking all install paths)
- 1 rebase conflict encountered and reconciled (README, task 006)
- 2 GitHub credential/push failures encountered and resolved via user token refresh
- 1 repo-visibility change (private → public) coordinated with the user after a permission-scope gap blocked doing it directly
