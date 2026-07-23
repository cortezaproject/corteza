---
kind: folder
covers: '.'
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/routes.js
touched-by: []
tests: []
---

# System admin views

## Intention

Admin UI for the "system" server component: identity, access, integration and
platform configuration. Each child folder is one resource area; every folder
carries its own recursive INTENT.md. No loose files live directly here.

## Map

- User, UserGroup, Role — identity: accounts, groups, roles and memberships
- AuthClient, Permissions — OAuth2 clients; system-wide permission matrix
- Connection, DataSource — integration connections; DAL (database) connections
- ApiGateway, Queue — integration gateway routes (+ profiler); messaging queues
- Label, Template, Application — cross-app labels; render templates; app registry
- LLMProvider, CodeSnippets — LLM provider configs; code snippet settings
- Settings, Email — auth/system settings (+ auth/ providers); SMTP settings
- ActionLog — read-only audit log viewer
