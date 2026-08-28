# How a turn reports

The one statement of the reply format for `/dev-task` and `/dev-change`. Both
skills point here; neither restates it.

Every turn ends in this block — mid-task turns as well as the final one. The
block is the reply; it does not get wrapped in prose above or below.

## The block

A status line, then eight sections in this order, **every one of them, every
time**. A section with nothing in it gets a single `—`. Predictable position is
the point: an empty `⚠️ NOT DONE` is information, a missing one is a question.

```
ordinary · step 8/9 · verifying

🎯 TASK
Restore control for deleted resources

🔍 FOUND
- delete is soft everywhere; only roles had a Restore control — RoleList.vue:88

🔧 CHANGED  3 files +48 −12
- CResourceList.vue:112 — restore action
- locale/en/human-webapp/list.yaml — new key

✅ VERIFIED
- repro re-run in browser: deleted row returns to the list
- vitest unify — 214 passed; teeth proven (reverted the guard → 3 fail)

⚠️ NOT DONE
- archived roles — non-goal

📊 CONFIDENCE  High
reproduced at the user's layer before and after; tests proven to bite

❓ OPEN
- —

➡️ NEXT
commit the three files — asked below

🔗 wt-alpha · http://localhost:5176 · ~/Human/human-worktrees/wt-alpha · not merged (branch wt-alpha)
```

### Status line

`lane · step N/total · what is running right now`. `/dev-task` counts steps out
of 9, `/dev-change` out of 7. It is one line and it is never omitted: a
mid-task turn has to say where in the flow it is without being read.

### The footer

**One line, last, every time: where the human can go and look.** They should
never have to ask which port this work is on, or reconstruct it from a slot
number.

```
🔗 <checkout> · <webapp URL> · <path> · <landed?>
```

- **The checkout this turn's work lives in**, whether that is a worktree or the
  primary. Add a second line only for another checkout the turn actually
  touched — never a listing of every worktree on the machine.
- **Say where it stands against main, every time.** A worktree is
  `merged to main` or `not merged (branch <name>)`; the primary is `on main` or
  `on <branch>`. Never leave the human to guess whether the work has landed.
- **Say when it is not serving.** A URL that answers nothing is worse than no
  URL: `http://localhost:5176 (down — 'worktree.sh up wt-alpha')`.
- **The path is there to be diffed**, so give the real one:
  `git -C <path> diff $(git -C <path> merge-base main HEAD)` shows a worktree's
  own work, committed and uncommitted, however far main has moved since.
  `git diff main...<branch>` shows nothing until it commits, and a plain
  `diff main` reports main's newer commits as deletions.
- Add the API URL only when the turn's own verification went through it.

`dev/agent/worktree.sh list` has the ports; the URL is `http://localhost:<vite>`.

### What each section holds

| Section       | Holds                                                           |
| ------------- | --------------------------------------------------------------- |
| 🎯 TASK       | What this turn was for, in the product's terms. One line.       |
| 🔍 FOUND      | What the investigation **established**. Not what was searched.  |
| 🔧 CHANGED    | Files touched, with `+n −n` on the heading. Nothing else.       |
| ✅ VERIFIED   | What was **run** and what it said. Never "should work".         |
| ⚠️ NOT DONE   | Skipped, deferred, failed, and declared non-goals. Plainly.     |
| 📊 CONFIDENCE | The grade and the rubric line that produced it. See below.      |
| ❓ OPEN       | Anything unresolved. Every entry becomes an interview question. |
| ➡️ NEXT       | The one decision now due — asked, not narrated.                 |

## Length

**Five bullets per section, one line each, no paragraphs anywhere.**

The cap is a real limit, not a target to squeeze under. A section that wants a
sixth bullet is a turn that should have been two turns, or a body of detail
that belongs in a file rather than the reply. Say which, and move on.

The section headings do the explaining. A sentence introducing a section, a
sentence summarising it afterwards, and a closing paragraph restating the block
are all the wall of text growing back.

## The queue

**Anything in `⚠️ NOT DONE` that should outlive the turn goes to the backlog as
you write it** — `dev/agent/backlog.sh add "…" --why … --files … --task …`.
Deferring without queueing is how a thing gets forgotten, which is the whole
reason the queue exists.

Which tier is not a judgement call, it is a test:

- **Local** — doing it would change a file this task is already changing, or it
  follows directly from this change. This is the default, and it is the only
  thing read back to you.
- **`--global`** — everything else, including the second, real problem you found
  while looking at something else. Filed, named here, and out of the way.

Not everything: a declared non-goal is not a todo, and neither is something
ruled against. Queue what someone would want done later.

The local todo is read back at exactly one moment — `➡️ NEXT`, where the open
items become interview options once the current task is fully finished. The
global pool is read when the human asks for it and at no other time.

## Confidence

One grade for the turn, chosen by this rubric — not by how the work feels.

| Grade      | Means                                                                                                                                 |
| ---------- | ------------------------------------------------------------------------------------------------------------------------------------- |
| **High**   | Reproduced before **and** re-run after, at the layer the user hits it; the touched package's tests are green and were proven to bite. |
| **Medium** | Run and green, but not exercised where the user hits it — or some path in the change never executed.                                  |
| **Low**    | Read and reasoned about; not run. Or run against something that might be stale.                                                       |

Rules on the grade:

- **It grades this turn's claims**, not the quality of the code.
- **Name what would raise it**, in the same breath — "Medium: the archived-role
  path never ran; running it makes this High." A grade with no route up is a
  mood.
- **Never High with anything in ❓ OPEN or ⚠️ NOT DONE that could falsify a
  ✅ VERIFIED bullet.** Contradiction outranks the rubric.
- **Downgrade on** — a subagent's claim not re-run; `dev_server_status` not
  checked before believing a live result; a test that passed first time and was
  never broken to prove it asserts.
- **A subagent's grade is not yours.** Take it one step down until you have
  re-run the claim yourself.
- **The grade decides whether work lands unasked**, so grade the evidence and
  not the mood. High in a worktree merges to main without a question; every
  reason to round a Medium up is a reason to look harder at what is missing.

## The interview

**Anything needed from the human is an `AskUserQuestion`, always.** Never a
question in the prose, never "let me know", never "shall I proceed". The human
should be able to answer the turn without reading it.

- **The test is whether the human could usefully weigh in — not whether you are
  stuck.** A decision, a preference, guidance, a choice among the moves now
  available: all of those are interviews. It is a low bar and it is deliberately
  not "shall I proceed".
- **When nothing is genuinely open, do not manufacture a question.** `➡️ NEXT`
  names the next move in one line and the turn ends there.
- **Every ❓ OPEN entry becomes a question.** If there are more than four, ask
  the four that block the most and say in `➡️ NEXT` what is being held back.
- **Recommendation first**, marked as such, with the reason and the cost.
- **An option is a sentence**, saying what it means and what it costs. A bare
  label is not an option.
- **Put the evidence in the question** — what was measured, what failed, what
  step 5 turned up. The question stands on its own.
- **Order by what blocks the most**, and batch related decisions into one round.
