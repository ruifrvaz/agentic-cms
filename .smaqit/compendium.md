# Project Compendium

## Content Schema

**How does the draft / work-in-progress content state work?**

An optional `status: draft | final | archived` field lives in `docs/` page frontmatter (`doc.md` template default is `final`). Setting it to `draft` marks a page as not yet first-class content — a brainstorm or half-formed note that will be iterated on before it graduates. Draft pages live at `docs/<topic>/drafts/<item>.md`, near their eventual home, and are deliberately NOT wired into `wiki/index.md` or `wiki/log.md` while draft: no index entry, no log entry. `content-lint`'s orphan check needs no special-casing for this — it only ever walks `ac-index list`, which never contains drafts by construction, so they're structurally unreachable by that check. `ac-index check`'s drift detection explicitly excludes `docs/*/drafts/*.md` from its `unindexed_pages` glob for the same reason.

`status:` is a live, first-class frontmatter field — like `type`, `tags`, `sources`, `refs`, and `classification` — not a commented-out template example, since it's knowable at page-creation time (a mode flag, not a judgment call) and `CONTENT.md` requires mechanical operations to go through the `.agentic-cms/scripts/` toolkit rather than hand-editing frontmatter. Concretely:

- Create a draft: `ac-page new doc docs/<topic>/drafts/<item>.md --title "<T>" --topic <topic> --status draft` (the `doc.md` template's `status: {{STATUS}}` placeholder gets filled; omitting `--status` defaults to `final`).
- Promote a draft: `ac-page promote docs/<topic>/drafts/<item>.md docs/<topic>/<item>.md` — moves the file, sets `status: final`, touches `updated:`, and refuses if the destination already exists (same refuse-on-conflict behavior as `ac-page new`). Follow with the normal `content-manage-item` register+log step (`ac-index add`, `ac-log append`) against the new path — that's the first time the item enters the wiki.
- `content-list` reports drafts in their own section, sourced from `ac-inventory`'s `drafts` field (per-topic path + count).

See also: how does archiving retired content work?

---

**How does archiving retired content work?**

`status: archived` marks a first-class page as retired without deleting it — archive is not delete. Archived pages move to `docs/<topic>/archive/<item>.md`, the structural mirror of `drafts/`. The key difference from drafts: an archived page *was* first-class, so unlike a draft it stays in `wiki/index.md`, re-filed under a `## Archived` section (auto-created if the index predates it), so `content-query` can still find it and knows it's retired.

- Archive an item: `ac-page archive docs/<topic>/<item>.md docs/<topic>/archive/<item>.md` — moves the file, sets `status: archived`, touches `updated:`. Then re-file the index entry (`ac-index remove` + `ac-index add archived ...`), append an `archive` log entry, and run `ac-links check` to catch (and fix) any inbound links from active pages now pointing at the old path — that's a deliberate forcing function, not a false positive.
- Un-archive: `ac-page promote` from `archive/` back up (mirrors draft promotion), plus the reverse index re-filing and an `unarchive` log entry.
- `content-lint`'s orphan check and stale-claims check both skip the `## Archived` section — having no inbound links, or not staying current, is an archived page's expected steady state, not a defect.
- `ac-inventory` reports archived items per topic, same shape as `drafts`.

See also: how does the draft / work-in-progress content state work?

---

## Content Classification

**Why does agentic-cms need confidentiality classification at all?**

An agentic CMS isn't a passive filing cabinet — the agent actively reads across everything dropped into `raw/`, synthesizes it into `docs/`, and cross-links it into a compounding `wiki/`. That synthesis is the entire value proposition, and it's also exactly the mechanism by which one sensitive detail — a customer's PII buried in an old email, a credential pasted into a meeting note — can silently bleed into a summary, an index one-liner, or an exported deck, resurfacing somewhere far from where a human reviewer would have caught it. The `classification: C0-C3` system (the CIA triad's confidentiality axis) makes that risk explicit and enforced instead of assumed, rather than relying on an agent's implicit judgment alone.

See also: how is content classification (confidentiality rating) enforced?

---

**How is content classification (confidentiality rating) enforced?**

Every `docs/`/`wiki/` page carries an optional `classification: C0 | C1 | C2 | C3` frontmatter field (standard CIA confidentiality axis: C0 Public, C1 Internal — the default when absent, C2 Confidential, C3 Restricted), rated by the agent at write time against a rubric in `CONTENT.md`. All detection logic — enum validity, staleness (via a `classified-hash` stamped at rating time), and heuristic floor patterns (credential-shaped strings imply at least C3; PII/currency-shaped content implies at least C2) — lives in exactly one place: `.agentic-cms/scripts/ac-classify` (`check`/`sweep`/`hook` subcommands). Every enforcement point is a thin caller of that one engine, never a re-implementation:

- **Agent hooks** (Claude Code `.claude/settings.json`, Codex `.codex/hooks.json`) fire on every content write, auto-raise floor violations, and block on stale ratings.
- **Git pre-commit gate** (`.agentic-cms/hooks/pre-commit`, wired into `.git/hooks/pre-commit` by `init`) materializes the staged tree via `git checkout-index` and runs the same engine — blocks floor violations and confidentiality leaks into `wiki/index.md`/`wiki/log.md` summaries, warns (non-blocking) on stale or unrated pages.
- **Write-path skills** (`content-manage-item`, `content-import`, `content-add-notes`, `content-research`) rate before writing and include `ac-classify check` in their verify tail — the guaranteed fallback wherever a blocking hook isn't installed or available (e.g. GitHub Copilot CLI, whose hook has no documented blocking capability, so isn't wired in).
- **`content-lint`** runs `ac-classify sweep` as the periodic catch-all for anything the other layers never saw — pre-existing drift, or a commit made with `git commit --no-verify`.

The ratchet rule: an agent may *raise* a page's classification; only the user may lower one. No mechanical path in the toolkit ever lowers a rating — `ac-page classify <path> <level>` will technically accept any level, but skill instructions and the hook's auto-raise logic only ever move upward.

See also: why does agentic-cms need confidentiality classification at all?; how does the git pre-commit classification gate get bypassed?; is classification enforcement available outside a project where `agentic-cms init` was run?

---

**How does the git pre-commit classification gate get bypassed?**

`git commit --no-verify` (or `-n`) is git's own built-in flag — it makes git skip invoking `.git/hooks/pre-commit` entirely before the script ever runs. This isn't something the classification gate can prevent from the inside; it's true of any git pre-commit hook, from any tool. Concrete ways a commit can land without the gate running:

- Deliberate use of `--no-verify`/`-n`.
- The hook was never installed in that particular clone — git hooks live only in a working copy's local `.git/hooks/` and are never copied by `clone`/`pull`/`fetch`, so any checkout that hasn't run `agentic-cms init`/`update` has no gate at all, `--no-verify` or not.
- Commits made outside the local git CLI — a GitHub web-based file edit, a PR merge button, or CI/automation committing via the GitHub API.
- A repository with `core.hooksPath` set to a custom location: `init` only reports that case and never writes there, so the gate stays uninstalled until someone wires it in manually.

Because a pre-commit gate can always be bypassed this way, `content-lint`'s `ac-classify sweep` is the deliberate backstop — it re-checks the whole content base independent of how or whether any commit happened, catching what slipped through. Bypassing the *commit-time* gate also doesn't bypass the agent-hook or skill-verify-tail layers, since those trigger on the file write itself, independent of the commit event.

See also: how is content classification (confidentiality rating) enforced?

---

**Is classification enforcement available outside a project where `agentic-cms init` was run?**

No — as of the current implementation, every classification enforcement layer (the toolkit script, both agent hooks, the git pre-commit gate, the skill verify tails) is installed per-project by `agentic-cms init`, and only exists in that specific target directory. Nothing is written to a user-global location (`~/.claude/skills/`, `~/.claude/settings.json`, or similar), and there is no global git hook mechanism in use. A project that has never run `init` on itself has zero coverage, with no fallback.

See also: how is content classification (confidentiality rating) enforced?

---

## Installation

**Does `install.sh` (or `go install`) support installing from a private GitHub repo?**

No — both installation paths assume `ruifrvaz/agentic-cms` is public. `install.sh` calls the unauthenticated GitHub Releases API and downloads release assets directly; both return `404` on a private repo without a credential. `go install .../agentic-cms@latest` has the same dependency in a different form: the public Go module proxy can't fetch a private module without the installing machine having `GOPRIVATE` and its own git credentials configured.

If private-repo distribution is ever needed, the fix isn't GitHub Packages — a private package requires the same kind of per-user credential distribution as a private release, so it relocates the problem rather than removing it. The correct approach would be an env-var-supplied token (`Authorization: Bearer`) sent on both the GitHub API request and the asset download — via the API's `/releases/assets/{id}` endpoint with `Accept: application/octet-stream`, not the plain `browser_download_url`: forwarding a Bearer token through that URL's redirect to GitHub's signed blob storage returns `404` rather than a clean auth error.

---

## Tooling

**Why do `.agents/`, `.claude/`, `.codex/`, `.github/agents/`, or `.github/skills/` sometimes reappear in this project after being deleted?**

This project uses `smaqit-extensions`' global-only installation: agents and skills live at `~/.agents/skills/`, `~/.claude/skills/`, `~/.claude/agents/`, `~/.claude/commands/`, `~/.codex/agents/`, and `~/.copilot/agents/` — not inside the repository. The only `smaqit-extensions` artifacts that belong in this project are `.smaqit/tasks/`, `.smaqit/history/`, `.smaqit/user-testing/`, `.smaqit/references/`, and `.github/workflows/post-merge-release.yml`.

If a full project-scoped mirror of agents/skills reappears (hundreds of files under those five paths), it's a known bug in `smaqit-extensions update` (not `init`, which was verified clean) — tracked upstream as `smaqit-extensions` task 033. It's safe to delete those five directories again; nothing in this project depends on them being present locally.

---
