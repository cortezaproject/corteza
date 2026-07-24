---
name: intent-task
description: Standard workflow for any code change in this repo — locate governing intent docs, confirm a plan with the human, implement intent-first, reconcile docs, sync, test, commit. Use for features, fixes, and refactors on intent-covered code.
---

# /intent-task <task description>

The standardized change workflow. Follow the phases in order; do not skip the
confirmation phase — every task gets one, that is the team contract.

## Phase 1 — Orient

1. Identify the files the task will touch (search, don't guess).
2. Read their governing intent docs: sibling `<name>.intent.md`, else nearest
   folder doc; ALWAYS also the area constitution if one exists (e.g.
   `sections/project/project.intent.md`).
3. Classify the task against the docs:
   - Touches a **locked contract** → the plan must call it out; the human must
     explicitly rule on the change (it is an intent change, not a code change).
   - Touches a **WIP zone** → interview first; do not build on unruled ground.
   - Contradicts a doc → STOP and report; never silently "fix" either side.
   - Plain change within recorded intent → note which docs stay unchanged.

## Phase 2 — Confirm (mandatory, no exceptions)

Present via AskUserQuestion: what will change (files + behavior), which intent
docs are affected (updated vs confirmed-unchanged), whether any locked/WIP
territory is involved, test plan. Concrete options, recommendation first.
Only proceed on approval.

## Phase 3 — Implement (intent-first)

1. If intent changes: edit the intent docs FIRST (they are the design).
2. Implement the code change.
3. PostToolUse hooks will name governing docs on every edit — reconcile as you
   go, not in a batch at the end.

## Phase 4 — Verify & close

1. prettier (FE) / gofmt (Go) on changed files only.
2. Run the touched package's unit tests, plus the e2e specs mapped to the
   change: `node .intent/intent.mjs affected <changed files>` prints them;
   run via `cd client/web/unify && npx playwright test <specs>` (requires the
   dev stack running + `.env.e2e`; if unavailable, say so in the report
   instead of skipping silently). Report results honestly.
3. `node .intent/intent.mjs sync <changed files>` then
   `node .intent/intent.mjs check` — must be green (never pipe through
   tail/head; read exit codes).
4. Commit per repo conventions (atomic; docs/bugfix/cleanup separate; no AI
   trailers). Update `.intent/TODO.md` in the same commit when the work
   resolves or creates an item.
5. Report: outcome, doc changes, any new WIP/DRIFT markers, open questions.
