---
kind: file
covers: translator.ts
backfilled: true
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/sections/compose/components/Translator/CTranslatorButton.vue
  - client/web/unify/src/sections/compose/components/Translator/CTranslatorDialog.vue
tests: []
---

# translator store

## Intention

Decouples the many "translate this resource" buttons from the single
CTranslatorDialog mounted once in ComposeHost: a button `open()`s the store
with a `TranslatorConfig`, the dialog renders whatever config is current.

## State owned

`visible` flag and the active `TranslatorConfig`: resource ID, per-key
`titles`, optional `highlightKey`, and — crucially — the caller-supplied
`fetcher()` / `updater(changes)` callbacks plus optional `keyPrettifier`.

## API surface consumed

None directly; all server I/O goes through the callbacks the opener provides
(module/page/namespace translators each bring their own endpoints).

## Consumers

CTranslatorButton (and the Module/Page/Namespace translator wrappers) open it;
so do the page Builder and module Edit views directly, for per-block and
per-field translation where there is no button to hang a wrapper on.
CTranslatorDialog reads it and calls the callbacks.

## Invariants

- `close()` clears `config` — the dialog must never act on a stale config,
  and callbacks must not be retained across openings.
- One dialog instance app-wide (in ComposeHost); openers must not mount their
  own.
