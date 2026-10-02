---
title: MCP tools
description: Every tool Human's MCP server offers, with its risk level and parameters.
outline: [2, 2]
---

<!-- This file is auto-generated from the MCP tool registry by server/tests/mcp/docs_test.go. -->

# MCP tools

Human is an MCP (Model Context Protocol) server: an AI client that speaks MCP
can connect to it, sign in as a Human user and call the tools on this page. The
tools are the same operations Human's own [agents](/platform/agents) use. A
tool only does what the signed-in user may do — every call goes through the
same permission checks as the webapp, so a tool that reads records returns only
the records that user can read.

## Endpoint and scope

The MCP endpoint is `/api/mcp` on the server's API. Two variants list a
narrower set of tools, for a client that works on one side of the system:

| Endpoint | Lists |
| --- | --- |
| `/api/mcp` | Every tool. |
| `/api/mcp/configuring` | Tools that change what the system _is_: modules, pages, charts, TAQ and workflow definitions, roles. |
| `/api/mcp/usage` | Tools that work with what the system _holds_, or run it: records, executions, chatbots, reports. |

The variants only shorten the list; they are not a permission boundary.

Add `?maxRisk=read`, `?maxRisk=write` or `?maxRisk=destructive` to
any of these to cap a session at a risk level. Tools above the cap are left out
of the list, and calling one by name is refused. This is a safety catch against
accidents, such as a session pointed at a production instance deleting a
namespace; permissions are what protect the data.

## Risk levels

Every tool declares one of three risk levels, shown as a badge on each entry
below and in the webapp's tool picker:

| Badge | Level | Meaning |
| --- | --- | --- |
| <Badge type="tip" text="Reads" /> | `read` | Makes no change. It can still return sensitive data, which permissions govern. |
| <Badge type="warning" text="Writes" /> | `write` | Creates or changes something. Running a TAQ, workflow or agent is always a write. |
| <Badge type="danger" text="Deletes" /> | `destructive` | Removes something. Most deletes can be undone with the matching `undelete` tool, but a deleted resource is hidden from everything until then. |

## Finding and loading tools

To keep the tool list small, the server lists each tool with a one-sentence
summary and without per-parameter descriptions. Two more tools, which belong to
the MCP server rather than to any part of Human, return the full documentation
shown on this page:

- `human_tool_load` — Get the full documentation for tools you can already name, without searching. The tool list summarises every tool to one line and drops the per-parameter detail; this returns it in full. Call this before using any tool whose summary says to, and any time a call was rejected for arguments you are unsure about. Use human_tool_search instead when you know what you want to do but not what it is called.
- `human_tool_search` — Find the right tool for what you are trying to do, and get its full documentation. Every tool on this server is already listed and callable, but the listing shows only a one-line summary of each and omits the per-parameter detail; this returns the complete description and parameter documentation for the tools that match. Use it when you are unsure which tool does the job, and whenever a tool's summary tells you to load it first. Query by resource or action, in as many words as you like: "page", "delete role", "create a user", "workflow". Tools matching more of your words rank first, so an extra word narrows the ranking rather than emptying the result. If you get nothing back, try a different word.

## Compose

Namespaces, modules, records, pages, page layouts and charts.

### `compose_chart_create` {#compose_chart_create}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Create a chart in a namespace. After creating, place it on a page with a Chart block: {"kind":"Chart","options":{"chartID":"&lt;created ID>"}}.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `config` | string | Yes | JSON chart configuration. Canonical shape, counting records grouped by a Select field:<br>{"colorScheme":"tableau.Tableau10","reports":[{"moduleID":"&lt;id from compose_module_lookup>","filter":"","dimensions":[{"field":"&lt;grouping field>","modifier":"(no grouping / buckets)","conditions":{}}],"metrics":[{"field":"count","type":"doughnut"}]}]}<br><br>A metric needs "field" and "type", and needs "aggregate" as well whenever "field" is anything other than "count":<br>- "field":"count" counts records and takes no aggregate.<br>- any other field is a numeric module field being reduced, so it MUST carry "aggregate": one of SUM, MAX, MIN, AVG. Summing an amount is {"field":"amount","aggregate":"SUM","type":"bar"}. Leaving it out is the single most common way to end up with a chart that renders "Metrics aggregate not defined" instead of data.<br>- "type": pie, doughnut, bar, line, funnel, gauge, radar or scatter. This is also what shapes the chart, so a funnel is made by naming it here. A gauge additionally needs bands on its dimension: "meta":{"steps":[{"value":0},{"value":50},{"value":100}]}.<br><br>A dimension is what the metric is grouped by. "field" is a module field — a Select field groups well — and "modifier" is one of "(no grouping / buckets)" (the field's own values, and the default when omitted), DATE, WEEK, MONTH, QUARTER or YEAR, the last five for bucketing a date field.<br><br>moduleID is the numeric ID from compose_module_lookup; a handle here leaves the chart with nothing to query.<br><br>"colorScheme" is optional and takes a "&lt;family>.&lt;Name>" key from the webapp's tables — brewer, office or tableau, e.g. "tableau.Tableau10". Omit it for the default palette. The number in a name is its swatch count and is part of the name, so guessing it is the usual way to get this wrong: "tableau.ClassicOrangeBlue13" exists, "tableau.ClassicOrangeBlue7" does not. An unknown name is refused here, with the near matches, because the webapp resolves it to no palette at all and the chart then draws its legend with no visible series. |
| `name` | string | Yes | Chart name |
| `namespace` | string | Yes | Namespace name, handle, slug, or ID (as string to prevent precision loss) |
| `handle` | string |  | URL-friendly identifier. Use snake_case (lowercase letters, digits, underscores); a hyphen is the subtraction operator wherever an identifier is parsed. |

### `compose_chart_delete` {#compose_chart_delete}

<Badge type="danger" text="Deletes" /> Listed under `/api/mcp/configuring`.

Delete a chart by name, handle, or ID. Remove Chart blocks referencing it from pages first — they render empty otherwise. The delete is soft and compose_chart_undelete reverses it, but note the chart ID first: a deleted chart can no longer be found by name or handle.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `chart` | string | Yes | Chart name, handle, or ID (as string to prevent precision loss) |
| `namespace` | string | Yes | Namespace name, handle, slug, or ID (as string to prevent precision loss) |

### `compose_chart_lookup` {#compose_chart_lookup}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/configuring` and `/api/mcp/usage`.

List charts in a namespace, or look one up by name, handle, or ID. Provide 'chart' to fetch a single chart in full, including its config; omit it to list every chart in the namespace. A listing is trimmed to chartID, name and handle — it never carries the config, so fetch the chart itself before changing its configuration. Chart blocks on pages reference charts by chartID.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `namespace` | string | Yes | Namespace name, handle, slug, or ID (as string to prevent precision loss) |
| `chart` | string |  | Chart name, handle, or ID (as string to prevent precision loss). Omit to list all charts in the namespace. |
| `limit` | string |  | Maximum charts to return when listing, default 50, capped at 200. |
| `pageCursor` | string |  | Cursor from a previous response, to fetch the next page. |

### `compose_chart_undelete` {#compose_chart_undelete}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Restore a soft-deleted chart, reversing compose_chart_delete. A delete only marks the chart — its name, handle and full report configuration were retained — so it comes back unchanged and Chart blocks that reference its chartID render again, without the page needing an edit. Requires the numeric chartID: deleted charts are excluded from every lookup path, so compose_chart_lookup can no longer resolve one by name or handle, and it exposes no includeDeleted-style filter. Use the ID compose_chart_delete reported, or one from a compose_chart_lookup taken before the delete. Calling this on a chart that is not deleted is accepted and changes nothing.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `chartID` | string | Yes | ID of the deleted chart (as string to prevent precision loss). A name or handle will not work — deleted charts are not resolvable by either. |
| `namespace` | string | Yes | Namespace name, handle, slug, or ID (as string to prevent precision loss) |

### `compose_chart_update` {#compose_chart_update}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Update an existing chart. A field you omit is left unchanged; passing an empty string for 'name' or 'handle' clears it. 'config' replaces the whole configuration rather than merging into it, so send the complete config, not only the part you are changing — call compose_chart_lookup with 'chart' first if you need the current one.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `chart` | string | Yes | Chart name, handle, or ID (as string to prevent precision loss) |
| `namespace` | string | Yes | Namespace name, handle, slug, or ID (as string to prevent precision loss) |
| `config` | string |  | JSON chart configuration. Canonical shape, counting records grouped by a Select field:<br>{"colorScheme":"tableau.Tableau10","reports":[{"moduleID":"&lt;id from compose_module_lookup>","filter":"","dimensions":[{"field":"&lt;grouping field>","modifier":"(no grouping / buckets)","conditions":{}}],"metrics":[{"field":"count","type":"doughnut"}]}]}<br><br>A metric needs "field" and "type", and needs "aggregate" as well whenever "field" is anything other than "count":<br>- "field":"count" counts records and takes no aggregate.<br>- any other field is a numeric module field being reduced, so it MUST carry "aggregate": one of SUM, MAX, MIN, AVG. Summing an amount is {"field":"amount","aggregate":"SUM","type":"bar"}. Leaving it out is the single most common way to end up with a chart that renders "Metrics aggregate not defined" instead of data.<br>- "type": pie, doughnut, bar, line, funnel, gauge, radar or scatter. This is also what shapes the chart, so a funnel is made by naming it here. A gauge additionally needs bands on its dimension: "meta":{"steps":[{"value":0},{"value":50},{"value":100}]}.<br><br>A dimension is what the metric is grouped by. "field" is a module field — a Select field groups well — and "modifier" is one of "(no grouping / buckets)" (the field's own values, and the default when omitted), DATE, WEEK, MONTH, QUARTER or YEAR, the last five for bucketing a date field.<br><br>moduleID is the numeric ID from compose_module_lookup; a handle here leaves the chart with nothing to query.<br><br>"colorScheme" is optional and takes a "&lt;family>.&lt;Name>" key from the webapp's tables — brewer, office or tableau, e.g. "tableau.Tableau10". Omit it for the default palette. The number in a name is its swatch count and is part of the name, so guessing it is the usual way to get this wrong: "tableau.ClassicOrangeBlue13" exists, "tableau.ClassicOrangeBlue7" does not. An unknown name is refused here, with the near matches, because the webapp resolves it to no palette at all and the chart then draws its legend with no visible series. |
| `handle` | string |  | New handle |
| `name` | string |  | New name |

### `compose_module_create` {#compose_module_create}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Create a new module in a namespace. A module defines a data structure (like a table) with typed fields. As a developer acting on behalf of the user, proactively add config where appropriate: enable duplicate detection for modules storing contacts/leads/customers (match on email or phone), enable recordRevisions for important transactional data, and set privacy disclosure for modules holding personal information.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `handle` | string | Yes | URL-friendly identifier (lowercase letters, digits, and underscores only) |
| `name` | string | Yes | Display name for the module |
| `namespace` | string | Yes | Namespace name, handle, slug, or ID (as string to prevent precision loss) |
| `config` | string |  | JSON object for module-level configuration. Supports: recordDeDup (duplicate detection), recordRevisions (audit trail), privacy (data sensitivity). Rules live under a "rules" array — a bare rule object is accepted and silently stored as {}. Example: {"recordDeDup":{"rules":[{"name":"unique-email","strict":true,"constraints":[{"attribute":"email","modifier":"ignore-case\|case-sensitive\|fuzzy-match\|sounds-like","multiValue":"one-of\|equal"}]}]},"recordRevisions":{"enabled":true},"privacy":{"usageDisclosure":"text","sensitivityLevelID":"123"}} |
| `detail` | string |  | How much of the stored module to echo back. "summary" (the default) confirms what was written — handle, name, and each field's name, kind, label, required/multi, select options and value expression. "full" adds every field's ID, timestamps and DAL storage config, which is thousands of tokens and is only useful when you need the field IDs. Issues with a field expression are reported either way. |
| `fields` | string |  | JSON array of field definitions. Each field: {"name":"fieldName","kind":"String","label":"Field Label","isRequired":false,"isMulti":false,"options":{},"defaultValue":[{"name":"fieldName","value":"default"}],"expressions":{"value":"","sanitizers":[],"validators":[],"formatters":[]}}. Supported kinds: String, Number, Bool, DateTime, Select, Email, Url, File, User, Record, Geometry — these are the whole set, and a kind outside it is not a field type the webapp can render. Kind-specific options: Select→{"options":[{"value":"a","text":"A"}]}, Record→{"moduleID":"123","labelField":"name","queryFields":["name"],"selectType":"default"} (labelField is what the picker and every viewer show; without it they fall back to the module's first field. recordLabelField is only the second-level label for when labelField itself points at a Record field), Number→{"precision":2,"format":"0,0.00","prefix":"","suffix":""} (precision rounds what is STORED; display comes from format, so precision alone drops the decimals it kept), DateTime→{"onlyDate":false,"onlyTime":false}, Bool→{"trueLabel":"Yes","falseLabel":"No"}, Geometry→{"center":[46.05,14.51],"zoom":7}. "defaultValue" fills the field when a record is created without it; the "name" key is optional and is stored empty, because the field already says which field it is for. Reserved field names, refused because the record already carries them as system fields: new, old, id, ID, recordID, tenantID, projectID, namespaceID, moduleID, revision, meta, ownedBy, createdAt, createdBy, createdByAgent, updatedAt, updatedBy, deletedAt, deletedBy; so are the query keywords NULL, TRUE, FALSE, IS, LIKE, NOT, AND, OR, XOR, IN, BETWEEN, DESC, ASC, INTERVAL in any casing, which record filters could never refer to — name an identifier field for what it identifies instead, e.g. IssueID. EXPRESSIONS. Each slot gets a different scope. String literals need DOUBLE quotes — a single-quoted string is a syntax error. "value" (the field value expression) sees the record's own fields by their BARE names (stage != "hired" && stage != "rejected"), plus "new" and "old" as whole records (new.values.stage, new.recordID, old.values.stage; on a create every "old" field is null). There is no "record" and no bare "values" — record.values.stage fails every save with 'unknown parameter record.values'. A value expression OVERWRITES whatever the caller sent for that field, so do not send it. "isRequired" on such a field means the expression must produce a value, not that a caller must supply one. "validators":[{"test":"…","error":"…"}] — test sees "value" (the value being saved, as a string), "oldValue" and "values.&lt;field>". READ THIS ONE TWICE: test names the condition under which the value is REJECTED. test "value >= 0 && value &lt;= 5" REFUSES 3 and stores 7 — the exact opposite of what it reads like, while showing an error message asserting the range. Write the rule you want as its rejection: "value &lt; 0 \|\| value > 5". "sanitizers" and "formatters" are transforms, not tests: each sees only "value" and its RESULT REPLACES the value — trim(value), toUpper(value). A sanitizer runs before validation and is stored; a formatter runs on the way out. The write result reports what cannot work under "issues"; no "issues" key is the clean result. |

### `compose_module_delete` {#compose_module_delete}

<Badge type="danger" text="Deletes" /> Listed under `/api/mcp/configuring`.

Delete a module by name, handle, or ID. The delete is soft: nothing is erased — the module stops appearing in lookups and its records become unreachable through it, but both are retained. Pages, charts and record fields that point at this module will no longer resolve it, so check what references it before deleting. compose_module_undelete reverses this, but note the module ID first — a deleted module can no longer be found by name or handle.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `module` | string | Yes | Module name, handle, or ID (as string to prevent precision loss) |
| `namespace` | string | Yes | Namespace name, handle, slug, or ID (as string to prevent precision loss) |

### `compose_module_lookup` {#compose_module_lookup}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/configuring` and `/api/mcp/usage`.

Look up modules in a namespace, or get one module's full definition. A module is a data structure (like a table) with typed fields. To list ALL modules in a namespace, omit the 'module' argument; the list is a slim view — module identity plus each field's name and kind. To get everything else about one module — field labels, options, default values, expressions and module config — provide its name, handle or ID. Never guess module or field names: list first if you are unsure what exists, then fetch the one you need. Use 'limit' and 'pageCursor' when a namespace holds more modules than one response can carry.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `namespace` | string | Yes | Namespace name, handle, slug, or ID (as string to prevent precision loss) |
| `detail` | string |  | How much of each field to return. "summary" (the default) gives name, kind, label, required/multi, select options and any value expression — what you need to read or filter data. "full" adds field IDs, timestamps and storage config, which you only need when editing the module itself. Summary is a fraction of the size; a namespace of modules at full detail is tens of kilobytes. |
| `limit` | string |  | Maximum modules to return when listing, default 50, capped at 200. |
| `module` | string |  | Module name, handle, or ID (as string to prevent precision loss). Omit to list all modules in the namespace. |
| `pageCursor` | string |  | Cursor from a previous response, to fetch the next page. |

### `compose_module_undelete` {#compose_module_undelete}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Restore a soft-deleted module, reversing compose_module_delete. A delete only sets a marker: the module definition, its fields and every record in it were retained, so the module comes back exactly as it was and the pages, charts and Record fields that reference it resolve again. Records that were deleted individually stay deleted — use compose_record_undelete for those. Requires the numeric moduleID, and nothing else will do: a deleted module is excluded from every lookup path, so compose_module_lookup can no longer resolve it by name or handle, and that tool exposes no includeDeleted-style filter. Take the ID from what compose_module_delete reported, or from a compose_module_lookup made before the delete. Calling this on a module that is not deleted is accepted and changes nothing.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `moduleID` | string | Yes | ID of the deleted module (as string to prevent precision loss). A name or handle will not work — deleted modules are not resolvable by either. |
| `namespace` | string | Yes | Namespace name, handle, slug, or ID (as string to prevent precision loss) |

### `compose_module_update` {#compose_module_update}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Update an existing module's name, handle, fields, or configuration. Arguments you omit are left unchanged. Fields are the exception to the usual replace-wholesale rule: they are MERGED by field name — fields you pass are added or updated, existing fields you do not mention are kept, and passing an empty array removes nothing. To drop fields, name them in 'removeFields'. 'config', by contrast, replaces the whole module config when provided.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `module` | string | Yes | Module name, handle, or ID (as string to prevent precision loss) |
| `namespace` | string | Yes | Namespace name, handle, slug, or ID (as string to prevent precision loss) |
| `config` | string |  | JSON object for module-level configuration. Replaces the existing config. Supports: recordDeDup (duplicate detection), recordRevisions (audit trail), privacy (data sensitivity). Example: {"recordDeDup":{"rules":[{"name":"unique-email","strict":true,"constraints":[{"attribute":"email","modifier":"ignore-case"}]}]},"recordRevisions":{"enabled":true},"privacy":{"usageDisclosure":"Used for customer contact only"}} |
| `detail` | string |  | How much of the stored module to echo back. "summary" (the default) confirms what was written — handle, name, and each field's name, kind, label, required/multi, select options and value expression. "full" adds every field's ID, timestamps and DAL storage config, which is thousands of tokens and is only useful when you need the field IDs. Issues with a field expression are reported either way. |
| `fields` | string |  | JSON array of fields to add or update. Same format as compose_module_create. Existing fields not listed are preserved. EXPRESSIONS. Each slot gets a different scope. String literals need DOUBLE quotes — a single-quoted string is a syntax error. "value" (the field value expression) sees the record's own fields by their BARE names (stage != "hired" && stage != "rejected"), plus "new" and "old" as whole records (new.values.stage, new.recordID, old.values.stage; on a create every "old" field is null). There is no "record" and no bare "values" — record.values.stage fails every save with 'unknown parameter record.values'. A value expression OVERWRITES whatever the caller sent for that field, so do not send it. "isRequired" on such a field means the expression must produce a value, not that a caller must supply one. "validators":[{"test":"…","error":"…"}] — test sees "value" (the value being saved, as a string), "oldValue" and "values.&lt;field>". READ THIS ONE TWICE: test names the condition under which the value is REJECTED. test "value >= 0 && value &lt;= 5" REFUSES 3 and stores 7 — the exact opposite of what it reads like, while showing an error message asserting the range. Write the rule you want as its rejection: "value &lt; 0 \|\| value > 5". "sanitizers" and "formatters" are transforms, not tests: each sees only "value" and its RESULT REPLACES the value — trim(value), toUpper(value). A sanitizer runs before validation and is stored; a formatter runs on the way out. The write result reports what cannot work under "issues"; no "issues" key is the clean result. |
| `handle` | string |  | New handle for the module. Use snake_case (lowercase letters, digits, underscores); a hyphen is the subtraction operator wherever an identifier is parsed. |
| `name` | string |  | New display name for the module |
| `removeFields` | string |  | JSON array of field names to remove, e.g. ["fieldA","fieldB"] |

### `compose_namespace_create` {#compose_namespace_create}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Create a new namespace. A namespace is a top-level container for modules and records in Human Compose.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `name` | string | Yes | Display name for the namespace |
| `slug` | string | Yes | URL-friendly identifier (lowercase letters, digits, and underscores only). Use snake_case: a hyphen is the subtraction operator wherever an identifier is parsed, so underscores keep the slug safe to reuse in filters and expressions. |
| `enabled` | boolean |  | Whether the namespace is enabled (default: true) |

### `compose_namespace_delete` {#compose_namespace_delete}

<Badge type="danger" text="Deletes" /> Listed under `/api/mcp/configuring`.

Delete a namespace by name, handle, slug, or ID. The delete is soft: the namespace stops appearing in lookups but is retained, and the modules, records, pages and charts inside it are left in place rather than removed. Any agent tool allow-list entry pointing at the namespace is dropped.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `namespace` | string | Yes | Namespace name, handle, slug, or ID (as string to prevent precision loss) |

### `compose_namespace_lookup` {#compose_namespace_lookup}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/configuring` and `/api/mcp/usage`.

Look up namespaces. Call this whenever the user asks what they have, what exists, what's set up, or anything about the current state of their data. Also call this to resolve a namespace before any other operation. Provide 'namespace' to fetch that one namespace in full; omit it to list namespaces as {namespaceID, name, slug}. If the namespace you name is not found, the list is returned instead so you can pick the right one without another call.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `limit` | string |  | Maximum namespaces to return when listing, default 50, capped at 200. |
| `namespace` | string |  | Namespace name, handle, slug, or ID (as string to prevent precision loss). Omit to list namespaces instead. |
| `pageCursor` | string |  | Cursor from a previous response, to fetch the next page. |

### `compose_namespace_undelete` {#compose_namespace_undelete}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Restore a soft-deleted namespace, reversing compose_namespace_delete. The delete only set a marker: the namespace, its modules, records, pages and charts were all retained, so it comes back exactly as it was and appears in compose_namespace_lookup again. Agent allow-list entries the delete dropped are not restored; re-grant them with system_agent_update. Requires the numeric namespaceID, and nothing else will do: a deleted namespace is excluded from every lookup path, so neither name nor slug resolves it. Take the ID from what compose_namespace_delete reported, or from a compose_namespace_lookup made before the delete. Calling this on a namespace that is not deleted is accepted and changes nothing.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `namespaceID` | string | Yes | ID of the deleted namespace (as string to prevent precision loss). A name or slug will not work — deleted namespaces are not resolvable by either. |

### `compose_namespace_update` {#compose_namespace_update}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Update an existing namespace's name, slug, or enabled state. Only the fields you send are written; fields you omit keep their current value, and a field sent as an empty string is cleared.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `namespace` | string | Yes | Namespace name, handle, slug, or ID (as string to prevent precision loss) |
| `enabled` | boolean |  | Whether the namespace should be enabled |
| `name` | string |  | New display name |
| `slug` | string |  | New URL-friendly identifier. Use snake_case (lowercase letters, digits, underscores); a hyphen is the subtraction operator wherever an identifier is parsed. |

### `compose_page_block_schema` {#compose_page_block_schema}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/configuring`.

Get the options structure for a page block kind — field names and types as a zero-value skeleton (not examples from live pages), under "options", plus "required": the options the block cannot render without. Leaving a required one out is not an error; the block is created and draws an empty panel, and compose_page_create says so in its note. Call this before creating blocks of an unfamiliar kind, and pass the result's field names into the block's "options" object in compose_page_create or compose_page_update. Semantics the skeleton cannot express: Metric items use metricField "count" with empty operation for record counts, or a numeric field with operation sum/avg/min/max; Chart blocks reference an existing chart resource by chartID, created with compose_chart_create; RecordList/Record moduleID/fields take IDs and field names from compose_module_lookup. This reads a static schema — it touches no namespace and no data.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `kind` | string |  | Block kind: Record, RecordList, Chart, Automation, Content, Metric, Progress, Comment, Calendar, RecordOrganizer, SocialFeed, ChatbotInbox. Omit it to list the kinds this server supports. |

### `compose_page_create` {#compose_page_create}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Create a new page in a namespace — the page holds what its blocks ARE (kind, options, title) and a LAYOUT holds where they go, so create seeds the page's primary layout from the blocks you send. A page is a screen in the namespace's navigation; it holds blocks that render records, charts and content. There are two distinct page types:

1. Record list page — shows all records in a table. Do NOT set the module parameter at page level. Add a RecordList block with moduleID in its options.
2. Record detail page — the form for viewing or editing a single record. Set the module parameter at page level. Add a Record block with the fields to display. Only one record detail page can exist per module, and creating a second one for the same module is rejected.

The layout grid is 48 columns wide, cell height 10px. Full width is w=48, half is 24, a quarter is 12 — w=12 is a QUARTER of this grid, not the half it would be on a 12-column one.

Blocks sit SIDE BY SIDE by stepping x and keeping y: a row of four tiles is [0,0,12,20], [12,0,12,20], [24,0,12,20], [36,0,12,20]. Stepping y instead puts every block in its own row with the page empty beside it, which is the layout to avoid — a dashboard reads as rows of tiles, and only tabular or form blocks earn the full 48.

The xywh you send is kept as sent. Omit xywh (or send w=0) on a block and it is placed for you: it flows into the current row at a width that suits its kind — Metric and Progress a quarter, Chart, Calendar, Content, Comment, SocialFeed and Automation a half, RecordList, Record and RecordOrganizer the full width — and wraps to a new row when the row fills, always below any block you positioned yourself. Mixing the two is fine; laying out the whole page yourself gives the better result.

Blocks CLIP their content when too short and fail silently as blank UI — give Metric blocks h>=20 and RecordList/Chart blocks h>=30.

Call compose_page_block_schema with the block kind to get its options before creating blocks. Blocks are numbered on create — the returned page carries the assigned blockIDs, which compose_page_update needs to change a block later.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `namespace` | string | Yes | Namespace name, handle, slug, or ID (as string to prevent precision loss) |
| `title` | string | Yes | Page title |
| `blocks` | string |  | JSON array of page blocks. Grid is 48 columns wide. A dashboard is rows of blocks: a row of four tiles above a full-width list is [{"kind":"Metric","title":"Total","xywh":[0,0,12,20],"options":{...}},{"kind":"Metric","title":"Open","xywh":[12,0,12,20],"options":{...}},{"kind":"Metric","title":"Urgent","xywh":[24,0,12,20],"options":{...}},{"kind":"Metric","title":"Overdue","xywh":[36,0,12,20],"options":{...}},{"kind":"Chart","title":"By month","xywh":[0,20,24,30],"options":{...}},{"kind":"Chart","title":"By owner","xywh":[24,20,24,30],"options":{...}},{"kind":"RecordList","title":"All tasks","xywh":[0,50,48,30],"options":{...}}] — four tiles across at y=0, two half-width charts across at y=20, the list full width below them. Omit xywh on a block to have it placed for you. Call compose_page_block_schema first for kind-specific options. |
| `config` | string |  | JSON object for page configuration. Example: {"navItem":{"expanded":true}}. A navigation icon is not something to set here — the webapp draws one as an image, so a font-awesome or inline-svg icon renders as a broken image and is rejected. |
| `description` | string |  | Page description |
| `handle` | string |  | URL-friendly identifier (lowercase letters, digits, and underscores; at least 2 characters) |
| `meta` | string |  | JSON object for page meta. Example: {"allowPersonalLayouts":true} |
| `module` | string |  | Module name, handle, or ID (as string to prevent precision loss). Set ONLY for record detail pages (the single-record form). Do not set for record list pages — put the module in the RecordList block options instead. |
| `parent` | string |  | Parent page title, handle, or ID (as string to prevent precision loss). Omit for a root-level page. |
| `visible` | boolean |  | Show page in navigation (default: true) |

### `compose_page_delete` {#compose_page_delete}

<Badge type="danger" text="Deletes" /> Listed under `/api/mcp/configuring`.

Delete a page by title, handle, or ID. The delete is soft: nothing is erased — the page stops appearing in lookups and in navigation, but it is retained. Records are not touched at all, because a page only displays data that lives in modules. Child pages are handled by 'strategy', which defaults to refusing the delete when the page has children, so decide explicitly what should happen to them. compose_page_undelete reverses this one page at a time, but note the page ID first: a deleted page can no longer be found by title or handle.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `namespace` | string | Yes | Namespace name, handle, slug, or ID (as string to prevent precision loss) |
| `page` | string | Yes | Page title, handle, or ID (as string to prevent precision loss) |
| `strategy` | string |  | How to handle child pages: "abort" (default — fail if the page has any undeleted children), "cascade" (delete the whole subtree), "rebase" (move the children up to this page's parent, then delete it), "force" (delete only this page and leave the children pointing at it — they then show up as root-level pages) |

### `compose_page_layout_create` {#compose_page_layout_create}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Add a layout to a page. A page holds what its blocks ARE — kind, options, title. A LAYOUT holds where they go, as {blockID, xywh}, and it is the layout the page renders and the builder draws. Every page is created with a primary layout seeded from its blocks, so these tools are for rearranging that layout, changing its record toolbar, or adding further layouts — not for giving a page its first one. A second layout is an alternative arrangement of the same page — a different subset of blocks, or a different record toolbar — shown when its visibility expression passes. Layouts are evaluated in weight order and the first whose expression passes is the one rendered, so an unconditional layout should sort last or it will always win.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `namespace` | string | Yes | Namespace name, handle, slug, or ID (as string to prevent precision loss) |
| `page` | string | Yes | Page title, handle, or ID (as string to prevent precision loss) |
| `blocks` | string |  | JSON array placing the page's blocks: [{"blockID":"1","xywh":[0,0,24,20]}]. blockID comes from compose_page_lookup; a block the layout does not name is not drawn in it. The grid is 48 columns wide and a cell is 10px tall — full width is w=48, half 24, a quarter 12 — and blocks sit side by side by stepping x and keeping y. Blocks CLIP when too short: give Metric h>=20 and RecordList/Chart h>=30. |
| `buttons` | string |  | JSON object turning record-toolbar buttons on or off, each with an optional label override: {"new":{"enabled":true},"edit":{"enabled":true,"label":"Amend"},"submit":{"enabled":true},"delete":{"enabled":false},"clone":{"enabled":true},"back":{"enabled":true}}. Only meaningful on a record page. A button you omit keeps its current setting. |
| `handle` | string |  | URL-friendly identifier, unique among the page's layouts (lowercase letters, digits, underscores) |
| `title` | string |  | Layout title, shown in the builder's layout picker |
| `visibility` | string |  | JSON object deciding when this layout is the one rendered: {"expression":"record.values.status == 'closed'","roles":["&lt;roleID>"]}. The expression is evaluated against the record being viewed, so it only makes sense on a record page. Layouts are tried in weight order and the first match wins. |
| `weight` | number |  | Sort order among the page's layouts; lower is evaluated first. |

### `compose_page_layout_delete` {#compose_page_layout_delete}

<Badge type="danger" text="Deletes" /> Listed under `/api/mcp/configuring`.

Delete a layout from a page. The page and its blocks are untouched — only the arrangement goes. Deleting a page's only layout leaves it with nothing to render and is rejected; delete the page instead, with compose_page_delete.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `layout` | string | Yes | Layout handle or ID (as string), from compose_page_layout_lookup |
| `namespace` | string | Yes | Namespace name, handle, slug, or ID (as string to prevent precision loss) |
| `page` | string | Yes | Page title, handle, or ID (as string to prevent precision loss) |

### `compose_page_layout_lookup` {#compose_page_layout_lookup}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/configuring`.

List a page's layouts, or fetch one whole. A page holds what its blocks ARE — kind, options, title. A LAYOUT holds where they go, as {blockID, xywh}, and it is the layout the page renders and the builder draws. Every page is created with a primary layout seeded from its blocks, so these tools are for rearranging that layout, changing its record toolbar, or adding further layouts — not for giving a page its first one. Provide 'layout' to fetch that one with its blocks and config; omit it to list the page's layouts as {pageLayoutID, handle, weight, title}. Call this before compose_page_layout_update: the update needs the layout's ID, and a block is placed by the blockID compose_page_lookup reports.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `namespace` | string | Yes | Namespace name, handle, slug, or ID (as string to prevent precision loss) |
| `page` | string | Yes | Page title, handle, or ID (as string to prevent precision loss) |
| `layout` | string |  | Layout handle or ID (as string). Omit to list the page's layouts. |
| `limit` | string |  | Maximum layouts to return when listing, default 50, capped at 200. |
| `pageCursor` | string |  | Cursor from a previous response, to fetch the next page. |

### `compose_page_layout_update` {#compose_page_layout_update}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Change a layout: move or resize its blocks, retitle it, change its record toolbar or its visibility expression. Every parameter is optional and only what you send is changed — omitting 'blocks' leaves the arrangement alone, and omitting 'buttons' leaves the toolbar alone. Sending 'blocks' REPLACES the arrangement wholesale, so send every block that should be placed, not only the ones that moved.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `layout` | string | Yes | Layout handle or ID (as string), from compose_page_layout_lookup |
| `namespace` | string | Yes | Namespace name, handle, slug, or ID (as string to prevent precision loss) |
| `page` | string | Yes | Page title, handle, or ID (as string to prevent precision loss) |
| `blocks` | string |  | JSON array placing the page's blocks: [{"blockID":"1","xywh":[0,0,24,20]}]. blockID comes from compose_page_lookup; a block the layout does not name is not drawn in it. The grid is 48 columns wide and a cell is 10px tall — full width is w=48, half 24, a quarter 12 — and blocks sit side by side by stepping x and keeping y. Blocks CLIP when too short: give Metric h>=20 and RecordList/Chart h>=30. |
| `buttons` | string |  | JSON object turning record-toolbar buttons on or off, each with an optional label override: {"new":{"enabled":true},"edit":{"enabled":true,"label":"Amend"},"submit":{"enabled":true},"delete":{"enabled":false},"clone":{"enabled":true},"back":{"enabled":true}}. Only meaningful on a record page. A button you omit keeps its current setting. |
| `handle` | string |  | New URL-friendly identifier |
| `title` | string |  | New layout title |
| `visibility` | string |  | JSON object deciding when this layout is the one rendered: {"expression":"record.values.status == 'closed'","roles":["&lt;roleID>"]}. The expression is evaluated against the record being viewed, so it only makes sense on a record page. Layouts are tried in weight order and the first match wins. |
| `weight` | number |  | Sort order among the page's layouts; lower is evaluated first. |

### `compose_page_lookup` {#compose_page_lookup}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/configuring`.

Look up pages in a namespace, or fetch one page in full. Provide 'page' (title, handle, or ID) to get that single page with its complete block layout — this is the only mode that returns blocks, and it is what you call before changing them, because compose_page_update needs the existing blockIDs. Omit 'page' to list the namespace's pages as {pageID, title, handle, parentID, moduleID, visible, weight}, ordered by navigation weight and paged with 'limit' and 'pageCursor'. parentID is the page's parent, 0 for a root-level page, so the hierarchy is reconstructible from a listing. Set 'tree' to true to get that same slim view nested as a hierarchy instead. A tree is always returned whole and ignores 'limit' and 'pageCursor': a hierarchy cannot be cut into pages without severing subtrees from their parents, so paging it would be meaningless. Use the flat, paged listing for a namespace with many pages. Neither listing carries blocks, config or meta — fetch the single page for those.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `namespace` | string | Yes | Namespace name, handle, slug, or ID (as string to prevent precision loss) |
| `limit` | string |  | Maximum pages to return when listing, default 50, capped at 200. Ignored in tree mode. |
| `page` | string |  | Page title, handle, or ID (as string to prevent precision loss). Omit to list the namespace's pages instead. |
| `pageCursor` | string |  | Cursor from a previous response, to fetch the next page of the listing. Ignored in tree mode. |
| `tree` | boolean |  | List the pages nested as a hierarchy instead of a flat list. Ignored when 'page' is given, and not paged. |

### `compose_page_remove_blocks` {#compose_page_remove_blocks}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Remove one or more blocks from a page, by blockID. This is the only way to take a block off a page: compose_page_update MERGES the blocks it is given by blockID, so it can add a block or overwrite one but never drop one. Use that tool to change a block and this one to delete it. This takes a block off the page, it does not delete the page — that is compose_page_delete — and it touches nothing a block pointed at: the module, records or chart a block rendered are left exactly as they were, because a block is only a view onto them. A removed block is also dropped from every layout of the page, because a block the page no longer has cannot be placed anywhere; no 'layout' argument is needed for that. Call compose_page_lookup with 'page' first to read the current blocks and their blockIDs; blockIDs are assigned per page as blocks are created, so they mean nothing on another page. Every ID you pass must exist on this page — an unknown blockID is rejected and nothing at all is removed, so a stale layout fails loudly instead of half applying. Blocks you do not list keep their position, which means removing one leaves a gap in the grid rather than reflowing the rest; reposition the survivors with compose_page_update if the layout should close up. The updated page is returned, so the remaining blocks and their blockIDs come back in the response.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `blockIDs` | string | Yes | JSON array of blockIDs to remove, as strings to prevent precision loss, e.g. ["1","3"]. Every ID must exist on the page; read them from compose_page_lookup with 'page'. |
| `namespace` | string | Yes | Namespace name, handle, slug, or ID (as string to prevent precision loss) |
| `page` | string | Yes | Page title, handle, or ID (as string to prevent precision loss) |

### `compose_page_reorder` {#compose_page_reorder}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Set the navigation order of the pages under one parent. The pages you list are weighted in the order given, so pass the COMPLETE ordered list of that parent's children: any child you leave out is pushed behind the ones you list, in an unspecified order. Omit 'parent' to order the root-level pages — note that this mode re-weights every other page in the namespace as well, so only use it with the full list of root pages. Call compose_page_lookup first to get the pageIDs and the current order; this tool changes ordering only, never the parent of a page — use compose_page_update's 'parent' to move a page.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `namespace` | string | Yes | Namespace name, handle, slug, or ID (as string to prevent precision loss) |
| `pageIDs` | string | Yes | JSON array of page IDs in the desired order, as strings to prevent precision loss, e.g. ["123","456","789"] |
| `parent` | string |  | Parent page title, handle, or ID (as string to prevent precision loss). Omit to reorder root-level pages. |

### `compose_page_undelete` {#compose_page_undelete}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Restore a soft-deleted page, reversing compose_page_delete. A delete only marks the page — its blocks, config, meta and navigation weight were all retained — so the page returns to navigation exactly as it was and no block needs rebuilding. Requires the numeric pageID: deleted pages are excluded from every lookup path, so compose_page_lookup can no longer resolve one by title or handle, its tree mode does not show it, and it exposes no includeDeleted-style filter. Use the ID compose_page_delete reported, or one from a compose_page_lookup taken before the delete. This restores one page and never its subtree: a page deleted with the 'cascade' strategy needs one call per child, and a child restored while its parent is still deleted shows up as a root-level page until the parent is restored too. Calling this on a page that is not deleted is accepted and changes nothing.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `namespace` | string | Yes | Namespace name, handle, slug, or ID (as string to prevent precision loss) |
| `pageID` | string | Yes | ID of the deleted page (as string to prevent precision loss). A title or handle will not work — deleted pages are not resolvable by either. |

### `compose_page_update` {#compose_page_update}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Update an existing page — the page holds what its blocks ARE and a LAYOUT is what places them, so a block added here is also placed on a layout: name it with 'layout', or omit that on a page that has only one. An argument you omit is left unchanged. Sending an empty string clears the value: 'parent' moves the page to the root, 'module' unlinks the module, 'handle' and 'description' are emptied. 'title' is the exception — a page without one is unusable, so an empty title is ignored rather than applied. 'blocks' is the documented exception to replacing wholesale: blocks are MERGED by blockID. A block carrying a blockID overwrites the existing block with that ID, a block without one is appended and is assigned a fresh blockID, and existing blocks you do not mention are kept — so send only the blocks you are adding or changing. An unknown blockID is rejected. Because the merge never removes, this tool cannot take a block away: use compose_page_remove_blocks for that. A block sent without xywh keeps the position it already has, so changing one block's options never moves it; a new block with no xywh is placed below the existing layout. 'config' and 'meta' each replace their whole object. Call compose_page_lookup with 'page' first when you need the current blocks or their blockIDs.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `namespace` | string | Yes | Namespace name, handle, slug, or ID (as string to prevent precision loss) |
| `page` | string | Yes | Page title, handle, or ID (as string to prevent precision loss) |
| `blocks` | string |  | JSON array of page blocks. Merged by blockID — include blockID to update an existing block, omit blockID to add a new one. Blocks are never removed here; omitting one keeps it, and compose_page_remove_blocks is what deletes one. Grid is 48 columns wide; a full-width block uses xywh [0,0,48,20] and a quarter-width tile [0,0,12,20]. xywh is read only for blocks NEW to the page; on an existing block it is accepted, echoed back and ignored, because placement lives on the layout. Moving or resizing an existing block is compose_page_layout_update. |
| `config` | string |  | JSON object for page configuration. Replaces existing config. Example: {"navItem":{"expanded":true}}. A navigation icon is not something to set here — the webapp draws one as an image, so a font-awesome or inline-svg icon renders as a broken image and is rejected. |
| `description` | string |  | New description. Pass an empty string to clear it. |
| `handle` | string |  | New handle. Pass an empty string to clear it. Use snake_case (lowercase letters, digits, underscores); a hyphen is the subtraction operator wherever an identifier is parsed. |
| `layout` | string |  | Layout handle or ID (as string) that blocks NEW to the page are placed on. Only read when 'blocks' introduces a block. Omit it on a page with a single layout and that layout is used; on a page with several, omitting it is refused rather than guessed. Existing blocks are never re-placed. |
| `meta` | string |  | JSON object for page meta. Replaces existing meta. Example: {"allowPersonalLayouts":true} |
| `module` | string |  | Module name, handle, or ID (as string to prevent precision loss) for record detail pages. Pass an empty string to unlink the module. |
| `parent` | string |  | New parent page title, handle, or ID (as string to prevent precision loss). Pass an empty string to move the page to the root. |
| `title` | string |  | New title. An empty string is ignored — a page cannot be left without a title. |
| `visible` | boolean |  | Show page in navigation |

### `compose_record_create` {#compose_record_create}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/usage`.

Create a new record. If you do not know the field names, call compose_module_lookup first to get them.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `module` | string | Yes | Module name, handle, or ID (as string to prevent precision loss) |
| `namespace` | string | Yes | Namespace name, handle, slug, or ID (as string to prevent precision loss) |
| `values` | string | Yes | JSON object of field name to value: {"title":"Kickoff","attendees":42,"done":true}. A MULTI-VALUE field takes an array, and each element becomes one of the record's values in the order given: {"tags":["red","blue"]}. A field whose value is itself structured takes an object, stored as its JSON — a Geometry point is {"geo":{"coordinates":[46.05,14.51]}} (latitude first). Anything else is refused rather than guessed at. A field you leave out is stored as no value at all, which is a different state from a false or an empty one and does not match a query for it: omitting a Bool rather than sending false means "done = false" finds none of those records, and the same goes for a prefilter or a Metric block filter built on that field. Send every field a filter or chart will group on, including the false ones. |

### `compose_record_delete` {#compose_record_delete}

<Badge type="danger" text="Deletes" /> Listed under `/api/mcp/usage`.

Delete a record by ID. Requires a record ID — use compose_record_lookup with a filter to find it if unknown. The delete is soft: the record stops appearing in lookups but is retained, and compose_record_undelete brings it back.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `module` | string | Yes | Module name, handle, or ID (as string to prevent precision loss) |
| `namespace` | string | Yes | Namespace name, handle, slug, or ID (as string to prevent precision loss) |
| `recordID` | string | Yes | Record ID (as string to prevent precision loss) |

### `compose_record_draft` {#compose_record_draft}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/usage`.

Put a proposed record in front of the user as an editable form, without saving anything. Use this instead of compose_record_create or compose_record_update when the user should review or complete the values first — a new entry they will fill in, or a change they want to check. The user saves from the form; the draft itself writes nothing. With 'recordID' it proposes changes to that record: the fields named in 'values' replace the record's own and the rest are shown as they are. 'values' may be partial or omitted.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `module` | string | Yes | Module name, handle, or ID (as string to prevent precision loss) |
| `namespace` | string | Yes | Namespace name, handle, slug, or ID (as string to prevent precision loss) |
| `recordID` | string |  | Record ID to propose changes to (as string to prevent precision loss). Omit to draft a new record. |
| `values` | string |  | JSON object of field name to value: {"title":"Kickoff","attendees":42,"done":true}. A MULTI-VALUE field takes an array, and each element becomes one of the record's values in the order given: {"tags":["red","blue"]}. A field whose value is itself structured takes an object, stored as its JSON — a Geometry point is {"geo":{"coordinates":[46.05,14.51]}} (latitude first). Anything else is refused rather than guessed at. A field you leave out is stored as no value at all, which is a different state from a false or an empty one and does not match a query for it: omitting a Bool rather than sending false means "done = false" finds none of those records, and the same goes for a prefilter or a Metric block filter built on that field. Send every field a filter or chart will group on, including the false ones. |

### `compose_record_lookup` {#compose_record_lookup}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/usage`.

Look up a record by ID, or list and filter records in a module. Provide 'recordID' to fetch one; omit it and use 'filter' to search by field values. Do NOT call this before creating a record — only use it when the user explicitly asks to search or check for existing records. Records are returned in full, including their values, so use 'limit' and narrow with 'filter' rather than listing a whole module: a large module will exceed the result size limit and the call will fail. A Record or User value is stored as the bare ID of what it points at. The response carries 'refs', a dictionary of every such ID to the name a person would see — read the name from there rather than looking each reference up one by one.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `module` | string | Yes | Module name, handle, or ID (as string to prevent precision loss) |
| `namespace` | string | Yes | Namespace name, handle, slug, or ID (as string to prevent precision loss) |
| `filter` | string |  | Filter expression when no recordID is given, e.g. "name = 'John'" or "status = 'open'". Field names, comparisons, AND/OR — not SQL: there are no subqueries and no joins. A Record field holds the target's ID, so filter it by ID ("card = '510764704641581057'") after looking the target up; a path like "card.name = 'Bolt'" is resolved for you only for '=' and LIKE, and any other operator on a path is an error. A Select field is compared by the value it stores, never by the label shown in its place ("stage = 'offer'", not "stage = 'Offer'") — the wrong one is refused, with the legal values listed. Escape an apostrophe inside a literal with a BACKSLASH — "name = 'Urza\\'s Saga'". Doubling it the way SQL does is not an escape here and quietly matches nothing. |
| `limit` | string |  | Maximum records to return, default 50, capped at 200. |
| `pageCursor` | string |  | Cursor from a previous response, to fetch the next page. |
| `recordID` | string |  | Record ID (as string to prevent precision loss). Omit to list or filter instead. |
| `recordIDs` | string |  | Comma-separated record IDs to fetch in one call, e.g. "101,102,103" — use this instead of one call per ID, and instead of an IN expression in 'filter', which the query language does not support. Takes precedence over 'filter'. |
| `sort` | string |  | Order the records, e.g. "created_at DESC" or "stage, amount DESC". Names a field on this module, or one of recordID, createdAt, updatedAt, ownedBy — not a dotted path through a Record reference. Ordering happens in the database, so use it with 'limit' to ask for a top-N directly instead of paging the whole module and ordering the rows yourself. A 'pageCursor' carries the sort that produced it: pass the same expression or none. |

### `compose_record_report` {#compose_record_report}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/usage`.

Aggregate records server-side: sums, averages, extremes and counts, optionally grouped by a field. Use this for ANY question that is answered by a number over more than a handful of records — a total, an average, a count, a breakdown. Reading the records with compose_record_lookup and adding them up yourself is a wrong answer waiting to happen, and it silently drops everything past the page you were given. Returns {"rows": [...], "units": {...}}: one row per group, each with the metrics you asked for plus 'count', the number of records in that group. With no 'dimension' you get a single row for the whole set, where dimension_0 is "\*". 'units' gives the prefix or suffix the module puts on each aggregated field — use it when stating the figure, and do not supply a currency or unit it does not name. Grouping by a Record or User field returns 'refs', a dictionary of each group's ID to its name — read the labels from there rather than looking the groups up one by one. Do NOT list the dimension field in 'metrics': it comes back as dimension_0 on every row.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `module` | string | Yes | Module name, handle, or ID (as string to prevent precision loss) |
| `namespace` | string | Yes | Namespace name, handle, slug, or ID (as string to prevent precision loss) |
| `dimension` | string |  | A single field name ON THIS MODULE to group by, e.g. "rarity" — not a dotted path through a Record reference: "card.rarity" is not a field and the call fails. To group by something held on a referenced record, aggregate that module instead, or denormalise the field. Omit for one row covering every record. Date fields can be bucketed with DATE(field), and a chart's modifiers (QUARTER, YEAR) are not available here — group by the raw field and combine the rows yourself if you need coarser buckets. |
| `filter` | string |  | Filter expression narrowing which records are aggregated, e.g. "rarity = 'Mythic'". Same syntax as compose_record_lookup's filter. |
| `limit` | string |  | Keep only the first N groups after sorting. Every group is still aggregated; this trims the answer. |
| `metrics` | string |  | Comma-separated aggregate expressions over numeric fields, e.g. "SUM(line_value) AS total, AVG(price) AS avg_price". Functions: SUM, AVG, MIN, MAX, COUNT. The alias after AS is the key in the result; without one the expression itself is the key. 'count' is always returned and needs no metric. |
| `sort` | string |  | Order the groups by one result key, e.g. "total DESC" — the alias you gave a metric, or 'count', or 'dimension_0' for the group itself. Use it with 'limit' to ask for a top-N directly instead of ranking the groups yourself. |

### `compose_record_undelete` {#compose_record_undelete}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/usage`.

Restore a soft-deleted record, reversing compose_record_delete. Deleting a record in Human only marks it deleted, so nothing was lost and the record comes back with its values intact. Requires the record ID. A deleted record is left out of every filtered listing and compose_record_lookup offers no includeDeleted-style filter — but fetching it directly by 'recordID' with compose_record_lookup does still return it, with 'deletedAt' set, which is how you confirm you have the right record before restoring it. Failing that, use the ID compose_record_delete reported. This restores one record per call. It writes to the record — the revision counter moves and undelete automation runs — so do not call it speculatively on a record that is not deleted.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `module` | string | Yes | Module name, handle, or ID (as string to prevent precision loss). The module itself must not be deleted — restore it first with compose_module_undelete if it is. |
| `namespace` | string | Yes | Namespace name, handle, slug, or ID (as string to prevent precision loss) |
| `recordID` | string | Yes | Record ID of the deleted record (as string to prevent precision loss) |

### `compose_record_update` {#compose_record_update}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/usage`.

Update an existing record. Requires a record ID — use compose_record_lookup with a filter to find it if unknown. Send only the fields you are changing: every other field on the record keeps its value. A field you name is replaced entire, so a multi-value field takes the whole list it should end up with, and sending a field as an empty string clears it. Pass replace=true to write the value set whole instead, clearing every field you leave out.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `module` | string | Yes | Module name, handle, or ID (as string to prevent precision loss) |
| `namespace` | string | Yes | Namespace name, handle, slug, or ID (as string to prevent precision loss) |
| `recordID` | string | Yes | Record ID (as string to prevent precision loss) |
| `values` | string | Yes | JSON object of field name to value: {"title":"Kickoff","attendees":42,"done":true}. A MULTI-VALUE field takes an array, and each element becomes one of the record's values in the order given: {"tags":["red","blue"]}. A field whose value is itself structured takes an object, stored as its JSON — a Geometry point is {"geo":{"coordinates":[46.05,14.51]}} (latitude first). Anything else is refused rather than guessed at. A field you leave out is stored as no value at all, which is a different state from a false or an empty one and does not match a query for it: omitting a Bool rather than sending false means "done = false" finds none of those records, and the same goes for a prefilter or a Metric block filter built on that field. Send every field a filter or chart will group on, including the false ones. On update only the fields named here change; the rest of the record is left alone. A named field is replaced entire, so adding one value to a multi-value field still means reading the record and writing the full list back. |
| `replace` | boolean |  | Write the value set whole, clearing every field not named in 'values'. Off by default — leave it off unless you have read the record and are sending all of it. |

## System

Users, roles, groups, permissions, applications, agents, chatbots, reminders and other instance-wide settings.

### `system_agent_create` {#system_agent_create}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Create an AI agent. An agent's configuration is four objects — meta, behavior, execution and access. Send only the ones you are changing; each REPLACES that whole section, so read the agent first and send an edited copy of the section rather than a fragment of it. A new agent is 'active' unless you say otherwise, and with no tools in 'access' it answers from its prompt alone, reaching nothing — grant it something in the same call. Name a provider and a model in 'execution.model' — system_llm_provider_lookup lists both — unless this instance has exactly one LLM provider, which is resolved at run time and leaves 'execution.model' empty until then. To let people talk to it in a widget, create a chatbot with a conversation scenario pointing at this agent — system_chatbot_create.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `access` | string |  | JSON object: what the agent may reach. {"allow":[{"namespaceID":"&lt;id>","moduleIDs":[]}],"tools":[{"name":"compose_record_lookup","permission":"always","description":"Read leads","allow":[{"namespaceID":"&lt;id>","moduleIDs":["&lt;id>"]}]}],"taqs":[{"id":"&lt;taqID>","description":"Escalate"}],"workflows":[{"id":"&lt;workflowID>","description":"Notify"}]}. Access is deny-by-default: an agent whose 'tools' is EMPTY can call no tool at all — it still answers from its prompt and knowledge, but cannot look anything up or change anything. Granting a tool does not widen anything — an agent runs as the person who invoked it, so RBAC is the ceiling on every call and the agent can never do what that user could not. 'allow' on the access object itself scopes the whole agent to namespaces and modules however its tools were granted; that is the setting to use for "this agent is only for namespace X". Granting an agent system_agent_update lets it re-grant itself any tool, so treat that entry the way you would a permission change. Every name is checked against the registry when you write it, and an unknown one is refused with the near matches: the runtime resolves the allow-list as a whole, so a single typo would stop the agent running at all rather than cost it one tool. An entry's "moduleIDs" narrows it to those modules; leave it EMPTY to mean every module in that namespace, now and in future. Prefer empty unless you actually need to withhold a module: an enumerated list has to be edited on every tool entry each time a module is added, and until it is the agent cannot see the new module and nothing says so. Each entry may carry a "permission": "always" runs the tool unannounced, "ask" stops the run and puts it to the user (the exec call comes back with status "awaiting_approval" and the pending call; send it again with that tool in "approvedTools" to carry on), and "deny" refuses it whatever the scope says. Omit it and the mode follows the tool's risk — reading always, anything that writes asks. An entry names ONE tool via "name", or a whole set via "group" plus "maxRisk" — {"group":"usage","maxRisk":"read","allow":[{"namespaceID":"&lt;id>","moduleIDs":[]}]} grants every read-only data tool in that namespace and picks up tools added later. "group" is "usage" (data and execution) or "configuring" (schema and definitions); "maxRisk" is "read", "write" or "destructive" and defaults to "read". Set one or the other, never both. Prefer a group: naming tools one at a time is what makes an agent tedious to set up and stale afterwards. |
| `behavior` | string |  | JSON object: what the agent is told and what it may draw on. {"systemPrompt":"You triage support tickets.","guardrails":["no-pii"],"knowledgeBases":["&lt;knowledgeBaseID>"]}. The platform context — what namespaces, modules and records are, and how to work with them — is not a setting: the server adds it to the prompt of any agent granted a tool in access.tools. |
| `execution` | string |  | JSON object: which model runs it and how far it may go. {"model":{"llmProviderID":"&lt;id>","model":"claude-sonnet-5","temperature":0.2},"limits":{"maxIterations":10,"timeout":"5m","softLimitRatio":0.8,"contextWindow":200000,"outputTokens":8192}}. 'timeout' is a Go duration string. 'contextWindow' is NOT the model's window: it is a cumulative token budget for the whole conversation, checked after each iteration and across turns, and 'softLimitRatio' is the fraction of it at which the agent is told to wrap up. Omitting 'llmProviderID' is resolved when the agent RUNS, not when it is written: the instance's sole active provider is used, and where there are several the run fails naming them. So an agent created without one reads back with an empty 'model' object and is not necessarily broken — system_llm_provider_lookup says which case this instance is in, and is where an ID and a model name come from. Sending a 'temperature' makes the server call the provider to check it, which also validates 'llmProviderID' and the model name — an unrecognised model fails with the provider's own error. Omit temperature and none of that is checked: the model name is stored as given. |
| `handle` | string |  | URL-friendly identifier, unique across agents. Use snake_case (lowercase letters, digits, underscores). |
| `meta` | string |  | JSON object: the agent's labels. {"short":"Support triage","description":"Reads new tickets and files them","sidebarRoles":["&lt;roleID>"]}. 'short' is the name a person sees and is required — an agent has no separate name field. 'sidebarRoles' limits who sees it in the webapp sidebar and is not an authorization check. |
| `namespace` | string |  | Compose namespace this agent is for, by handle, slug or ID. Sets 'access.allow' to that namespace, so whatever tools you grant are confined to it. It does NOT grant any tool: an agent scoped here and granted nothing still reaches nothing. Ignored when 'access' already says something. It is not stored — what is stored is the scope it produced. |
| `status` | string |  | Lifecycle status; defaults to "active". |

### `system_agent_delete` {#system_agent_delete}

<Badge type="danger" text="Deletes" /> Listed under `/api/mcp/configuring`.

Delete an agent. The delete is soft: the configuration is retained and system_agent_undelete brings it back intact, but the agent stops appearing in lookups and stops being invocable. Note the agentID before calling — a deleted agent is excluded from every listing, so a handle no longer resolves and undelete needs the ID. A chatbot scenario pointing at this agent is NOT updated and will fail to find it, so check system_chatbot_lookup for scenarios naming it first.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `agent` | string | Yes | Agent handle, or ID as a string to prevent precision loss. |

### `system_agent_exec` {#system_agent_exec}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/usage`.

Run an agent and return what it said. This is how you find out whether an agent you just wrote actually works — storing one proves nothing, and most of what goes wrong with an agent (no model, a prompt that ignores its tools, a tool it was never granted) is invisible until it runs. 'input' is the message to send. Pass 'conversationID' from a previous response to carry on the same conversation; leave it out and a new one is started, and its ID is in the response. Read 'toolCalls' in the response, not just 'output'. An agent that answered plausibly having called nothing is the usual failure: it made the answer up from its prompt. An empty 'toolCalls' on a question that needed data is the tell. 'status' is "complete", or "awaiting_approval" when the run stopped to ask. In that case 'pendingApproval' names the tool it wants to use and the arguments it wants to use it with; send the same call again with that tool's name in 'approvedTools' to let it carry on. A tool the agent was granted with permission "ask" stops the run this way every time. The agent runs as YOU. Its tool allow-list narrows what it may do; it never widens it, so the agent cannot reach anything you could not reach yourself. An agent cannot call this tool to start another agent. Chaining agents is what a TAQ or workflow step is for — 'agentPrompt' in the construct library, 'agentRun' in the workflow registry.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `agent` | string | Yes | Agent handle, or ID as a string to prevent precision loss. |
| `approvedTools` | string |  | JSON array of tool names the agent may now use, e.g. ["compose_record_create"]. Send this with the same input and conversationID to resume a run that came back "awaiting_approval". An approval covers this call only. |
| `context` | string |  | JSON object of extra facts the agent's prompt can read, e.g. {"recordID":"511"}. Treated as data, not instructions. |
| `conversationID` | string |  | Conversation to continue, as a string. From a previous response. Omit to start a new one. |
| `input` | string |  | The message to send. Omit it only when resuming a run with 'approvedTools'. |

### `system_agent_lookup` {#system_agent_lookup}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/configuring`.

Read the AI agents configured on this instance — an agent is a model, a system prompt and an allow-list of tools it may call. Provide 'agent' to fetch one whole, with all five configuration sections; omit it to list them as {agentID, handle, short, status}. Call this before system_agent_update: an update replaces whole sections, so you need the current one to send back an edited copy. An agent is not a chatbot — a chatbot is the widget and the conversation, and it delegates to an agent; see system_chatbot_lookup.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `agent` | string |  | Agent handle, or ID as a string to prevent precision loss. Omit to list. An agent has no name to resolve by — 'short' is a label, not an identifier. |
| `handle` | string |  | Filter the listing to an exact handle. |
| `includeDeleted` | boolean |  | Include soft-deleted agents in the listing. They are excluded by default, which is why undelete needs an ID. |
| `limit` | string |  | Maximum agents to return when listing, default 50, capped at 200. |
| `pageCursor` | string |  | Cursor from a previous response, to fetch the next page. |
| `query` | string |  | Free-text search across the agents. |
| `status` | string |  | Filter the listing by status, e.g. "active". New agents are created active. |

### `system_agent_undelete` {#system_agent_undelete}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Restore a soft-deleted agent, reversing system_agent_delete. Nothing was erased, so the agent returns with its prompt, model, limits and tool allow-list exactly as they were. Requires the numeric agentID: a deleted agent is excluded from the default listing, so its handle resolves to nothing. Use the ID system_agent_delete reported, or call system_agent_lookup with includeDeleted first. Calling this on an agent that is not deleted is accepted and changes nothing.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `agentID` | string | Yes | ID of the deleted agent, as a string to prevent precision loss. A handle will not work. |

### `system_agent_update` {#system_agent_update}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Change an existing agent. An agent's configuration is four objects — meta, behavior, execution and access. Send only the ones you are changing; each REPLACES that whole section, so read the agent first and send an edited copy of the section rather than a fragment of it. An argument you omit is left unchanged, so this is safe to call with one section. What is NOT safe is sending a section you built from nothing: 'behavior' without the knowledge bases the agent had drops them, because the section replaces rather than merges. Call system_agent_lookup with 'agent' first. The revision number is the server's and increments on every update; it is not a parameter and sending one changes nothing.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `agent` | string | Yes | Agent handle, or ID as a string to prevent precision loss. |
| `access` | string |  | JSON object: what the agent may reach. {"allow":[{"namespaceID":"&lt;id>","moduleIDs":[]}],"tools":[{"name":"compose_record_lookup","permission":"always","description":"Read leads","allow":[{"namespaceID":"&lt;id>","moduleIDs":["&lt;id>"]}]}],"taqs":[{"id":"&lt;taqID>","description":"Escalate"}],"workflows":[{"id":"&lt;workflowID>","description":"Notify"}]}. Access is deny-by-default: an agent whose 'tools' is EMPTY can call no tool at all — it still answers from its prompt and knowledge, but cannot look anything up or change anything. Granting a tool does not widen anything — an agent runs as the person who invoked it, so RBAC is the ceiling on every call and the agent can never do what that user could not. 'allow' on the access object itself scopes the whole agent to namespaces and modules however its tools were granted; that is the setting to use for "this agent is only for namespace X". Granting an agent system_agent_update lets it re-grant itself any tool, so treat that entry the way you would a permission change. Every name is checked against the registry when you write it, and an unknown one is refused with the near matches: the runtime resolves the allow-list as a whole, so a single typo would stop the agent running at all rather than cost it one tool. An entry's "moduleIDs" narrows it to those modules; leave it EMPTY to mean every module in that namespace, now and in future. Prefer empty unless you actually need to withhold a module: an enumerated list has to be edited on every tool entry each time a module is added, and until it is the agent cannot see the new module and nothing says so. Each entry may carry a "permission": "always" runs the tool unannounced, "ask" stops the run and puts it to the user (the exec call comes back with status "awaiting_approval" and the pending call; send it again with that tool in "approvedTools" to carry on), and "deny" refuses it whatever the scope says. Omit it and the mode follows the tool's risk — reading always, anything that writes asks. An entry names ONE tool via "name", or a whole set via "group" plus "maxRisk" — {"group":"usage","maxRisk":"read","allow":[{"namespaceID":"&lt;id>","moduleIDs":[]}]} grants every read-only data tool in that namespace and picks up tools added later. "group" is "usage" (data and execution) or "configuring" (schema and definitions); "maxRisk" is "read", "write" or "destructive" and defaults to "read". Set one or the other, never both. Prefer a group: naming tools one at a time is what makes an agent tedious to set up and stale afterwards. |
| `behavior` | string |  | JSON object: what the agent is told and what it may draw on. {"systemPrompt":"You triage support tickets.","guardrails":["no-pii"],"knowledgeBases":["&lt;knowledgeBaseID>"]}. The platform context — what namespaces, modules and records are, and how to work with them — is not a setting: the server adds it to the prompt of any agent granted a tool in access.tools. |
| `execution` | string |  | JSON object: which model runs it and how far it may go. {"model":{"llmProviderID":"&lt;id>","model":"claude-sonnet-5","temperature":0.2},"limits":{"maxIterations":10,"timeout":"5m","softLimitRatio":0.8,"contextWindow":200000,"outputTokens":8192}}. 'timeout' is a Go duration string. 'contextWindow' is NOT the model's window: it is a cumulative token budget for the whole conversation, checked after each iteration and across turns, and 'softLimitRatio' is the fraction of it at which the agent is told to wrap up. Omitting 'llmProviderID' is resolved when the agent RUNS, not when it is written: the instance's sole active provider is used, and where there are several the run fails naming them. So an agent created without one reads back with an empty 'model' object and is not necessarily broken — system_llm_provider_lookup says which case this instance is in, and is where an ID and a model name come from. Sending a 'temperature' makes the server call the provider to check it, which also validates 'llmProviderID' and the model name — an unrecognised model fails with the provider's own error. Omit temperature and none of that is checked: the model name is stored as given. |
| `handle` | string |  | New handle. Pass an empty string to clear it. |
| `meta` | string |  | JSON object: the agent's labels. {"short":"Support triage","description":"Reads new tickets and files them","sidebarRoles":["&lt;roleID>"]}. 'short' is the name a person sees and is required — an agent has no separate name field. 'sidebarRoles' limits who sees it in the webapp sidebar and is not an authorization check. |
| `status` | string |  | New lifecycle status. |

### `system_application_create` {#system_application_create}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Create an application, adding an entry to this instance's app selector. 'name' is the administrative name used in listings; the 'unify' block is what users actually see, and an application with no unify URL has nothing to open. A new application is disabled and unlisted unless you say otherwise: pass enabled true and unify.listed true for it to reach users. It is also placed last in the ordering — this tool sets no weight, use system_application_reorder to position it. With unify.kind "custom" the application is its own HTML document rather than a link to a section: leave unify.url out and it is filled in with app/&lt;applicationID> once the ID exists, and system_application_source_set is what gives it a body. Call system_application_lookup on an existing application first if you are unsure what a working unify block looks like on this instance.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `name` | string | Yes | Administrative name of the application, shown in admin listings. |
| `enabled` | boolean |  | Whether the application is active. Defaults to false, which hides it from every user regardless of the unify block. |
| `unify` | string |  | App selector configuration, as a JSON object. Keys: "name" (label in the selector, falls back to the application name), "listed" (boolean, whether users see it), "url" (where it opens), "config" (free-form configuration string), "icon" and "logo" (URLs), "iconID" and "logoID" (attachment IDs as strings), "kind" ("custom" for an application whose UI is one HTML document, shown by the app view at /app/&lt;applicationID>; with no "url" of your own the server sets it to app/&lt;applicationID>). Example: {"name":"Reports","listed":true,"url":"/compose/ns/reports"} |

### `system_application_delete` {#system_application_delete}

<Badge type="danger" text="Deletes" /> Listed under `/api/mcp/configuring`.

Delete an application, removing it from the app selector for every user. The delete is soft: the record is kept, is visible again by passing includeDeleted to system_application_lookup, and is restored by system_application_undelete. Prefer system_application_update with enabled false when you only want to take the application out of circulation — that leaves it in the admin listing where someone can find it again.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `application` | string | Yes | Exact application name, or an application ID as a string to prevent precision loss. |

### `system_application_flag` {#system_application_flag}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Put a flag on an application. Flags are free-form labels the app selector reads, "pinned" being the usual one; check what an existing application carries with system_application_lookup before inventing a name. 'mode' is required and decides who the flag is for, so state it deliberately: "own" sets a flag only you see, while "global" writes shared state that every user sees and needs a separate, higher permission. You cannot flag on another user's behalf in either mode. Flagging something already flagged in that mode fails rather than doing nothing.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `application` | string | Yes | Exact application name, or an application ID as a string to prevent precision loss. |
| `flag` | string | Yes | Flag name to set, e.g. "pinned". |
| `mode` | string | Yes | "own" flags it for you alone; "global" flags it for every user and requires the global-flag permission. There is no default — pick one. |

### `system_application_lookup` {#system_application_lookup}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/configuring`.

Look up one application by name or ID, or list and filter applications. An application is an entry in this instance's app selector: an administrative name, an on/off switch, an ordering weight, and an 'unify' block holding the label, URL, icon, logo and configuration the selector renders. Provide 'application' to fetch one, which returns the full record including the unify block; omit it to list, which returns a compact form (applicationID, name, enabled, weight, flags, deletedAt) with the unify block left out because it is bulky and useless for picking one out of a list — fetch by name or ID when you need it. 'application' matches an exact name or an ID, never a partial name; use 'query' to search by fragment. Deleted applications are hidden unless you set 'includeDeleted', which is how you find the ID that system_application_undelete requires. Weights are listed in ascending order of appearance; call this before system_application_reorder to get the current order and the complete set of IDs.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `application` | string |  | Exact application name, or an application ID as a string to prevent precision loss. Omit to list instead of fetching one. |
| `flags` | string |  | List mode only. Return only applications carrying any of these flags. JSON array of strings, e.g. ["pinned"]. Matches your own flags plus global ones, exactly as the app selector sees them. |
| `includeDeleted` | boolean |  | Include soft-deleted applications alongside live ones. Deletes are reversible with system_application_undelete. |
| `limit` | string |  | Maximum results, default 50, capped at 200. |
| `pageCursor` | string |  | Cursor from a previous response, to fetch the next page. |
| `query` | string |  | List mode only. Case-insensitive substring match against the application name. |

### `system_application_reorder` {#system_application_reorder}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Set the order applications appear in in the app selector. The applications you list are weighted in the order given. Applications are a single flat list for the whole instance, so pass the COMPLETE ordered list: anything you leave out is pushed behind what you listed, keeping its current relative order. Call system_application_lookup first for the full set of applicationIDs and the current weights. This changes ordering only — nothing else about an application — and it is the only way to set weight, which the create and update tools deliberately do not expose.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `order` | string | Yes | JSON array of application IDs in the desired order, as strings to prevent precision loss, e.g. ["123","456","789"] |

### `system_application_source_get` {#system_application_source_get}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/configuring`.

Read a custom application's HTML back, exactly as stored, with the meta the app view reads instead of the document: hash, byte size, and the namespace and modules the app is allowed to query. Neither system_application_lookup nor the listing carries the source — this is the only way to see it. Call it before patching with system_application_source_set, so 'old_string' is copied from what is actually stored rather than from what you believe you wrote. Needs 'access' on the application, the same grant that lets a user open it.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `application` | string | Yes | Exact application name, or an application ID as a string to prevent precision loss. |

### `system_application_source_set` {#system_application_source_set}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Store the HTML of a custom application — the whole document the app view renders inside its sandbox. The application must already carry unify.kind "custom"; any other application is refused. Send it one of two ways: 'source' replaces the document outright, or 'old_string' with 'new_string' patches the stored one, which is how you edit an app without resending it. A patch is applied only when 'old_string' matches exactly once — zero or several matches are refused with the count, so include enough surrounding text to be unambiguous. 'namespace' and 'modules' declare the data the app may read: the bridge refuses every module not named here. A patch that sends neither keeps the declaration already stored. The document is plain HTML with inline script. JSX, ES modules and a React import are refused, and so are fetch, XMLHttpRequest and WebSocket — the sandbox's CSP blocks all of them, so an app using them fails in front of a user instead of here. Load the custom_app skill with system_skill_lookup before writing one.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `application` | string | Yes | Exact application name, or an application ID as a string to prevent precision loss. |
| `modules` | string |  | JSON array of module handles the app may query, as strings, e.g. ["Lead","Deal"]. This is the allowlist the bridge enforces; a module left out is refused at runtime. Omit in patch mode to keep the stored list. |
| `namespace` | string |  | Handle of the compose namespace the app reads from. Omit in patch mode to keep the stored one. |
| `new_string` | string |  | Patch mode: what old_string becomes. Pass an empty string to delete the matched text. Required whenever old_string is given. |
| `old_string` | string |  | Patch mode: the exact text to replace, whitespace included. It must occur exactly once in the stored source — read it with system_application_source_get first. |
| `source` | string |  | The whole HTML document, replacing whatever is stored. Capped at 256 KB; aim well under that. Cannot be combined with old_string. |
| `writes` | string |  | JSON array of the declared module handles the app may also create and change records in, e.g. ["Lead"]; every one must be in modules too. Leave it out for a read-only app. Each viewer is asked once, in Human, before the app's first change, and changes run with that viewer's own permissions. Omit in patch mode to keep the stored list. |

### `system_application_undelete` {#system_application_undelete}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Restore a soft-deleted application, putting it back in the app selector with all of its settings intact. This one takes an ID rather than a name, because name lookup only ever sees live applications: call system_application_lookup with includeDeleted set to find the applicationID of the deleted entry first.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `applicationID` | string | Yes | Application ID as a string, to prevent precision loss. Get it from system_application_lookup with includeDeleted set. |

### `system_application_unflag` {#system_application_unflag}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Take a flag off an application, reversing system_application_flag. 'mode' is required and must match the flag you are removing: "own" removes only the flag you set yourself and never touches a global one, while "global" removes the flag for every user and needs the global-flag permission. One quirk worth knowing: if an application carries a global flag and you have also set the same flag in "own" mode, unflagging in "own" mode suppresses the global flag for you alone and leaves everyone else's view untouched. Removing a flag that is not set fails rather than doing nothing.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `application` | string | Yes | Exact application name, or an application ID as a string to prevent precision loss. |
| `flag` | string | Yes | Flag name to remove, e.g. "pinned". |
| `mode` | string | Yes | "own" removes your own flag; "global" removes the shared flag for every user and requires the global-flag permission. There is no default — pick one. |

### `system_application_update` {#system_application_update}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Update an application. Omit a field to leave it unchanged. 'unify' is merged key by key: the keys you send are applied, the ones you leave out keep their current value, and sending a key with an empty value clears it — so you can flip listed without resending the URL. This tool never changes the ordering weight; use system_application_reorder for that. Setting unify.kind to "custom" on an application with no unify.url points the URL at app/&lt;applicationID>, which is where the app view serves the stored HTML. To take an application away from users without deleting it, set enabled false — that is reversible and keeps every setting.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `application` | string | Yes | Exact application name, or an application ID as a string to prevent precision loss. |
| `enabled` | boolean |  | Whether the application is active. Omit to leave unchanged. |
| `name` | string |  | New administrative name. Omit to leave unchanged; it cannot be set to an empty string. |
| `unify` | string |  | App selector configuration to merge in, as a JSON object. Same keys as system_application_create, "kind" included. Only the keys present are touched, so {"listed":false} hides the application from the selector and leaves the URL, icon and logo alone. |

### `system_auth_client_create` {#system_auth_client_create}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Create an OAuth2 auth client — the identity an external application, integration or script uses to obtain tokens from this instance. THIS IS THE ONE TOOL THAT HANDS BACK A CREDENTIAL. The server generates the client secret during creation and this returns it, once. Nothing reads it back afterwards: there is no tool that exposes or regenerates a secret, so a secret that is lost means creating a new client. Put it where it is going before you go on, and do not repeat it anywhere it does not need to be. Which grant to use decides what else you must supply. "authorization_code" is the flow where a person signs in and is sent back to 'redirectURI', which is then required. "client_credentials" has no person in it: the client authenticates as itself and acts as one nominated user, so 'impersonateUser' is required and everything the client does is done with that user's permissions. Choosing a user with more access than you have is how a client ends up more powerful than the person who made it — pick the narrowest account that can do the job. 'scope' is what the token may be used for and is almost always "profile api": the API middleware requires the 'api' scope, so a client without it authenticates and is then refused by every endpoint. The client is enabled on creation unless you say otherwise. 'validFrom' and 'expiresAt' bound when it works at all and are the clean way to issue a credential that stops working on its own. Nothing here sets the client's permitted, prohibited or forced roles; those are a person's job in the admin UI, and system_auth_client_update does not touch them either.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `name` | string | Yes | Display name, shown in listings. An auth client must have one. |
| `description` | string |  | What this client is for and who runs it. Worth writing: a credential nobody can attribute is a credential nobody dares revoke. |
| `enabled` | boolean |  | Whether the client may obtain tokens. Defaults to true. Pass false to create it dormant and enable it later with system_auth_client_update. |
| `expiresAt` | string |  | Client cannot be used after this time. RFC3339. Omit for a credential that never expires on its own. |
| `handle` | string |  | Short URL-safe identifier used in OAuth2 requests, e.g. "reporting_service". Optional, but a client without one can only ever be referenced by its numeric ID. |
| `impersonateUser` | string |  | The user a client_credentials client acts as: user ID as a string (to prevent precision loss), handle or email. REQUIRED when validGrant is client_credentials, and rejected by the server otherwise. Everything the client does carries this user's permissions. |
| `redirectURI` | string |  | Absolute URL the person is sent back to after authorizing. Required in practice for authorization_code; meaningless for client_credentials. |
| `scope` | string |  | Space-separated OAuth2 scopes the client may request. Use "profile api" unless you have a reason not to — without the 'api' scope every API call is refused after a successful login. |
| `trusted` | boolean |  | A trusted client skips the consent screen a person is otherwise shown when authorizing it. Set it only for clients this instance itself owns. |
| `validFrom` | string |  | Client cannot be used before this time. RFC3339, e.g. 2026-08-03T09:00:00Z. Omit for no lower bound. |
| `validGrant` | string |  | OAuth2 grant type: "authorization_code" (a person signs in; needs redirectURI) or "client_credentials" (no person; needs impersonateUser). Omit to leave it unset, which means the client cannot obtain a token until system_auth_client_update sets one. |

### `system_auth_client_delete` {#system_auth_client_delete}

<Badge type="danger" text="Deletes" /> Listed under `/api/mcp/configuring`.

Delete an auth client. Every application, integration or script authenticating with it stops being able to obtain tokens, so establish what uses it before calling. The delete is soft: the client is retained with its secret intact, is listed again by passing includeDeleted to system_auth_client_lookup, and is restored by system_auth_client_undelete. If you only want to stop the client working while keeping it obviously present, set enabled=false with system_auth_client_update instead. The instance's default auth client cannot be deleted.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `authClient` | string | Yes | Auth client ID as a string, to prevent precision loss, or the client's handle. |

### `system_auth_client_lookup` {#system_auth_client_lookup}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/configuring`.

Look up an OAuth2 auth client by ID or handle, or list and filter auth clients. An auth client is the identity an external application, integration or script uses to obtain a token from this instance. Provide 'authClient' to fetch one, which returns its full configuration; omit it to list, which returns a compact form (ID, handle, name, enabled, isDefault, validFrom, expiresAt, deletedAt). No response from this tool ever includes a client secret. A client's secret is shown once, by system_auth_client_create, and there is no tool that reads it back or regenerates it — a client whose secret has been lost is replaced rather than recovered, or its secret is read by a person in the admin UI. Soft-deleted clients are hidden unless you set 'includeDeleted'.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `authClient` | string |  | Auth client ID as a string, to prevent precision loss, or the client's handle. Omit to list instead. |
| `handle` | string |  | List only the client whose handle is exactly this. Ignored when 'authClient' is given; unlike a lookup by handle, an unknown value returns an empty list rather than an error. |
| `includeDeleted` | boolean |  | Include soft-deleted auth clients. They are hidden by default, so this is how you find one to restore with system_auth_client_undelete. |
| `limit` | string |  | Maximum results, default 50, capped at 200. |
| `pageCursor` | string |  | Cursor from a previous response, to fetch the next page. |

### `system_auth_client_undelete` {#system_auth_client_undelete}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Restore a soft-deleted auth client. It becomes usable again with the same ID, handle and secret as before, so applications configured against it start working without being reconfigured. Find the client first with system_auth_client_lookup and includeDeleted set — a deleted client is invisible to a plain lookup. A handle works here even though a handle search normally skips deleted clients: this tool searches them deliberately.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `authClient` | string | Yes | Auth client ID as a string, to prevent precision loss, or the client's handle. |

### `system_auth_client_update` {#system_auth_client_update}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Update an auth client's configuration: its handle, name and description, whether it is enabled and trusted, its OAuth2 grant type, redirect URI, scope and validity window. Omit a field to leave it unchanged; pass it empty to clear it. This cannot read, set or regenerate the client secret, so an update never disturbs the credential the client is already using. A secret is shown once, by system_auth_client_create, and never again. It also leaves the client's security settings alone (impersonated user, permitted, prohibited and forced roles): those grant privileges and are set when the client is created or by a person in the admin UI, which is why setting 'validGrant' to client_credentials is rejected unless an impersonation user is already configured. The instance's default auth client — the one 'isDefault' marks in a lookup — cannot have its handle changed and cannot be disabled.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `authClient` | string | Yes | Auth client ID as a string, to prevent precision loss, or the client's handle. |
| `description` | string |  | New description of what the client is for. Empty clears it. |
| `enabled` | boolean |  | Whether the client may be used to obtain tokens at all. Disabling is the reversible alternative to deleting; the default client cannot be disabled. |
| `expiresAt` | string |  | Client is unusable after this time. RFC3339. Empty removes the expiry. |
| `handle` | string |  | New handle: a short URL-safe identifier used in OAuth2 requests. Empty clears it. Cannot be changed on the default client. |
| `name` | string |  | New display name. Cannot be cleared: an auth client must have a name, so an empty value is rejected. |
| `redirectURI` | string |  | Absolute URL the user is sent back to after authorizing. Empty clears it. |
| `scope` | string |  | Space-separated OAuth2 scopes the client may request, e.g. "profile api". Empty clears it. |
| `trusted` | boolean |  | Trusted clients skip the consent step users are otherwise shown when authorizing. Grant this only to clients this instance owns. |
| `validFrom` | string |  | Client is unusable before this time. RFC3339, e.g. 2026-08-03T09:00:00Z. Empty removes the restriction. |
| `validGrant` | string |  | OAuth2 grant type, typically "authorization_code" or "client_credentials". Empty clears it. client_credentials needs security settings this tool does not set. |

### `system_chatbot_create` {#system_chatbot_create}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Create a chatbot. A chatbot's configuration is three objects — handoff, styling and scenarios. Send only the ones you are changing; each REPLACES that whole section, so read the chatbot first and send an edited copy rather than a fragment. Two things decide whether it does anything. It needs a 'conversation' scenario naming an agentID, or it has nothing to answer with. And it needs 'allowedOrigins', because that list is what actually restricts where the widget may run — a chatbot with none is usable from any site that has its widget key. The widget key is generated here and returned; you cannot choose it.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `name` | string | Yes | Human-readable name, shown in the admin list. |
| `allowedOrigins` | string |  | JSON array of origins the widget may be embedded on, e.g. ["https://example.com"]. This is the access control on the widget, not decoration: an empty list restricts nothing. |
| `enabled` | boolean |  | Whether the widget answers. A disabled chatbot keeps its configuration and serves nobody. |
| `handle` | string |  | URL-friendly identifier, unique across chatbots. Use snake_case (lowercase letters, digits, underscores). |
| `handoff` | string |  | JSON object: passing the conversation to a person. {"enabled":true,"automation":{"onRequested":{"taqID":"&lt;id>"},"onAccepted":{"taqID":"&lt;id>"}}}. The hooks fire when a visitor asks for a human and when one picks the conversation up. |
| `scenarios` | string |  | JSON array: what the chatbot can do, in order. [{"id":"triage","name":"Triage","type":"conversation","agentID":"&lt;agentID>","automation":{"before":{},"after":{}}}]. A scenario of type "conversation" MUST carry an agentID — the whole save is rejected otherwise, not just that scenario. Find the agent with system_agent_lookup. This is a collection and always replaces: send every scenario the chatbot should have, not only the new one. |
| `sessionTTL` | string |  | How long a visitor's session survives, as a duration string, e.g. "30m". |
| `styling` | string |  | JSON object: how the widget looks. {"logoAttachmentID":"&lt;id>","logoURL":"https://…","fontFamily":"Inter","fontSizes":{},"colors":{},"launcher":{}}. An attachment must already exist — it is uploaded through the REST endpoint /chatbots/{id}/upload-asset, and MCP has no way to transfer a file. |

### `system_chatbot_delete` {#system_chatbot_delete}

<Badge type="danger" text="Deletes" /> Listed under `/api/mcp/configuring`.

Delete a chatbot. The delete is soft: the configuration is retained and system_chatbot_undelete brings it back, but the widget stops answering everywhere it is embedded. Note the chatbotID before calling — a deleted chatbot is excluded from every listing, so its handle no longer resolves and undelete needs the ID. The agent a scenario pointed at is untouched; deleting the widget does not delete what answered through it.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `chatbot` | string | Yes | Chatbot handle, or ID as a string to prevent precision loss. |

### `system_chatbot_lookup` {#system_chatbot_lookup}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/configuring`.

Read the chatbots configured on this instance — a chatbot is the embeddable widget and the conversation around it, and it delegates the thinking to an agent through a scenario. Provide 'chatbot' to fetch one whole, with its handoff, styling and scenarios; omit it to list them as {chatbotID, handle, name, enabled}. Call this before system_chatbot_update: an update replaces whole sections. A chatbot is not an agent — see system_agent_lookup for the model and prompt behind it. Conversation transcripts are deliberately not reachable here.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `chatbot` | string |  | Chatbot handle, or ID as a string to prevent precision loss. Omit to list. A chatbot's name is not resolvable — only its handle is. |
| `handle` | string |  | Filter the listing to an exact handle. |
| `includeDeleted` | boolean |  | Include soft-deleted chatbots in the listing. They are excluded by default, which is why undelete needs an ID. |
| `limit` | string |  | Maximum chatbots to return when listing, default 50, capped at 200. |
| `pageCursor` | string |  | Cursor from a previous response, to fetch the next page. |
| `query` | string |  | Free-text search across the chatbots. |

### `system_chatbot_regenerate_key` {#system_chatbot_regenerate_key}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Issue the chatbot a new widget key. The previous key stops working the moment this returns, so every page embedding the widget breaks until its embed code is updated with the new key — which this returns. Use it when a key has leaked somewhere it should not be, or when retiring an embed you no longer control. It is not part of ordinary editing: nothing else about the chatbot changes, and an update never rotates the key on its own. The key identifies the widget rather than authorising it — 'allowedOrigins' is what restricts where it may run, so rotating a key on a chatbot that allows every origin buys very little.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `chatbot` | string | Yes | Chatbot handle, or ID as a string to prevent precision loss. |

### `system_chatbot_undelete` {#system_chatbot_undelete}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Restore a soft-deleted chatbot, reversing system_chatbot_delete. Nothing was erased, so it returns with its scenarios, styling and the same widget key — pages already embedding it start working again without being changed. Requires the numeric chatbotID: a deleted chatbot is excluded from the default listing, so its handle resolves to nothing. Use the ID system_chatbot_delete reported, or call system_chatbot_lookup with includeDeleted first. Calling this on a chatbot that is not deleted is accepted and changes nothing.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `chatbotID` | string | Yes | ID of the deleted chatbot, as a string to prevent precision loss. A handle will not work. |

### `system_chatbot_update` {#system_chatbot_update}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Change an existing chatbot. A chatbot's configuration is three objects — handoff, styling and scenarios. Send only the ones you are changing; each REPLACES that whole section, so read the chatbot first and send an edited copy rather than a fragment. An argument you omit is left unchanged, so this is safe to call with one section. 'scenarios' is a collection and always replaces: send every scenario the chatbot should end up with, because sending only the new one deletes the rest. Call system_chatbot_lookup with 'chatbot' first. A 'conversation' scenario missing its agentID rejects the whole call. The widget key cannot be set here — it is carried forward untouched. Rotate it with system_chatbot_regenerate_key.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `chatbot` | string | Yes | Chatbot handle, or ID as a string to prevent precision loss. |
| `allowedOrigins` | string |  | JSON array of origins the widget may be embedded on. Replaces the list wholesale; an empty array removes every restriction. |
| `enabled` | boolean |  | Whether the widget answers. |
| `handle` | string |  | New handle. Pass an empty string to clear it. |
| `handoff` | string |  | JSON object: passing the conversation to a person. {"enabled":true,"automation":{"onRequested":{"taqID":"&lt;id>"},"onAccepted":{"taqID":"&lt;id>"}}}. The hooks fire when a visitor asks for a human and when one picks the conversation up. |
| `name` | string |  | New name. An empty string is ignored — a chatbot without one is unidentifiable in the admin list. |
| `scenarios` | string |  | JSON array: what the chatbot can do, in order. [{"id":"triage","name":"Triage","type":"conversation","agentID":"&lt;agentID>","automation":{"before":{},"after":{}}}]. A scenario of type "conversation" MUST carry an agentID — the whole save is rejected otherwise, not just that scenario. Find the agent with system_agent_lookup. This is a collection and always replaces: send every scenario the chatbot should have, not only the new one. |
| `sessionTTL` | string |  | How long a visitor's session survives, as a duration string, e.g. "30m". |
| `styling` | string |  | JSON object: how the widget looks. {"logoAttachmentID":"&lt;id>","logoURL":"https://…","fontFamily":"Inter","fontSizes":{},"colors":{},"launcher":{}}. An attachment must already exist — it is uploaded through the REST endpoint /chatbots/{id}/upload-asset, and MCP has no way to transfer a file. |

### `system_llm_provider_lookup` {#system_llm_provider_lookup}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/configuring`.

List the LLM providers configured on this instance, or fetch one. A provider is the connection to a model vendor — its type (anthropic, openai, mistral, …), its endpoint and its default model — and it is what 'execution.model.llmProviderID' on an agent names. Read this before system_agent_create or system_agent_update. The ID is not guessable and there is no other way to discover one, so an agent gets its provider from here. An agent that names none runs on the instance's sole active provider; where there is more than one, an unnamed provider is an agent that fails when it is first run rather than when it is written — this is where you find out which case you are in. Set 'models' to fetch the model names the provider itself advertises, which is where 'execution.model.model' comes from. That calls the vendor over the network with the stored credential, so it is slower than the rest of this tool and fails when the credential is wrong — which is itself worth knowing before an agent is built on it. 'models' needs a single provider named. No response here contains an API key. A provider carries only the ID of the credential it uses, and nothing on this surface reads a credential. Providers are created and their keys set by a person in the admin UI; there is no tool for either.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `limit` | string |  | Maximum results, default 50, capped at 200. |
| `llmProvider` | string |  | Provider ID as a string (to prevent precision loss), or its handle. Omit to list instead. |
| `models` | boolean |  | Also return the model names this provider advertises, fetched from the vendor. Requires 'llmProvider'. Slower, and it surfaces a bad credential as an error rather than a silent failure later. |
| `pageCursor` | string |  | Cursor from a previous response, to fetch the next page. |
| `provider` | string |  | List only providers of this vendor type, e.g. "anthropic", "openai", "mistral". Ignored when 'llmProvider' is given. |
| `status` | string |  | List only providers in this state, typically "active". A provider that is not active is not eligible to run an agent. Ignored when 'llmProvider' is given. |

### `system_permission_grant` {#system_permission_grant}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Give a role permissions. This is the second half of access control: system_role_create makes an empty container and system_role_member_add puts people in it, but until a rule is granted the role allows nothing. YOU CAN ONLY GRANT WHAT YOU HOLD. Every rule is checked against your own access on that exact resource, evaluated the same way the server evaluates it when you act yourself. An operation you cannot perform is refused, by name, and the whole call is rejected rather than partly applied. This is structural, not a policy: there is no argument, no role and no resource that lets this tool hand out more than the person driving it already has. You also need the 'grant' permission on the component — holding an operation is not the same as being allowed to delegate it. Get the resource strings and operation names from system_permission_schema; they are not guessable and a wrong one is refused rather than stored. Check what you can give with system_permission_lookup, which returns 'yourAccess' for a resource. A rule replaces any existing rule for the same role, resource and operation. Rules for resources you did not name are untouched — unlike system_role_clone_rules, which replaces a role's whole rule set. Effect is immediate for new sessions, but an already signed-in user keeps the access their session was built with until they sign in again. 'deny' is a rule like any other and outranks an allow inherited from anywhere else, so use it deliberately. To take a rule away rather than override it, use system_permission_revoke.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `role` | string | Yes | Role to grant to: ID as a string (to prevent precision loss), handle or name. Create one with system_role_create. |
| `rules` | string | Yes | JSON array of rules. One rule: {"resource":"corteza::compose:module/511/512","operation":"record.create","access":"allow"}<br>resource is the full RBAC resource string, exactly as system_permission_schema prints it — component, resource type and every path segment. A short form like "module/\*" is refused by the server, and so is a path with the wrong number of segments: a module is "corteza::compose:module/&lt;namespaceID>/&lt;moduleID>", a record is "corteza::compose:record/&lt;namespaceID>/&lt;moduleID>/&lt;recordID>". A component-wide rule has an empty path and a trailing slash: "corteza::compose/".<br>A path segment may be "\*" to cover every resource at that level; a rule written that way is only accepted if you hold the operation at that same breadth.<br>operation is one of the operations system_permission_schema lists for that resource type. Operations are per resource type, not global: "record.create" lives on the module, not on the record.<br>access is "allow" or "deny". Omit it and it is "allow". To remove a rule, use system_permission_revoke — writing "inherit" here is refused, so that taking access away is never something this tool does by accident.<br>Rules for different components may be mixed in one call; they are grouped and applied per component. |

### `system_permission_lookup` {#system_permission_lookup}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/configuring`.

Read the permission rules a role holds, and what you yourself may grant. Pass 'role' — an ID, handle or name — for that role's rules; omit it to sweep every role you can see, which is how you find out who already has access to something. 'resource' narrows to one resource string, and also accepts a prefix: 'corteza::compose:module/511' returns the rules on every module in that namespace. A rule's 'access' is 'allow' or 'deny'. A permission with no rule at all is not listed: it inherits, which for most things means denied. So an empty result means the role holds nothing here, not that something failed. When you pass 'resource', the answer also carries 'yourAccess' — the operations you personally hold on it. That list is exactly the ceiling on what system_permission_grant will let you give away, so reading it first turns a refusal into a decision. Reading rules needs the 'grant' permission on the component the rules belong to; a component you cannot manage is left out of the answer and named under 'componentsSkipped' rather than silently omitted.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `limit` | string |  | Maximum rules returned, default 50, capped at 200. |
| `resource` | string |  | Full RBAC resource string, or a prefix of one. Use system_permission_schema for the shapes. Omit for every resource. |
| `role` | string |  | Role ID as a string (to prevent precision loss), handle or name. Omit to cover every role you can see. |

### `system_permission_revoke` {#system_permission_revoke}

<Badge type="danger" text="Deletes" /> Listed under `/api/mcp/configuring`.

Take permission rules off a role. The rule is deleted rather than set to 'deny': the permission goes back to inheriting, which for most things means the role no longer allows it, but a rule inherited from elsewhere — a wildcard rule, or the Authenticated role — takes over again. If you need the operation blocked outright, grant 'deny' with system_permission_grant instead. THE SAME CEILING APPLIES AS ON GRANT: you can only revoke an operation you yourself hold on that resource, and you need the component's 'grant' permission. Revoking is a privilege change like any other and is not a way around the boundary. Read the role's rules with system_permission_lookup first — revoking a rule that was never there succeeds and changes nothing, which reads as if it worked. Nothing here deletes a role or removes a member; that is system_role_delete and system_role_member_remove.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `role` | string | Yes | Role to revoke from: ID as a string (to prevent precision loss), handle or name. |
| `rules` | string | Yes | JSON array of rules to remove: [{"resource":"corteza::compose:module/511/512","operation":"record.create"}]. An 'access' key is ignored here — revoking always means back to inherit. Resource and operation must match the granted rule exactly; system_permission_lookup prints them. |

### `system_permission_schema` {#system_permission_schema}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/configuring`.

List every permission that can be granted: each resource type, the shape of its resource string, and the operations valid on it. This is a fixed reference list built into Human — read-only, the same on every instance, and not something a caller adds to. Read it before system_permission_grant or system_permission_revoke. Both take a full RBAC resource string and an operation, neither is guessable, and a wrong one is refused by the server rather than stored — so this is the only place to get them right. Two mistakes it prevents: a shortened resource ('module/\*' instead of 'corteza::compose:module/\*/\*'), which the server rejects outright; and an operation looked for on the wrong resource type — 'record.create' is an operation on the module, not on the record, and letting a role into an application needs 'access' and 'read' together, on the application. The 'resource' field of each entry is the wildcard form covering every resource of that type. Replace the '\*' segments with real IDs to scope a rule to one thing: 'corteza::compose:module/511/512' is one module, 'corteza::compose:module/511/\*' is every module in namespace 511. Around 40 resource types across three components, returned whole by default; 'component' and 'resourceType' narrow it.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `component` | string |  | Narrow to one component: "system", "compose" or "automation". Accepts the full form ("corteza::compose") too. Omit for all three. |
| `resourceType` | string |  | Narrow to one resource type. Either the full form ("corteza::compose:module") or the bare tail ("module"). Case-insensitive. A value matching nothing answers with every known type rather than an empty result. |

### `system_reminder_create` {#system_reminder_create}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/usage`.

Create a reminder. 'resource' identifies what the reminder is about and 'payload' carries whatever the consumer needs to render it — both are free-form, so mirror an existing reminder's shape rather than inventing one; call system_reminder_lookup first if unsure. Assigning to a user other than yourself requires the assign-reminder permission and fails without it.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `resource` | string | Yes | Resource the reminder refers to, e.g. "compose:record/1/2/3". |
| `assignedTo` | string |  | Assignee user ID as a string. Defaults to you; assigning to anyone else needs permission. |
| `payload` | string |  | The reminder's contents, as a JSON object or as the JSON text of one. Invalid JSON is refused. |
| `remindAt` | string |  | When to surface the reminder. RFC3339. Omit for a reminder with no schedule. |

### `system_reminder_delete` {#system_reminder_delete}

<Badge type="danger" text="Deletes" /> Listed under `/api/mcp/usage`.

Delete a reminder. The delete is soft — the record is retained and can be seen again by passing includeDeleted to system_reminder_lookup — but there is no undelete tool, because the reminder service exposes no undelete operation. Treat it as final. If you only want to stop a reminder surfacing, use system_reminder_dismiss instead: that is reversible.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `reminderID` | string | Yes | Reminder ID as a string. |

### `system_reminder_dismiss` {#system_reminder_dismiss}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/usage`.

Dismiss a reminder so it stops surfacing, keeping the record. Reversible with system_reminder_undismiss. Prefer this over system_reminder_delete when the user is done with a reminder rather than wanting it gone.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `reminderID` | string | Yes | Reminder ID as a string. |

### `system_reminder_lookup` {#system_reminder_lookup}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/usage`.

Look up a reminder by ID, or list and filter reminders. Provide 'reminderID' to fetch one, which returns the full reminder including its payload; omit it to list, which returns a compact form (ID, resource, remindAt, assignedTo, dismissedAt) — fetch by ID when you need the payload. You only ever see reminders assigned to you: the service filters by assignee regardless of what you pass in 'assignedTo', so this tool cannot read another user's reminders. Dismissed reminders are included unless you set 'excludeDismissed'.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `assignedTo` | string |  | Assignee user ID as a string. Reads are scoped to you regardless; this narrows, it cannot widen. |
| `excludeDismissed` | boolean |  | Omit reminders that have been dismissed. |
| `includeDeleted` | boolean |  | Include soft-deleted reminders. Deletes are soft and there is no undelete tool, so this is the only way to see them. |
| `limit` | string |  | Maximum results, default 50, capped at 200. |
| `pageCursor` | string |  | Cursor from a previous response, to fetch the next page. |
| `reminderID` | string |  | Reminder ID as a string, to prevent precision loss. Omit to list instead. |
| `resource` | string |  | Filter by the resource string the reminder is attached to, e.g. "compose:record". |
| `scheduledFrom` | string |  | Only reminders due at or after this time. RFC3339, e.g. 2026-08-03T09:00:00Z. |
| `scheduledOnly` | boolean |  | Return only reminders that have a remindAt set. |
| `scheduledUntil` | string |  | Only reminders due at or before this time. RFC3339. |

### `system_reminder_snooze` {#system_reminder_snooze}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/usage`.

Push a reminder's due time out and record that it was snoozed, incrementing its snooze count. Use this rather than system_reminder_update when the user is deferring a reminder — the count is what tells you a reminder keeps being put off.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `remindAt` | string | Yes | New due time, RFC3339. |
| `reminderID` | string | Yes | Reminder ID as a string. |

### `system_reminder_undismiss` {#system_reminder_undismiss}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/usage`.

Reverse a dismissal so the reminder surfaces again. Use this when a user changes their mind about a reminder they dismissed, or dismissed by accident. Has no effect on a reminder that was never dismissed, and cannot recover a deleted one — system_reminder_delete has no undo.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `reminderID` | string | Yes | Reminder ID as a string. |

### `system_reminder_update` {#system_reminder_update}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/usage`.

Update a reminder. Omit a field to leave it unchanged; pass it empty to clear it. 'payload' replaces the stored object wholesale rather than merging into it. To change only the schedule, prefer system_reminder_snooze — it also records that the reminder was snoozed, which this tool does not.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `reminderID` | string | Yes | Reminder ID as a string. |
| `payload` | string |  | Replacement contents, as a JSON object or as the JSON text of one. Empty clears it. Not merged. |
| `remindAt` | string |  | New due time, RFC3339. Empty clears the schedule. |
| `resource` | string |  | New resource string. Empty clears it. |

### `system_role_archive` {#system_role_archive}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Archive a role: it stops being applied, so its members stop receiving what it grants, and it drops out of default listings — but nothing is removed and system_role_unarchive puts it straight back. Use this to retire a role you may want later. Archiving is not deleting: system_role_delete marks the role deleted and hides it from everything, archiving parks it. It is also not suspending a user — it affects every member of the role at once, so to stop one person holding the role use system_role_member_remove. System roles are refused.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `role` | string | Yes | Role ID as a string, handle, or name. |

### `system_role_clone_rules` {#system_role_clone_rules}

<Badge type="danger" text="Deletes" /> Listed under `/api/mcp/configuring`.

Copy every permission rule from one role onto another, replacing the target's rules. This is the widest-reaching write in this tool surface: the target ends up with exactly the source's access, so every member of the target immediately holds everything members of the source hold. Never use it to 'start from' a powerful role, and confirm with a human before pointing it at any role you do not fully understand. It is destructive for the target as well: the target's existing rules are discarded, not merged with the source's, and there is no undo. The source role is left untouched. Requires permission to manage permissions across the instance, which is a higher bar than editing a role. To review either side first, read both roles with system_role_lookup.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `role` | string | Yes | Source role to copy permission rules FROM. ID as a string, handle, or name. Unchanged by this call. |
| `targetRole` | string | Yes | Target role to copy permission rules ONTO. ID as a string, handle, or name. Its own rules are discarded first. |

### `system_role_create` {#system_role_create}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Create a role. This is the first half of granting access: a role is a permission container, and every user or user group later added to it receives everything the role is permitted to do. Creating one requires the create-role permission and is recorded in the action log. A new role starts with no permissions and no members, so on its own it grants nothing — it becomes real when permissions are assigned to it (system_role_clone_rules copies another role's entire permission set onto it) and members are added with system_role_member_add. 'handle' is the stable identifier other tools accept in place of an ID; a role without one can only be referenced by ID or name. Handle and name must each be unique.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `name` | string | Yes | Human-readable name, e.g. "Support agents". Must be unique. |
| `description` | string |  | What the role is for. Free text; say what it grants, since nothing else records that. |
| `handle` | string |  | Stable identifier, e.g. "support-agents". Letters, digits, dash and underscore, at least 2 characters. Must be unique. |

### `system_role_delete` {#system_role_delete}

<Badge type="danger" text="Deletes" /> Listed under `/api/mcp/configuring`.

Delete a role. Everyone who held access through this role loses it: a deleted role stops being applied, so every permission it granted stops taking effect for every one of its members at once. The delete is soft — the role and its permission rules are retained, system_role_lookup with includeDeleted lists it, and system_role_undelete restores it with its rules intact. System roles (bypass, authenticated, anonymous) are refused. If you want the role out of use but expect to bring it back, prefer system_role_archive. If you want one user to stop holding the role, use system_role_member_remove rather than deleting the role for everybody.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `role` | string | Yes | Role ID as a string, handle, or name. |

### `system_role_lookup` {#system_role_lookup}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/configuring`.

Look up a role by ID, handle or name, or list and filter roles. Provide 'role' to fetch one, which returns the whole role plus isSystem and isClosed; omit it to list, which returns a compact form (roleID, name, handle, isSystem, isClosed, archivedAt, deletedAt). Read this before any role write. A role is a permission container: whatever it is granted is held by everyone in it, so knowing which role you are about to touch matters. isSystem marks the roles the instance itself runs on — the bypass, authenticated and anonymous roles — which refuse delete, archive and undelete and cannot be renamed; isClosed marks roles whose membership cannot be changed. Use system_role_member_list to see who is in a role, and the 'member' filter here to see which roles a user holds. Deleted and archived roles are excluded unless you pass includeDeleted or includeArchived, and a deleted role can then only be fetched by ID, because handle and name lookups skip deleted roles.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `handle` | string |  | Exact handle match. |
| `includeArchived` | boolean |  | Include archived roles, which are otherwise hidden. |
| `includeDeleted` | boolean |  | Include soft-deleted roles. This is how you find the roleID that system_role_undelete needs. |
| `limit` | string |  | Maximum results, default 50, capped at 200. |
| `member` | string |  | Only roles this user is a member of. User ID as a string, handle or email. Cannot be combined with userGroup. |
| `name` | string |  | Exact name match. |
| `pageCursor` | string |  | Cursor from a previous response, to fetch the next page. |
| `query` | string |  | Free-text match against handle and name. |
| `role` | string |  | Role ID as a string, handle, or name. Omit to list instead. |
| `userGroup` | string |  | Only roles this user group is a member of. Group ID as a string or handle. Cannot be combined with member. |

### `system_role_member_add` {#system_role_member_add}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Add a user or a user group to a role, granting that user — or every member of that group — everything the role is permitted to do. This is a privilege change, not bookkeeping: read the role first and do not add anyone to a role whose permissions you have not inspected. Pass exactly one of 'user' or 'userGroup'. Requires permission to manage members on the role. Closed roles and context roles refuse membership changes; system_role_lookup shows isClosed. Reversible with system_role_member_remove.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `role` | string | Yes | Role ID as a string, handle, or name. This is the role whose permissions the member will hold. |
| `user` | string |  | User ID as a string, handle, or email address. Pass this or userGroup, not both. |
| `userGroup` | string |  | User group ID as a string, or handle. Adds the whole group, so every current and future member of it holds the role. |

### `system_role_member_list` {#system_role_member_list}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/configuring`.

List the members of a role — the users and user groups that hold whatever the role grants. Each entry is a kind (user or userGroup) and an ID; resolve those with the user or user group tools if you need names. Returns every member in one response with no paging, so call it for one role at a time. Requires read access to the role, and is refused for closed roles and for context roles, whose membership is computed from an expression rather than stored. To go the other way — which roles one user holds — list roles with system_role_lookup and its 'member' filter.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `role` | string | Yes | Role ID as a string, handle, or name. |

### `system_role_member_remove` {#system_role_member_remove}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Remove a user or a user group from a role, withdrawing everything that role granted them. Any other role they hold is unaffected, so this does not necessarily leave them without access — check system_role_lookup with the 'member' filter if that is the goal. Pass exactly one of 'user' or 'userGroup'. Removing a user does not remove them from a user group that is itself a member of the role: if their access comes through the group, remove the group or remove them from it. Requires permission to manage members on the role; closed and context roles refuse membership changes. Reversible with system_role_member_add.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `role` | string | Yes | Role ID as a string, handle, or name. |
| `user` | string |  | User ID as a string, handle, or email address. Pass this or userGroup, not both. |
| `userGroup` | string |  | User group ID as a string, or handle. Removes the whole group's membership of the role. |

### `system_role_unarchive` {#system_role_unarchive}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Return an archived role to use. Its permissions start applying to its members again, so this widens access by exactly as much as the role grants — read the role with system_role_lookup before calling it. Has no effect on a role that was never archived, and cannot recover a deleted one: use system_role_undelete for that. System roles are refused.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `role` | string | Yes | Role ID as a string, handle, or name. Archived roles need includeArchived to show up in system_role_lookup. |

### `system_role_undelete` {#system_role_undelete}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Restore a soft-deleted role, and with it every permission that role granted to its members. This re-grants access — it does not merely make the role visible again — so check what the role is before restoring it. Takes 'roleID' rather than a handle or name because a deleted role is only reachable by ID: handle and name lookups skip deleted roles. Get the ID from system_role_lookup with includeDeleted. System roles are refused.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `roleID` | string | Yes | Role ID as a string. Handles and names do not resolve to deleted roles. |

### `system_role_update` {#system_role_update}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Rename a role or change its description. Handle and name are what other tools and stored configuration use to refer to a role, so changing them can break anything that names the old value. System roles are refused a rename: the bypass, authenticated and anonymous roles reject any change to handle or name, and only a description-only edit goes through. Check isSystem with system_role_lookup first rather than calling this and retrying. Omit a field to leave it unchanged; pass it empty to clear it. This tool changes neither the role's permissions nor its members — use system_role_clone_rules for permissions and system_role_member_add or system_role_member_remove for members.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `role` | string | Yes | Role ID as a string, handle, or name. |
| `description` | string |  | New description. Empty clears it. Other role metadata, including any context expression, is left untouched. |
| `handle` | string |  | New handle. Empty clears it, which leaves the role referenceable only by ID or name. Must stay unique. |
| `name` | string |  | New name. Empty clears it. Must stay unique. |

### `system_skill_lookup` {#system_skill_lookup}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/configuring` and `/api/mcp/usage`.

Read Human's own guidance on how to use its tools — the rules a tool's parameter list cannot state, written per topic rather than per tool. Call it with 'tool' BEFORE the first call to an unfamiliar tool: it returns the skills that apply to that tool, and skipping it is how a call succeeds and is quietly wrong. A page block added without reading page_layout, for instance, is stored on the page and placed on no layout, so nothing renders it. Pass 'skill' to read one whole, or neither argument to list what exists as {name, description, tools}. This is guidance, not configuration: nothing here changes the instance, and it is the same text for every caller. For a single tool's own full documentation, use human_tool_load instead — the two are complementary, and a skill spans the several tools a job actually takes.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `skill` | string |  | Name of one skill to read whole, e.g. "page_layout". Omit to list them. |
| `tool` | string |  | Tool name, e.g. "compose_page_update". Returns every skill that applies to that tool, each with its full text. A tool with no skill returns an empty list, which is not an error. |

### `system_theme_lookup` {#system_theme_lookup}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/configuring`.

Read the webapp's theme colours. Returns each theme — light, dark, and general — with its colours as a plain object, which is what system_theme_update expects back. Call this before updating: an update merges into what is already there, so seeing the current palette is how you avoid changing a colour you did not mean to. Colour names: primary, secondary, success, warning, danger, black, white, light, extra-light, body-bg, sidebar-bg, topbar-bg. Values are hex, e.g. "#09344E".

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `theme` | string |  | Which theme to read: "light", "dark", or "general". Omit for all of them. |

### `system_theme_update` {#system_theme_update}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Change the webapp's theme colours. Colours are merged, not replaced: pass only what should change and the rest of the palette is left alone. This restyles the webapp for EVERY user of this instance, not just the caller, and takes effect on their next page load — it is an instance-wide change, so confirm the colours with the user before calling. Light and dark are separate palettes and neither inherits from the other: a colour changed in one is unchanged in the other, and changing only 'light' leaves anyone using dark mode seeing the old palette. Colour names: primary, secondary, success, warning, danger, black, white, light, extra-light, body-bg, sidebar-bg, topbar-bg. Values are hex, e.g. "#09344E".

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `colors` | string | Yes | JSON object of colour name to hex value, e.g. {"primary":"#09344E","body-bg":"#F4F4F5"}. Only the names given are changed. |
| `theme` | string | Yes | Which theme to change: "light", "dark", or "general". |

### `system_user_create` {#system_user_create}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Create a user account. Only the email address is required; when you omit 'handle' one is derived from the name and email. The account is created with its email marked confirmed and with no password, so the person cannot sign in until a credential exists — setting one is deliberately not available through this tool surface, so send them through the normal sign-in or password-reset flow. To put the new user in a user group afterwards call system_user_group_member_add rather than system_user_update: only that call also refreshes the permission graph.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `email` | string | Yes | Email address. Must be a valid address and unique among users that are not deleted. This is what the person signs in with. |
| `handle` | string |  | URL-safe handle, letters/digits/._- starting with a letter, at least 2 characters. Omit to have one generated from the name and email. |
| `name` | string |  | Display name, e.g. "Ada Lovelace". Optional but used to derive the handle and the avatar initials. |

### `system_user_delete` {#system_user_delete}

<Badge type="danger" text="Deletes" /> Listed under `/api/mcp/configuring`.

Delete a user account. The delete is soft — the record is kept and system_user_undelete restores it — but the account disappears from every listing and its access tokens are revoked immediately, so treat this as removing the person. If you only want to stop someone signing in while keeping the account visible and easily reversible, use system_user_suspend instead. Before deleting, note that undelete needs the numeric userID, which afterwards is only obtainable from system_user_lookup with includeDeleted set.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `user` | string | Yes | User reference: numeric ID as a string, a handle, or an email address. |

### `system_user_group_create` {#system_user_group_create}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Create a user group. Groups form a hierarchy that permissions flow down, so where you attach a new group decides what its members will be able to do. Every group except the instance's single root group must name at least one parent in 'parents'; creating one without a parent is rejected as an invalid parent reference. Call system_user_group_lookup first to find the parent you want. This creates an empty group — members are added afterwards, one at a time, with system_user_group_member_add.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `handle` | string | Yes | Unique handle, at least 2 characters, starting with a letter and ending alphanumeric; letters, digits, underscore, hyphen and dot in between. This is how the group is referenced everywhere. |
| `description` | string |  | Longer description of what the group is for. |
| `parents` | string |  | JSON array of parent groups, e.g. ["engineering"] or [{"parent":"engineering","name":"reporting-line"}]. Each 'parent' is a group handle or ID as a string. Give every entry a distinct 'name' when there is more than one parent — duplicate or missing names are rejected. Omit only for the root group. |
| `short` | string |  | Short human label for the group, shown where a handle would read badly. |

### `system_user_group_delete` {#system_user_group_delete}

<Badge type="danger" text="Deletes" /> Listed under `/api/mcp/configuring`.

Delete a user group. The delete is soft and reversible with system_user_group_undelete, but it takes the group out of the permission graph straight away, so everyone in it loses whatever that group and its ancestors granted them. Users are not deleted and still point at the group, so undeleting restores their permissions. Child groups are not deleted either and are left pointing at a group that is no longer in the graph — re-parent them first with system_user_group_update. A group used as the security user group of an auth client cannot be deleted; that is reported as a permission error, so check the auth clients before assuming you lack rights.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `userGroup` | string | Yes | User group handle, or ID as a string to prevent precision loss. |

### `system_user_group_lookup` {#system_user_group_lookup}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/configuring`.

Look up a user group by ID or handle, or list and filter user groups. User groups are the hierarchy that permissions are granted through: a group inherits from its parents, and every member of a group gets what the group and its ancestors are allowed to do. Provide 'userGroup' to fetch one group in full, including its parent links, description and labels; omit it to list groups as {userGroupID, handle, short, isRoot, parentIDs, archivedAt, deletedAt}. A user group has no name field — the handle is its identifier and 'short' is the only human label it carries. Members are never returned here, by either mode; call system_user_group_member_list for those. To find which group a given user is in, read the 'userGroupID' field on the user instead — this tool cannot filter by member. Archived and deleted are two distinct states and both are excluded by default, so a group you archived or deleted will not appear, and will not resolve by handle either, until you set 'includeArchived' or 'includeDeleted'.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `includeArchived` | boolean |  | Include archived groups. Archiving is separate from deleting; a group can be one, both or neither. |
| `includeDeleted` | boolean |  | Include soft-deleted groups. Needed to find a group before calling system_user_group_undelete. |
| `limit` | string |  | Maximum groups to return when listing, default 50, capped at 200. |
| `pageCursor` | string |  | Cursor from a previous response, to fetch the next page. |
| `query` | string |  | Substring match against the handle. Only the handle is searched — descriptions are not indexed. |
| `userGroup` | string |  | User group handle, or ID as a string to prevent precision loss. Omit to list groups instead. |

### `system_user_group_member_add` {#system_user_group_member_add}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Add a user to a user group. This changes what that user is allowed to do: group membership is mirrored into the permission graph, so the user immediately gains everything this group and its ancestors grant. Treat it as a privilege change and confirm the group is the one intended before calling. A user belongs to exactly one group at a time, so this moves them: whatever group they were in before, they are no longer in, and they lose the permissions that came with it. There is no companion removal tool because the service exposes no removal operation — to take a user out of a group, add them to a different one.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `user` | string | Yes | User handle, email, or ID as a string to prevent precision loss. |
| `userGroup` | string | Yes | Target user group handle, or ID as a string to prevent precision loss. |

### `system_user_group_member_list` {#system_user_group_member_list}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/configuring`.

List the users who are members of one user group, as {userID, handle, name}. Membership is deliberately kept out of system_user_group_lookup because it is the heavy part of a group, so this is the only way to read it. Direct members only — users in child groups are not included, even though they inherit the group's permissions. This returns every member in one response and cannot be paged, so a very large group may exceed the result size limit. Email addresses and other user detail are not returned; look the user up by ID if you need them.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `userGroup` | string | Yes | User group handle, or ID as a string to prevent precision loss. |

### `system_user_group_undelete` {#system_user_group_undelete}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Restore a soft-deleted user group and put it back into the permission graph, which restores what its members could do. Find the group first with system_user_group_lookup and 'includeDeleted' set — a deleted group does not resolve by handle or ID without it. This is the way back from system_user_group_delete; it has no effect on a group that was never deleted.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `userGroup` | string | Yes | User group handle, or ID as a string to prevent precision loss. |

### `system_user_group_update` {#system_user_group_update}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Update a user group's handle, labels or parent links. Omit a field to leave it unchanged; pass it as an empty string to clear it. 'parents' replaces the existing parent links wholesale rather than merging, so send the full set you want. Re-parenting is not cosmetic: the group moves in the permission graph and its members immediately gain or lose whatever the old and new ancestors granted. This cannot update the root group — the service requires at least one parent on every update and the root group has none, so an update of it fails with a missing-parent error whatever you send. Members are not touched here; use system_user_group_member_add.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `userGroup` | string | Yes | User group handle, or ID as a string to prevent precision loss. |
| `description` | string |  | New description. Empty clears it. |
| `handle` | string |  | New handle. Must stay unique across all groups. |
| `parents` | string |  | Replacement JSON array of parent groups, e.g. ["engineering"] or [{"parent":"engineering","name":"reporting-line"}]. Replaces every existing parent link. Omit to leave the current parents alone. |
| `short` | string |  | New short label. Empty clears it. |

### `system_user_lookup` {#system_user_lookup}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/configuring`.

Look up one user by ID, handle or email address, or list and filter the user accounts on this Human instance. Provide 'user' to fetch one, which returns the full record including meta and labels; omit it to list, which returns a compact form (userID, email, handle, name, kind, suspendedAt, deletedAt) — fetch by ref when you need more than that. Listing hides deleted and suspended accounts unless you ask for them, and returns only ordinary users unless you set 'kind' or 'allKinds'. This is also the only way to obtain the numeric userID that system_user_undelete requires.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `allKinds` | boolean |  | Return ordinary and system accounts together, ignoring 'kind'. |
| `email` | string |  | Exact email address. Unlike 'query' this must match in full. |
| `handle` | string |  | Exact handle. Unlike 'query' this must match in full. |
| `includeDeleted` | boolean |  | Also return soft-deleted users. Off by default. Set this to find the numeric userID that system_user_undelete needs. |
| `includeSuspended` | boolean |  | Also return suspended users. Off by default, so a suspended account is invisible to a plain list. |
| `kind` | string |  | User kind. Empty (the default) means ordinary user accounts; "sys" means built-in system accounts. Those are the only two values. |
| `limit` | string |  | Maximum results, default 50, capped at 200. |
| `pageCursor` | string |  | Cursor from a previous response, to fetch the next page. |
| `query` | string |  | Free-text search, matched as a substring against email, username, handle and name. |
| `role` | string |  | Only users who are members of this role. Accepts a role ID as a string, a handle, or a role name. |
| `user` | string |  | User reference: numeric ID as a string, a handle, or an email address — anything containing '@' is read as an email. Usernames are not resolved. Omit to list instead. |
| `userGroup` | string |  | Only users belonging to this user group. Accepts a group ID as a string, a handle, or a group name. |
| `username` | string |  | Exact username. A legacy field that is empty on most accounts; prefer 'handle' or 'email'. |

### `system_user_set_email_confirmed` {#system_user_set_email_confirmed}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Mark a user's email address as confirmed, or withdraw that confirmation. This flips the stored flag only: it neither sends a confirmation mail nor validates one, and it does not change the address itself — use system_user_update for that. Use it to unblock someone stuck in the sign-up confirmation flow, or to force a re-confirmation after an address is corrected. Accounts made with system_user_create are already confirmed, so a fresh user rarely needs this.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `confirmed` | boolean | Yes | True marks the address confirmed, false marks it unconfirmed. Must be given explicitly. |
| `user` | string | Yes | User reference: numeric ID as a string, a handle, or an email address. |

### `system_user_suspend` {#system_user_suspend}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Suspend a user account: the person can no longer sign in and their existing access tokens are revoked, but the account, its data, and its group and role membership are all kept. Reversible with system_user_unsuspend. Prefer this over system_user_delete whenever the person may come back or the account merely needs freezing — suspend is about access, delete is about removal. A suspended account drops out of system_user_lookup listings unless you pass includeSuspended, though it stays resolvable by ID, handle or email.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `user` | string | Yes | User reference: numeric ID as a string, a handle, or an email address. |

### `system_user_undelete` {#system_user_undelete}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Restore a soft-deleted user, reversing system_user_delete. This takes the numeric user ID and nothing else: handle and email lookups do not resolve deleted accounts, so a reference that worked before the delete will fail here. Get the ID from system_user_lookup with includeDeleted set. The restore is refused if the email, handle or username of the deleted account now collides with a live user, or if it would take the instance over its user limit.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `userID` | string | Yes | Numeric user ID as a string, to prevent precision loss. Obtain it from system_user_lookup with includeDeleted set. |

### `system_user_unsuspend` {#system_user_unsuspend}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Lift a suspension so the user can sign in again, reversing system_user_suspend. Accepts the same references as suspend, because a suspended account is still resolvable by ID, handle or email. This cannot recover a deleted account — use system_user_undelete for that — and it fails if reactivating the user would take the instance over its configured user limit.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `user` | string | Yes | User reference: numeric ID as a string, a handle, or an email address. |

### `system_user_update` {#system_user_update}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Update a user's email address, display name or handle. Omit a field to leave it unchanged; pass 'name' or 'handle' as an empty string to clear them. 'email' cannot be cleared and must stay a valid, unique address: it is the identity the person signs in with and the address notifications go to, so changing it changes both. It does not send or require a new confirmation mail — use system_user_set_email_confirmed if you need the confirmation flag changed. Kind, avatar and user group are not editable here; to move someone between groups call system_user_group_member_add, which also refreshes the permission graph.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `user` | string | Yes | User reference: numeric ID as a string, a handle, or an email address. |
| `email` | string |  | New email address. Must be valid and unique. Cannot be cleared. |
| `handle` | string |  | New handle. Empty clears it. Must be unique among users that are not deleted. |
| `name` | string |  | New display name. Empty clears it. |

## Automation

TAQs, workflows, triggers and the catalogues they are built from.

### `automation_event_type_lookup` {#automation_event_type_lookup}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/configuring`.

List the automation events Human can react to: every valid eventType and resourceType pair, with the properties each event carries and the constraints a trigger may narrow it by. This is a fixed reference list built into Human — read-only, the same on every instance, and not something a caller can add to. Read it before calling automation_trigger_create or automation_trigger_update: those need an eventType and resourceType that exist together, this is the only place they are written down, and a pair that is not in this list is stored without complaint and then never fires. It is also how you find what a workflow will receive when the event fires — the 'properties' of an entry are the variables put into the workflow's scope — and which constraint names the workflow editor will accept for that event. The whole catalogue is around 120 entries and is returned in full when you pass nothing, which is usually what you want. Pass 'resourceType' to narrow it: an exact value like "compose:record" returns just that resource's events, while a bare prefix like "compose" returns that resource and all of its sub-resources. There is no paging and no per-entry fetch — one call returns everything there is to know.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `resourceType` | string |  | Narrow to one resource, e.g. "compose:record" for exactly that one or "compose" for it and every compose sub-resource. Case-insensitive. Omit for the whole catalogue. |

### `automation_taq_construct_lookup` {#automation_taq_construct_lookup}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/configuring`.

List the construct library — the vocabulary a TAQ (Trigger Action Query) is built from. Two catalogues in one call: 'functions', every ref a TAQ step of kind function or iterator may name, each with its parameters, their exact type names and which are required; and 'triggers', every legal resourceType/eventType pair a TAQ trigger may use. There are around 19 functions and 22 triggers, so the default is to return both whole. Read this before automation_taq_create or automation_taq_update. A step whose 'ref' is not in this catalogue is still stored and still returns success — it comes back with a 'function.unknown' issue, 'runnable' false, and automation_taq_exec then answers 'manager: executable not found'. The same goes for a trigger pair that is not listed: it stores without complaint and never fires. This is NOT the workflow function registry. A TAQ and a workflow are different engines with different vocabularies, and most of the 93 workflow functions do not exist here. For a workflow use automation_workflow_function_lookup instead — mixing the two is the most common way to author an automation that stores cleanly and never runs. Bind a step's arguments by 'argumentName', spelled exactly as listed here, and give each one a 'type' spelled exactly as it appears in that parameter's 'types' array. Note the workflow registry keys its parameters by 'name' instead, which is another reason a ref copied across engines does not work. The webapp's form layout for each entry ('segments') is omitted; nothing an author needs is in it and it is four fifths of the payload.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `catalogue` | string |  | Return one catalogue only: "functions" for step refs, "triggers" for resourceType/eventType pairs. Omit for both. |
| `ref` | string |  | Narrow the functions to one. An exact ref is matched first, then any ref containing it, case-insensitively. A ref that matches nothing answers with the full list of refs rather than an empty result. Ignored when catalogue is "triggers". |
| `resourceType` | string |  | Narrow the triggers to one resource, e.g. "compose:record" for exactly that or "compose" for it and every sub-resource. Case-insensitive. Ignored when catalogue is "functions". |

### `automation_taq_create` {#automation_taq_create}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Create a TAQ (Trigger Action Query automation) — its triggers, steps and paths in one call. A TAQ is a graph: triggers say what starts it, steps say what it does, paths connect them. Find out what a step can do before writing one. automation_taq_construct_lookup lists every function a step's 'ref' may name, each with its parameters, their exact type names and whether they are required, and every resourceType/eventType pair a trigger may use. Do not use automation_workflow_function_lookup — that is the workflow function registry and most of its entries are unavailable to a TAQ. The smallest TAQ that runs is one trigger, one function step, and no paths at all: a lone unconnected trigger and a lone unconnected step are wired together for you, and termination is added automatically. Read the result, do not assume it worked. A TAQ whose graph fails validation is still stored and still succeeds, and says what is wrong under 'issues' — a code, a severity, and the exact step, parameter and legal types. Any issue at all, of any severity, stops the TAQ being registered with the runtime: 'runnable' comes back false and automation_taq_exec answers 'manager: executable not found', which means the TAQ has issues and not that the ID is wrong. No 'issues' key at all is the clean result. Fix issues with automation_taq_update. 'enabled' defaults to true, which arms the TAQ immediately — its triggers begin listening for real events. Pass false to author without arming, but a disabled TAQ also refuses automation_taq_exec, so you cannot test one until you enable it. Storing is not proof of working. Run it with automation_taq_exec and then read automation_taq_execution_trace, because exec reports status only: a 'completed' status with no working frames in the trace — nothing, or only frames of kind 'trigger' and 'termination' — is a failure, it means no step ran. A termination frame records the flow reaching that end, which is how you tell which branch arm ran, but it is not work done. A frame carries the arguments it resolved under 'args' and its results under 'output', so reading 'args' is how you confirm an expression evaluated to what you meant — an iterator's resolved 'query' is right there, and a filter that came out wrong is visible at a glance. It proves the step ran with those arguments and nothing more, so where the step has a visible effect, check for that effect separately. An iterator that matched no records still records its own frame; only the absence of body frames after it shows that it looped zero times.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `handle` | string | Yes | URL-friendly identifier, unique among TAQs in the same project. Use snake_case (lowercase letters, digits, underscores); a hyphen is the subtraction operator wherever an identifier is parsed. This is what automation_taq_lookup searches and what automation_taq_exec resolves, so a TAQ without one can only ever be reached by its numeric ID. |
| `name` | string | Yes | Human-readable name, shown in listings. |
| `description` | string |  | What this TAQ is for. |
| `enabled` | boolean |  | Whether the TAQ is live. Defaults to true. False leaves its triggers unregistered AND makes automation_taq_exec refuse it, so it cannot be tested while disabled. |
| `paths` | string |  | JSON array of edges: [{"parentID":"&lt;trigger or step ID>","childID":"&lt;step ID>"}]. IDs are strings of digits.<br>Send [] for the common case. With exactly one unconnected trigger and exactly one unconnected step, the two are wired together for you; more than one of either and you get a graph.ambiguousEntry issue instead.<br>Order carries meaning, not presentation. A non-gateway step's second outbound path is its error handler rather than a second successor; an iterator's paths are its body first, then its exit.<br>A gatewayExclusive or gatewayInclusive step needs at least two outbound paths: each carries a "condition", except exactly one which carries none and is the else branch, always tested last. A condition is an AST object, not an expression string: {"ref":"eq","args":[{"symbol":"foo"},{"value":{"@type":"String","@value":"bar"}}]}. Operators: and, or, not, isNull, isNotNull, eq, ne, lt, gt, lte, gte. Nothing checks a condition when it is written — a wrong operator, symbol or scope surfaces only when the TAQ runs. |
| `steps` | string |  | JSON array of steps — what the TAQ does. One step:<br>{"stepID":"1","handle":"notify","kind":"function","ref":"notificationSend","meta":{"short":"Notify the owner"},"arguments":[{"argumentName":"recipient","value":"507566326668132353","type":"ID"},{"argumentName":"title","value":"Hello","type":"String"}]}<br>stepID is a string of digits you choose, unique across steps AND triggers; paths reference steps by it.<br>kind is one of: function, iterator, gatewayExclusive, gatewayInclusive, termination, error. Those are the only six, and they are NOT the workflow step kinds — a TAQ has no "expressions" step. A kind outside the six is rejected before anything is written.<br>ref, on a function or an iterator, names a construct-library function. List them with automation_taq_construct_lookup — each with its parameters, their types and whether they are required; it is the authority on which refs exist. Do NOT use automation_workflow_function_lookup: that is the workflow function registry, and most of its entries do not exist here. A ref that is not in the construct library is stored anyway and comes back as a "function.unknown" issue that names the near matches.<br>arguments bind by argumentName, which must equal one of that function's parameter argumentNames exactly; every required parameter must be present. "type" is compared literally against that parameter's "types" array and must be spelled as listed there (ID, String, Integer, Boolean, ComposeRecord, …) — omitting it is the same as sending "" and fails. Supply a literal with "value", copy an earlier step's output with "source" plus "scope" (the producing step's handle), or compute one with "expr".<br>An "expr" is evaluated INSIDE its own "scope", and without one that is the trigger's: so in an iterator body the record variable is still the trigger's record, and reaching the loop item means giving the argument the iterator's handle as "scope". One expression sees one scope, so a formula needing both the trigger's record and the loop item cannot be written as a single expr — split it across arguments, or derive the value with a module field value expression instead (compose_module_update, "fields[].expressions.value" — its scope is the record's own fields by bare name, documented on that parameter and in the field_expressions skill).<br>In an expression a Record-ref field holds the target's ID as a string, not a record: a filter over them reads "card = " + record.values.card (string concatenation). Writing record.values.card.recordID yields nothing, so the query matches nothing and an iterator over it runs ZERO times in silence — no error, no issue, and the execution still reports completed. The resolved query is recorded in that step's frame under "args", which is where a filter that came out wrong becomes obvious.<br>A step's results cannot be set: they are derived from the function definition and anything you send is overwritten. An error step's message likewise cannot be set through the API.<br>A termination step is optional — every leaf step is wired to an auto-injected one. |
| `triggers` | string |  | JSON array of triggers — what starts the TAQ. One trigger:<br>{"triggerID":"1001","handle":"onAgent","enabled":true,"resourceType":"automation:trigger:agentic","eventType":"onAgentic","meta":{"short":"Invoked by an agent"},"inputSchema":[{"name":"subject","type":"String","required":true}]}<br>resourceType and eventType are a pair from the trigger catalogue: automation_taq_construct_lookup with catalogue "triggers", 22 entries. "automation:trigger:agentic" with "onAgentic" is the pair automation_taq_exec runs, so include one of those if you want to be able to run the TAQ on demand; the compose:record and system:user pairs fire on real events instead.<br>triggerID is a string of digits. Omit it on create and one is minted for you and returned. Supply it when a path references the trigger, and echo the existing triggerID on update — a trigger sent without one is replaced by a new trigger with a new ID, orphaning any path that pointed at the old one.<br>handle names the trigger: it is what automation_taq_exec takes as entryPoint, and what a step reads from as a scope.<br>inputSchema declares the parameters automation_taq_exec accepts under "input"; it is the only place those names are written down. |

### `automation_taq_delete` {#automation_taq_delete}

<Badge type="danger" text="Deletes" /> Listed under `/api/mcp/configuring`.

Delete a TAQ (Trigger Action Query automation). The delete is soft — its triggers, steps and paths are all retained — and automation_taq_undelete restores the whole definition. Deleting stops the TAQ running: its triggers are unregistered, so nothing it was listening for fires it any more. Prefer disabling a TAQ over deleting it when you only want to pause it, because a disabled TAQ still appears in automation_taq_lookup with includeDisabled, whereas a deleted one cannot be found at all — not by handle, not by search — and restoring it needs an ID you must have kept.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `taq` | string | Yes | TAQ ID as a string (to prevent precision loss), or handle — the same reference automation_taq_lookup and automation_taq_update take. |
| `taqID` | string |  | Accepted as an alias for 'taq'. |

### `automation_taq_exec` {#automation_taq_exec}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/usage`.

Run a TAQ and wait for it to finish. Returns the execution result — executionID, status, error, timings — not the values the TAQ produced; pass the returned executionID to automation_taq_execution_trace to see each step's input and output. This runs the automation for real, doing whatever its steps do. There is no dry run, so read the TAQ with automation_taq_lookup first if you are unsure what it will do. 'entryPoint' names the trigger to start from, by trigger handle; omitted, the TAQ's first trigger is used. 'input' must satisfy that trigger's input schema, which automation_taq_lookup returns: keys are matched to the schema case-insensitively, and a missing required parameter fails the call before the TAQ starts.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `taq` | string | Yes | TAQ ID as a string (to prevent precision loss), or handle. |
| `entryPoint` | string |  | Handle of the trigger to start from. Defaults to the TAQ's first trigger. |
| `input` | string |  | JSON object of input parameters, shaped by the entry-point trigger's input schema. Omit when the trigger takes no input. |

### `automation_taq_execution_trace` {#automation_taq_execution_trace}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/usage`.

Get the step-by-step trace of one TAQ run: a stack frame per executed step under 'frames', with its handle, kind, arguments, input, output, timings and error. This is the tool for diagnosing why a run failed or produced what it did. Any step that recorded an error is also listed under 'failedSteps', with the exact message, so a failure nested in a long frame is not something you have to scroll for. A failed step does mark the run failed — but the action that triggered it still succeeds: a record write that fires a failing automation saves and returns normally, so the failure reaches nobody unless someone reads the executions or this trace. Get the executionID from automation_taq_executions, or from the result of automation_taq_exec. A long run traces every step it took, so a trace can exceed the result size limit.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `executionID` | string | Yes | Execution ID as a string (to prevent precision loss), from automation_taq_executions or automation_taq_exec. |
| `taq` | string | Yes | TAQ ID as a string (to prevent precision loss), or handle. |

### `automation_taq_executions` {#automation_taq_executions}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/usage`.

List past runs of one TAQ: executionID, status, any error, start and end time and duration. Use it to check whether a TAQ ran and whether it succeeded, then pass an executionID to automation_taq_execution_trace for the step-by-step detail of one run. Every retained execution of the TAQ is returned — the service takes no filter and no paging — so a heavily used TAQ can exceed the result size limit. Newest first, so the run you just caused is the first entry.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `taq` | string | Yes | TAQ ID as a string (to prevent precision loss), or handle. |

### `automation_taq_lookup` {#automation_taq_lookup}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/configuring`.

List TAQs (Trigger Action Query automations), or fetch one whole. Omit 'taq' to list: a list row is a compact form — automationID, handle, name, enabled — enough to pick a TAQ, not enough to run one, because a TAQ carries its entire trigger, step and path graph and one graph per row would swamp the result. Pass 'taq' — an ID or a handle — to fetch that single TAQ in full, including its triggers with their entry-point handles and input schemas. Do this before automation_taq_exec: the trigger's input schema is the only place the expected input parameter names are written down. 'query' matches the handle only — not the name — as a case-insensitive substring. A query containing spaces that matches nothing is retried lowercased with spaces turned into underscores, because handles are slugs. Disabled TAQs are omitted from the list unless you set 'includeDisabled'; deleted ones are always omitted.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `includeDisabled` | boolean |  | Also list disabled TAQs. Ignored when 'taq' is given. |
| `limit` | string |  | Maximum results, default 50, capped at 200. |
| `pageCursor` | string |  | Cursor from a previous response, to fetch the next page. |
| `query` | string |  | Case-insensitive substring of the handle. Ignored when 'taq' is given. |
| `taq` | string |  | TAQ ID as a string (to prevent precision loss), or handle. Omit to list instead. 'taqID', 'automation' and 'automationID' are accepted as aliases, so a near-miss name fetches the one TAQ rather than silently returning a list. |

### `automation_taq_undelete` {#automation_taq_undelete}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Restore a deleted TAQ (Trigger Action Query automation). Deleting a TAQ is soft — its triggers, steps and paths are all kept — so restoring puts the whole definition back, re-registers its triggers, and makes it runnable again if it is enabled. You cannot search for a deleted TAQ: automation_taq_lookup omits deleted TAQs from its list and does not find them by handle either, so you must already hold the TAQ's numeric ID — the 'automationID' automation_taq_lookup reported before the delete, or an ID from an audit trail. That ID does still resolve: pass it as 'taq' to automation_taq_lookup to see the TAQ, with 'deletedAt' set, before you restore it. Restoring a TAQ that is not deleted succeeds and changes nothing.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `taqID` | string | Yes | TAQ ID as a string (to prevent precision loss). A handle does not work here — a deleted TAQ cannot be resolved by handle. |

### `automation_taq_update` {#automation_taq_update}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Update a TAQ (Trigger Action Query automation): rename it, enable or disable it, or replace its trigger, step and path graph. This is also the tool for fixing the 'issues' automation_taq_create reported. A field you omit is left alone, so a rename does not disturb the graph and a graph change does not disturb the name. 'triggers', 'steps' and 'paths' each replace their whole collection rather than merging into it: send the complete array, not only the part you changed, and send [] only when you mean to remove every one. A TAQ with no steps stores clean and then executes as 'completed' having run nothing. Read the TAQ first with automation_taq_lookup — it returns exactly the shape these params accept, so the safe edit is lookup, change, send back. Echo each trigger's existing triggerID when you resend triggers, or the trigger is replaced by a new one with a new ID and any path pointing at the old one is orphaned. Anything these params do not carry is preserved from the stored TAQ: run-as user, owner, labels and scope are never cleared by omitting them. The write succeeds even when the graph is invalid. 'issues' says what is wrong, and any issue at all unregisters the TAQ from the runtime — 'runnable' comes back false and automation_taq_exec answers 'manager: executable not found' until the graph is clean. Disabling stops the triggers firing and also makes automation_taq_exec refuse the TAQ. To take one out of service reversibly, prefer enabled:false over automation_taq_delete — a deleted TAQ can no longer be found by handle at all.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `taq` | string | Yes | TAQ ID as a string (to prevent precision loss), or handle. Find it with automation_taq_lookup. |
| `description` | string |  | New description. Omit to leave it alone; an empty string clears it. |
| `enabled` | boolean |  | Whether the TAQ is live. Omit to leave it alone. False unregisters its triggers AND makes automation_taq_exec refuse it. |
| `handle` | string |  | New handle. Omit to leave it alone; an empty string clears it, after which the TAQ is reachable only by its numeric ID. |
| `name` | string |  | New name. Omit to leave it alone; it cannot be cleared, because a TAQ must have a name. |
| `paths` | string |  | Replaces every path. Omit to leave the paths alone. JSON array of edges: [{"parentID":"&lt;trigger or step ID>","childID":"&lt;step ID>"}]. IDs are strings of digits.<br>Send [] for the common case. With exactly one unconnected trigger and exactly one unconnected step, the two are wired together for you; more than one of either and you get a graph.ambiguousEntry issue instead.<br>Order carries meaning, not presentation. A non-gateway step's second outbound path is its error handler rather than a second successor; an iterator's paths are its body first, then its exit.<br>A gatewayExclusive or gatewayInclusive step needs at least two outbound paths: each carries a "condition", except exactly one which carries none and is the else branch, always tested last. A condition is an AST object, not an expression string: {"ref":"eq","args":[{"symbol":"foo"},{"value":{"@type":"String","@value":"bar"}}]}. Operators: and, or, not, isNull, isNotNull, eq, ne, lt, gt, lte, gte. Nothing checks a condition when it is written — a wrong operator, symbol or scope surfaces only when the TAQ runs. |
| `steps` | string |  | Replaces every step. Omit to leave the steps alone. JSON array of steps — what the TAQ does. One step:<br>{"stepID":"1","handle":"notify","kind":"function","ref":"notificationSend","meta":{"short":"Notify the owner"},"arguments":[{"argumentName":"recipient","value":"507566326668132353","type":"ID"},{"argumentName":"title","value":"Hello","type":"String"}]}<br>stepID is a string of digits you choose, unique across steps AND triggers; paths reference steps by it.<br>kind is one of: function, iterator, gatewayExclusive, gatewayInclusive, termination, error. Those are the only six, and they are NOT the workflow step kinds — a TAQ has no "expressions" step. A kind outside the six is rejected before anything is written.<br>ref, on a function or an iterator, names a construct-library function. List them with automation_taq_construct_lookup — each with its parameters, their types and whether they are required; it is the authority on which refs exist. Do NOT use automation_workflow_function_lookup: that is the workflow function registry, and most of its entries do not exist here. A ref that is not in the construct library is stored anyway and comes back as a "function.unknown" issue that names the near matches.<br>arguments bind by argumentName, which must equal one of that function's parameter argumentNames exactly; every required parameter must be present. "type" is compared literally against that parameter's "types" array and must be spelled as listed there (ID, String, Integer, Boolean, ComposeRecord, …) — omitting it is the same as sending "" and fails. Supply a literal with "value", copy an earlier step's output with "source" plus "scope" (the producing step's handle), or compute one with "expr".<br>An "expr" is evaluated INSIDE its own "scope", and without one that is the trigger's: so in an iterator body the record variable is still the trigger's record, and reaching the loop item means giving the argument the iterator's handle as "scope". One expression sees one scope, so a formula needing both the trigger's record and the loop item cannot be written as a single expr — split it across arguments, or derive the value with a module field value expression instead (compose_module_update, "fields[].expressions.value" — its scope is the record's own fields by bare name, documented on that parameter and in the field_expressions skill).<br>In an expression a Record-ref field holds the target's ID as a string, not a record: a filter over them reads "card = " + record.values.card (string concatenation). Writing record.values.card.recordID yields nothing, so the query matches nothing and an iterator over it runs ZERO times in silence — no error, no issue, and the execution still reports completed. The resolved query is recorded in that step's frame under "args", which is where a filter that came out wrong becomes obvious.<br>A step's results cannot be set: they are derived from the function definition and anything you send is overwritten. An error step's message likewise cannot be set through the API.<br>A termination step is optional — every leaf step is wired to an auto-injected one. |
| `triggers` | string |  | Replaces every trigger. Omit to leave the triggers alone. JSON array of triggers — what starts the TAQ. One trigger:<br>{"triggerID":"1001","handle":"onAgent","enabled":true,"resourceType":"automation:trigger:agentic","eventType":"onAgentic","meta":{"short":"Invoked by an agent"},"inputSchema":[{"name":"subject","type":"String","required":true}]}<br>resourceType and eventType are a pair from the trigger catalogue: automation_taq_construct_lookup with catalogue "triggers", 22 entries. "automation:trigger:agentic" with "onAgentic" is the pair automation_taq_exec runs, so include one of those if you want to be able to run the TAQ on demand; the compose:record and system:user pairs fire on real events instead.<br>triggerID is a string of digits. Omit it on create and one is minted for you and returned. Supply it when a path references the trigger, and echo the existing triggerID on update — a trigger sent without one is replaced by a new trigger with a new ID, orphaning any path that pointed at the old one.<br>handle names the trigger: it is what automation_taq_exec takes as entryPoint, and what a step reads from as a scope.<br>inputSchema declares the parameters automation_taq_exec accepts under "input"; it is the only place those names are written down. |

### `automation_trigger_create` {#automation_trigger_create}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Create a trigger, making a workflow run whenever a given event fires. This changes when the named workflow runs, so only create one when you intend that workflow to start responding to the event. It takes effect immediately: the trigger is registered as this call returns and the workflow can fire on the very next matching event, with no restart and no separate activation step. 'eventType' and 'resourceType' must be a pair that actually exists — call automation_event_type_lookup for the catalogue and copy a pair from it. A pair that does not exist is stored without complaint and simply never fires. A trigger created here carries NO constraints, so it fires for every event of that type on that resource type: an afterUpdate trigger on compose:record runs the workflow for every record update in every module of every namespace. Constraints, which narrow a trigger to a particular namespace, module or field value, cannot be set through this tool — create constrained triggers in the workflow editor in the Human webapp instead. Two cases are stored but silently never registered, so nothing fires and no error is returned: an onInterval or onTimestamp trigger (and system:sink and system:queue events) on a workflow with no run-as user, and any enabled trigger on a workflow marked as a sub-workflow. Check the workflow with automation_workflow_lookup first if either might apply. A disabled workflow, or one with unresolved issues, likewise registers nothing.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `eventType` | string | Yes | Event that fires the workflow, e.g. "onManual" or "afterCreate". Must be a value automation_event_type_lookup reports for this resourceType. |
| `resourceType` | string | Yes | Resource the event happens on, e.g. "compose:record" or "system:user". Must be paired with 'eventType' as automation_event_type_lookup reports it. |
| `workflow` | string | Yes | Workflow to run, as an ID string (to prevent precision loss) or a handle. This is the workflow that will start responding to the event. |
| `description` | string |  | Human-readable note on what this trigger is for. Shown in the workflow editor. |
| `enabled` | boolean |  | Whether the trigger fires. Defaults to true. A disabled trigger is stored but never registered. |
| `input` | string |  | JSON object of fixed input variables merged into the workflow's scope on every run. Omit when the event's own properties are all the workflow needs. |
| `stepID` | string |  | ID as a string of the workflow step to start at. Omit only when the workflow has exactly one starting step; otherwise starting it fails at run time. It also decides where the trigger is drawn: the Start node is placed one row above this step and connected to it. |

### `automation_trigger_delete` {#automation_trigger_delete}

<Badge type="danger" text="Deletes" /> Listed under `/api/mcp/configuring`.

Delete a trigger, so its event stops starting its workflow. The workflow itself is untouched and still exists: deleting every trigger only means nothing starts it automatically, and it can still be run on demand with automation_workflow_exec. The delete is soft and reversible — automation_trigger_undelete restores it, and automation_trigger_lookup with 'includeDeleted' still lists it. Prefer disabling the trigger with automation_trigger_update when you only want to pause it, since that keeps it visible in the workflow editor. The trigger stops firing immediately — it is unregistered from the running instance as well as marked deleted — and automation_trigger_undelete brings it back.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `triggerID` | string | Yes | Trigger ID as a string (to prevent precision loss). Find it with automation_trigger_lookup. |

### `automation_trigger_lookup` {#automation_trigger_lookup}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/configuring`.

List automation triggers, or fetch one whole. A trigger binds a workflow to an event: when an event of its eventType fires on its resourceType, the workflow runs. Omit 'triggerID' to list: a list row is a compact form — triggerID, workflowID, eventType, resourceType, enabled, deletedAt — enough to pick a trigger out of a list. Pass 'triggerID' to fetch that one trigger in full, which additionally returns its constraints (the conditions that decide whether it actually fires), its fixed input and its metadata; those are left out of list rows because they are bulky and rarely what tells two rows apart. Narrow the list with 'workflow' to answer "what makes this workflow run?", or with 'eventType' and 'resourceType' to answer "what runs when this happens?". Both match exactly, not as substrings, and the values that exist are listed by automation_event_type_lookup — guessing them returns an empty list rather than an error. Disabled triggers are omitted unless you set 'includeDisabled', deleted ones unless you set 'includeDeleted'. Fetching by 'triggerID' ignores both flags and returns the trigger either way, so this is also how you inspect a deleted trigger before restoring it with automation_trigger_undelete.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `eventType` | string |  | Exact event type to match, e.g. "onManual" or "afterUpdate". Valid values come from automation_event_type_lookup. |
| `includeDeleted` | boolean |  | Also list deleted triggers. This is how you find the ID for automation_trigger_undelete. |
| `includeDisabled` | boolean |  | Also list disabled triggers, which stay stored but never fire. |
| `limit` | string |  | Maximum results, default 50, capped at 200. |
| `pageCursor` | string |  | Cursor from a previous response, to fetch the next page. |
| `resourceType` | string |  | Exact resource type to match, e.g. "compose:record". Valid values come from automation_event_type_lookup. |
| `triggerID` | string |  | Trigger ID as a string (to prevent precision loss). Omit to list instead. Triggers have no handle; they are identified by ID only. |
| `workflow` | string |  | Workflow ID as a string (to prevent precision loss), or workflow handle. Lists only that workflow's triggers. Ignored when 'triggerID' is given. |

### `automation_trigger_undelete` {#automation_trigger_undelete}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Restore a deleted trigger, putting its workflow back on the event it used to respond to. The delete was soft, so the event type, resource type, constraints and input all come back as they were. Find the ID with automation_trigger_lookup: unlike most resources here, deleted triggers are reachable — pass 'includeDeleted' to list them, and passing a deleted trigger's ID as 'triggerID' returns it in full with 'deletedAt' set, so you can check what you are about to restore. Restoring a trigger that is not deleted succeeds and changes nothing. The trigger starts firing again immediately: it is re-registered on the running instance as part of the restore.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `triggerID` | string | Yes | Trigger ID as a string (to prevent precision loss), from automation_trigger_lookup with includeDeleted set. |

### `automation_trigger_update` {#automation_trigger_update}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Change an existing trigger: re-point it at another workflow or event, enable or disable it, or edit its note and fixed input. Omit a field to leave it unchanged; pass 'description' or 'input' empty to clear it. 'workflow', 'eventType' and 'resourceType' can be changed but not cleared — an empty value is refused rather than stored, because a trigger with no workflow can never be loaded again and so becomes impossible to edit, delete or restore. Constraints and labels are preserved untouched; neither can be set here. New 'eventType'/'resourceType' values must be a pair automation_event_type_lookup reports, and an invalid pair is stored without complaint and simply never fires. The change takes effect immediately: the trigger is re-registered on the running instance, so a trigger you disable here stops firing at once and one you re-point starts firing on its new event.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `triggerID` | string | Yes | Trigger ID as a string (to prevent precision loss). Find it with automation_trigger_lookup. |
| `description` | string |  | New note. Empty clears it. |
| `enabled` | boolean |  | Enable or disable the trigger. See the note above: disabling does not stop it firing until the workflow is saved again or the server restarts. |
| `eventType` | string |  | New event type, from automation_event_type_lookup. Cannot be cleared. |
| `input` | string |  | Replacement JSON object of fixed input variables. Replaces the stored set wholesale rather than merging. Empty clears it. |
| `resourceType` | string |  | New resource type, from automation_event_type_lookup. Cannot be cleared. |
| `stepID` | string |  | ID as a string of the workflow step to start at. Empty resets it, which requires the workflow to have exactly one starting step. |
| `workflow` | string |  | Move the trigger to this workflow: an ID string or a handle. Cannot be cleared. |

### `automation_workflow_create` {#automation_workflow_create}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Create an automation workflow: a graph of steps joined by paths, which runs when a trigger fires it or when automation_workflow_exec calls it. Workflows are the step-and-path automations; Trigger-Action-Query automations are a different resource with their own tools, and nothing here transfers to them. The smallest workflow that runs is one step and no paths: steps [{"stepID":"1","kind":"expressions","arguments":[{"target":"result","expr":"1 + 1","type":"Integer"}]}] with no 'paths' at all. Execution starts at the one step nothing points at and follows the paths out of it until it reaches a step with none. A workflow may have only ONE such starting step, and this tool refuses a graph with more than one rather than storing something that fails on its first run. A BROKEN WORKFLOW IS STILL STORED. Only a missing name, an invalid or already-taken handle, a run-as user that cannot be loaded and a stale 'updatedAt' are refused. Everything else — an unknown function ref, an argument whose type does not match the parameter, a step kind given arguments it does not take, a path pointing at a step that does not exist, an unparseable expression — comes back in the "issues" array of this tool's result, and the workflow is written anyway. A workflow with issues never runs and its triggers are never registered, so a call that "succeeds" with a non-empty "issues" is a failure to fix, not a warning to note. Read "issues" before reporting success. Before you write a 'function' or 'iterator' step, read the catalogue with automation_workflow_function_lookup — 93 entries, each with the 'ref' you put on the step, its parameters, and the exact type names each parameter accepts (a parameter declares 'types' as a list; your argument's single 'type' has to be one of them, spelled the same). A guessed ref stores a workflow that will not run. NOTHING FIRES THIS WORKFLOW ON ITS OWN. Triggers are a separate resource: call automation_trigger_create with this workflow's handle or ID afterwards, and give it 'stepID' when you want it to start somewhere other than the workflow's only starting step. Until then the workflow runs only when automation_workflow_exec is called. Verify with automation_workflow_exec once created — storing and running are different claims.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `handle` | string | Yes | Unique handle, e.g. "order_sync". Must start with a letter, be at least 2 characters, and use only letters, digits, underscores, dashes and dots — but prefer snake_case: a hyphen is the subtraction operator wherever an identifier is parsed. Handles are unique across the whole instance, not per project, and a taken one is refused. |
| `name` | string | Yes | Display name shown in the workflow list and editor. A workflow without one is refused. |
| `description` | string |  | What this workflow is for. Shown in the workflow editor. |
| `enabled` | boolean |  | Whether the workflow may run. Defaults to true. A disabled workflow keeps its definition but never fires, and automation_workflow_exec reports it as not found. It is also left out of the automation_workflow_lookup listing unless includeDisabled is set — but it still resolves by handle, so it can be enabled again without its ID. |
| `paths` | string |  | JSON array of connections between steps: [{"parentID":"1","childID":"2","expr":"","meta":{"visual":{"id":"e1","parent":"1","points":[],"style":""}}}]. "parentID" and "childID" are the QUOTED stepIDs from 'steps'. Omit paths entirely for a single-step workflow — a step with no outbound path ends the run. ARRAY ORDER IS SEMANTICS, not presentation, wherever a step branches: the FIRST path out of an iterator is the loop body and the SECOND is the exit; the SECOND path out of an error-handler is the catch branch; an "excl" gateway tests the "expr" of each path leaving it in array order and takes the first that passes, so the fallback path — the one with an empty "expr" — must be last. "expr" is read only on paths leaving an "excl" or "incl" gateway; on any other path it is ignored. |
| `runAs` | string |  | User ID as a string (to prevent precision loss) whose permissions the workflow's steps run with. Omit to run as whoever or whatever started it. Interval and timestamp triggers only register on a workflow that has one. A user that cannot be loaded fails the whole call. |
| `scope` | string |  | JSON object of variables every run starts with, e.g. {"retries":3}. Input passed to automation_workflow_exec, and input carried by a trigger, is merged over this. |
| `steps` | string |  | JSON array of steps. Give every step a "stepID" of your own choosing, as a QUOTED string, unique within this workflow — paths reference it. Shape: {"stepID":"1","kind":"expressions","ref":"","arguments":[],"results":[],"meta":{"name":"","visual":{"id":"1","parent":"1","value":"Label","xywh":[240,120,180,64]}}}. An argument or a result is {"target":"...","expr":"...","value":...,"type":"..."}: "target" is the name (on a function step, the parameter the argument feeds; on an expressions step or on any result, the scope variable it writes), "expr" is an expression evaluated against the workflow scope, "value" is a literal used when there is no "expr", and "type" is the expression type — "Any" when omitted, except on a function or iterator argument, where it must match one of that parameter's declared types exactly. Step kinds and what each one expects: "expressions" — no ref, one or more arguments, at most one outbound path, e.g. {"target":"total","expr":"1 + 1","type":"Integer"}. "function" — "ref" names a construct from automation_workflow_function_lookup (93 of them; NOT automation_taq_construct_lookup, which is the smaller registry Trigger-Action-Query automations use and which a workflow must not be built from), arguments match that construct's parameters by "target", results copy its outputs into scope as {"target":"&lt;variable>","expr":"&lt;result name>","type":"&lt;type>"}, at most one outbound path. "iterator" — like function but the ref must be an iterating construct (composeRecordsEach and friends), and it needs exactly two outbound paths. "gateway" — "ref" is "fork" (run every branch), "join" (wait for branches to meet), "excl" (take the first matching branch) or "incl" (take every matching branch); no arguments and no results. "termination" — no ref, no arguments, no outbound path; ends the workflow, and takes ONE inbound path: the editor refuses a connection into a termination that already has one, and applies that rule to stored paths as well as drawn ones, so two arms aimed at the same termination make each other invalid and BOTH vanish from the canvas. Give every branch arm its own termination. "error" — a single {"target":"message","type":"String","value":"..."} argument and no outbound path; fails the run. "error-handler" — no ref, one or two outbound paths. "delay" — exactly one argument, either {"target":"timestamp","type":"DateTime"} or {"target":"offset","type":"Duration"}. "prompt" — "ref" names a prompt construct such as "notification". "exec-workflow" — no ref, a required {"target":"workflow","type":"Handle"} (or type "ID") argument naming another workflow, and an optional {"target":"scope","type":"Vars"}. "break" and "continue" — inside an iterator body, nothing else set. "debug" — logs the whole scope. "visual" — a canvas note that never runs. "meta.visual.xywh" is [x, y, width, height] on the workflow editor's canvas. Omit it and the graph is laid out for you, which is what you usually want: a node is 200x80, a step sits 160 below the one it follows, and a branch opens a column 280 to its right. The first path out of a step carries its column straight down and every later path takes a new column, so a gateway's arms and an error handler's catch branch separate on their own; an iterator inverts that, sending its body (the first path) into the new column and carrying its exit (the second) down. Set the field only to place a node deliberately — a step that carries coordinates is left exactly where you put it, and the computed graph is laid out clear of it. |
| `subWorkflow` | boolean |  | Mark this workflow as one that only other workflows call, through their 'exec-workflow' steps. Triggers on a sub-workflow are stored but never registered, so set this only when nothing should fire it directly. |
| `trace` | boolean |  | Keep a full step-by-step trace of every run. Useful while building; it makes every run store more, so turn it off once the workflow works. |

### `automation_workflow_delete` {#automation_workflow_delete}

<Badge type="danger" text="Deletes" /> Listed under `/api/mcp/configuring`.

Delete an automation workflow. The delete is soft — its step and path graph is kept — and automation_workflow_undelete restores it intact. Deleting unregisters its triggers, so nothing fires the workflow any more and automation_workflow_exec will not find it. Any TAQ or trigger pointing at this workflow stops working; check automation_trigger_lookup with this workflow first if you are not sure what depends on it. Prefer disabling the workflow when you only want to pause it — a disabled workflow still resolves by handle, a deleted one is not findable by handle at all and restoring it needs an ID you must have kept.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `workflow` | string | Yes | Workflow ID as a string (to prevent precision loss) or handle. Find it with automation_workflow_lookup. |

### `automation_workflow_exec` {#automation_workflow_exec}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/usage`.

Run a workflow by ID or handle, wait for it to finish, and return the variables it produced. This performs the workflow's real side effects — it is not a way to inspect what a workflow does; use automation_workflow_lookup for that. The workflow must be enabled: executing a disabled workflow fails rather than queueing. 'input' supplies the workflow's input variables and overrides any input the workflow's own trigger defines. For Trigger-Action-Query automations use automation_taq_exec instead.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `workflow` | string | Yes | Workflow ID as string (to prevent precision loss) or handle |
| `input` | string |  | JSON object of input variables |

### `automation_workflow_function_lookup` {#automation_workflow_function_lookup}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/configuring`.

List the workflow function registry — every ref a workflow step of kind 'function' or 'iterator' may name, with its parameters, their exact type names and which are required, plus the results it puts into scope for later steps. There are around 93, returned whole by default. Read this before automation_workflow_create or automation_workflow_update. A step naming a ref that is not here is stored and returns success, and the workflow then fails to convert into anything runnable — a guessed ref is not a guess that gets corrected, it is a workflow that silently does nothing. This is NOT the TAQ construct library. Most of these functions are unavailable to a TAQ, and a TAQ step naming one of them stores with a 'function.unknown' issue and never runs. For a TAQ use automation_taq_construct_lookup instead. A workflow function's parameters are keyed by 'name' — the construct library keys its own by 'argumentName' — so an argument list copied from one engine to the other is wrong even where the ref happens to exist in both. 'kind' says how the step is wired: a 'function' runs once, an 'iterator' runs its body once per item. 'labels' records where the registry advertises the function.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `ref` | string |  | Narrow to one function. An exact ref is matched first, then any ref containing it, case-insensitively. A ref that matches nothing answers with the full list of refs rather than an empty result. |

### `automation_workflow_lookup` {#automation_workflow_lookup}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/configuring`.

Look up one automation workflow, or list the workflows on this instance. Provide 'workflow' (ID or handle) to fetch a single workflow in full, including its step and path graph; omit it to list workflows as summaries — workflow ID, handle, name and enabled flag — and then fetch the one you want. List mode never returns the graph, so a listing cannot tell you what a workflow does. Narrow the list with 'query': it matches the handle as a case-insensitive substring and does not search the display name; a query containing spaces that finds nothing is retried once as a slug ("My Workflow" is retried as "my_workflow"). 'query' is ignored when 'workflow' is given. Workflows are the step-and-path automations; Trigger-Action-Query automations are a separate resource — use automation_taq_lookup for those.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `includeDisabled` | boolean |  | Include disabled workflows in the listing. Off by default, matching the rest of the product. Naming one workflow by handle always finds it, enabled or not; this flag only affects the listing. |
| `limit` | string |  | Maximum workflows to return in list mode, default 50, capped at 200. |
| `pageCursor` | string |  | Cursor from a previous list response, to fetch the next page. |
| `query` | string |  | Case-insensitive substring match on the workflow handle. Applies to list mode only. |
| `workflow` | string |  | Workflow ID as string (to prevent precision loss) or handle. Omit to list all. |

### `automation_workflow_undelete` {#automation_workflow_undelete}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Restore a deleted automation workflow. Deleting a workflow is soft — its step and path graph is kept — so restoring brings the workflow back intact, re-registers its triggers and makes it executable again if it is enabled. You cannot search for a deleted workflow: automation_workflow_lookup leaves deleted workflows out of its list and does not find them by handle, so you need the workflow's numeric ID from before the delete, or from an audit trail. That ID still resolves: pass it as 'workflow' to automation_workflow_lookup to see what you are about to restore — a deleted workflow comes back with 'deletedAt' set. Restoring a workflow that is not deleted succeeds and changes nothing. To restore a Trigger-Action-Query automation use automation_taq_undelete instead.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `workflowID` | string | Yes | Workflow ID as a string (to prevent precision loss). A handle does not work here — a deleted workflow cannot be resolved by handle. |

### `automation_workflow_update` {#automation_workflow_update}

<Badge type="warning" text="Writes" /> Listed under `/api/mcp/configuring`.

Change an existing workflow: rename it, enable or disable it, or replace its step and path graph. This is a read-modify-write — the workflow is loaded, only the arguments you pass are applied, and its stored 'updatedAt' is echoed back so an edit someone else made in the meantime is refused instead of being silently overwritten. An argument you omit is left exactly as it is, so renaming or disabling a workflow is safe to send on its own and needs no graph. But 'steps' and 'paths' REPLACE what is stored when you do pass them — they are never merged, and passing an empty array [] is a request to delete every step or every path, not a way to say "unchanged". Send the complete graph: call automation_workflow_lookup with 'workflow' first and edit what it returns. Everything you cannot set here — the owner, how many sessions are kept, labels — is preserved. Changing the graph takes effect at once: the workflow is re-converted and re-registered on the running instance, so the next event or the next automation_workflow_exec uses the new steps. A workflow may have only ONE step that nothing points at, and this tool refuses a graph with more than one rather than storing something that fails on its first run. A BROKEN WORKFLOW IS STILL STORED. Only a missing name, an invalid or already-taken handle, a run-as user that cannot be loaded and a stale 'updatedAt' are refused. Everything else — an unknown function ref, an argument whose type does not match the parameter, a step kind given arguments it does not take, a path pointing at a step that does not exist, an unparseable expression — comes back in the "issues" array of this tool's result, and the workflow is written anyway. A workflow with issues never runs and its triggers are never registered, so a call that "succeeds" with a non-empty "issues" is a failure to fix, not a warning to note. Read "issues" before reporting success. Fixing the issues and calling this tool again clears them. Triggers are a separate resource and are not touched here: use automation_trigger_update to change when the workflow fires, and check automation_trigger_lookup for this workflow after you renumber steps, since a trigger pinned to a stepID that no longer exists fails at run time.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `workflow` | string | Yes | Workflow ID as a string (to prevent precision loss) or handle. Find it with automation_workflow_lookup; a disabled workflow resolves by handle too, though it is left out of the listing unless includeDisabled is set. |
| `description` | string |  | New description. An empty string clears it. |
| `enabled` | boolean |  | Enable or disable the workflow. The definition is kept either way. A disabled workflow is left out of the automation_workflow_lookup listing unless includeDisabled is set, but it still resolves by handle, so this is reversible with the handle alone. |
| `handle` | string |  | New handle. Must be unique across the instance. An empty string clears it, after which the workflow can only be referenced by ID. |
| `name` | string |  | New display name. Cannot be emptied — a workflow with no name is refused. |
| `paths` | string |  | Replaces the stored paths wholesale. JSON array of connections between steps: [{"parentID":"1","childID":"2","expr":"","meta":{"visual":{"id":"e1","parent":"1","points":[],"style":""}}}]. "parentID" and "childID" are the QUOTED stepIDs from 'steps'. Omit paths entirely for a single-step workflow — a step with no outbound path ends the run. ARRAY ORDER IS SEMANTICS, not presentation, wherever a step branches: the FIRST path out of an iterator is the loop body and the SECOND is the exit; the SECOND path out of an error-handler is the catch branch; an "excl" gateway tests the "expr" of each path leaving it in array order and takes the first that passes, so the fallback path — the one with an empty "expr" — must be last. "expr" is read only on paths leaving an "excl" or "incl" gateway; on any other path it is ignored. |
| `runAs` | string |  | User ID as a string (to prevent precision loss) whose permissions the steps run with. An empty string clears it, and interval or timestamp triggers on this workflow then stop registering. |
| `scope` | string |  | JSON object of variables every run starts with. Replaces the stored set wholesale; {} clears it. |
| `steps` | string |  | Replaces the stored steps wholesale. JSON array of steps. Give every step a "stepID" of your own choosing, as a QUOTED string, unique within this workflow — paths reference it. Shape: {"stepID":"1","kind":"expressions","ref":"","arguments":[],"results":[],"meta":{"name":"","visual":{"id":"1","parent":"1","value":"Label","xywh":[240,120,180,64]}}}. An argument or a result is {"target":"...","expr":"...","value":...,"type":"..."}: "target" is the name (on a function step, the parameter the argument feeds; on an expressions step or on any result, the scope variable it writes), "expr" is an expression evaluated against the workflow scope, "value" is a literal used when there is no "expr", and "type" is the expression type — "Any" when omitted, except on a function or iterator argument, where it must match one of that parameter's declared types exactly. Step kinds and what each one expects: "expressions" — no ref, one or more arguments, at most one outbound path, e.g. {"target":"total","expr":"1 + 1","type":"Integer"}. "function" — "ref" names a construct from automation_workflow_function_lookup (93 of them; NOT automation_taq_construct_lookup, which is the smaller registry Trigger-Action-Query automations use and which a workflow must not be built from), arguments match that construct's parameters by "target", results copy its outputs into scope as {"target":"&lt;variable>","expr":"&lt;result name>","type":"&lt;type>"}, at most one outbound path. "iterator" — like function but the ref must be an iterating construct (composeRecordsEach and friends), and it needs exactly two outbound paths. "gateway" — "ref" is "fork" (run every branch), "join" (wait for branches to meet), "excl" (take the first matching branch) or "incl" (take every matching branch); no arguments and no results. "termination" — no ref, no arguments, no outbound path; ends the workflow, and takes ONE inbound path: the editor refuses a connection into a termination that already has one, and applies that rule to stored paths as well as drawn ones, so two arms aimed at the same termination make each other invalid and BOTH vanish from the canvas. Give every branch arm its own termination. "error" — a single {"target":"message","type":"String","value":"..."} argument and no outbound path; fails the run. "error-handler" — no ref, one or two outbound paths. "delay" — exactly one argument, either {"target":"timestamp","type":"DateTime"} or {"target":"offset","type":"Duration"}. "prompt" — "ref" names a prompt construct such as "notification". "exec-workflow" — no ref, a required {"target":"workflow","type":"Handle"} (or type "ID") argument naming another workflow, and an optional {"target":"scope","type":"Vars"}. "break" and "continue" — inside an iterator body, nothing else set. "debug" — logs the whole scope. "visual" — a canvas note that never runs. "meta.visual.xywh" is [x, y, width, height] on the workflow editor's canvas. Omit it and the graph is laid out for you, which is what you usually want: a node is 200x80, a step sits 160 below the one it follows, and a branch opens a column 280 to its right. The first path out of a step carries its column straight down and every later path takes a new column, so a gateway's arms and an error handler's catch branch separate on their own; an iterator inverts that, sending its body (the first path) into the new column and carrying its exit (the second) down. Set the field only to place a node deliberately — a step that carries coordinates is left exactly where you put it, and the computed graph is laid out clear of it. |
| `subWorkflow` | boolean |  | Mark the workflow as one that only other workflows call. Turning this on unregisters its triggers. |
| `trace` | boolean |  | Keep a full step-by-step trace of every run. |

## Discovery

Full-text search across the instance.

### `discovery_search` {#discovery_search}

<Badge type="tip" text="Reads" /> Listed under `/api/mcp/usage`.

::: info
Only available when [`DISCOVERY_ENABLED`](/reference/environment#DISCOVERY_ENABLED), [`DISCOVERY_BASE_URL`](/reference/environment#DISCOVERY_BASE_URL) and [`DISCOVERY_JWT_SECRET`](/reference/environment#DISCOVERY_JWT_SECRET) are all set.
:::

Search or list records you have access to. Use this when the user asks to find, list, or show existing records. Leave 'query' empty to list records, or give a specific field value (a name, email, phone) to search within them. Scope the search with 'namespace', and optionally 'module', to avoid searching everything you can read. Omitting 'namespace' searches every namespace you have access to, which is rarely what you want and can be slow. When running as a configured agent, only the namespaces in your DISCOVERY ACCESS section are permitted and anything else is denied. Results are capped and cannot be paged — narrow the query rather than asking for more.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `limit` | string |  | Maximum records to return, default 50, capped at 200. There is no cursor: this search cannot be paged. |
| `module` | string |  | Module name, handle, or ID (as string to prevent precision loss). Requires 'namespace'. Omit to search all modules in the namespace. |
| `namespace` | string |  | Namespace name, handle, slug, or ID (as string to prevent precision loss). Omit to search every namespace you can access. |
| `query` | string |  | A specific value to search for inside record fields (e.g. a name, email, phone). Leave empty to list records. |
