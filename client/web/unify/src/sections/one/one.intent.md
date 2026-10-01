---
kind: folder
covers: recursive
owner: fe
depends-on:
  - lib/vue/src/components/navigation/CAppList.vue
  - server/provision/101_applications/0200_applications.yaml
touched-by:
  - client/web/unify/src/sections/index.js
tests:
  - client/web/unify/src/sections/one/one.access.test.js
---

# One section

## Intention

A full-page application catalog at `/one`: the instance logo, a search box, and
every application the user may open as a grid of cards. It is the landing page
an installation had before the webapps merged into one shell, kept so an
installation that still routes people there keeps working.

`home` is the landing page a new installation gets, and it lists the same
applications in its left column. `one` is therefore deliberately redundant, and
ships switched off — see Access.

## Contract (section values declared by index.js)

- `id: 'one'`; the single route `one` at `/one` is tagged `meta.section: 'one'`.
- `app: 'one/'` — the registry application gating entry.
- `sidebar: null`; `topbar.hideAppSelector` — the page is the app selector, so the topbar's button would open a menu of what is already on screen.

## Access

The provisioned `one` application is `enabled: false` and `unify.listed: false`,
and `one.access.test.js` holds it that way. Turning it on is a deliberate act by
an installation that wants this catalog; until then the section is routable but
unreachable, and `/one` renders the disabled screen.

## Map

- `index.js` — section contract: id, gate, single route, topbar override.
- `views/AppList.vue` — the catalog screen. Not a `fileTier` view, so it is governed here rather than by a sidecar.
- `one.access.test.js` — pins the gate to the provisioned application and the shipped-off default.

## Data touched

- `CAppList` (lib/vue) in its `grid` variant — this section is that variant's only consumer.
- `useBrandLogo('main')` (lib/vue) for the header logo, so it follows the active theme.
- Locale keys `one.apps.*` in the merged `human-webapp` bundle.

## When changing this

- The route path is fixed by the registry url `one/`: the app menu resolves `unify.url` to a section, so renaming one without the other makes the tile inert.
- Enabling the application by default is a product decision, not a cleanup — the test will fail first and is meant to.
