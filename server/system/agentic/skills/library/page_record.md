---
name: page_record
description: How to create a record detail page (the single-record form for a module).
triggers:
  - compose_page_create
  - compose_page_update
---

# Record Detail Page

A record detail page shows a single record. Users reach it by clicking a row in the record list. It is always a child of the list page and never appears in the sidebar.

## When to create one

Every time you create a record list page, immediately follow it with a record detail page for the same module. Do not skip this step.

## How to create it

Call `compose_page_create` with:
- `title` — the module name, singular (e.g. "Sleep Log")
- `module` — the module name, handle, or ID (this is what makes it a detail page)
- `visible: false` — it must NOT appear in the sidebar
- `parent` — the ID of the list page you just created
- `blocks` — a single `Record` block. Do NOT put the module in the block options, it comes from the page:

```json
[{
  "kind": "Record",
  "title": "Sleep Log"
}]
```

Omit `fields` in the Record block options to show all module fields. Set `fields` only if the user asks to show specific fields.

## Record list needs a detail page to be usable

Without a detail page, users can see the table but cannot open, edit, or create individual records. Always create the detail page right after the list page.
