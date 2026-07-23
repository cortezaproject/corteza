---
kind: file
covers: Editor.vue
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/composables/useUnsavedGuard.ts
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# Node Editor view

## Intention

Create or edit a federation node (name, base URL, contact) and produce the
pairing URI that the other instance needs to complete the handshake.

## UX capabilities

- Validated form: name and baseURL required, baseURL must parse as a URL.
- Edit mode adds: generate-URI dialog (`nodeGenerateUri`, with copy-to-clipboard), per-node permissions button, and delete (hidden once deleted).
- Unsaved-changes guard (deep diff vs loaded copy); create redirects into edit mode; fetch failure redirects back to the list.

## Routes

`federation.nodes.create` at `/federation/nodes/new`; `federation.nodes.edit` at `/federation/nodes/:nodeID` (param watched — reload on change). Back navigates to `federation.nodes`.

## When changing this

- The generated URI is the input to the pair dialog on the other instance — keep it consistent with the server's pairing protocol.
- Payload sends only name/baseURL/contact; pairing state is server-owned.
- Unlike sibling editors, actions are not gated by per-item `can*` flags — only the permissions resource `corteza::federation:node/<id>` applies.
