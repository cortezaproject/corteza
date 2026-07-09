# DevDocs

## Codegen

The following parts of the system are or should be generated.
Most generated code focuses on the core part and the rest is scaffolding for the hand-written code

These are the note-worthy parts:

- REST layer: low-level HTTP stuff is fully generated, REST<->SVC interaction is scaffolded.
- Service layer: service functions are scaffolded, the actual low-level stuff is hand-written. The CRUD stuff is always generated (except if a thing is opted out). Custom functions are noted in .cue files and also generated.
- Resources: system resources as well as sub-resources (user meta, chart config, ...) are generated from .cue files.
  This also generates a bunch of supporting functions such as calculating diff, cloning, ...

- Supporting stuff:
  - EnvoyX (import export stuff)
  - RBAC rules and related access control stuff

### NG and OG codegen

Old codegen is based on yaml files, new one on cue files.
A bunch of things was migrated to NG but some bits are left on OG.

Largest bits left on OG:

- Generic resource helpers
- REST handlers and request processing
- Action log events

### How to migrate remained from OG to NG

- Action log events need to be moved to the .cue schema; nothing special
- Generic resource helpers can be moved to the .cue schema; no big deal here also
- Rest handlers and request processing need:
  - Expand the .cue schema to support additional REST endpoints outside of standard CRUD
  - Replace the lib/js codegen with a server-side codegen which hands off API clients to the lib/js package

## Resources

Resources (at least the core part of them) are fully generated in the .cue.

## Services

The scaffolding of each service is and must be fully generated from .cue
The scaffolding needs to cover access control, tenancy, and project access.

## REST Layer

Th
