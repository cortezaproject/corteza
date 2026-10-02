---
name: module_building
description: Rules for creating and changing Compose modules — the field shape, the closed list of kinds, what locks once records exist, and the config a developer applies unasked.
triggers:
  - compose_module_lookup
  - compose_module_create
  - compose_module_update
  - compose_module_delete
  - compose_module_undelete
---

# Modules

A module defines the structure of a record, like a table: named, typed fields.
Fields are always referred to by name. Resolve the namespace first with
`compose_namespace_lookup`; never ask the user for an ID.

## Creating a module

- Call `compose_module_create` directly with the name, handle and fields the
  user described. Do not ask whether an existing module could serve instead.
- Handles are identifiers: lowercase letters, digits and underscores.
- `detail` decides what the write echoes back: `summary` (the default) confirms
  what was written; `full` adds field IDs, timestamps and storage config, which
  is thousands of tokens and only needed for the IDs.

## The field shape

Each field is `{name, kind, label, isRequired, isMulti, options}`. `kind` is
one of String, Number, Bool, DateTime, Select, Email, Url, File, User, Record,
Geometry — nothing else is accepted. The web app shows friendlier words (Text,
Checkbox, Date and Time, Dropdown, Location); the kind is the word above.

`options` depends on the kind:

- Select: `{"options": [{"value": "open", "text": "Open"}, …]}`. Filters and
  records compare the **value**, never the text.
- Record: `{"moduleID": "<module>", "labelField": "name"}` — the module the
  field points at and the field shown for it.
- Number: `precision`, `format`; DateTime: `onlyDate`, `onlyTime`.
- A multi-value field is `isMulti: true`; its values are arrays on records.

Expressions on a field (`value`, `validators`, `sanitizers`, `formatters`) have
their own rules in the `field_expressions` skill; a module is stored even when
one is broken, and the write reports it under `issues`.

## Changing a module

`compose_module_update` merges `fields` by name — fields you send are added or
updated, the rest are kept, and removing one means naming it in `removeFields`.
`config` replaces whole. Once a module holds records, a field's **kind and name
are locked**: retyping or renaming one is refused, so add a new field and move
the data instead.

## Auto-applied configuration

You act as a developer on the user's behalf; apply what a developer would
without being asked:

- Modules holding contacts, leads, customers, suppliers, or anything identified
  by email or phone: duplicate-detection rules on those fields (`ignore-case`
  for email; `case-sensitive` or `fuzzy-match` for phone).
- Modules holding transactions, orders, invoices, contracts, or anything whose
  change history matters: record revisions on.
- Modules holding personal information (names, addresses, health or financial
  data): a privacy disclosure saying how the data is used.

Do not add config to simple lookup or reference modules (a product category
list).
