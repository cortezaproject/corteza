---
kind: folder
covers: recursive
owner: fe
depends-on:
  - lib/vue
touched-by:
  - client/web/unify/src/main.js
tests: []
---

# Plugins / bootstrap sequence

## Intention

Own the entire app-setup pipeline (`setupAndAuthenticate`): auth handshake, API
plugin registration, settings, Pinia, router, i18n, theming, PrimeVue services.

## Contract

- Auth first; on `Unauthenticated` the auth flow starts and setup aborts silently — every other error propagates.
- API plugins registered here (`$SystemAPI`, `$ComposeAPI`, `$AutomationAPI`, `$DiscoveryAPI`, `$FederationAPI`) are the only API access path for the app.
- i18n loads the single merged `human-webapp` bundle in the user's preferred language; missing translations fall back to `{}`, never block boot.
- PrimeVue theme preset comes from user meta + studio themes setting; CSS layer order `tailwind-base, primevue, tailwind-utilities` is load-bearing for styling.

## When changing this

- Order is a contract (see `../intent doc`). Settings must be initialized before PrimeVue theming.
- New global plugins/services belong here, not in `main.js` or `App.vue`.
