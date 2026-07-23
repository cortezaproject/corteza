---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/js/src/api-clients
  - lib/js/src/compose
  - lib/js/src/system
  - lib/js/src/shared
touched-by: []
tests:
  - lib/js/src/corredor/args.test.ts
  - lib/js/src/corredor/ctx.test.ts
  - lib/js/src/corredor/exec.test.ts
  - lib/js/src/corredor/helpers/compose.test.ts
  - lib/js/src/corredor/helpers/system.test.ts
  - lib/js/src/corredor/helpers/shared.test.ts
---

# Corredor scripting support

## Intention

Runtime support for user-written automation scripts executed by the corredor
server (and client-script bundles): script execution context, argument
casting, and the high-level `Compose`/`System` helper classes scripts call
(`Compose.saveRecord(...)`, `System.findUserByEmail(...)`, …).

## Map

- `exec.ts` — `Exec(script, args, ctx)`: wraps sync/async script `exec` fns, normalizes results.
- `ctx.ts` — `Ctx`: what a script sees as `this` — casted args, logger, API clients, helper instances (uppercase = Human classes/helpers, lowercase = raw).
- `args.ts` / `args-human.ts` — `Args`/`ArgsProxy`: cast raw event args into Human classes per `HumanTypes` mapping.
- `helpers/compose.ts`, `helpers/system.ts` — scripting-facing CRUD/permission/notification helpers over the API clients; `helpers/shared.ts` — common plumbing (extractID, permission updater, list unwrapping).
- `shared.ts` — `BaseArgs` contract every script invocation receives.

## Contract (what user scripts may rely on)

- Helper methods accept flexible identifiers (object, ID string, handle) via `extractID`-style resolution and return resource classes, not raw payloads.
- `Ctx` naming convention is API: uppercase properties are safe, casted Human types; renaming any breaks deployed user scripts.

## When changing this

- This is a public scripting API executed against customer-authored scripts — treat every rename/signature change as breaking.
- Helper behavior must track API-client regeneration (endpoints come from server rest.yaml).
