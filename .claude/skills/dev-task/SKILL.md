---
name: dev-task
description: The task-level workflow — triage, write the task, scout, clarify, investigate, clarify again, implement, verify independently, present. Use when starting any piece of work that is not a one-liner. Calls /dev-change for the per-change mechanics and /intent-task when locked or WIP ground is involved.
---

# /dev-task

The procedure for taking a piece of work from a sentence to something landed.

Two rules govern everything below:

- **Questions are for decisions, not status.** Ask when two readings of the task
  would produce materially different work. Never ask "shall I proceed" — that is
  the human doing your job.
- **Nothing is true because a report said so.** Not a subagent's transcript, not
  a green test, not your own earlier reasoning. Claims get re-run.
- **A statement about what the software does is a claim, and a claim needs the
  file open.** Not the call site skimmed, not the function's name read, not the
  argument list inferred from. Twice in one afternoon a defect was asserted from
  unread code and both were wrong: a toast helper was assumed to _replace_ the
  server's message when it _prefixes_ it — the proof was in the bug report's own
  text — and a project lock was blamed for a failure it had nothing to do with.
  Reading the file takes a minute; a wrong diagnosis costs an investigation.

## How to ask

Every question round in this skill is an interview, and a question is not an
interview just because it uses AskUserQuestion.

- **Explain the question and each option well enough to decide on.** An option
  is a sentence about what it means and what it costs, never a bare label. The
  reader should not have to ask what an option implies.
- **Recommendation first**, marked as such, with the reason.
- **Order by what blocks the most.** The pressing decision leads; the cosmetic
  one can wait or go unasked.
- **Put the evidence in the question.** What was measured, what failed, what the
  investigation turned up. A question that withholds the finding forces a guess.
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

## 3. Reproduce it — before you have a theory

**For a defect this is required, and it comes before the questions.**

Make it fail in front of you, at the layer the user hit it. Stand up whatever
state that needs, prefixed `agent-`, and note what you created so step 8 can
remove it.

- **Match the instrument to the layer.** The MCP tools reach the REST API and
  never load the webapp: an MCP-green result says nothing about a frontend
  defect. API and service behaviour → MCP or `dev/agent/api.sh`. Anything the
  user clicks → the browser. Logic in between → a unit test.
- **`dev_server_status` first** — a reproduction against a stale binary is
  worse than none, because it looks like evidence.
- **Failing to reproduce is a result, not a blocker.** It means the first
  question is "what were you doing when this happened", and asking that beats
  spending a deep investigation on a guess.

This step exists because it was missing. Two defects in one afternoon were
diagnosed from reading and both diagnoses were wrong — one blamed a project lock
for a failure it had nothing to do with, and a live reproduction would have
killed it in a minute, since it would have needed a project that did not exist.

Keep the reproduction. Step 7 re-runs the same steps against the fix, in the
same environment — that is the only thing that closes the loop.

## 4. Clarify round one — what are we building

Questions that change **what** gets built. Scope, behaviour, which of two
readings is meant, what "done" looks like in the product.

Do not ask about implementation here; you have not investigated enough to have
an opinion worth acting on, and asking anyway spends the human's attention on a
decision that will be re-opened at step 5.

## 5. Investigate properly

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

## 6. Clarify round two — how, and how we will know

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

Write the criteria down in the reply. They are the contract for step 8.

## 7. Implement

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

## 8. Verify and review — re-run the claims

Review means exercising the claims, not reading the diff.

- **Re-run the step-3 reproduction against the fix**, in the same environment.
  A fix that was never shown to stop the original failure is a hypothesis.
- **Prove the environment is running the fix before you believe the re-run.**
  This is where a server change gets a false verdict, and the false verdict is
  a _red_ one: the reproduction still fails, which looks exactly like a fix that
  did not work, and the honest response — go rewrite it — makes a correct fix
  worse. `dev_server_status` now decides staleness from the running process's
  start time, so trust it again here; a green unit test says nothing about which
  binary is answering. gin only rebuilds when its proxy is hit
  (`curl -s localhost:3001/api/`) — requests to the API port never trigger it,
  so a server that is never proxied never picks your change up at all.
- Run `/dev-change`'s verify sequence: tests, format, intent drift, affected
  specs.
- **Check against the step-5 criteria**, one at a time, out loud.
- **Prove the tests have teeth**: break what they cover, watch them fail, put it
  back. A green suite says nothing about a test that does not assert — that
  exact case is why this step exists.
- **Verify at the layer the bug lives in.** The MCP tools reach the REST API
  directly and never load the webapp, so an MCP-green result says nothing about
  a frontend defect — issue #30 deleted a module happily through MCP while the
  edit view could not, because the view never sent the request. Ask which layer
  is on trial before choosing the instrument: MCP for service and API behaviour,
  a browser for anything the user clicks, a unit test for logic in between.
- **Re-run anything a subagent claimed.** Their transcripts are evidence, not
  proof; an agent's live exercise has been right about the result and wrong
  about the reason.
- For substantial work, consider a cold adversarial pass: a fresh agent, no
  context, told to break it. It has been worth more than any amount of reading.

Anything that fails here goes back to step 7 — or further, per the loop rules.

## 9. Present, and record what surprised you

Report what changed, what you verified and how, and what you deliberately did
not do. Failures and skipped work get stated plainly, with the output.

**Remove what you created.** Everything `agent-` prefixed from step 3 goes, and
unprefixed data is never yours to touch. A leftover namespace sat on the dev
server for an afternoon because this was nobody's job.

Then **write down what surprised you**, in the repo, where the next session will
find it: a wrong assumption in a brief, an endpoint that is not what its name
suggests, a check that reported success without checking. Each of these cost
real time once; recording them is minutes and stops the next session paying
again.

Finally, offer the next decision. If more questions are needed, loop:

- **New information changes the plan** → back to step 5.
- **The plan was wrong** → back to step 6.
- **The task was wrong** → back to step 1.

The loop is bounded by the step-5 criteria: when they are met, the work is done.
If they keep moving, that is the task changing, and it should go back to step 1
openly rather than being absorbed as another iteration.

## Committing

Once the work is done and verified, commit it. Do not leave ready changes
sitting in the tree waiting to be asked for — the human should not have to
chase a finished change into history. Use `dev_commit_create`; see
`/dev-change` for the message and atomicity rules.

"Verified" means step 8 actually happened: the reproduction re-run against the
fix, the step-5 criteria checked one at a time, the touched package's tests
green. Work that has not cleared that bar is not done, and a commit is not what
makes it so.

Three things this does not license:

- **Pushing.** A local commit is reversible; publishing is not. Pushes, PRs and
  anything else outward-facing stay the human's call, every time.
- **Sweeping the tree.** Name the files belonging to the change in `files`.
  Another session's dirty file, or the human's own half-finished edit, is not
  yours to commit — that has happened, and it buries someone else's work in
  your history.
- **Committing through a conversation.** Work the human is still turning over —
  changing the values, the shape, the scope — is not done however green it is.
  Let it settle, then commit once, rather than once per turn.
