---
name: dev-change
description: The mechanics of making one change — orient, check contracts, baseline, edit, verify, prove the test has teeth, commit — driven by the human-dev MCP tools. This is the inner loop /dev-task calls once a plan is agreed; invoke it directly for a change small enough to need no planning. For anything touching a locked contract or a WIP zone, use /intent-task.
---

# /dev-change

The `human-dev` MCP server provides the primitives; this skill is the order to
use them in. If a tool below is not available, the server is not connected —
tell the human to restart their MCP client, and fall back to the shell
equivalents (`git status`, `go test`, `gofmt`).

## How to report

**The format is `.claude/reporting.md`** — read it before your first report and
follow it verbatim: the status line, the eight sections in order, the five-bullet
cap, the confidence rubric, and the rule that anything needed from the human is
an `AskUserQuestion` rather than a sentence. It holds whether this skill was
invoked directly or called by `/dev-task`.

Count the status line's steps out of 7: `dev-change · step 3/7 · baseline`.

What this skill adds to it:

- **`🔍 FOUND` carries step 2's answer verbatim** when a contract or WIP marker
  is in scope. A contract restated in your own words has stopped being the
  contract, and paraphrasing it into a bullet is how it stops being one.
- **A change whose test was never broken to prove it asserts is Medium**, per
  step 6, however green the suite is.
- **`⚠️ NOT DONE` holds the second problem** you found and did not fold in —
  and it goes to `backlog.sh` at the same time, or it is forgotten.
- **The footer names the checkout and its URL** — after a change, that is where
  the human goes to look at it.
- **A change that needs the server** belongs in its own worktree
  (`dev/agent/worktree.sh new`), not in a checkout another session is editing.

## 1. Orient

`dev_branch_status` — branch, base, how far ahead, what is already dirty.

Know this _before_ you edit: a tree holding unrelated changes decides whether
you can commit staged files or must name paths explicitly to keep the commit
atomic.

## 2. Read the contracts before editing, not after

`dev_intent_governing` with the files you intend to change.

Read what it returns; do not skim it. It gives locked contracts and WIP markers
**verbatim** because a contract restated in your own words has stopped being the
contract.

- **A locked contract in scope** → stop. That is an intent change, not a code
  change, and the human rules on it first. Switch to `/intent-task`.
- **A WIP marker in scope** → ask before building on it. The ground is
  undecided by definition.
- `covers: governs` binds the file directly. `covers: area-context` is the
  area's constitution — it does not formally cover the path, but its contracts
  still hold.
- "No intent doc governs this path" is a real answer, not a failure. Much of
  `server/` is not yet enforced.

## 3. Baseline before you touch anything

`dev_test_run` on the package you are about to change.

If it is already failing, you need to know now: pre-existing breakage is not
yours to debug, and you cannot tell it from your own once you have edited.

## 4. Make the change

Ordinary Read/Edit. The tools do not do the thinking.

Edit at the layer that owns the concept: `server/` and `lib/` are the source of
truth and the app imports them one-way. Check upstream for an existing rule
before writing one in `client/web/unify` — otherwise you ship a second copy that
drifts from the first.

Scope discipline: do what was asked. If you find a second, real problem, finish
the first, then queue it — `dev/agent/backlog.sh add "…" --why … --files …` —
and name it in `⚠️ NOT DONE`. Do not fold it in silently, and do not leave it
as a sentence in a report nobody re-reads.

## 5. Verify

In this order, because each one can invalidate the next:

1. `dev_test_run` — the touched package. Returns failures only; a green run is
   one line.
2. `dev_format_run` — no arguments formats what you changed. **Never format a
   directory**: it sweeps up files that were already unformatted and turns a
   small diff into an unreviewable one.
3. `dev_intent_check` — reports your drift separately from the repo's
   pre-existing baseline. Reconcile the docs listed under `yours`. Leave the
   baseline alone.
4. `dev_intent_affected` — e2e specs go to the human to run
   (`npx playwright test <specs>`); never run them yourself, they are slow.
   Anything under `unitTests` you run with `dev_test_run`.

## 6. Prove the test has teeth

**Do not skip this. It is the step that catches the failure nobody else does.**

If you added or changed a test, break the thing it covers — revert your fix,
flip a condition — and confirm the suite fails, then put it back.

**A green suite proves nothing about a test that does not assert** — a test that
logs instead of failing reports `passed: true`, correctly and meaninglessly. No
tool can tell you this; you have to break something and watch. Treat "it passed
first time" as a question, not an answer.

## 7. Commit — once the work is verified

**The rule is CLAUDE.md § Commit convention.** Read it there; it is not
restated here, and it is not what `dev_commit_create` decides — the tool
enforces the mechanical half (message shape, no AI attribution, formatting
before staging) and warns on the rest.

What this skill adds is only what "verified and nothing open" means here:
steps 4–6 actually ran — tests green, the change exercised at its own layer,
the test proven to have teeth — and nothing you turned up in them is still
waiting. If the human asks for a commit anyway, say plainly what is still open
in `⚠️ NOT DONE`, then do as they asked.

## What this skill is not

It is the inner loop, not the whole procedure. It has no clarification rounds
and no planning: it assumes you already know what to build. For work that needs
deciding first, `/dev-task` wraps this with the question rounds and calls it at
step 6. For a change on locked or WIP ground, use `/intent-task`.
