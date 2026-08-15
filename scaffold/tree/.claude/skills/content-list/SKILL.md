---
name: content-list
description: List and summarize the contents of the CMS — topics, items, wiki pages, recent activity. Use for "what's in the content base", "list topics", "what changed recently", or a health overview.
license: MIT
compatibility: Requires bash and python3 (uses the .agentic-cms/bin toolkit)
allowed-tools: Read Grep Glob Bash(.agentic-cms/bin/*)
---

# content-list — inventory and status

Read `CONTENT.md` at the project root first if you haven't this session. This
skill is almost fully deterministic — three commands, then formatting.

## Steps

1. **Gather** (deterministic):
   ```sh
   .agentic-cms/bin/ac-inventory
   .agentic-cms/bin/ac-index check
   .agentic-cms/bin/ac-log tail 5
   ```
2. **Report** concisely from the JSON:
   - Topics with item counts (`topics`)
   - Drafts by topic (`drafts`) — work in progress, kept separate from
     first-class content; suggest promoting via `content-new-item` when one
     looks ready
   - Wiki pages by category (`wiki.entities/concepts/sources`)
   - Raw sources not yet ingested (`raw_uningested`) — suggest `content-import`
   - Last 5 log entries (`recent_log`)
   - Drift (`ac-index check`: `dead_entries`, `unindexed_pages`) — drafts never
     appear here by design, not a drift signal
3. **Offer to fix drift** if `ac-index check` was not clean:
   ```sh
   .agentic-cms/bin/ac-index remove <dead-path>              # per dead entry
   .agentic-cms/bin/ac-index add <section> <path> "<summary>" # per unindexed page
   .agentic-cms/bin/ac-log append lint "index repair" "<what was fixed>"
   ```
   (Write the one-line summaries yourself — that's the only judgment here.)

## Rules

- Read-only except for offered index repairs. Do not modify content pages.
- Do not log pure listings; log only repairs.
