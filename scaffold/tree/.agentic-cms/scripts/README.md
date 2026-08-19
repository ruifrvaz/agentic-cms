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
            [--status draft|final] [--classification C0|C1|C2|C3]
  templates: doc entity concept source topic  (from .agentic-cms/templates/)
  → {"ok":true,"path":...,"template":...,"title":...,"slug":...,"classification":...}
  Refuses to overwrite. Fills {{TITLE}} {{TOPIC}} {{RAW_PATH}} {{STATUS}}
  {{CLASSIFICATION}} {{CLASSIFIED_HASH}} {{DATE}}.
  --status defaults to final; pass --status draft for docs/<topic>/drafts/ items.
  --classification defaults to C1; stamps classified-hash from the page's
  initial body so ac-classify check reports it as freshly rated.
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
ac-page classify <file.md> <C0|C1|C2|C3> [--ack-floor]
  → {"ok":true,"path":...,"classification":...,"classified-hash":...,"updated":...[,"classification-ack":...]}
  Sets classification (inserting the field if absent), restamps
  classified-hash from the file's current body, touches updated:. Purely
  mechanical — does not enforce the ratchet rule (raise-only); that is a
  skill-instruction and hook-logic policy, not a tool constraint.
  --ack-floor additionally stamps classification-ack: bound to the same
  body hash — the USER's recorded decision that a heuristic floor hit is
  a false positive and the rating stands. Agents never pass it on their
  own judgment, only relaying an explicit user instruction. Any body
  change invalidates the ack (hash mismatch); a plain classify without
  the flag withdraws a standing ack.
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
  operations: init new new-item import research notes query lint export archive unarchive classify
ac-log tail [n]           → {"ok":true,"entries":[{date,operation,subject}]}
```

### ac-links — link integrity
```
ac-links check            → {"ok":true,"clean":bool,"links_checked":n,"broken":[{page,line,target}]}
  Accepts links relative to the page OR to the project root.
```

### ac-inventory — full content-base inventory
```
ac-inventory              → {"ok":true,"topics":{...},"drafts":{...},"archived":{...},
                              "classification":{...},"wiki":{...},"raw_uningested":[],"recent_log":[...]}
  drafts: {<topic>: {"count":n,"files":[...]}} — status: draft items in docs/<topic>/drafts/,
    not yet wired into the index (see CONTENT.md's Drafts convention)
  archived: same shape — status: archived items in docs/<topic>/archive/, still
    indexed under the Archived section (see CONTENT.md's Archive convention)
  classification: {"C0":n,"C1":n,"C2":n,"C3":n,"unrated":n} — a frontmatter
    tally only (counts the classification: value already on each page); for
    staleness/floor/bleed *detection* use ac-classify check|sweep instead —
    this field is the cheap read-only count, that is the audit
```

### ac-classify — classification detection engine
```
ac-classify check <path...>
  → {"ok":true,"clean":bool,"pages":[{path,classification,valid,unrated,
                                       stale,floor,floor_violation,acked[,implied_level]}]}
  Per-page: is the classification value legal, is it unrated (no field —
  defaults to C1), is it stale (classified-hash doesn't match the current
  body — content changed since last rated), and does the body's content
  trip a heuristic floor pattern (credential-shaped → C3; email/SSN/
  currency-shaped in any form — symbol-prefixed/suffixed, ISO-coded,
  word-form → C2) above the current rating. The floor's contract is
  recall (0% false negatives over its enumerated shapes); precision is
  the ack's job, never a narrower pattern. A user floor-ack
  (classification-ack: matching the current body hash — see ac-page
  classify --ack-floor) reports acked:true and suppresses the
  floor_violation while the body is unchanged; the floor itself stays
  visible. clean is false only on an invalid value or a floor violation —
  staleness and unrated are advisory, not blocking, by design (see
  CONTENT.md's Classification section for the block-vs-warn split).
ac-classify sweep
  → {"ok":true,"clean":bool,"pages":[...],"bleed":[{page,location,line,text,detected}]}
  check across every docs/**/*.md, wiki/entities/*.md, wiki/concepts/*.md,
  wiki/sources/*.md, plus the bleed check: any wiki/index.md or
  wiki/log.md line that itself trips a floor pattern while referencing a
  C2+ page (independent of that page's own body) — the C2+ index/log
  bleed rule.
ac-classify hook
  Reads a Claude Code/Codex PostToolUse JSON payload from stdin (auto-
  detects tool_input.file_path/path shape), filters to docs/**/wiki/**
  markdown paths only — everything else is a silent {"ok":true,"skip":true}.
  A floor violation auto-applies (ac-page classify to the implied level +
  an ac-log classify entry) and returns {"ok":true,"auto_raised":true,...}.
  A stale rating blocks: exit 2, {"decision":"block","reason":...} on
  stdout, the reason also on stderr. Otherwise {"ok":true,"clean":true}.
  This is the ONLY engine — every caller (skill verify tails, both agent
  hooks, the git pre-commit gate, content-lint) invokes check/sweep/hook;
  none re-implements detection.
```
Wired in by the scaffold: `.claude/settings.json` and `.codex/hooks.json`
(both PostToolUse → `ac-classify hook`), and `.agentic-cms/hooks/pre-commit`
(a versioned script `agentic-cms init` wires into `.git/hooks/pre-commit`
non-destructively — see the root README's install notes). The pre-commit
gate is delta-scoped: it blocks only on violations in the files staged in
that commit (plus bleed in a staged wiki/index.md or wiki/log.md,
whichever page leaked), and summarizes pre-existing drift elsewhere as one
non-blocking warning — content-lint owns that backlog. Copilot CLI's
hook is not wired in this version (no documented blocking capability as of
this writing) — the skill-tail `ac-classify check` calls are the guaranteed
fallback there.

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
