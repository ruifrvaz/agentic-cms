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
approval.
<!-- agentic-cms:end -->
