---
name: page_record_list
description: How to create a record list page (the table view for a module).
triggers:
  - compose_page_create
  - compose_page_update
---

# Record List Page

A record list page shows all records of a module in a table. It is the entry point — users land here first, then open individual records from it.

## When to create one

Any time the user asks to set up pages for a module, or to "create a page" for a module, always create both a record list page AND a record detail page. Never create one without the other.

## How to create it

Call `compose_page_create` with:
- `title` — the module name, usually pluralised (e.g. "Sleep Logs")
- `visible: true` — it appears in the sidebar navigation
- Do NOT set `module` at the page level
- `blocks` — a single `RecordList` block:

```json
[{
  "kind": "RecordList",
  "title": "Sleep Logs",
  "options": {
    "module": "Sleep Log"
  }
}]
```

The server resolves `"module"` from a name to an ID automatically — you do not need to look up the ID first.

After creating the list page, immediately create the matching record detail page. Save the list page ID — the detail page will be set as its child.
