# Project Compendium

## Content Schema

**How does the draft / work-in-progress content state work?**

An optional `status: draft | final` field lives in `docs/` page frontmatter (`doc.md` template default is `final`). Setting it to `draft` marks a page as not yet first-class content — a brainstorm or half-formed note that will be iterated on before it graduates. Draft pages live at `docs/<topic>/drafts/<item>.md`, near their eventual home, and are deliberately NOT wired into `wiki/index.md` or `wiki/log.md` while draft: no index entry, no log entry. `content-lint`'s orphan check needs no special-casing for this — it only ever walks `ac-index list`, which never contains drafts by construction, so they're structurally unreachable by that check. `ac-index check`'s drift detection explicitly excludes `docs/*/drafts/*.md` from its `unindexed_pages` glob for the same reason.

`status:` is a live, first-class frontmatter field — like `type`, `tags`, `sources`, and `refs` — not a commented-out template example, since it's knowable at page-creation time (a mode flag, not a judgment call) and `CONTENT.md` requires mechanical operations to go through the `.agentic-cms/bin/` toolkit rather than hand-editing frontmatter. Concretely:

- Create a draft: `ac-page new doc docs/<topic>/drafts/<item>.md --title "<T>" --topic <topic> --status draft` (the `doc.md` template's `status: {{STATUS}}` placeholder gets filled; omitting `--status` defaults to `final`).
- Promote a draft: `ac-page promote docs/<topic>/drafts/<item>.md docs/<topic>/<item>.md` — moves the file, sets `status: final`, touches `updated:`, and refuses if the destination already exists (same refuse-on-conflict behavior as `ac-page new`). Follow with the normal `content-new-item` register+log step (`ac-index add`, `ac-log append`) against the new path — that's the first time the item enters the wiki.
- `content-list` reports drafts in their own section, sourced from `ac-inventory`'s `drafts` field (per-topic path + count).

---

## Tooling

**Why do `.agents/`, `.claude/`, `.codex/`, `.github/agents/`, or `.github/skills/` sometimes reappear in this project after being deleted?**

This project uses `smaqit-extensions`' global-only installation: agents and skills live at `~/.agents/skills/`, `~/.claude/skills/`, `~/.claude/agents/`, `~/.claude/commands/`, `~/.codex/agents/`, and `~/.copilot/agents/` — not inside the repository. The only `smaqit-extensions` artifacts that belong in this project are `.smaqit/tasks/`, `.smaqit/history/`, `.smaqit/user-testing/`, `.smaqit/references/`, and `.github/workflows/post-merge-release.yml`.

If a full project-scoped mirror of agents/skills reappears (hundreds of files under those five paths), it's a known bug in `smaqit-extensions update` (not `init`, which was verified clean) — tracked upstream as `smaqit-extensions` task 033. It's safe to delete those five directories again; nothing in this project depends on them being present locally.

---
