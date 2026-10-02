# Intent system gotchas

Working facts about the intent tooling (`.intent/intent.mjs`, `dev_intent_check`) that `.intent/SPEC.md` and the intent skills do not state.

## Build the sync file list by hand

`node .intent/intent.mjs sync <files>` records each named file's current hash as reconciled with its governing doc. Never feed it `git status` output: another session often edits the same checkout, and its unsynced files are in that list too. `dev_intent_check`'s `yours` is not your list either — it is every drifted file in the working tree, whoever changed it.

**Why:** a sync over a derived list marks someone else's in-flight files reconciled against docs nobody read.

**How to apply:** list the files the task actually touched and pass exactly those.

## Sync scope widens silently

`sync` with no paths syncs every enforced file, and a folder path syncs everything beneath it. Every sync also drops lock entries for files that no longer exist, repo-wide. There is no dry run. A scoped sync changes a handful of lock lines; a diff in the hundreds means it went repo-wide.

**How to apply:** check the `.intent/intent.lock.json` diff before committing; if it is wrong and the lock was clean beforehand, `git checkout -- .intent/intent.lock.json` and re-sync the right list. Splitting a change into two commits means one sync per commit, each with its own files.

## Full check fails on the baseline

`intent.mjs check` with no option covers every enforced file and fails on hundreds of `(never synced)` files nobody has backfilled; that is not a signal about the task. Positional paths are ignored: `check <paths>` still prints the full baseline. Only `--changed` (working tree plus untracked) and `--staged` narrow it.

**How to apply:** judge a task by `check --changed`. To compare full runs, compare the line count against the primary checkout, and grep the output for the task's own paths — short names like `CSidebar` match `CSidebarSearchNav` too.

## Capture the exit code

`check` and `coverage` print their failures to stderr and exit 1; a clean run prints one OK line to stdout. Piping through `tail` or `head` can hide a failure block or drop it entirely.

**How to apply:** read the whole output and the exit status, never a truncated pipe.
