---
name: namespace_building
description: Rules for creating, looking up, and updating Compose namespaces.
triggers:
  - compose_namespace_lookup
  - compose_namespace_create
  - compose_namespace_update
  - compose_namespace_delete
---

# Namespaces

A namespace is a self-contained application that groups modules, pages, and automations. Reference by handle or numeric ID.

Before creating any module, you must know which namespace to put it in. Call `compose_namespace_lookup` first to see what namespaces exist. If there is only one, use it. If there are several, pick the one whose name or handle best matches the user's request — use your judgement. Only ask the user if you genuinely cannot tell from context.

## Creating a namespace

- Call `compose_namespace_create` directly. Never suggest that an existing namespace could serve the same purpose. Never ask if the user wants to use something else instead.
- Use the name and slug the user specifies directly — do not check whether a similar namespace already exists before creating.

## Resolving a namespace

- Prefer handles over numeric IDs.
- If a tool call is denied with "namespace is required but was not specified", you forgot to resolve the namespace first — call `compose_namespace_lookup`, then retry.
