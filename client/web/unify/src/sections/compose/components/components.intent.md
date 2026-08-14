---
kind: folder
covers: .
backfilled: true
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/sections/compose/sidebar/ComposeSidebar.vue
tests: []
---

# Compose components

## Intention

Shared, non-route components of the compose section. Each family folder carries
its own intent doc; this doc indexes them and directly governs the three loose
sidebar files.

## Map

- Admin — module/page/chart configuration panels used by the admin editor views
- Chart — chart rendering (ChartRenderer) and report/chart-type config editors
- Common — RecordListFilter, the shared record-list filter builder
- ModuleFields — module field configuration dialog (per-kind option panels)
- Modules — ModuleImporter dialog
- Namespaces — NamespaceImporter dialog and namespace translator button
- PageBlocks — page block runtime components, their configurators, block registry, gridstack Grid
- Public — record Importer/Exporter dialogs for record lists
- Record — RecordModal, URL-query-driven record page in a dialog
- Reminders — reminder sidebar/list/edit/toast components over the reminder store
- Translator — resource-translation button/dialog/form trio over the translator store
- CSidebarNamespaceSwitcher.vue — namespace Select for the compose sidebar; navigates on switch, links to namespace manage/edit
- CSidebarNavigation.vue — sidebar search box plus two trees rendered together, each when it has items: the public page tree and the admin nav (modules/pages/charts). The current route only changes the search placeholder, never which tree shows
- CSidebarNamespaceNav.vue — the namespace editor's sidebar: a local name/short-name filter over every listable namespace, disabled ones included and badged, each row opening that namespace's editor. The "Namespaces" header is the link to the chooser (chevron expands, label navigates), so no separate entry for it. Disabled namespaces belong here precisely because this is where they get configured — the switcher, which navigates into a namespace, filters them out

## Data touched

Sidebar files read namespace, page, and module stores from `@planetcrust/human-vue`.
The admin nav builds `_id/_parentId/_label/_route` items; the page tree instead
points the shared tree at the page shape via `*-key` props (`pageID`/`selfID`/`title`).

## When changing this

New component families get their own folder + `<name>.intent.md`; keep this index
in sync. Record pages in the sidebar tree deliberately route to the new-record
creator (recordID '0').
