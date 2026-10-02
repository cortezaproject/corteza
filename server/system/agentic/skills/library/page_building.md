---
name: page_building
description: Block options reference for Compose pages — the option names each block kind actually reads, and what happens when one is wrong.
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

## Options keys are exact

An option name the block does not read is **kept and ignored**. The page saves,
the tool reports success, and the block draws an empty panel — nothing errors
until a person opens the page. `compose_page_block_schema` returns the real
names for a kind; it is generated from the same structs the webapp reads, so it
is the authority and this page is a summary of it.

Only two spellings are forgiven: `module` and `chart` are folded to `moduleID`
and `chartID` before the block is stored. Nothing else is.

## Block options reference

**RecordList**: `{"moduleID": "candidate", "fields": [{"name": "full_name"}], "perPage": 20, "prefilter": "stage = 'applied'", "presort": "createdAt DESC"}`
— `moduleID` is required and takes a module name, handle or ID; the server resolves it. `fields` picks the columns (omit for all).

**Record**: `{"fields": [{"name": "full_name"}]}`
— No module here: a Record block takes it from the page's own `module`. Omit `fields` to show all.

**Chart**: `{"chartID": "by_stage"}`
— Required, and takes a chart name, handle or ID. Create the chart first with `compose_chart_create`, then reference it here.

**Content**: `{"body": "<p>HTML here</p>"}`

**Metric**: `{"metrics": [{"label": "Total candidates", "moduleID": "candidate", "metricField": "count", "operation": "", "filter": ""}]}`
— One tile per entry in `metrics`. `metricField: "count"` with an empty `operation` counts records; any other value is a numeric field being reduced and needs `operation` set to `sum`, `avg`, `min` or `max`. `label` is what the tile is called — without it the tile reads "Unnamed metric". There is no `field` and no `reduce`.

**Calendar**: `{"defaultView": "month", "feeds": [{"resource": "compose:record", "startField": "start_at", "endField": "end_at", "titleField": "title", "options": {"moduleID": "interview"}}]}`
— The module lives under the feed's own `options`, not beside `startField`.

**Comment**: `{"moduleID": "candidate", "contentField": "body", "titleField": "full_name"}`
— `moduleID` and `contentField` are both required. The field holding the text is `contentField`; there is no `commentField`.

**Progress**: `{"value": {"moduleID": "candidate", "field": "score", "operation": "avg"}, "minValue": {"default": 0}, "maxValue": {"default": 100}}`
— All three are objects, not numbers. Each reads a value from a module (`moduleID` + `field` + `operation`) or takes a fixed `default`. Only `value` is required.

**RecordOrganizer**: `{"moduleID": "candidate", "groupField": "stage", "labelField": "full_name", "positionField": "sort_order"}`

**Automation**: `{"buttons": [{"label": "Run", "enabled": true, "automationID": "<taq ID>", "triggerHandle": "manual"}]}`
— Every button needs one target: `automationID` plus `triggerHandle` (a TAQ), `workflowID` plus `stepID`, or `script`. A button with none renders disabled.

`compose_page_block_schema` with no `kind` lists every kind this server has — Tabs, IFrame, File, Navigation, RecordRevisions, Geometry, AgentChat and ChatbotInbox included — with the options each cannot render without.

## What "required" means here

Leaving a required option out is **not** an error. The block is created and
draws an empty panel, and `compose_page_create` says so in its returned note —
read it. The note checks that the option is present, not that its contents make
sense, so a Metric with a `metrics` array of unusable tiles reports clean.

## Adding and editing blocks

To add a new block: call `compose_page_update` with only the new block in `blocks`, without a `blockID`. Existing blocks are preserved automatically.

To edit an existing block: look the page up first, include the block's `blockID` exactly as returned. Never invent a `blockID`.

Where a block sits is the layout's business, not the page's: the grid, the
default sizes and `xywh` are in the `page_layout` skill.

When the user asks to add a block, look the page up first and check what's actually there. Only the response from that lookup counts — not memory, not prior turns.
