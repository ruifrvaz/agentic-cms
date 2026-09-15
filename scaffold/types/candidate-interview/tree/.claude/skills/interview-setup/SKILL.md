---
name: interview-setup
description: Bootstrap a candidate-interview engagement in a fresh CMS — creates the six type topics with their READMEs, the company/candidate entity skeletons, the culture-values and role-family concept pages, the company and role research items, the Rounds table, and ingests whatever the user has already dropped in raw/ (candidate-preparation guide, CV). Use when the user says "new engagement", "set up the interview KB", "start preparing for <company>", or invokes /interview-setup.
license: MIT
compatibility: candidate-interview type (CONTENT.md, Type section); greenfield only — the six topics must not exist yet; uses the .agentic-cms/scripts toolkit and the type templates in .agentic-cms/templates/; delegates web research to the content-researcher subagent
allowed-tools: Read Write Edit Grep Glob Bash(.agentic-cms/scripts/*) Task
---

# interview-setup — bootstrap an engagement

Read `CONTENT.md` — the Type section in particular — first if you haven't this
session. All `ac-*` commands live in `.agentic-cms/scripts/`, run from the
project root, and return JSON — check `"ok"` after every call (contract:
`.agentic-cms/scripts/README.md`). Every page is rated at write time against
the Classification section; the Type section lists this type's defaults.

## Steps

1. **Gather the engagement facts** in one exchange, if not already given:
   the target company (name, sector, HQ), the target role (title, team,
   scope), what already sits in `raw/` (a candidate-preparation guide, the
   CV, a job spec, recruiter mail), and the rounds as far as they are known
   (or "unknown — derive from the guide"). Record company and role in
   `AGENTS.md`'s Domain Context placeholders.

2. **Guard — greenfield only:**
   ```sh
   .agentic-cms/scripts/ac-inventory
   ```
   If any of `company`, `candidate`, `interview-prep`, `interview-practice`,
   `offer`, `poc` already exists under `topics`, stop and say so: extend the
   existing engagement with `interview-round` and `content-manage-item`
   instead of re-running setup.

3. **Create the six topics**, each from its own template, then write the
   "About this topic" paragraph from the engagement facts:
   ```sh
   .agentic-cms/scripts/ac-page new topic-company            docs/company/README.md            --title "Company"            --topic company            --classification C0
   .agentic-cms/scripts/ac-page new topic-candidate          docs/candidate/README.md          --title "Candidate"          --topic candidate
   .agentic-cms/scripts/ac-page new topic-interview-prep     docs/interview-prep/README.md     --title "Interview Prep"     --topic interview-prep
   .agentic-cms/scripts/ac-page new topic-interview-practice docs/interview-practice/README.md --title "Interview Practice" --topic interview-practice
   .agentic-cms/scripts/ac-page new topic-offer              docs/offer/README.md              --title "Offer"              --topic offer
   .agentic-cms/scripts/ac-page new topic-poc                docs/poc/README.md                --title "PoC — Interview Demo" --topic poc
   ```
   Register each: `ac-index add topics docs/<topic>/README.md "<one-line
   summary>"` and `ac-log append new "<topic>"`. Leave the Rounds table in
   `docs/interview-prep/README.md` with only its header row for now.

4. **Seed the wiki skeletons** — only pages the engagement facts justify:
   ```sh
   .agentic-cms/scripts/ac-page new entity-company        wiki/entities/the-company.md    --title "<Company name>"   --classification C0
   .agentic-cms/scripts/ac-page new entity-candidate      wiki/entities/the-candidate.md  --title "<Candidate name>"
   .agentic-cms/scripts/ac-page new concept-culture-values wiki/concepts/culture-values.md --title "<Company> Culture and Values" --classification C0
   .agentic-cms/scripts/ac-page new concept-role-family   wiki/concepts/<family-slug>.md  --title "<Job family>"
   ```
   The culture page only if the company publishes a values statement; the
   role-family page only if the role title maps to an industry job family
   (see the `role-research` template); `entity-poc` / `entity-flagship-project`
   only when the user plans a PoC or has a standout public project. Fill
   "What it is" / "Key facts" from what you were told; leave the rest for
   research. `ac-index add entities|concepts ...` each one.

5. **Research and file the company and role items** — delegate web
   research to the `content-researcher` subagent (company site, engineering
   blog, careers page, public postings for the role and adjacent titles,
   the industry definition of the job family), then create and write:
   ```sh
   .agentic-cms/scripts/ac-page new company-overview     docs/company/company-overview.md          --title "Company Overview"                   --topic company        --classification C0
   .agentic-cms/scripts/ac-page new engineering-culture  docs/company/engineering-culture.md       --title "Engineering Culture and Organization" --topic company      --classification C0
   .agentic-cms/scripts/ac-page new culture-values       docs/company/culture-values.md            --title "Culture and Values"                 --topic company        --classification C0
   .agentic-cms/scripts/ac-page new hiring-process       docs/interview-prep/hiring-process.md     --title "Hiring Process"                     --topic interview-prep --classification C0
   .agentic-cms/scripts/ac-page new role-research        docs/interview-prep/role-research.md      --title "Role Research"                      --topic interview-prep
   ```
   A company page carrying financial figures trips the C2 currency floor —
   rate it C2, or leave it C0/C1 only if the **user** acks the floor. Any
   technology the research names gets its own `wiki/concepts/` or
   `wiki/entities/` page (Type section, "Company stack in the wiki") — never
   inline stack facts into prose only. Register, log (`research`), and list
   every item in its topic README.

6. **Ingest what is already in `raw/`:**
   - Candidate-preparation guide →
     `ac-page new candidate-preparation-guide docs/interview-prep/candidate-preparation-guide.md --title "Candidate Preparation Guide" --topic interview-prep --raw-path raw/<file> --classification C2`
     plus `ac-page new source-preparation-guide wiki/sources/<slug>.md --title "..." --raw-path raw/<file> --classification C2`.
     Its stage table becomes the **Rounds table** in
     `docs/interview-prep/README.md`: one row per round — kebab-case `key`
     (no stage numbers), name, duration, format, Plan `_pending_`, Debrief
     `_pending_`. Reconcile `role-research` and `hiring-process` against it
     (the guide is authoritative where they differ; note corrections
     explicitly).
   - CV → `ac-page new candidate-cv docs/candidate/cv.md --title "CV — <name>" --topic candidate --raw-path raw/assets/<file> --classification C2`
     plus `source-candidate-cv`; update the candidate entity's Key facts.
   - Then `public-footprint` (C1; research the candidate's public profiles)
     and, last, `role-fit-analysis` (C2) synthesized from CV + guide +
     company pages, with its round-by-round table mirroring the Rounds
     table. Use `import` / `research` log operations accordingly.
   If `raw/` is empty, skip this step and say which sources would unlock
   the rest (the guide above all).

7. **Rounds**: for each round the user wants planned now, run
   `interview-round <key>` (it creates the plan page and panelist pages and
   fills the row's Plan link). Rounds without a known panel stay as rows
   with `_pending_` links.

8. **Verify and report:**
   ```sh
   .agentic-cms/scripts/ac-index check && .agentic-cms/scripts/ac-links check && .agentic-cms/scripts/ac-classify sweep
   ```
   All three must report `"clean": true`. Then summarize what was created,
   what is pending (rounds without plans, an empty `offer` topic, a PoC
   plan not yet written — `poc-plan` / `poc-learning-log` via
   `content-manage-item`), and the first recommended next step.

## Rules

- Never invent company, role, or panel facts — research them or ask;
  leave a section empty rather than plausible.
- Type classification defaults apply (Type section): PII, panel research,
  interview specifics, private strategy are C2; public research C0/C1.
  Index and log one-liners for C2 pages stay opaque.
- Ratchet: you may raise a rating; only the user may lower one or ack a
  floor.
- Don't create pages for topics with nothing to say: `offer` stays a bare
  README until an offer exists; `poc` until a PoC is planned;
  `interview-practice` is filled only by `interview-rehearsal`.
- Round keys, tags, and file names never carry stage numbers — the Rounds
  table is the naming authority.
