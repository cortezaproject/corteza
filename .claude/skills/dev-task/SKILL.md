---
name: dev-task
description: The task-level workflow — triage, write the task, scout, clarify, investigate, clarify again, implement, verify independently, present. Use when starting any piece of work that is not a one-liner. Calls /dev-change for the per-change mechanics and /intent-task when locked or WIP ground is involved.
---

# /dev-task

The procedure for taking a piece of work from a sentence to something landed.

Two rules govern everything below:

- **Questions are for decisions, not status.** Ask when two readings of the task
  would produce materially different work. Never ask "shall I proceed" — that is
  the human doing your job. Use AskUserQuestion with concrete options,
  recommendation first.
- **Nothing is true because a report said so.** Not a subagent's transcript, not
  a green test, not your own earlier reasoning. Claims get re-run.

## 0. Triage — pick a lane, say which

State the lane in one line and proceed. The human can override in one word.

| Lane            | When                                                                                      | Flow                                                          |
| --------------- | ----------------------------------------------------------------------------------------- | ------------------------------------------------------------- |
| **trivial**     | Typo, rename, a one-liner whose cause is already known and whose blast radius is one file | → `/dev-change`, no questions                                 |
| **ordinary**    | A bug with a clear cause, or a small addition to an existing pattern                      | Scout → **one** question round → implement → verify → present |
| **substantial** | New surface, unclear cause, several files, or anything touching a contract                | The full flow below                                           |

Getting this wrong in the cheap direction is the dangerous one: if scouting
turns up a locked contract or a second cause, re-triage upward and say so.

## 1. Write the task

Before investigating, write down three things. Two sentences each is plenty.

- **Outcome** — what is true when this is done, in the product's terms, not the
  code's.
- **Non-goals** — what this deliberately does not change. This is what stops
  scope creep later, and it is cheapest to write now.
- **Unknowns** — what you would need to know to be confident. These become the
  first question round.

If the human wrote the task, restate it in these three parts and let them
correct the restatement. A misread task is the most expensive error available,
and it is free to catch here.

## 2. Scout — shallow, and bounded

Enough context to ask good questions. Not enough to form a plan.

A few files, the obvious entry points, whatever the human named. **Stop when you
can ask a sharper question than you could five minutes ago** — that is the whole
purpose of this step. Reading the whole subsystem here means asking questions
you have already answered, badly.

Use `Explore` or a read-only subagent for breadth; keep the depth for step 4.

## 3. Clarify round one — what are we building

Questions that change **what** gets built. Scope, behaviour, which of two
readings is meant, what "done" looks like in the product.

Do not ask about implementation here; you have not investigated enough to have
an opinion worth acting on, and asking anyway spends the human's attention on a
decision that will be re-opened at step 5.

## 4. Investigate properly

Now go deep, on the _clarified_ task.

Three things are mandatory, because each one changes the shape of the work
rather than its content:

1. **Contracts** — `dev_intent_governing` on the files you expect to touch. A
   locked contract in scope means **stop**: it is an intent change, the human
   rules on it, and the work moves to `/intent-task`. A WIP marker means the
   ground is undecided — ask before building on it.
2. **Baseline** — `dev_test_run` on the packages involved. Know what green looks
   like before you change anything. This repo has pre-existing breakage, and
   inheriting someone else's failure costs an hour.
3. **Prior art** — has this been solved, half-solved, or deliberately not solved
   nearby? A `@todo`, a brief, or a comment explaining why not is worth more
   than any amount of fresh reasoning.

Fan out read-only subagents for breadth here. Give each a specific question, not
a topic.

## 5. Clarify round two — how, and how we will know

Questions that change **how** it gets built: approach, trade-offs, what to do
about anything surprising step 4 turned up.

**End this round by agreeing done-criteria, concretely.** Not "it works" —
_which_ test, _which_ live check, what the failure would look like if it were
still broken. This is the single highest-value thing in the whole flow, because:

- It is decided before any code exists, so it cannot be quietly bent to fit what
  was built.
- Step 7 checks against it rather than against your own judgement.
- It forces the question "how would I know if this were wrong", which is the
  question that catches the failures nothing else catches.

Write the criteria down in the reply. They are the contract for step 7.

## 6. Implement

Main context by default.

**Fan out to subagents when** the work splits into slices that touch disjoint
files and each slice is describable in a precise brief. Model: Opus 5 unless the
slice is a mechanical sweep.

**A brief is only good enough if it says what is already known to be wrong.**
Briefs written this way have had agents come back and correct them — which is
the point, and only possible when the brief is specific enough to be falsified.

**Consolidation is yours, and it is a named step, not a hope.** Independent
agents solve shared problems independently: two agents on separate files wrote
the same JSON-decoding helper twice, in the same package, on the same afternoon.
After every fan-out, read the diffs together and factor what duplicated.

Fan out only on disjoint files. Anything needing judgement, or touching
locked/WIP ground, stays in the main context.

## 7. Verify and review — re-run the claims

Review means exercising the claims, not reading the diff.

- Run `/dev-change`'s verify sequence: tests, format, intent drift, affected
  specs.
- **Check against the step-5 criteria**, one at a time, out loud.
- **Prove the tests have teeth**: break what they cover, watch them fail, put it
  back. A green suite says nothing about a test that does not assert — that
  exact case is why this step exists.
- **Re-run anything a subagent claimed.** Their transcripts are evidence, not
  proof; an agent's live exercise has been right about the result and wrong
  about the reason.
- For substantial work, consider a cold adversarial pass: a fresh agent, no
  context, told to break it. It has been worth more than any amount of reading.

Anything that fails here goes back to step 6 — or further, per the loop rules.

## 8. Present, and record what surprised you

Report what changed, what you verified and how, and what you deliberately did
not do. Failures and skipped work get stated plainly, with the output.

Then **write down what surprised you**, in the repo, where the next session will
find it: a wrong assumption in a brief, an endpoint that is not what its name
suggests, a check that reported success without checking. Each of these cost
real time once; recording them is minutes and stops the next session paying
again.

Finally, offer the next decision. If more questions are needed, loop:

- **New information changes the plan** → back to step 4.
- **The plan was wrong** → back to step 5.
- **The task was wrong** → back to step 1.

The loop is bounded by the step-5 criteria: when they are met, the work is done.
If they keep moving, that is the task changing, and it should go back to step 1
openly rather than being absorbed as another iteration.

## Committing

Never autonomously. Commit when the human asks, once the work is confirmed
working, through `dev_commit_create` — see `/dev-change` for the message and
atomicity rules.
