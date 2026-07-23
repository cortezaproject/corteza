---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/composables/useInternalLink.ts
touched-by: []
tests: []
---

# app/

## Intention

App-level full-screen states shown instead of the app itself.

## Map

- `CAppDisabled.vue` — full-viewport "this app is disabled" screen with a
  back-to-home link (origin URL, SPA-navigated via useInternalLink).

## Contracts consumers rely on

Prop-less; strings come from the `general.appDisabled.*` i18n keys, so the
hosting app's locale bundle must provide them.

## When changing this

Renders before/without the normal app shell — keep it dependency-light (no
stores, no injected APIs).
