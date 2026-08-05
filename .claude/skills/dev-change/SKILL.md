---
name: dev-change
description: The mechanics of making one change — orient, check contracts, baseline, edit, verify, prove the test has teeth, commit — driven by the human-dev MCP tools. This is the inner loop /dev-task calls once a plan is agreed; invoke it directly for a change small enough to need no planning. For anything touching a locked contract or a WIP zone, use /intent-task.
---

# /dev-change

The `human-dev` MCP server provides the primitives; this skill is the order to
use them in. If a tool below is not available, the server is not connected —
tell the human to restart their MCP client, and fall back to the shell
equivalents (`git status`, `go test`, `gofmt`).

## 1. Orient

`dev_branch_status` — branch, base, how far ahead, what is already dirty.

A tree with unrelated changes already in it is worth knowing about _before_ you
edit, because it decides whether you can commit with staged files or must name
paths explicitly to keep the commit atomic.

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

If it is already failing, you need to know now — this repo has pre-existing
breakage (`federation/service` does not build; `lib/js` rollup is broken), and
half an hour spent debugging someone else's failure is the cost of skipping
this.

## 4. Make the change

Ordinary Read/Edit. The tools do not do the thinking.

Scope discipline: do what was asked. If you find a second, real problem, finish
the first, then say what you found — do not fold it in silently.

## 5. Verify

In this order, because each one can invalidate the next:

1. `dev_test_run` — the touched package. Returns failures only; a green run is
   one line.
2. `dev_format_run` — no arguments formats what you changed. **Never format a
   directory**: this repo has files that were already unformatted before you
   arrived, and reformatting them has turned a small diff into an unreviewable
   one twice.
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

This exists because of a real case in this repo: `inputguard`'s
`TestKnownFalsePositives` used `t.Logf` when an input stopped being blocked, so
removing three detection patterns left the suite green. `dev_test_run` reported
`passed: true`, correctly and meaninglessly. **A green suite proves nothing
about a test that does not assert.** No tool can tell you this; you have to
break something and watch.

The same shape shows up everywhere in this codebase's history: a check that
reports success without having checked. Treat "it passed first time" as a
question, not an answer.

## 7. Commit — only when the human asked for one

`dev_commit_create` enforces how a commit is made. It does not decide whether
one should happen: that is the human's call, every time.

- One imperative line, under 72 characters, capitalised, no trailing period.
- No AI attribution anywhere in the message — no co-author trailer, no
  generated-with line, no robot emoji. The tool refuses all three.
- Pass `files` explicitly to keep a commit atomic when the tree holds unrelated
  work; omit it to commit what is already staged.
- Formatting runs before staging, so what lands is what was formatted.
- Mixing docs with code produces a **warning**, not a refusal. Read it: split
  the commit unless the pieces genuinely belong together.

Separate atomic commits — a product fix and a tooling fix are two commits, even
when you made them in the same sitting.

## What this skill is not

It is the inner loop, not the whole procedure. It has no clarification rounds
and no planning: it assumes you already know what to build. For work that needs
deciding first, `/dev-task` wraps this with the question rounds and calls it at
step 6. For a change on locked or WIP ground, use `/intent-task`.
