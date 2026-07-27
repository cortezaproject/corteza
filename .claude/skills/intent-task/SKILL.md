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

   When the Agent tool is available, do 1–2 with parallel read-only scouts
   (model: sonnet), launched in one message: one locates the touched files,
   one gathers the governing docs. Scouts return paths, verbatim doc excerpts
   (locked-contract / WIP / DRIFT blocks in full), and `file:line` evidence —
   never classifications or prose paraphrases of contracts. Without the Agent
   tool, do it inline as before.

3. Classify the task against the docs — always the main agent's job, from the
   evidence, never delegated:
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

1. If intent changes: edit the intent docs FIRST (they are the design) — in
   the main context, never via subagent.
2. Implement the code change. Main context is the default. For purely
   mechanical multi-file sweeps (renames, i18n strings, boilerplate) over
   disjoint files you MAY fan out parallel implementation agents
   (model: sonnet), one per independent slice; their prompts must name the
   governing intent docs as constraints and end with prettier/gofmt on the
   files they changed. Anything needing judgment — or touching locked/WIP
   territory — stays in the main context.
3. Run `node .intent/intent.mjs check --changed` as you edit to see which
   governing docs need reconciling — reconcile as you go, not in a batch at
   the end. Edits made by subagents are yours to reconcile: review their
   diffs against the governing docs before Phase 4.

## Phase 4 — Verify & close

1. prettier (FE) / gofmt (Go) on changed files only.
2. Run the touched package's unit tests. Do NOT run e2e specs — they are slow
   and human-run for now. Instead, list the affected specs in the report:
   `node .intent/intent.mjs affected <changed files>` prints them; the human
   runs `npx playwright test <specs>` themselves.

   Parallelize with background Bash, not agents: launch unit tests as
   background commands and run the format→sync→check chain (steps 1 and 3,
   which must stay sequential with each other) while they run. Read every
   exit code yourself; a background failure is still a failure.

3. `node .intent/intent.mjs sync <changed files>` then
   `node .intent/intent.mjs check` — must be green (never pipe through
   tail/head; read exit codes). Sync only after formatting — prettier changes
   hashes.
4. Commit per repo conventions (atomic; docs/bugfix/cleanup separate; no AI
   trailers). Update `.intent/TODO.md` in the same commit when the work
   resolves or creates an item.
5. Report: outcome, doc changes, any new WIP/DRIFT markers, open questions.
