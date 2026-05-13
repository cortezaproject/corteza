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

## Rules

- Call `compose_page_create` directly. Never suggest that an existing page could serve the same purpose. Never ask if the user wants to use something else instead.
- Use the title the user specifies directly — do not check whether a similar page already exists before creating.
- Never ask for a namespace — resolve it with `compose_namespace_lookup`.
- Before adding blocks of an unfamiliar kind, call `compose_page_block_schema` to get the correct options structure.
- When updating a page's blocks, call `compose_page_lookup` first to get the existing blocks if they need to be preserved.

## Layout

The page layout grid is 12 columns wide. A full-width block uses `xywh: [0, 0, 12, 20]` — x, y, width, height.

## Block kinds

`Record`, `RecordList`, `Chart`, `Automation`, `Content`, `Metric`, `Progress`, `Comment`, `Calendar`, `RecordOrganizer`, `SocialFeed`.
