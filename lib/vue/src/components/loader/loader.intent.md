---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on: []
touched-by: []
tests: []
---

# loader/

## Intention

Full-screen boot loader shown while an app initializes (auth, settings,
stores), branded with the instance logo.

## Map

- `CLoaderLogo.vue` — fixed overlay (z-9999) pulsing `logoUrl`; `show` prop
  toggles it.

## Contracts consumers rely on

Renders nothing but the overlay — apps keep it mounted and flip `show` off
when bootstrapping completes. No logo URL → blank overlay (still covers the
app), so callers resolve the logo from settings before or during boot.

## When changing this

Must stay renderable pre-bootstrap: no stores, no injected APIs, no i18n.
