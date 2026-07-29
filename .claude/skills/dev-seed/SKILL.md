---
name: dev-seed
description: Seed, reset, and harvest demo/test data on the local dev server using versioned envoy fixtures (dev/fixtures). Use when a task needs known test data, a clean slate, or when good agent-built structures should be captured as fixtures.
---

# /dev-seed

Versioned fixtures live in `dev/fixtures/<slug>/` (`def.yaml` + record CSVs),
where `<slug>` is the compose namespace slug (always `agent-` prefixed).
The `agent-sandbox` fixture ships modules (Contacts, Tasks with a Record
ref), records, a chart, and pages.

## Commands

```sh
dev/agent/seed.sh                 # import all fixtures (skips already-seeded)
dev/agent/seed.sh --force agent-sandbox   # delete + re-import one fixture
dev/agent/cleanup.sh              # delete ALL agent-* compose namespaces
```

Record CSV imports are not idempotent, so `seed.sh` skips a fixture whose
namespace already exists — use `--force` for a clean slate.

## Authoring / extending fixtures

- Model a fixture on `dev/fixtures/agent-sandbox/def.yaml`; record CSVs are
  wired via `source:` blocks (key column → cross-module refs by key).
- Envoy YAML limits (learned the hard way): auth clients can't be fully
  configured (validGrant dropped, impersonateUser unresolved); field keys in
  the generic decoder match lowercased Go field names case-sensitively.
- Import must go through `seed.sh`, not plain `server_cli import`: it imports
  the fixture *directory* (CSV providers are only discovered on directory
  decode) and then re-registers DAL models on the live server (CLI import
  bypasses the running process — without this, record queries fail with
  "model does not exist").

## Harvesting

When a good structure is built live (via API or UI), capture it back into a
fixture instead of rebuilding next time:

```sh
cd server && ./build/gin-bin --env-file .env export compose-namespace <slug> 2>&1
```

Review the YAML, place it under `dev/fixtures/<slug>/`, commit.
