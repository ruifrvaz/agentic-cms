<!-- ac-type-anchor: after "  sources/            one summary page per ingested raw source" -->
exercises/            type-specific: runnable coding drills (see Type section)
<!-- ac-type-anchor: replace "  templates/          markdown templates for every page type" -->
  templates/          markdown templates for every page type — the base
                      kinds (doc, entity, concept, source, topic) plus the
                      type-specific templates listed in the Type section
<!-- ac-type-anchor: after "  VERSION             scaffold version installed by agentic-cms" -->
  TYPE.md             the installed content type and what it owns
<!-- ac-type-anchor: replace "  skills/             content-* skills (workflows)" -->
  skills/             content-* skills (workflows) + type-specific skills
<!-- ac-type-anchor: after "- **Filenames**: kebab-case, descriptive, no dates in names (`transformer-architecture.md`)." -->
  **Exception (type-defined)**: the `interview-practice` topic uses numbered
  session chains of `type: rehearsal` pages sharing one sequential `NNN` per
  session — `NNN-interview.md` (verbatim mock-interview transcript),
  `NNN-debrief.md` (its scored debrief), and optionally `NNN-refresher.md`
  (gap-targeted course) and `NNN-exam.md` (multiple-choice exam generated from
  the same gaps). Templates: `.agentic-cms/templates/rehearsal-{interview,
  debrief,refresher,exam}.md`. See the Type section.
<!-- ac-type-anchor: replace "  type: doc | entity | concept | source | note" -->
  type: doc | entity | concept | source | note | rehearsal
<!-- ac-type-anchor: before "## Operations" -->
## Type: candidate-interview

This CMS instance runs the **candidate-interview** content type: a personal
knowledge base supporting one job-interview recruitment process end to end,
from the candidate's side. Everything in this section is type-defined on top
of the base schema above; the base schema still applies in full. The type is
company-, role-, and stack-agnostic by design — nothing below names a company,
a role, or a technology; those are per-engagement content.

### Topics

A fresh instance is empty. `interview-setup` creates the six topics below
(each `docs/<topic>/README.md` from its `topic-<topic>` template); they are
the type's vocabulary and every item lands in exactly one of them.

| Topic | Holds | Item templates |
|---|---|---|
| `company` | research on the target company | `company-overview`, `engineering-culture`, `culture-values` |
| `candidate` | the candidate's own material | `candidate-cv`, `public-footprint`, `role-fit-analysis` |
| `interview-prep` | the hiring process and one plan/debrief pair per round | `hiring-process`, `role-research`, `candidate-preparation-guide`, `round-plan`, `round-debrief` |
| `interview-practice` | mock-interview session chains (`NNN-*`) | `rehearsal-interview`, `rehearsal-debrief`, `rehearsal-refresher`, `rehearsal-exam` |
| `offer` | outcome/leveling/compensation correspondence, decoded | `offer-outcome` |
| `poc` | plan and learning log of a technical proof-of-concept built for the loop | `poc-plan`, `poc-learning-log` |

Wiki templates: entities `entity-company`, `entity-candidate`,
`entity-flagship-project`, `entity-poc`, `entity-panelist` (one page per
interviewer — a real person in a live engagement, so C2); concepts
`concept-culture-values` (the company's published values) and
`concept-role-family` (the industry job family the target role maps to, if
any); sources `source-candidate-cv`, `source-preparation-guide`,
`source-recruiter-email`. Generic `doc`/`entity`/`concept`/`source` templates
remain available for anything the type didn't anticipate.

### Rounds

`docs/interview-prep/README.md` carries a **Rounds** table — one row per
round of the real process, with a `key`, `name`, `duration`, `format`
summary, and links to that round's plan and debrief pages. It is the single
naming authority for rounds: skills select a round by its `key`, tags use
round keys (never stage numbers), and the round's plan page carries a
`## Panel` section linking the `wiki/entities/` panelist pages the rehearsal
skill plays in character. `interview-round` adds a round (plan page, panelist
entities, table row) and later files its as-run debrief.

### Practice loop

`interview-rehearsal` runs a mock round against the plan pages and files the
`NNN-interview.md` + `NNN-debrief.md` pair; `interview-refresher` extends the
same chain with `NNN-refresher.md` (one unit per gap) and, once every unit is
checked, `NNN-exam.md` (multiple-choice, answer key in an HTML comment,
graded in place). The next rehearsal reads the latest debrief and graded
exam, so weaknesses carry forward mechanically. All chain pages are C2:
they contain the candidate's real answers and panel personas derived from
research on real people.

### Coding exercises

`exercises/NNN-<slug>/` (repo root) holds real, runnable code for
HackerRank-style drills matching a coding-round format — "complete an
unfinished codebase by fixing unit tests": a small, self-contained,
stdlib-only module with an intentionally incomplete/buggy implementation
plus a full unit test suite that only passes once fixed, verified against a
reference solution before being handed over (the reference is never shown).
Exercises are the type's one deliberate exception to "markdown only": they
are **not** CMS pages — no frontmatter, no index/log entries, no
classification — ordinary tracked code. The coding round's language is
per-engagement information like the company's stack; the type ships one
illustrative exercise, not a language curriculum. Write new exercises in
whatever language the real round uses. A debrief of each attempt is CMS
content and belongs in `interview-prep`/`interview-practice`, cross-linking
the exercise folder by path.

### Company stack in the wiki

The type ships **no technology-specific content** — no concept or entity
page for any cloud platform, framework, or tool. A specific company's stack
is per-engagement research, not part of the mechanism. When it surfaces
(candidate-prep guide, job spec, PoC work, rehearsal prep), capture it the
normal way: one `wiki/concepts/<technology>.md` per general technology or
practice the round prep depends on, one `wiki/entities/<product>.md` per
specific product the company uses, linked from the round plans, PoC pages,
and panelist pages that need them. Skills and templates refer to "the
target stack" generically and never hardcode a technology.

### Classification defaults for this type

Recruiter correspondence, compensation figures, interview specifics, panel
research on real people, the candidate's CV/PII, private strategy
(role-fit analysis), and every rehearsal-chain page are **C2**. A verbatim
offer letter or contract is **C3**. Public company research is C0/C1. Panel
research must close an identity-verification chain before a real person's
profile is filed (employer named on a professional profile, a code account
tied to a personal site, commit history matching a listed employer — never
name similarity alone); namesakes found during research are recorded as
excluded so they cannot re-enter later.

<!-- ac-type-anchor: before "## Greenfield vs brownfield" -->
**Type operations** (skills: `interview-setup`, `interview-round`,
`interview-rehearsal`, `interview-refresher`): bootstrap an engagement, add
and debrief rounds, run the practice loop — see the Type section.

<!-- ac-type-anchor: replace "  `docs/` and `wiki/` from scratch." -->
  `docs/` and `wiki/` from scratch. For this type, run `interview-setup` first —
  it lays down the six type topics so `content-new` is only needed for
  anything beyond them.
