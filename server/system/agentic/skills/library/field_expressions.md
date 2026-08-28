---
name: field_expressions
description: What a module field's value, validator, sanitizer and formatter expressions may reference, and the one place the obvious reading is backwards.
triggers:
  - compose_module_create
  - compose_module_update
  - compose_module_lookup
---

# Field expressions

A field can carry four kinds of expression under `expressions`. They do **not**
share a scope, and two of them are transforms rather than tests. Writing one
from the shape of another stores clean and misbehaves later, at record save.

String literals need **double quotes**. `'hired'` is a syntax error.

## `value` — the field value expression

Computed from the record being saved. In scope:

| Name          | What it is                                                                                      |
| ------------- | ----------------------------------------------------------------------------------------------- |
| `<fieldName>` | every field of this module, by its **bare** name                                                |
| `new`         | the record being saved — `new.values.<field>`, `new.recordID`, `new.ownedBy`, `new.createdAt` … |
| `old`         | the same shape for the record before this save; on a create every field is null                 |

There is no `record` and no bare `values`. `record.values.stage` fails **every**
save with `unknown parameter record.values`, and the error names the record, not
the module.

```
stage != "hired" && stage != "rejected"      ✅
new.values.stage != "hired"                  ✅
record.values.stage != "hired"               ❌ unknown parameter record.values
values.stage != "hired"                      ❌ unknown parameter values.stage
stage != 'hired'                             ❌ could not parse string
```

Two consequences worth knowing:

- The expression **overwrites** whatever a caller sent for that field. Do not
  send it.
- `isRequired` on such a field means the **expression** must produce a value. An
  expression that evaluates to null on a required field refuses the save.

## `validators[].test` — reads backwards

In scope: `value` (the value being saved, as a string), `oldValue`, and
`values.<field>` for the record's other fields.

**`test` names the condition under which the value is REJECTED.** It is not a
statement of what is allowed, even though it sits under `validators`, is called
`test`, and is paired with an `error` describing the valid range.

```
{"test": "value >= 0 && value <= 5", "error": "Score must be between 0 and 5"}
  score 3   → REFUSED "Score must be between 0 and 5"
  score 7   → stored
```

Write the rule as its rejection instead:

```
{"test": "value < 0 || value > 5", "error": "Score must be between 0 and 5"}
  score 3.5 → stored
  score 9   → REFUSED "Score must be between 0 and 5"
```

A wrong-way-round validator never errors and never gets noticed: it admits
exactly the values it was meant to block and refuses the ones it was meant to
allow, while displaying a message asserting the opposite.

## `sanitizers[]` and `formatters[]` — transforms, not tests

Each sees only `value`, and its **result replaces the value**. A sanitizer runs
before validation and what it produces is what gets stored; a formatter runs on
the way out.

```
"sanitizers": ["trim(value)"], "formatters": ["toUpper(value)"]
```

## Check the write result

`compose_module_create` and `compose_module_update` store the module even when an
expression cannot work, and report what is wrong under `issues` — the field, the
slot, a severity and the expression. **No `issues` key is the clean result.** An
issue of severity `error` fails on every record save until it is fixed.

The check dry-runs each expression with every field empty, so it cannot catch a
misspelled bare field name — that resolves to null rather than failing. Look up
the module and read the field names back if a computed value comes out empty.
