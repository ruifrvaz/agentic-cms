# ac toolkit — deterministic CMS operations

Executable scripts installed by agentic-cms. All commands: run from the project
root, emit a single JSON object on stdout, exit non-zero with
`{"ok": false, "error": "..."}` on failure. Agents: always check `"ok"` before
proceeding, and prefer these commands over hand-editing `wiki/index.md`,
`wiki/log.md`, or frontmatter — they enforce the schema mechanically.

Requires bash and python3 (stdlib only).

## Commands

### ac-page — pages and frontmatter
```
ac-page new <template> <dest.md> --title <T> [--topic <t>] [--raw-path <p>] [--status draft|final]
  templates: doc entity concept source topic  (from .agentic-cms/templates/)
  → {"ok":true,"path":...,"template":...,"title":...,"slug":...}
  Refuses to overwrite. Fills {{TITLE}} {{TOPIC}} {{RAW_PATH}} {{STATUS}} {{DATE}}.
  --status defaults to final; pass --status draft for docs/<topic>/drafts/ items.
ac-page meta <file.md>    → {"ok":true,"path":...,"frontmatter":{...}}
ac-page touch <file.md>   → sets updated: to today
ac-page promote <src.md> <dest.md>
  → {"ok":true,"path":...,"from":...,"status":"final","updated":...}
  Moves src to dest, sets status: final, touches updated:. Refuses to
  overwrite an existing dest. Used to graduate a draft out of drafts/,
  and to un-archive an item out of archive/.
ac-page archive <src.md> <dest.md>
  → {"ok":true,"path":...,"from":...,"status":"archived","updated":...}
  Mirror of promote: moves src to dest, sets status: archived (inserting
  the field if absent), touches updated:. Refuses to overwrite. Used to
  retire an item into docs/<topic>/archive/.
```

### ac-index — wiki/index.md maintenance
```
ac-index add <topics|entities|concepts|sources|archived> <path> <summary>
  → inserts "- [Title](path) — summary" under the section (title read from the page);
    a known section missing its header (older index.md) is appended automatically
ac-index remove <path>
ac-index list             → {"ok":true,"entries":[{section,title,path,summary}]}
ac-index check            → {"ok":true,"clean":bool,"dead_entries":[],"unindexed_pages":[]}
```

### ac-log — wiki/log.md journal
```
ac-log append <operation> <subject> [detail]
  operations: init new new-item import research notes query lint export archive unarchive
ac-log tail [n]           → {"ok":true,"entries":[{date,operation,subject}]}
```

### ac-links — link integrity
```
ac-links check            → {"ok":true,"clean":bool,"links_checked":n,"broken":[{page,line,target}]}
  Accepts links relative to the page OR to the project root.
```

### ac-inventory — full content-base inventory
```
ac-inventory              → {"ok":true,"topics":{...},"drafts":{...},"archived":{...},"wiki":{...},"raw_uningested":[],"recent_log":[...]}
  drafts: {<topic>: {"count":n,"files":[...]}} — status: draft items in docs/<topic>/drafts/,
    not yet wired into the index (see CONTENT.md's Drafts convention)
  archived: same shape — status: archived items in docs/<topic>/archive/, still
    indexed under the Archived section (see CONTENT.md's Archive convention)
```

### ac-search — term search (retrieval)
```
ac-search <term> [term...]  → {"ok":true,"files_matched":n,"results":[{file,hit_count,hits:[{line,terms,text}]}]}
  Terms OR-ed, case-insensitive, ranked by hits, wiki/ and docs/ only.
```

## Composition example

```sh
.agentic-cms/scripts/ac-page new doc docs/ai/agents.md --title "Agents" --topic ai
.agentic-cms/scripts/ac-index add topics docs/ai/agents.md "Agent architectures"
.agentic-cms/scripts/ac-log append new-item "ai/agents"
.agentic-cms/scripts/ac-index check   # verify clean afterwards
```
