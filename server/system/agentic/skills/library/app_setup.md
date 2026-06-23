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

1. Create the namespace
2. Create all modules
3. Create all pages

Never create pages before all modules exist. Pages reference modules — a page built before its module exists will be broken.

## Partial requests

When the user asks to add just one thing to an existing app, check dependencies first:

- **"Add a page for X"** — look up whether module X exists before creating the page. If it doesn't exist, create the module first, then the page.
- **"Add a module"** — look up the namespace first, then create the module. Ask whether pages are needed for it.
- **"Set up pages for my app"** — look up what modules already exist in the namespace before creating any pages. Build pages only for modules that exist.

## In general

Before creating any page that references a module, confirm that module exists. Do not assume it was created earlier in the conversation — look it up.
