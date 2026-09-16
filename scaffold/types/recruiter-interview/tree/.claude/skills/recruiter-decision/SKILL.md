---
name: recruiter-decision
description: Run a hiring-committee decision point — compare every candidate at a stage per competency across their scorecards, record who advances, is held, rejected, or offered, with the evidence; or make an offer to one candidate (level, terms, approvals, negotiation log, outcome) and close the roster row. Use when the user says "committee for <round>", "who advances", "decide on <slug>", "make an offer to <slug>", "<slug> accepted/declined", or invokes /recruiter-decision compare <key> | offer <slug>.
license: MIT
compatibility: recruiter-interview type (CONTENT.md, Type section); requires scorecards in docs/candidates/ and the competency framework; uses the .agentic-cms/scripts toolkit and the decision-record and candidate-offer templates
allowed-tools: Read Write Edit Grep Glob Bash(.agentic-cms/scripts/*)
---

# recruiter-decision — compare and decide, then offer

Read `CONTENT.md`'s Type section first if you haven't this session. All
`ac-*` commands live in `.agentic-cms/scripts/`, return JSON — check `"ok"`
after every call. Everything here is **C2** and refers to candidates by
**slug only**.

## Compare mode (`compare <key>`)

1. **Load:** the Roster (which slugs are at stage `<key>` or `decision`
   after it), every scorecard of those candidates
   (`docs/candidates/cNNN-*-scorecard.md`), their screening notes and
   references if any, the competency framework (weights, must-haves), and
   the hiring process's decision rules (quorum, tie-break).

2. **Build the comparison:** one row per slug, one column per competency
   assessed so far — the score from the relevant scorecard, "—" for no
   evidence — plus each candidate's conferred recommendation. Weighted
   totals are a check, not the decision; must-have competencies below the
   bar are disqualifying regardless of total.

3. **Record the decision:**
   ```sh
   .agentic-cms/scripts/ac-page new decision-record docs/decisions/<date>-<key>.md --title "Decision — <Round name>" --topic decisions --classification C2
   ```
   (The date in the filename is the exception the Type section allows for
   decision points: `YYYY-MM-DD-<key>.md` orders the committee's records.)
   Attendees (interviewer entities), the table, one subsection per slug —
   advance (to which round), hold (until what), reject (why, in
   evidence terms), offer (at which level) — dissent recorded, and the
   process observations (a competency no round evidenced, divergent
   scoring → `recruiter-round calibrate`, a question that never
   discriminates). `refs:` every scorecard compared, the framework; add the
   stage to `tags`.

4. **Apply and bookkeep:** update each slug's roster stage/status and hub
   (`ac-page touch`, re-rate C2); run `recruiter-candidate archive <slug>
   rejected` for rejections; then
   ```sh
   .agentic-cms/scripts/ac-index add topics docs/decisions/<date>-<key>.md "decision — <key> (C2)"
   .agentic-cms/scripts/ac-log append new-item "decisions/<date>-<key> (opaque; page is C2)"
   .agentic-cms/scripts/ac-index check && .agentic-cms/scripts/ac-links check && .agentic-cms/scripts/ac-classify check docs/decisions/<date>-<key>.md
   ```
   All must report `"clean": true`. Tell the user which candidate
   communications are now due (the process page's commitments).

## Offer mode (`offer <slug>`)

1. **Load:** the decision record that authorized the offer, the candidate's
   hub and screening notes (expectations), the job spec's level and band,
   the references outcome.

2. **Create the offer page:**
   ```sh
   .agentic-cms/scripts/ac-page new candidate-offer docs/candidates/cNNN-offer.md --title "cNNN — Offer" --topic candidates --classification C2
   ```
   Decision basis (link the record; level and why), terms as approved with
   the approval chain, the negotiation log as dated entries, outcome
   pending. A verbatim offer letter is C3 and stays out of this repository
   — reference where it lives. `refs:` the decision record, the hub; add
   the slug to `tags`.

3. **Bookkeep:** roster stage `offer`; hub Relations; index (opaque), log,
   the three checks clean.

4. **Close:** when the candidate answers, fill the outcome — accepted
   (start date; roster status `hired`, requisition state on the role
   entity) or declined (reason if disclosed; roster `declined`) — `ac-page
   classify ... C2` to re-stamp, `ac-log append notes`. If declined and the
   committee held a runner-up, point the user back to compare mode.

## Rules

- Decisions cite evidence from scorecards, never impressions that aren't
  on a scorecard; if the evidence isn't filed, file it first.
- Must-have competencies are gates; weighted totals never override them.
- Slug only, everywhere outside the candidate's own pages.
- Figures make every offer entry C2; the signed document is C3 and
  external.
- Ratchet: raise freely, never lower; floor acks are the user's alone.
