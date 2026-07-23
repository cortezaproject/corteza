---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/js
  - lib/vue/src/libs/url.ts
touched-by:
  - client/web/unify/src/plugins
tests: []
---

# App bootstrap plugins

## Intention

The Vue plugins every host app installs at boot: authentication, API clients,
i18n, settings, event bus, toasts, and the global PrimeVue/shared component
registry. Together they define the injection surface (`$Auth`, `$SystemAPI`,
`$ComposeAPI`, `$AutomationAPI`, `$toast`, `$eventBus`, `$Settings`) that
stores and components in this lib assume exists.

## Map

- auth.ts — OAuth2 default-client flow + refresh-token lifecycle; provides `$Auth` with `accessTokenFn` and the authenticated user; derives auth/callback URLs from `window.HumanAPI`/base tag.
- event-bus.ts — minimal typed pub/sub; provides `$eventBus` (used e.g. for realtime message re-broadcast).
- human-api.ts — one install-plugin per API client (System/Compose/Automation/Federation/Discovery); base URL from `window.HumanAPI`, token from `$Auth`.
- i18n.ts — creates vue-i18n (Composition mode) from a pre-loaded translations bundle, `en` fallback.
- primevue-components.ts — globally registers PrimeVue components plus shared C\* inputs/permission components so SFCs skip per-file imports.
- settings.ts — `Settings` service around `settingsCurrent()`: dot-path `get(k, d)` and attachment-URL resolution; provides `$Settings`.
- toast.ts — wraps PrimeVue toast into `$toast` helpers (`toastSuccess/Warning/Info/Danger`, `toastErrorHandler`).

## When changing this

- Install order in the host app matters: auth before API plugins (token fn),
  PrimeVue's own toast service before ToastPlugin (it wraps
  `globalProperties.$toast`), settings before theming.
- Everything is provided under string keys (`'$SystemAPI'`, …) — renaming a
  key breaks every `inject` in stores/components.
- auth.ts persists flow state in sessionStorage/localStorage keys (`auth.*`);
  changing them invalidates in-progress logins.
