---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/js/src/cast.ts
  - lib/js/src/guards.ts
touched-by:
  - lib/vue
  - client/web/unify
tests:
  - lib/js/src/compose/types/record.test.ts
  - lib/js/src/compose/types/module.test.ts
  - lib/js/src/compose/types/namespace.test.ts
  - lib/js/src/compose/types/page.test.ts
  - lib/js/src/compose/types/module-field/base.test.ts
---

# Compose resource classes

## Intention

Client-side model layer for the compose subsystem: `Record`, `Module`,
`Namespace`, `Page`, `PageLayout`, page blocks, module fields, charts, and
compose event constructors. Apps and corredor user scripts work with these
classes, never with raw API payloads.

## Map

- `types/record.ts` — `Record`: the central class; versatile ctor accepts Record/Module/partial/value combos in either argument order.
- `types/module.ts`, `types/module-field/` — `Module` + one field class per kind (String, Number, Record, User, …); field classes own default values and multi-value semantics.
- `types/namespace.ts`, `types/page.ts`, `types/page-layout.ts`, `types/page-block/` — page tree; one block class per page-block kind, registered in `page-block/index.ts`.
- `types/chart/` — chart config classes (generic, funnel, gauge, radar).
- `types/revision.ts` — record revision types.
- `events.ts` — compose event/trigger constructors for the eventbus; `helpers/idgen.ts` — `tempID-…` UID generator for not-yet-saved entities; `helpers/interpolate.ts` — `interpolateTemplate` evaluates a JS template literal against record/user variables.

## Contract (what apps may rely on)

- `Record.values` is an **object keyed by field name** (`string | string[]` per field's `isMulti`); the raw API array shape (`[{name, value}]`) is normalized by the class both ways (`apply`/`toJSON`).
- A `Record` always has a `Module` with fields; module is frozen on set and cannot be swapped (moduleID/namespaceID mismatches throw).
- `Record.cleanValues` is the frozen snapshot from construction — dirty-checking in editors depends on it.
- IDs are strings (`NoID = '0'`); date fields are `Date` objects, cast from ISO8601 strings.
- `toJSON` emits the wire shape (values back to the raw array), so classes can be posted directly to the API.
- `interpolateTemplate(template, { record, user, recordID, ownerID, userID })` (exported from the barrel as `compose.interpolateTemplate`) evaluates `template` as a JS template literal via `new Function`; `template` must be trusted page-author configuration (a prefilter, URL, label, … typed into a page/field builder by an admin), never end-user input. The variables are passed as explicit function parameters rather than closed over from the enclosing scope: bindings referenced only inside the template string are invisible to bundlers and can be dropped or renamed when minified.

## When changing this

- Value normalization in `Record` is the contract half of memory's "raw API array vs object" split — viewers/editors across lib/vue and the apps assume the object shape.
- New field kinds need a class in `module-field/` **and** matching editor/viewer registration in `lib/vue`'s field registry.
- New page-block kinds must be registered in `page-block/index.ts` or pages containing them fail to deserialize.
