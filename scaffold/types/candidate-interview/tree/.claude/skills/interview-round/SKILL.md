---
name: interview-round
description: Add a round to the engagement — a round-plan page, one identity-verified wiki entity page per panelist, and the row in the Rounds table — or, once the real round has been held, file its as-run debrief and feed the implications forward into later rounds. Use when the user says "plan the <round>", "add a round", "the panel for <round> is X and Y", "debrief the <round>", or invokes /interview-round <key> [debrief].
license: MIT
compatibility: candidate-interview type (CONTENT.md, Type section); requires docs/interview-prep/README.md with a Rounds table (created by interview-setup); uses the .agentic-cms/scripts toolkit and the round-plan, round-debrief, and entity-panelist templates; delegates panelist research to the content-researcher subagent
allowed-tools: Read Write Edit Grep Glob Bash(.agentic-cms/scripts/*) Task
---

# interview-round — plan a round, then debrief it

Read `CONTENT.md`'s Type section first if you haven't this session. All
`ac-*` commands live in `.agentic-cms/scripts/`, return JSON — check `"ok"`
after every call. Two modes: **plan** (default; before the round) and
**debrief** (after the real round has been held).

## Plan mode

1. **Load grounding:** `docs/interview-prep/README.md` (the Rounds table),
   `candidate-preparation-guide.md` and `hiring-process.md` (the round's
   official definition, if any), `docs/candidate/role-fit-analysis.md` (gaps
   this round exposes), `docs/company/culture-values.md` (the behavioral
   rubric), and every earlier round's debrief page (lessons to carry
   forward — format pivots, framing corrections, panel dynamics).

2. **Resolve the round:** the `key` argument, else ask. Keys are kebab-case
   round names, never stage numbers. If the row already has a Plan link,
   stop — extend the existing plan with `content-add-notes`. If there is no
   row yet, add one (key, name, duration, format, `_pending_`, `_pending_`).

3. **Panel:** names and titles come from the user (recruiter's mail) — never
   guess a panel. For each panelist, delegate research to
   `content-researcher` with the identity-verification brief from the Type
   section (a professional profile naming the employer, a code account tied
   to a personal site, commit history matching a listed employer — name
   similarity alone proves nothing; account for name variants; enumerate
   namesakes to exclude; a thin footprint is a documented absence). Then:
   ```sh
   .agentic-cms/scripts/ac-page new entity-panelist wiki/entities/<slug>.md --title "<Name>" --classification C2
   .agentic-cms/scripts/ac-index add entities wiki/entities/<slug>.md "<opaque one-liner: seat/round only>"
   .agentic-cms/scripts/ac-log append research "<key> panel — <slug> (opaque; page is C2)"
   ```
   Fill the page with the seat, verified vs. inferred facts, predicted lens,
   excluded namesakes; add the round key and specialty to `tags`.

4. **Plan page** — rate C2 (it embeds panel research), then:
   ```sh
   .agentic-cms/scripts/ac-page new round-plan docs/interview-prep/<key>-plan.md --title "<Round name> — Panel & Plan" --topic interview-prep --classification C2
   ```
   Keep only the sections that match the round's format (the template says
   which); delete the rest. `## Panel` links every panelist page with the
   lens they bring. Strategy and predicted questions with answer anchors
   come from the guide's own framing, earlier debriefs, the role-fit gaps,
   and the panelists' lenses — the rehearsal skill draws ~70% of its
   questions from here. Reverse Q&A prepared per seat. Watch-outs. Set
   `refs:` (guide, prior debriefs, role-fit, culture page, panelist pages)
   and add the round key to `tags`.

5. **Register, cross-link, verify:** put the Plan link into the Rounds
   table row; `ac-index add topics docs/interview-prep/<key>-plan.md
   "<opaque one-liner>"`; `ac-log append new-item "interview-prep/<key>-plan
   (opaque; page is C2)"`; add the plan to each panelist's "Mentioned in"
   and `refs:` (`ac-page touch` each edited page); then
   ```sh
   .agentic-cms/scripts/ac-index check && .agentic-cms/scripts/ac-links check && .agentic-cms/scripts/ac-classify check docs/interview-prep/<key>-plan.md wiki/entities/<slug>.md
   ```
   All must report `"clean": true`. Suggest `interview-rehearsal <key>` as
   the next step.

## Debrief mode

1. **Create the record** from the candidate's account of the round, as
   soon after it as possible (recall decays):
   ```sh
   .agentic-cms/scripts/ac-page new round-debrief docs/interview-prep/<key>-debrief.md --title "<Round name> — As-Run Debrief" --topic interview-prep --classification C2
   ```
   Format as run vs. predicted; questions recalled with the honest read and
   the fuller answer for any incomplete one; reverse-Q&A intel; signals and
   implications; performance read. `refs:` the plan page.

2. **Feed forward:** fill the Debrief link in the Rounds table; append the
   implications as "Lessons carried from previous rounds" to the next
   round's plan page if it exists, and to `role-fit-analysis.md` as a dated
   note if a gap was confirmed live (`ac-page touch` each); update the
   panelist pages with the observed dynamic (who carried, who observed);
   any technology the panel named gets its `wiki/concepts/` or
   `wiki/entities/` page (Type section, "Company stack in the wiki") linked
   from the debrief. Later recall goes into the debrief's Notes as dated
   entries, never by rewriting the body.

3. **Register and verify:** `ac-index add topics ... "<opaque one-liner>"`,
   `ac-log append new-item "interview-prep/<key>-debrief (opaque; page is
   C2)"`, plus a `notes` entry per page fed forward; then the same three
   checks as plan mode, all `"clean": true`.

## Rules

- The plan is the prediction; the debrief is the record. Never rewrite a
  plan to match what happened — record the pivot in the debrief.
- A real person's page is filed only after the verification chain closes;
  an unresolved identity is recorded as unresolved, not guessed.
- Everything this skill writes is C2: opaque index/log one-liners, no
  names, figures, or quoted content in the bookkeeping layer.
- Ratchet: raise freely, never lower; floor acks are the user's alone.
- No stage numbers anywhere — the Rounds table's keys are the names.
