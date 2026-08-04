# MCP resource-name table

The authoritative mapping from REST resource to tool resource segment. A tool
author does **not** derive the segment — they look it up here. Without this,
~170 fanned-out agents invent ~170 answers to the same question.

Tool names are `{app}_{resource}_{op}`; see `CONVENTIONS.md` §6 for the op
vocabulary.

## Rules

1. **All snake_case.** `userGroup` → `user_group`, never `userGroup` or
   `usergroup`. Consequence: a tool name cannot be split mechanically into
   app/resource/op, because the resource segment contains underscores. Nothing
   parses names — the structural test validates against this table.
2. **Singular.** `queues` → `queue`. `lookup` already covers fetch-one and
   list, so a plural would imply a distinction that does not exist.
3. **The table wins over the REST name.** Where Human's vocabulary differs from
   the generated REST resource, the vocabulary wins. One case today:
   `ngAutomation` → `taq`.

## Table

42 codegen'd REST resources. This table fixes **names**, not coverage — how many
tools a resource has today is in `TOOLS.md`, which is generated from the
registry and asserted by a test. A hand-maintained count here drifted within a
day of being written.

### compose

| REST resource | Tool segment | Notes |
|---|---|---|
| `namespace` | `compose_namespace` | |
| `module` | `compose_module` | |
| `record` | `compose_record` | usage group; the rest of compose is configuring |
| `page` | `compose_page` | includes `block_schema`, a legacy 4-segment name |
| `pageLayout` | `compose_page_layout` | ⚠ prefix-shadows `compose_page_*` |
| `chart` | `compose_chart` | |
| `attachment` | `compose_attachment` | |
| `icon` | `compose_icon` | |

### automation

| REST resource | Tool segment | Notes |
|---|---|---|
| `ngAutomation` | `automation_taq` | **vocabulary divergence** — TAQ = Trigger Action Query, never "ngAutomation" or "Task Queue" in any description |
| `workflow` | `automation_workflow` | |
| `trigger` | `automation_trigger` | |
| `eventTypes` | `automation_event_type` | singularised |

Note `automation_<numeric id>` is reserved: the runtime mints per-TAQ tools
under that shape. No static tool may take a numeric resource segment.

### system

| REST resource | Tool segment | Notes |
|---|---|---|
| `user` | `system_user` | has `FindByAny` |
| `userGroup` | `system_user_group` | has `FindByAny` |
| `role` | `system_role` | has `FindByAny`; `role.Membership` is on the §8.6 deny-list |
| `authClient` | `system_auth_client` | `ExposeSecret`, `RegenerateSecret` and `Create` excluded under §8.6b — they disclose a credential, not for want of a check |
| `application` | `system_application` | |
| `template` | `system_template` | has `FindByAny` |
| `tenant` | `system_tenant` | out of scope while tenancy is single-instance |
| `reminder` | `system_reminder` | **exemplar** — soft delete, domain ops, no handle, caller scoping |
| `notification` | `system_notification` | |
| `report` | `system_report` | runs saved definitions only; ad-hoc analytics is `compose record report`, a separate undocumented endpoint |
| `queues` | `system_queue` | singularised |
| `agent` | `system_agent` | |
| `chatbot` | `system_chatbot` | `chatbotSession` and `chatbotPreview` are on the §8.6 deny-list |
| `knowledgeBase` | `system_knowledge_base` | |
| `aiConversation` | `system_ai_conversation` | |
| `connection` | `system_connection` | L4 builds on this |
| `configuredConnection` | `system_configured_connection` | L4 builds on this |
| `dalConnection` | `system_dal_connection` | distinct from `connection` |
| `dalSensitivityLevel` | `system_dal_sensitivity_level` | |
| `apigwRoute` | `system_apigw_route` | |
| `apigwFilter` | `system_apigw_filter` | |
| `project` | `system_project` | ⚠ prefix-shadows the 8 below |
| `projectTask` | `system_project_task` | |
| `projectFeature` | `system_project_feature` | |
| `projectBacklogItem` | `system_project_backlog_item` | |
| `projectIncident` | `system_project_incident` | |
| `projectReview` | `system_project_review` | |
| `projectPrivacy` | `system_project_privacy` | |
| `projectAiSystem` | `system_project_ai_system` | |
| `projectFriaScenario` | `system_project_fria_scenario` | |

`discovery_search` is not in this table: it is a legacy two-segment name in
package `system/agentic`, hardcoded in eight places in `executor.go`, and
grandfathered rather than renamed.

## Prefix shadowing

Two families produce names where one resource's prefix is another's:

- `compose_page_*` and `compose_page_layout_*`
- `system_project_*` and the eight `system_project_<sub>_*`

This is safe for dispatch — nothing parses names — but these are exactly the
pairs a model is most likely to confuse when selecting a tool. Descriptions in
both families must state what the tool is *not*: a page layout is not a page,
a project task is not a project.

## Not in this table

Non-resource REST surfaces with no codegen'd resource — `permissions`, `stats`,
`dml`, `expression`, `locales`, `sink`, `settings`, `data_privacy` — are out of
scope for the resource grammar. Where they need tools, the shape is decided
case by case and recorded here first.
