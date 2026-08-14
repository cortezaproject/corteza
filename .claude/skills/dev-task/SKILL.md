---
name: dev-task
description: The task-level workflow — triage, write the task, scout, clarify, investigate, clarify again, implement, verify independently, present. Use when starting any piece of work that is not a one-liner. Calls /dev-change for the per-change mechanics and /intent-task when locked or WIP ground is involved.
---

# /dev-task

The procedure for taking a piece of work from a sentence to something landed.

Three rules govern everything below:

- **Questions are for decisions, not status.** Ask when two readings of the task
  would produce materially different work. Never ask "shall I proceed" — that is
  the human doing your job.
- **Nothing is true because a report said so.** Not a subagent's transcript, not
  a green test, not your own earlier reasoning. Claims get re-run.
- **A statement about what the software does is a claim, and a claim needs the
  file open.** Not the call site skimmed, not the function's name read, not the
  argument list inferred from.

## How to ask

Every question round in this skill is an interview, and a question is not an
interview just because it uses AskUserQuestion.

- **Explain the question and each option well enough to decide on.** An option
  is a sentence about what it means and what it costs, never a bare label.
- **Recommendation first**, marked as such, with the reason.
- **Order by what blocks the most.** The pressing decision leads; the cosmetic
  one can wait or go unasked.
- **Put the evidence in the question.** What was measured, what failed, what the
  investigation turned up.
- **End the turn with the questions.** If decisions are open, the last thing in
  the turn is the interview — not a summary that buries them.
- Batch the round: several related decisions together beats dripping one at a
  time.

## 0. Triage — pick a lane, say which

State the lane in one line and proceed. The human can override in one word.

| Lane            | When                                                                                      | Flow                                                                      |
| --------------- | ----------------------------------------------------------------------------------------- | ------------------------------------------------------------------------- |
| **trivial**     | Typo, rename, a one-liner whose cause is already known and whose blast radius is one file | → `/dev-change`, no questions                                             |
| **ordinary**    | A bug with a clear cause, or a small addition to an existing pattern                      | Scout → reproduce → **one** question round → implement → verify → present |
| **substantial** | New surface, unclear cause, several files, or anything touching a contract                | The full flow below                                                       |

Under-triaging is the dangerous direction: if scouting turns up a locked
contract or a second cause, re-triage upward and say so.

## 1. Write the task

Before investigating, write down three things. Two sentences each is plenty.

- **Outcome** — what is true when this is done, in the product's terms, not the
  code's.
- **Non-goals** — what this deliberately does not change.
- **Unknowns** — what you would need to know to be confident. These become the
  first question round.

If the human wrote the task, restate it in these three parts and let them
correct the restatement.

## 2. Scout — shallow, and bounded

Enough context to ask good questions. Not enough to form a plan.

A few files, the obvious entry points, whatever the human named. **Stop when you
can ask a sharper question than you could five minutes ago.**

**Read down the dependency tree: `server/` and `lib/` before
`client/web/unify`.** The higher layer is the single source of truth; the app
usually inherits what it already decided.

Use `Explore` or a read-only subagent for breadth; keep the depth for step 5.

## 3. Reproduce it — before you have a theory

**For a defect this is required, and it comes before the questions.**

Make it fail in front of you, at the layer the user hit it. Stand up whatever
state that needs; `dev/agent/api.sh` and `mcp.py` record what you create in
`dev/agent/.state/created.jsonl`, and step 9's cleanup deletes exactly that. The
ledger is the only thing that decides — do **not** prefix things `agent-`
expecting it to matter.

- **Match the instrument to the layer.** API and service behaviour → MCP or
  `dev/agent/api.sh`. Anything the user clicks → the browser. Logic in between →
  a unit test. The MCP tools never load the webapp, so an MCP-green result says
  nothing about a frontend defect.
- **`dev_server_status` first** — a reproduction against a stale binary looks
  like evidence and is not.
- **Failing to reproduce is a result, not a blocker.** Ask what the human was
  doing when it happened.

Keep the reproduction. Step 8 re-runs the same steps against the fix, in the
same environment.

## 4. Clarify round one — what are we building

Questions that change **what** gets built. Scope, behaviour, which of two
readings is meant, what "done" looks like in the product.

Do not ask about implementation here — it gets re-opened at step 6.

## 5. Investigate properly

Now go deep, on the _clarified_ task. Three things are mandatory:

1. **Contracts** — `dev_intent_governing` on the files you expect to touch. A
   locked contract in scope means **stop**: it is an intent change, the human
   rules on it, and the work moves to `/intent-task`. A WIP marker means the
   ground is undecided — ask before building on it.
2. **Baseline** — `dev_test_run` on the packages involved. Know what green looks
   like before you change anything; pre-existing breakage is not yours to debug.
3. **Prior art** — has this been solved, half-solved, or deliberately not solved
   nearby? Look for a `@todo`, a brief, or a comment explaining why not.

Fan out read-only subagents for breadth here. Give each a specific question, not
a topic.

## 6. Clarify round two — how, and how we will know

Questions that change **how** it gets built: approach, trade-offs, what to do
about anything surprising step 5 turned up.

**End this round by agreeing done-criteria, concretely.** Not "it works" —
_which_ test, _which_ live check, what the failure would look like if it were
still broken. Ask "how would I know if this were wrong" and answer it in the
criteria.

Write the criteria down in the reply. They are the contract for step 8.

## 7. Implement

Main context by default.

**Fan out to subagents when** the work splits into slices that touch disjoint
files and each slice is describable in a precise brief. Model: Opus 5 unless the
slice is a mechanical sweep.

**A brief is only good enough if it says what is already known to be wrong** —
specific enough for the agent to come back and falsify it.

**Consolidation is yours, and it is a named step, not a hope.** After every
fan-out, read the diffs together and factor what duplicated.

Fan out only on disjoint files. Anything needing judgement, or touching
locked/WIP ground, stays in the main context.

## 8. Verify and review — re-run the claims

Review means exercising the claims, not reading the diff.

- **Re-run the step-3 reproduction against the fix**, in the same environment.
- **Prove the environment is running the fix before you believe the re-run.**
  `dev_server_status` decides staleness from the running process's start time.
  gin only rebuilds when its proxy is hit (`curl -s localhost:3001/api/`);
  requests to the API port never trigger it.
- Run `/dev-change`'s verify sequence: tests, format, intent drift, affected
  specs.
- **Check against the step-6 criteria**, one at a time, out loud.
- **Prove the tests have teeth**: break what they cover, watch them fail, put it
  back.
- **Verify at the layer the bug lives in** — pick the instrument as in step 3.
- **Re-run anything a subagent claimed.** Their transcripts are evidence, not
  proof.
- For substantial work, consider a cold adversarial pass: a fresh agent, no
  context, told to break it.

Anything that fails here goes back to step 7 — or further, per the loop rules.

## 9. Present, and record what surprised you

Report what changed, what you verified and how, and what you deliberately did
not do. Failures and skipped work get stated plainly, with the output. The
report is where a task ends: the issue tracker is never yours to close, comment
on, or ask about (CLAUDE.md § Commit convention).

**Remove what you created.** `dev_fixture_cleanup` (or `dev/agent/cleanup.sh`)
deletes what this session's ledger records and nothing else — data the session
did not create is off-limits whatever it is called, and data it did create goes
whatever it is named.

Then **write down what surprised you**, in the repo, where the next session will
find it: a wrong assumption in a brief, an endpoint that is not what its name
suggests, a check that reported success without checking.

Finally, offer the next decision. If more questions are needed, loop:

- **New information changes the plan** → back to step 5.
- **The plan was wrong** → back to step 6.
- **The task was wrong** → back to step 1.

The loop is bounded by the step-6 criteria: when they are met, the work is done.
If they keep moving, that is the task changing — go back to step 1 openly rather
than absorbing it as another iteration.

## Committing

**The rule is CLAUDE.md § Commit convention** — authorship, message, atomicity,
intent docs riding with their code, and never pushing. Use `dev_commit_create`.

What this skill adds is what its two conditions mean in terms of these steps:

- **Verified** = step 8 actually happened: the reproduction re-run against the
  fix, the step-6 criteria checked one at a time, the touched package's tests
  green.
- **Nothing still open** = step 9's report would raise nothing you then act on:
  no question waiting on the human, no intent doc you noticed had drifted, no
  check you named and skipped. Those belong in the commit, not after it. If
  reconciling a doc needs the human's ruling, that is an open thread — get the
  ruling, then commit the lot.

One thing worth stating twice: work the human is still turning over — changing
the values, the shape, the scope — is not done however green it is. Let it
settle, then commit once.
