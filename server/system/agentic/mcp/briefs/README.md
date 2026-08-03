# Per-resource briefs

One file per resource, written before its tools are. Each is the complete input
to whoever writes that resource's tools — read alongside `CONVENTIONS.md` (the
rules) and `RESOURCES.md` (the names), not instead of them.

## Why these exist

Writing `system/reminder` from the conventions alone required inventing roughly
ten things the spec had no answer for: which filter fields become params, what
to do when a service has no `FindByAny`, how a `*time.Time` param is shaped,
whether an assignee defaults to the caller, which service methods become domain
ops. Prose cannot carry that; it is per-resource fact.

A brief that omits something is how ~165 tools end up inconsistent, so each one
records what was actually read in the service, with the file references to
check it.

## Required sections

Every brief covers, in this order:

1. **Service** — the exact type and `Default*` var, and which methods become
   tools with their real signatures.
2. **Authorization (§8.6)** — where the `ac.Can*` calls are, or the ownership
   scope, or the absence of both. A resource that fails §8.6 gets no tool and
   the brief says so.
3. **Identifier strategy** — whether `FindByAny` exists. Only `user`,
   `userGroup`, `role`, `template`, `namespace` and `module` have one; every
   other resource needs an explicit decision, not an assumption.
4. **Filter fields** — which become params and which are deliberately left out.
5. **Undelete** — whether `UndeleteByID` exists (check `*.gen.go`, not only the
   hand-written file) and how a deleted item is found.
6. **Domain ops** — service methods that are not CRUD.
7. **Traps** — what is surprising about this resource. The most valuable
   section; it is where the hours go if it is missing.

## Status

`system/reminder` is the worked exemplar and has no brief — the code is the
brief. See `system/agentic/reminder_tools.go`.
