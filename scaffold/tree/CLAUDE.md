<!-- agentic-cms:begin -->
## Content management system

This project uses a markdown-based agentic CMS (installed by agentic-cms). Before
working with anything under `raw/`, `docs/`, or `wiki/`, read `CONTENT.md` — it
defines the layer rules, conventions, and workflows. Use the `content-*` skills in
`.claude/skills/` for all content operations, and the subagents in
`.claude/agents/` (content-researcher, content-importer, content-exporter) for
research, bulk import, and export tasks.

Hard rules: `raw/` is immutable; `wiki/index.md` and `wiki/log.md` are updated on
every content operation; schema changes to `CONTENT.md` require explicit user
approval; every page carries a `classification: C0-C3` rating you set at
write time (see CONTENT.md's Classification section) — never lower one
yourself, only raise.
<!-- agentic-cms:end -->
