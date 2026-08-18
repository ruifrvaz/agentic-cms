---
name: content-research
description: Research a topic on the web and file the findings into the CMS. Use when the user wants to gather external information — "research X", "find out about Y", "what's the latest on Z" — and keep the results in the content base.
license: MIT
compatibility: Requires bash, python3, and web access (uses the .agentic-cms/scripts toolkit)
allowed-tools: Read Write Edit Grep Glob WebSearch WebFetch Bash(.agentic-cms/scripts/*)
---

# content-research — web research into the content base

Read `CONTENT.md` at the project root first if you haven't this session. Toolkit
contract: `.agentic-cms/scripts/README.md` — check `"ok"` after every `ac-*` call.

## Steps

1. **Scope against what's known** (deterministic):
   ```sh
   .agentic-cms/scripts/ac-search <topic terms>
   .agentic-cms/scripts/ac-index list
   ```
   Read the matching pages; research the gaps, not the settled ground.
2. **Delegate** to the `content-researcher` subagent with: the research question,
   the relevant existing page paths, and instructions to return findings with
   source URLs. If subagents are unavailable, do the web research yourself.
3. **File the findings** — web research is almost always C0/C1 (public
   sources, synthesized), but rate deliberately rather than assuming: a
   research question about internal-only material, or findings that
   incorporate something the user shared, can carry a higher rating.
   - New material:
     ```sh
     .agentic-cms/scripts/ac-page new doc docs/<topic>/<item>.md --title "<Title>" --topic <topic> --classification <C0|C1|C2|C3>
     .agentic-cms/scripts/ac-index add topics docs/<topic>/<item>.md "<summary>"
     ```
     Write the findings in (judgment), URLs listed in `sources:`.
   - Updates to existing pages: edit in place, flag contradictions
     explicitly. If the update raises the page's sensitivity, re-rate
     (`ac-page classify <file> <higher-level>`, which also touches
     `updated:`); otherwise `.agentic-cms/scripts/ac-page touch <file>`.
   - Update affected wiki entity/concept pages and cross-links the same way.
4. **Log and verify**:
   ```sh
   .agentic-cms/scripts/ac-log append research "<question>" "<pages touched>"
   .agentic-cms/scripts/ac-index check && .agentic-cms/scripts/ac-links check && .agentic-cms/scripts/ac-classify check <pages touched>
   ```
   All three must report `"clean": true`.

## Rules

- Every researched claim carries its source URL.
- Distinguish clearly between what sources say and your own synthesis.
- One-off answers still get filed — explorations compound in the wiki.
