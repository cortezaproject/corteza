---
name: record_handling
description: Rules for creating, looking up, updating, and filtering Compose records.
triggers:
  - compose_record_lookup
  - compose_record_create
  - compose_record_update
  - compose_record_delete
---

# Records

A record is a single data entry in a module. All business data lives in records — adding a lead, scheduling a meeting, creating a task all map to records in the relevant module.

When a user asks you to do anything involving a record, start tool calls immediately — do not ask questions first.

## Workflow

If the namespace and module IDs are listed in the **ACCESSIBLE NAMESPACES AND MODULES** section, use those `id` values directly — do not call `compose_namespace_lookup`. Before creating a record, always call `compose_module_lookup` first to get the field names, then create with the correct fields.

If the namespace or module is not in that section:

1. Call `compose_namespace_lookup` with `namespace` set to the namespace name or handle. Never pass a numeric ID here.
2. Call `compose_module_lookup` with `namespace` and optionally `module` (name or handle). Never pass numeric IDs here either.
3. Ask the user only for values that match fields returned by step 2.
4. Perform the operation.

## Rules

- If a TAQ exists for the operation, use it instead of `compose_record_create`.
- All values the user provides are field values for the target record. Never treat a name, person, or any other value as a reference to look up in another module — put it directly into the field as given.
- Call `compose_record_create` directly. Never call `compose_record_lookup` before creating, for any reason.
- Use values the user gives you directly — do not verify them with extra tool calls.
- Never ask for namespace IDs, module IDs, or record IDs — resolve them with tools.
- To find an existing record by value, use `compose_record_lookup` with a `filter`. Use `discovery_search` only for full-text or fuzzy search across modules.

## Filter syntax

- Equality: `fieldName = 'value'`
- Numeric: `fieldName = 42`
- AND/OR: `status = 'open' AND assignee = 'john'`
- Contains: `name LIKE '%john%'`
- String literals use **single quotes**. Variable references use no quotes.
