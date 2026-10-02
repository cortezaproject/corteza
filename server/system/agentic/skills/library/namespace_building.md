---
name: namespace_building
description: Rules for creating, resolving, updating and restoring Compose namespaces — the slug is an identifier, a deleted namespace is reachable only by ID.
triggers:
  - compose_namespace_lookup
  - compose_namespace_create
  - compose_namespace_update
  - compose_namespace_delete
  - compose_namespace_undelete
---

# Namespaces

A namespace is the container one app is built in: its modules, records, pages
and charts. Every tool that takes `namespace` accepts the name, the slug or the
ID as a string.

## Resolving one

Before creating a module or a page, know which namespace it goes in. Call
`compose_namespace_lookup` once: with no argument it lists what exists; with a
name or slug it returns that one, and when the name is not found it returns the
list instead so you can pick without a second call. If there is one namespace,
use it. If there are several, pick the one whose name or slug matches the
request; ask the user only when nothing matches.

A call refused with `invalid_argument` and "namespace is required" means you
left the argument out; one refused with `not_found` means the name is not a
namespace on this instance — read the list the message carries.

## Creating one

Look up first, then create without asking: the lookup is what tells you whether
the app already exists under another name, and the slug has to be unique. Do
not ask the user whether an existing namespace could serve instead — if the
lookup shows the app they asked for, use it; otherwise create what they asked
for with the name and slug they gave.

The slug is an identifier, not a label: lowercase letters, digits and
underscores. It is what filters, expressions and URLs refer to, so a hyphen in
it is a subtraction wherever an identifier is parsed. A namespace starts
enabled.

## Changing and removing one

- `compose_namespace_update` writes only the fields you send. Renaming the slug
  changes every URL and filter that names it; disabling hides the namespace from
  users without deleting anything.
- `compose_namespace_delete` is soft: the namespace disappears from every lookup
  and its modules, records, pages and charts are left in place. Agent allow-list
  entries pointing at it are dropped.
- `compose_namespace_undelete` restores it. It takes the numeric ID and nothing
  else, because a deleted namespace cannot be resolved by name or slug — keep
  the ID the delete reported.
