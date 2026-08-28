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

| REST resource | Tool segment          | Notes                                            |
| ------------- | --------------------- | ------------------------------------------------ |
| `namespace`   | `compose_namespace`   |                                                  |
| `module`      | `compose_module`      |                                                  |
| `record`      | `compose_record`      | usage group; the rest of compose is configuring  |
| `page`        | `compose_page`        | includes `block_schema`, a legacy 4-segment name |
| `pageLayout`  | `compose_page_layout` | ⚠ prefix-shadows `compose_page_*`                |
| `chart`       | `compose_chart`       |                                                  |
| `attachment`  | `compose_attachment`  |                                                  |
| `icon`        | `compose_icon`        |                                                  |

### automation

| REST resource      | Tool segment                   | Notes                                                                                                           |
| ------------------ | ------------------------------ | --------------------------------------------------------------------------------------------------------------- |
| `ngAutomation`     | `automation_taq`               | **vocabulary divergence** — TAQ = Trigger Action Query, never "ngAutomation" or "Task Queue" in any description |
| `workflow`         | `automation_workflow`          |                                                                                                                 |
| `trigger`          | `automation_trigger`           |                                                                                                                 |
| `eventTypes`       | `automation_event_type`        | singularised                                                                                                    |
| `constructLibrary` | `automation_taq_construct`     | **vocabulary divergence** — named for what it is the vocabulary _of_. See below                                 |
| `function`         | `automation_workflow_function` | the workflow function registry, and not the construct library. See below                                        |

Note `automation_<numeric id>` is reserved: the runtime mints per-TAQ tools
under that shape. No static tool may take a numeric resource segment.

`constructLibrary` and `function` are two catalogues of step refs for two
different engines, and naming one from the other is the failure both authoring
tools warn about hardest — the definition stores, validates as far as anything
checks, and never runs. A single `automation_function_lookup` would have been
the one tool both callers reached for, so the segment carries the engine
instead: a TAQ author finds `automation_taq_construct_lookup` beside
`automation_taq_create`, a workflow author finds
`automation_workflow_function_lookup` beside `automation_workflow_create`, and
neither name reads as the other's. This is the prefix shadowing below used on
purpose.

### system

| REST resource          | Tool segment                   | Notes                                                                                                                      |
| ---------------------- | ------------------------------ | -------------------------------------------------------------------------------------------------------------------------- |
| `user`                 | `system_user`                  | has `FindByAny`                                                                                                            |
| `userGroup`            | `system_user_group`            | has `FindByAny`; there is no member removal to map — see below                                                             |
| `role`                 | `system_role`                  | has `FindByAny`; `role.Membership` is on the §8.6 deny-list                                                                |
| `authClient`           | `system_auth_client`           | `ExposeSecret`, `RegenerateSecret` and `Create` excluded under §8.6b — they disclose a credential, not for want of a check |
| `application`          | `system_application`           |                                                                                                                            |
| `template`             | `system_template`              | has `FindByAny`                                                                                                            |
| `tenant`               | `system_tenant`                | out of scope while tenancy is single-instance                                                                              |
| `reminder`             | `system_reminder`              | **exemplar** — soft delete, domain ops, no handle, caller scoping                                                          |
| `notification`         | `system_notification`          |                                                                                                                            |
| `report`               | `system_report`                | runs saved definitions only; ad-hoc analytics is `compose record report`, a separate undocumented endpoint                 |
| `queues`               | `system_queue`                 | singularised                                                                                                               |
| `agent`                | `system_agent`                 |                                                                                                                            |
| `chatbot`              | `system_chatbot`               | `chatbotSession` and `chatbotPreview` are on the §8.6 deny-list                                                            |
| `knowledgeBase`        | `system_knowledge_base`        |                                                                                                                            |
| `aiConversation`       | `system_ai_conversation`       |                                                                                                                            |
| `connection`           | `system_connection`            | L4 builds on this                                                                                                          |
| `configuredConnection` | `system_configured_connection` | L4 builds on this                                                                                                          |
| `dalConnection`        | `system_dal_connection`        | distinct from `connection`                                                                                                 |
| `dalSensitivityLevel`  | `system_dal_sensitivity_level` |                                                                                                                            |
| `apigwRoute`           | `system_apigw_route`           |                                                                                                                            |
| `apigwFilter`          | `system_apigw_filter`          |                                                                                                                            |
| `project`              | `system_project`               | ⚠ prefix-shadows the 8 below                                                                                               |
| `projectTask`          | `system_project_task`          |                                                                                                                            |
| `projectFeature`       | `system_project_feature`       |                                                                                                                            |
| `projectBacklogItem`   | `system_project_backlog_item`  |                                                                                                                            |
| `projectIncident`      | `system_project_incident`      |                                                                                                                            |
| `projectReview`        | `system_project_review`        |                                                                                                                            |
| `projectPrivacy`       | `system_project_privacy`       |                                                                                                                            |
| `projectAiSystem`      | `system_project_ai_system`     |                                                                                                                            |
| `projectFriaScenario`  | `system_project_fria_scenario` |                                                                                                                            |

`discovery_search` is not in this table: it is a legacy two-segment name in
package `system/agentic`, hardcoded in eight places in `executor.go`, and
grandfathered rather than renamed.

`system_theme` is not in the table either, for the opposite reason: there is no
`theme` REST resource to map from. The tools read and write one settings value,
`ui.studio.themes`, so the segment is named for what the caller is changing —
the webapp's theme — rather than derived from a resource that does not exist.
Any later settings-backed family should follow the same rule and be recorded
here, so the next author finds a precedent instead of inventing one.

A `system_user_group_member_remove` is absent because the operation does not
exist: `userGroup` has `memberAdd` and `memberList` and no removal, in the
service, in `rest.yaml` and in the JS client. A user carries one `UserGroupID`,
so removal would mean setting it to zero — and a group-less user has the org
tree turn Inherit into Deny, leaving them only what a role rule explicitly
allows. Moving a user is `member_add` on the new group. `system_role` is the
other shape and does have `member_remove`: a user may hold any number of roles.

`system_skill` follows that rule. It maps to no REST resource either: the
skills are markdown files embedded from `system/agentic/skills/library`, read
by the agentic runtime and — through `system_skill_lookup` — by any MCP client.
The segment is named for what the caller reads. It is a `lookup` only; a skill
is changed by editing the repository, not the instance.

## Prefix shadowing

Two families produce names where one resource's prefix is another's:

- `compose_page_*` and `compose_page_layout_*`
- `system_project_*` and the eight `system_project_<sub>_*`
- `automation_taq_*` and `automation_taq_construct_*`
- `automation_workflow_*` and `automation_workflow_function_*`

This is safe for dispatch — nothing parses names — but these are exactly the
pairs a model is most likely to confuse when selecting a tool. Descriptions in
both families must state what the tool is _not_: a page layout is not a page,
a project task is not a project.

## Not in this table

Non-resource REST surfaces with no codegen'd resource — `stats`, `dml`,
`expression`, `locales`, `sink`, `settings`, `data_privacy` — are out of scope
for the resource grammar. Where they need tools, the shape is decided case by
case and recorded here first.

`permissions` is the first of those to get tools, so its shape is recorded
here. The segment is `system_permission`, singular, under `system` even though
rules span three components — RBAC is one mechanism with one vocabulary, and a
`compose_permission_*` family beside a `system_permission_*` family would make
a caller pick a tool by the resource they happen to be granting on. The
component is a property of the resource string, resolved from it, never a tool
name. Four ops: `schema` (a domain op, the grantable catalogue, precedent
`compose_page_block_schema`), `lookup`, `grant` and `revoke` — the last two
domain ops rather than `update`/`delete`, because that is the product's own
word for them and because revoke is `destructive` while grant is `write`.

**The one rule that makes these tools safe to have: an agent may grant only
what it already holds.** The service layer's `Grant` checks a single
component-wide `grant` permission and then writes whatever rule it is handed —
that is what the admin UI needs and it is an escalation path for an agent, so
the tools add a per-rule check of the caller's own access on the same resource.
Any later tool that writes an RBAC rule inherits that rule; a tool that passes
a rule straight to `Grant` must not be written.

Federation is out of scope: its services initialise conditionally and after the
MCP registry is wired, and its REST routes do not load unless
`FEDERATION_ENABLED` is set. A federation resource is refused by name.
