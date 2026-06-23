---
name: page_building
description: Block options reference and layout rules for Compose pages.
triggers:
  - compose_page_lookup
  - compose_page_create
  - compose_page_update
  - compose_page_delete
  - compose_page_reorder
  - compose_page_block_schema
---

# Page Building

## Page types

- **Module pages** — always a list+detail pair. See the record_list and record page skills.
- **Dashboard pages** — overview screens with charts, metrics, and summaries. A single page, no module.
- **Content pages** — informational screens. A single page with Content blocks.

## Block options reference

Options keys must match exactly — wrong keys are silently ignored.

**RecordList**: `{"module": "<module name or ID>", "fields": [{"name": "fieldName"}, ...], "perPage": 20}`
— `module` is required. `fields` filters which columns show (omit for all). `prefilter` and `presort` accept filter/sort expressions.

**Record**: `{"fields": [{"name": "fieldName"}, ...]}`
— No `module` here. Module comes from the page. Omit `fields` to show all.

**Chart**: `{"chart": "<chart name or ID>"}`
— Key is `"chart"`, not `"chartID"`. You cannot create charts. If the user asks for a page with charts, tell them you cannot create charts yet and skip the Chart blocks.

**Content**: `{"body": "<p>HTML here</p>"}`

**Metric**: `{"metrics": [{"label": "Total Records", "moduleID": "<module name or ID>", "field": "fieldName", "reduce": "count"}]}`
— `label` is required on every metric item. Without it the block shows "Unnamed metric".

**Calendar**: `{"defaultView": "month", "feeds": [{"moduleID": "<module name or ID>", "startField": "start", "endField": "end", "titleField": "name"}]}`

**Comment**: `{"moduleID": "<module name or ID>", "commentField": "fieldName", "titleField": "fieldName"}`

**Progress**: `{"moduleID": "<module name or ID>", "field": "fieldName", "minValue": 0, "maxValue": 100}`

**RecordOrganizer**: `{"moduleID": "<module name or ID>", "groupField": "status", "labelField": "name"}`

**Automation**: `{"buttons": [{"label": "Run", "enabled": true, "workflowID": "<id>"}]}`

Module and chart references can be names or IDs — the server resolves them automatically.

## Adding and editing blocks

To add a new block: call `compose_page_update` with only the new block in `blocks`, without a `blockID`. Existing blocks are preserved automatically.

To edit an existing block: look the page up first, include the block's `blockID` exactly as returned. Never invent a `blockID`.

When the user asks to add a block, look the page up first and check what's actually there. Only the response from that lookup counts — not memory, not prior turns.
