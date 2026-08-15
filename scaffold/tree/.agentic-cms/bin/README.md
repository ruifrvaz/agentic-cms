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
ac-page new <template> <dest.md> --title <T> [--topic <t>] [--raw-path <p>]
  templates: doc entity concept source topic  (from .agentic-cms/templates/)
  → {"ok":true,"path":...,"template":...,"title":...,"slug":...}
  Refuses to overwrite. Fills {{TITLE}} {{TOPIC}} {{RAW_PATH}} {{DATE}}.
ac-page meta <file.md>    → {"ok":true,"path":...,"frontmatter":{...}}
ac-page touch <file.md>   → sets updated: to today
```

### ac-index — wiki/index.md maintenance
```
ac-index add <topics|entities|concepts|sources> <path> <summary>
  → inserts "- [Title](path) — summary" under the section (title read from the page)
ac-index remove <path>
ac-index list             → {"ok":true,"entries":[{section,title,path,summary}]}
ac-index check            → {"ok":true,"clean":bool,"dead_entries":[],"unindexed_pages":[]}
```

### ac-log — wiki/log.md journal
```
ac-log append <operation> <subject> [detail]
  operations: init new new-item import research notes query lint export
ac-log tail [n]           → {"ok":true,"entries":[{date,operation,subject}]}
```

### ac-links — link integrity
```
ac-links check            → {"ok":true,"clean":bool,"links_checked":n,"broken":[{page,line,target}]}
  Accepts links relative to the page OR to the project root.
```

### ac-inventory — full content-base inventory
```
ac-inventory              → {"ok":true,"topics":{...},"wiki":{...},"raw_uningested":[],"recent_log":[...]}
```

### ac-search — term search (retrieval)
```
ac-search <term> [term...]  → {"ok":true,"files_matched":n,"results":[{file,hit_count,hits:[{line,terms,text}]}]}
  Terms OR-ed, case-insensitive, ranked by hits, wiki/ and docs/ only.
```

## Composition example

```sh
.agentic-cms/bin/ac-page new doc docs/ai/agents.md --title "Agents" --topic ai
.agentic-cms/bin/ac-index add topics docs/ai/agents.md "Agent architectures"
.agentic-cms/bin/ac-log append new-item "ai/agents"
.agentic-cms/bin/ac-index check   # verify clean afterwards
```
