---
name: orchestrate
description: Take a list of issues — pasted, backlog ids, or the whole open queue — and work them in parallel lanes, each in its own worktree, batching every lane's questions into one interview and landing the results in a chosen order. Use when there are several independent pieces of work; use /dev-task for one.
---

# /orchestrate

Several issues, worked at once, without them tripping over each other.

Invoking this skill **is** the authorisation to spawn subagents; nothing else
in this repo is. It is not for a single task — that is `/dev-task`, and one
issue orchestrated is one issue done slowly.

Three rules shape everything below:

- **A subagent cannot ask the human anything.** Only the main thread can run an
  interview. Agents return their open questions as _data_ and you batch them.
- **Fan-out is the easy half; landing is the hard half.** N lanes means N
  branches into a main that moves under each of them.
- **Nothing is true because a lane said so.** Same rule as everywhere else: a
  claim you did not re-run is a claim.

## How to report

**The format is `.claude/reporting.md`** — the status line, the eight sections,
the five-bullet cap, the confidence rubric, the footer. Count phases out of 6:
`orchestrate · phase 3/6 · grouping`.

What this skill adds:

- **`🔧 CHANGED` is per lane**, one line each, with the lane's own confidence
  grade. Do not merge them into a single narrative.
- **`📊 CONFIDENCE` is the lowest lane's grade**, not an average. One shaky lane
  makes the run shaky, because they all land into the same main.
- **The footer names every live lane** and its URL — that is how the human looks
  at any of them.

## 1. Intake

Three inputs, all valid:

- **Pasted issues** — free text. Each becomes a lane candidate.
- **Backlog ids** — `dev/agent/backlog.sh show <id>`. These already record
  `--files`, which phase 3 needs.
- **Nothing** — `backlog.sh list` and work what it shows: the shared pool plus
  your own. Another session's unpromoted items are invisible on purpose — they
  are findings it is still holding opinions about, and its files are often live
  while it works them.

`backlog.sh list --orphaned` is worth a glance at intake: items whose session
has ended without promoting them. Nobody owns those, and nobody will.

Restate the list before doing anything, numbered. The human corrects the
restatement now, not after three worktrees exist.

## 2. Scope — read-only, and cheap

**Every issue needs a file list before anything can be grouped**, and a pasted
issue has none.

Fan out one read-only agent per issue whose entire job is to answer: _which
files would this change?_ No worktree, no server, no edits. `Explore` is the
right agent type. A backlog item with `--files` skips this.

Under-listing is the dangerous direction — a file missed here is two lanes
colliding in phase 5.

## 3. Group — by file overlap, nothing else

Issues whose file sets are **disjoint** may run concurrently. Issues that share
even one file go in the same lane and run **in sequence** inside it.

- Grouping is arithmetic on the phase-2 lists. Do it in your own head or in a
  script, never by asking an agent's opinion.
- **Three lanes at once** by default. More only if the human asks; each lane
  needing a stack costs a database clone and a first-run build.
- Extra lanes queue. Say in the report how many are waiting.

Report the grouping before provisioning. A human who disagrees with the
grouping wants to say so before the worktrees exist.

## 4. Provision — a stack only where one is needed

- **Needs to run** — touches `server/`, or anything the browser or e2e must
  exercise → `dev/agent/worktree.sh new <lane>` then `up`.
- **Does not** — docs, a pure refactor, a config change → a plain checkout, or
  work it in sequence in the primary.

**Never use the Agent tool's `isolation: "worktree"`.** It gives a checkout and
nothing else — no database, no ports, no server — so an agent in one can edit
but cannot verify, which is the false-green this whole toolkit exists to
remove. `worktree.sh` is the only worktree that can run the product.

Eight slots exist in total, shared with every other session on the machine.
Run `worktree.sh gc` first: a slot held by an abandoned worktree is a lane you
cannot start.

## 5. Execute — in waves, questions collected

Each lane's agent gets a brief that names the issue, its files, its worktree
path and ports, and **what is already known to be wrong**. It runs the
`/dev-task` flow inside its lane.

**Tell each lane to pass `worktree: <its lane name>` to every repo-acting MCP
tool.** Lanes share one MCP process, rooted at the orchestrator's checkout —
the primary. A lane that omits it tests the primary's tree and commits to the
primary's branch, and both come back saying they succeeded.

Every agent returns, as structured data:

- what changed, and the commit(s) on its branch
- its confidence grade, and the rubric line behind it
- **its open questions** — the decisions it could not make alone
- what it deliberately did not do

**Agents commit to their own branch and stop. No lane lands.** Phase 6 does
that, in an order you choose.

**Then batch the questions into one interview.** Four at a time, ordered by
what blocks the most, each carrying the evidence its lane found. A lane that
comes back with a question and no fix is a normal outcome — say so plainly
rather than sending the agent back to guess.

Re-run anything a lane claims before you repeat it as fact.

## 6. Land in order, then collect

Serialised, always. Each `worktree.sh land` fast-forwards onto a main the
previous one just moved.

- **Order by blast radius**: the lane the others build on goes first.
- **A refusal stops the run.** `land` refuses on an unclean rebase or an
  overlap with the primary's uncommitted work. Report its reason verbatim and
  hand the conflict to the human — never resolve one inside a subagent, and
  never work around a refusal.
- **A lane below High confidence does not land**, exactly as in `/dev-task`. It
  keeps its worktree and becomes an interview option.
- **`worktree.sh gc` last.** Report what it held back and why.

Then the report: one line per lane, the lowest grade as the run's grade, and
every surviving worktree named with its URL.
