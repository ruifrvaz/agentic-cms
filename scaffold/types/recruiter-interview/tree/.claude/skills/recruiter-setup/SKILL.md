---
name: recruiter-setup
description: Bootstrap a requisition in a fresh CMS — creates the five type topics with their READMEs, the company and role entity skeletons, the competency framework and one competency concept page per competency, the hiring process, the Rounds table, the employer context pages, and ingests whatever role documents the user has dropped in raw/ (job description, leveling guide, loop definition). Use when the user says "open a requisition", "set up the hiring KB", "we're hiring a <role>", or invokes /recruiter-setup.
license: MIT
compatibility: recruiter-interview type (CONTENT.md, Type section); greenfield only — the five topics must not exist yet; uses the .agentic-cms/scripts toolkit and the type templates in .agentic-cms/templates/; delegates web research to the content-researcher subagent
allowed-tools: Read Write Edit Grep Glob Bash(.agentic-cms/scripts/*) Task
---

# recruiter-setup — open a requisition

Read `CONTENT.md` — the Type section in particular — first if you haven't this
session. All `ac-*` commands live in `.agentic-cms/scripts/`, run from the
project root, and return JSON — check `"ok"` after every call (contract:
`.agentic-cms/scripts/README.md`). Every page is rated at write time against
the Classification section; the Type section lists this type's defaults.

## Steps

1. **Gather the requisition facts** in one exchange, if not already given:
   role title, team, level, location/remote policy; hiring manager and
   recruiter; the interviewers already known; the competencies the role is
   assessed on (or "derive from the job description"); the loop as far as
   it is designed (rounds, a take-home?); what sits in `raw/` (job
   description, leveling guide, loop definition — never candidate
   documents).

2. **Guard — greenfield only:**
   ```sh
   .agentic-cms/scripts/ac-inventory
   ```
   If any of `role`, `company`, `interview-design`, `candidates`,
   `decisions` already exists under `topics`, stop and say so: this is a
   one-requisition instance — open a new instance for a new requisition.

3. **Create the five topics**, each from its own template, then write the
   "About this topic" paragraph from the requisition facts:
   ```sh
   .agentic-cms/scripts/ac-page new topic-role             docs/role/README.md             --title "Role"             --topic role
   .agentic-cms/scripts/ac-page new topic-company          docs/company/README.md          --title "Company"          --topic company          --classification C0
   .agentic-cms/scripts/ac-page new topic-interview-design docs/interview-design/README.md --title "Interview Design" --topic interview-design
   .agentic-cms/scripts/ac-page new topic-candidates       docs/candidates/README.md       --title "Candidates"       --topic candidates
   .agentic-cms/scripts/ac-page new topic-decisions        docs/decisions/README.md        --title "Decisions"        --topic decisions
   ```
   Register each: `ac-index add topics docs/<topic>/README.md "<one-line
   summary>"` and `ac-log append new "<topic>"`. The Rounds and Roster
   tables keep only their header rows for now.

4. **Seed the wiki skeletons:**
   ```sh
   .agentic-cms/scripts/ac-page new entity-company         wiki/entities/the-company.md    --title "<Company name>"        --classification C0
   .agentic-cms/scripts/ac-page new entity-role            wiki/entities/the-role.md       --title "<Role title>"
   .agentic-cms/scripts/ac-page new concept-culture-values wiki/concepts/culture-values.md --title "<Company> Culture and Values" --classification C0
   .agentic-cms/scripts/ac-page new concept-role-family    wiki/concepts/<family-slug>.md  --title "<Job family>"
   ```
   The culture page only if the company publishes a values statement; the
   role-family page only if the title maps to an industry job family. Then
   **one `concept-competency` page per competency** — these are the hubs
   every scorecard will link:
   ```sh
   .agentic-cms/scripts/ac-page new concept-competency wiki/concepts/<competency-slug>.md --title "<Competency>"
   ```
   Fill definitions from the requisition facts; `ac-index add
   entities|concepts ...` each one.

5. **Ingest the role documents in `raw/`** (the only documents this type
   stores): for each,
   `ac-page new source-role-document wiki/sources/<slug>.md --title "..."
   --raw-path raw/<file>` (C2 if it carries compensation bands or names
   people), then derive the role items — rate first (a job spec with a
   band is C2 unless the **user** acks the floor):
   ```sh
   .agentic-cms/scripts/ac-page new job-spec             docs/role/job-spec.md             --title "Job Spec — <Role title>" --topic role --raw-path raw/<file> --classification <C0|C1|C2>
   .agentic-cms/scripts/ac-page new competency-framework docs/role/competency-framework.md --title "Competency Framework"     --topic role
   .agentic-cms/scripts/ac-page new hiring-process       docs/role/hiring-process.md       --title "Hiring Process"          --topic role
   ```
   The framework's coverage map must assess every competency at least
   once; the process's interview steps become the **Rounds table** rows in
   `docs/interview-design/README.md` — kebab-case `key` (never a stage
   number), name, duration, format, Guide `_pending_`, Interviewers
   `_pending_`. Every requirement in the job spec maps to a competency;
   any technology it names gets its `wiki/concepts/` or `wiki/entities/`
   page (Type section, "Company stack in the wiki"). If `raw/` is empty,
   author the three items from the requisition facts and say what a
   received document would sharpen.

6. **Employer context** — delegate research to `content-researcher`
   (company site, engineering blog, careers page) and create:
   ```sh
   .agentic-cms/scripts/ac-page new company-overview    docs/company/company-overview.md    --title "Company Overview"                     --topic company --classification C0
   .agentic-cms/scripts/ac-page new engineering-culture docs/company/engineering-culture.md --title "Engineering Culture and Organization" --topic company
   .agentic-cms/scripts/ac-page new culture-values      docs/company/culture-values.md      --title "Culture and Values"                   --topic company --classification C0
   ```
   Map the values to behavioral competencies on the culture pages.
   Register, log (`research`), list in the topic README.

7. **Candidate-facing guide** — last, once rounds are known:
   `ac-page new candidate-preparation-guide docs/role/candidate-preparation-guide.md
   --title "Candidate Preparation Guide" --topic role`; consistent with the
   Rounds table and the guides' assessed competencies, revealing no
   question banks or rubrics.

8. **Rounds**: for each round the user wants designed now, run
   `recruiter-round <key>` (it creates the interview guide and interviewer
   pages and fills the row). For a take-home, `recruiter-round take-home`.

9. **Verify and report:**
   ```sh
   .agentic-cms/scripts/ac-index check && .agentic-cms/scripts/ac-links check && .agentic-cms/scripts/ac-classify sweep
   ```
   All three must report `"clean": true`. Summarize what was created, what
   is pending (rounds without guides, an empty roster), and the next step —
   usually `recruiter-round` for the first round, then
   `recruiter-candidate` as applications arrive.

## Rules

- Never invent company, role, or people facts — research them or ask;
  leave a section empty rather than plausible.
- Candidate documents are never stored in `raw/` (Type section,
  Retention); this skill never touches candidate data at all.
- Type classification defaults apply; index and log one-liners for C2
  pages stay opaque. Ratchet: raise freely; only the user lowers or acks.
- Round keys, tags, and file names never carry stage numbers — the Rounds
  table is the naming authority.
- Every competency has exactly one concept page; every job-spec
  requirement maps to one — unmapped requirements are a framework gap to
  raise, not to paper over.
