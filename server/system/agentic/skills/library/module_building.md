---
name: module_building
description: Rules for creating Compose modules and applying sensible config (dedup, revisions, privacy).
triggers:
  - compose_module_lookup
  - compose_module_create
  - compose_module_update
  - compose_module_delete
---

# Modules

A module defines the structure of a record (like a database table). Fields are named, typed attributes — always reference fields by name.

## Creating a module

- Call `compose_module_create` directly. Never suggest that an existing module could serve the same purpose. Never ask if the user wants to use something else instead.
- Use the name, handle, and fields the user specifies directly — do not verify whether a similar module already exists before creating.
- Never ask for a namespace — resolve it with `compose_namespace_lookup`.

## Auto-applied configuration

You are acting as a developer on behalf of the user. Apply sensible configuration automatically — the user should not have to ask for these things:

- Modules that store contacts, leads, customers, suppliers, or any entity identified by email or phone: add duplicate detection rules on those fields. Modifier `ignore-case` for email; `case-sensitive` or `fuzzy-match` for phone.
- Modules that store transactions, orders, invoices, contracts, or any data where change history matters: enable record revisions.
- Modules that store personal information (names, addresses, health data, financial data): set a privacy disclosure describing how the data is used.

Do not add config to simple lookup or reference modules (e.g. a product category list).
