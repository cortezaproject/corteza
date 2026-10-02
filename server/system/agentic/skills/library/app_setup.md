---
name: app_setup
description: Order of operations when building or extending an application — namespace, modules, then pages, whether doing a full setup or adding one thing.
triggers:
  - compose_namespace_create
  - compose_namespace_lookup
  - compose_module_create
  - compose_module_lookup
  - compose_page_create
  - compose_page_lookup
---

# Application Setup Order

## Full setup

When building an application from scratch, complete each step fully before moving to the next:

1. Resolve or create the namespace (`namespace_building`).
2. Create all modules, Record fields last so the module they point at exists (`module_building`).
3. Create charts, which name a module by its numeric ID (`chart_building`).
4. Create pages: a list and a record page per module (`page_record_list`, `page_record`), then dashboards that place the charts (`page_building`, `page_layout`).
5. Automations last, once the fields they read exist (`taq_authoring`).

Never create pages before all modules exist. Pages reference modules — a page built before its module exists will be broken. A module has one record page; a second is refused.

A refused call carries a code: `not_found` means resolve the resource with its lookup first, `invalid_argument` means read the tool's parameter documentation, `forbidden` means the caller's roles do not allow it and no retry will.

## Partial requests

When the user asks to add just one thing to an existing app, check dependencies first:

- **"Add a page for X"** — look up whether module X exists before creating the page. If it doesn't exist, create the module first, then the page.
- **"Add a module"** — look up the namespace first, then create the module. Ask whether pages are needed for it.
- **"Set up pages for my app"** — look up what modules already exist in the namespace before creating any pages. Build pages only for modules that exist.

## In general

Before creating any page that references a module, confirm that module exists. Do not assume it was created earlier in the conversation — look it up.
