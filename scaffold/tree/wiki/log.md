---
title: Operations Log
type: doc
created: {{DATE}}
updated: {{DATE}}
---

# Log

Append-only journal of every operation performed on this content base. Entry format:

    ## [YYYY-MM-DD] <operation> | <subject>

Greppable with: `grep "^## \[" wiki/log.md | tail -5`

---

## [{{DATE}}] init | scaffold installed by agentic-cms
