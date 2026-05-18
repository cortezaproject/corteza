---
name: page_building
description: Rules and reference for building Compose pages.
triggers:
  - compose_page_lookup
  - compose_page_create
  - compose_page_update
  - compose_page_delete
  - compose_page_reorder
  - compose_page_block_schema
---

# Page Building

A page is a screen in the application's navigation. It can contain blocks that display records, charts, content, and more.

There are two distinct page types:

1. **Record list page** — shows all records in a table. Do NOT set the module parameter at page level. Add a RecordList block with moduleID in its options.
2. **Record detail page** — the form for viewing or editing a single record. Set the module parameter at page level. Add a Record block with the fields to display. Only one record detail page can exist per module.

## Creating a page

- Call `compose_page_create` directly. Never suggest that an existing page could serve the same purpose. Never ask if the user wants to use something else instead.
- Use the title the user specifies directly — do not check whether a similar page already exists before creating.
- Never ask for a namespace — resolve it with `compose_namespace_lookup`.

## Adding a block to an existing page

Follow this sequence — never skip a step:

1. Call `compose_page_lookup` with the page handle/title. The response is the single source of truth for what's on the page. Treat anything not in this response as not existing.
2. Compute `nextY = max(block.y + block.h)` across the existing blocks. This is the y for your new block — anything lower will overlap.
3. Call `compose_page_update` with `blocks` containing only the new block. Do not include `blockID` — that's how the handler knows it's new. Existing blocks you don't pass are kept.

A typical new full-width block looks like `{"kind": "Content", "xywh": [0, <nextY>, 12, 4], "options": {...}}`.

When the user asks to add a block, ADD IT. Do not claim it already exists unless the lookup response from THIS turn contains a block of the matching kind AND content. Past conversations, your own prior messages, or assumptions are not evidence — only the JSON returned by the most recent `compose_page_lookup` counts.

## Editing an existing block

- Look the page up to get the existing block's `blockID`.
- Call `compose_page_update` with that block in `blocks`, including the `blockID` exactly as returned by lookup. Never invent or guess a `blockID`.

## Layout

Grid is 12 columns wide. Full-width = `w: 12`. Heights are roughly `h: 4` for a small Content/Metric block, `h: 20` for a list. Blocks must not overlap on the y-axis — overlapping blocks get rendered side-by-side at half-width.

## Block options

`compose_page_block_schema` returns only the field names. Use these option formats:

- `Content`: `{"body": "<p>HTML content here</p>"}` — body holds raw HTML.
- `RecordList`: `{"moduleID": "<id>"}` — at minimum the module to list.
- `Record`: `{"moduleID": "<id>"}` — the module whose record is being shown.
- `Chart`: `{"chartID": "<id>"}` — references an existing chart definition.
- `Metric`: `{"moduleID": "<id>", "metrics": [...]}` — see schema for shape.

## Block kinds

`Record`, `RecordList`, `Chart`, `Automation`, `Content`, `Metric`, `Progress`, `Comment`, `Calendar`, `RecordOrganizer`, `SocialFeed`.
