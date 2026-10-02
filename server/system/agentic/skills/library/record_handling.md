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

1. Call `compose_namespace_lookup` with `namespace` set to the name, slug or ID (IDs always as strings).
2. Call `compose_module_lookup` with `namespace` and optionally `module` (name, handle or ID).
3. Ask the user only for values that match fields returned by step 2.
4. Perform the operation.

## Values

- Values are keyed by field name. A multi-value field takes an array; a
  structured value (a Geometry point) takes an object; Record and User fields
  take IDs as strings.
- A field you leave out is stored as **no value**, which is not `false` and not
  `""`: `done = false` does not find records whose `done` was never sent. Send
  every field a filter or a chart will group on.
- `compose_record_update` changes only the fields you send; a named field is
  replaced entire, so adding one value to a multi-value field means sending the
  whole list. `replace=true` writes the set whole and clears the rest.
- A duplicate-detection warning comes back on a **successful** save; read it.

## Rules

- If a TAQ exists for the operation, use it instead of `compose_record_create`.
- All values the user provides are field values for the target record. Never treat a name, person, or any other value as a reference to look up in another module — put it directly into the field as given.
- Call `compose_record_create` directly. Never call `compose_record_lookup` before creating, for any reason.
- Use values the user gives you directly — do not verify them with extra tool calls.
- Never ask for namespace IDs, module IDs, or record IDs — resolve them with tools.
- To find an existing record by value, use `compose_record_lookup` with a `filter`. `discovery_search` does full-text and fuzzy search across modules, but it is only registered when this instance has discovery configured — if it is not in your tool list, a `filter` is the only way.

## Filter syntax

- Equality: `fieldName = 'value'`
- Numeric: `fieldName = 42`
- AND/OR: `status = 'open' AND assignee = 'john'`
- Contains: `name LIKE '%john%'`
- Never set: `field IS NULL`; cleared to empty: `field = ''` — two different
  states, and `= ''` errors on a Number or DateTime column.
- Membership in a multi-value field: `tags IN ('red', 'blue')`.
- String literals use **single quotes**. Variable references use no quotes.
- A Select field is compared by its stored **value**, never by the label a
  person sees: `stage = 'offer'`, not `stage = 'Offer'`. The two differ by
  design, and the wrong one is refused with the legal values listed.
  `compose_module_lookup` shows every option.
