---
title: Expressions
description: The expression language used in calculated fields, field validation and contextual roles.
---

# Expressions

Human has a small expression language for the places where a value has to be
computed or a condition tested. An expression is a single formula: it reads
the variables the place provides, and produces one value.

```
(price * quantity) * 1.22
```

## Where expressions are used

| Where                                             | The expression…                                     | Variables                                                                                                        |
| ------------------------------------------------- | --------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------- |
| **Field value expression** on a module field      | calculates the field's value when a record is saved | every field of the record by name; `new` (the record being saved); `old` (the previous version, empty on create) |
| **Value sanitization** on a module field          | rewrites the value before it is validated           | `value`                                                                                                          |
| **Value validation** on a module field            | rejects the value when it comes out **true**        | `value`, `oldValue`, `values` (all field values of the record)                                                   |
| **Contextual role** (Admin Area → System → Roles) | decides whether the role applies to a resource      | `resource` (the resource being accessed), `userID`, `agentID`                                                    |
| **While** loop in a TAQ                           | decides whether to run the loop again               | the TAQ's values: what the trigger carries and the results of earlier steps                                      |

::: tip Validation reads the other way round
A validation expression describes when a value is **wrong**. To allow values
from 0 to 5, write `value < 0 || value > 5`, with the message to show when it
comes out true.
:::

Branch conditions in TAQs are not written as expressions: they are built in
the TAQ builder.

## Values

| Kind    | Examples                     |
| ------- | ---------------------------- |
| Number  | `42`, `3.14`                 |
| String  | `"text"`                     |
| Boolean | `true`, `false`              |
| Array   | `[1, 2, 3]`, `["vip", "eu"]` |

There is no `nil` literal. To handle a value that may be missing, see
[Missing values](#missing-values).

## Variables and properties

Variables come from the place the expression runs in (see the table above);
you cannot declare your own. Reach into an object with a dot or with brackets:

```
lead.values.totalCost / 10      // 25, when totalCost is 250
lead["values"]["totalCost"]     // 250
```

## Operators

| Operator                    | Meaning                          | Example                         | Result                      |
| --------------------------- | -------------------------------- | ------------------------------- | --------------------------- |
| `+` `-` `*` `/`             | arithmetic                       | `price * quantity`              | `42`                        |
| `%`                         | remainder                        | `10 % 3`                        | `1`                         |
| `**`                        | power                            | `2 ** 3`                        | `8`                         |
| `+`                         | joins strings                    | `"Hello, " + trim(name)`        | `"Hello, Ada"`              |
| `==` `!=` `<` `<=` `>` `>=` | comparison                       | `value > oldValue`              | `true`                      |
| `&&` `\|\|` `!`             | and, or, not                     | `!(value > 5)`                  | `false`                     |
| `? :`                       | if-then-else                     | `quantity > 3 ? "bulk" : "one"` | `"bulk"` when quantity is 4 |
| `in`                        | array contains                   | `"vip" in tags`                 | `true`                      |
| `=~` `!~`                   | matches / does not match a regex | `name =~ "^Ada"`                | `true`                      |

Parentheses group as usual: `(price * quantity) * 1.22`.

## Functions

Functions are called by name with arguments in parentheses, and can be nested:

```
toUpper(trim(name))
round(price / 3, 2)
format("%s has %d items", "cart", quantity)
```

[Expression functions](./functions) lists every function, with an example.

## Missing values

Reading a variable or property that does not exist is an **error**, not an
empty value. Check first with `has()`, and use the if-then-else operator, which
only evaluates the branch it takes:

```
has(lead.values, "discount") ? lead.values.discount : 0
```

## Types

Values passed around inside Human, such as records, users and modules, have
named types with known properties. [Expression types](./types) lists them.
