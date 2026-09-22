---
kind: folder
covers: recursive
owner: fe
depends-on:
  - lib/js/src/eventbus
  - lib/js/src/corredor
  - lib/vue/src/stores/usePageStore.js
touched-by:
  - client/web/unify/src/plugins/corredor.js
  - client/web/unify/src/sections/compose/components/PageBlocks/Shared/AutomationButtons.vue
  - lib/vue/src/components/corredor/CManualScriptButtons.vue
tests:
  - lib/vue/src/corredor/ui-hooks.test.ts
  - lib/vue/src/corredor/compose-ui.test.ts
  - lib/vue/src/corredor/bundle-loader.test.ts
  - lib/vue/src/corredor/script-display.test.ts
---

# Corredor in the browser

## Intention

What a web app needs to run Corredor automation scripts the way Corteza's apps
did: a bus scripts are dispatched on, a registry of the scripts a page may offer
as buttons, the context a client script executes with, and the loader that
brings a Corredor-built bundle in.

## Map

- `script-bus.ts` — `ScriptBusPlugin`: one lib/js `EventBus` per app as `$ScriptBus` (also provided under `ScriptBusKey`), with the well-known resource/event pairs. Distinct from `$eventBus`, the app's own on/off/emit bus.
- `ui-hooks.ts` — `UIHooks`/`Button` + `UIHooksPlugin` (`$UIHooks`): every `onManual` trigger of a registered script becomes a button carrying its uiProps (`app`, `page`, `slot`, `label`, `variant`); `Find(resourceType, page, slot, app)` answers what a spot may show. Built for a set of apps, since one webapp hosts what were Corteza's compose and admin apps.
- `compose-ui.ts` — `ComposeUIHelper`: what a compose client script reaches as `ComposeUI` — `gotoRecordViewer`/`gotoRecordEditor` (the module's record page, route `page.record`; the editor is `?edit=1`), `success`/`warning` toasts.
- `script-display.ts` — `constraintChips(constraints)`: how a trigger constraint reads on screen (`<name> <op> <value>`, equality when the script named no op); shared by the admin inventory and the builder.
- `compose-ctx.ts` / `webapp-ctx.ts` — script execution contexts over lib/js `Ctx`, cloned per event with `withArgs`; `logger.ts` — console-backed logger in pino's shape.
- `bundle-loader.ts` — `loadClientScripts` fetches `<bundle>-client-scripts.js` from the system API, evaluates it and registers each script's triggers on the bus with the ctx; `registerServerScripts` wraps `onManual` server scripts as handlers forwarding to the API and registers every script with `$UIHooks`.

## When changing this

- Names a script author sees — `ComposeUI` methods, ctx properties, the bundle global `<bundle>ClientScripts` — are the Corteza scripting API; renaming one breaks deployed extensions.
- Constraint names must match the server's compose events (`namespace`, `module`, …); the dispatcher supplies the matcher, since a record's module carries no namespace here.
- A missing or broken bundle never stops the app: the loader warns and registers nothing.
