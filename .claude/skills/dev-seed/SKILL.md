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
dev/agent/cleanup.sh              # delete namespaces THIS session created
dev/agent/cleanup.sh --all        # every recorded session's namespaces
```

Record CSV imports are not idempotent, so `seed.sh` skips a fixture whose
namespace already exists — use `--force` for a clean slate.

## Authoring / extending fixtures

- A fixture is `def.yaml` (namespace + modules + record CSVs via `source:`
  blocks) **plus `ui.json`** (charts + pages, applied by
  `dev/agent/pagebuild.py` — spec format in its header). Model new fixtures
  on `dev/fixtures/agent-sandbox/`.
- **Never put pages/charts in def.yaml**: envoy import leaves block refs
  unresolved (module/chart handles instead of IDs) → broken pages in the UI.
  Presentation always goes through ui.json/pagebuild.
- Envoy YAML limits (learned the hard way): auth clients can't be fully
  configured (validGrant dropped, impersonateUser unresolved); field keys in
  the generic decoder match lowercased Go field names case-sensitively.
- Import must go through `seed.sh`, not plain `server_cli import`: it stages
  the fixture *directory* without ui.json (CSV providers are only discovered
  on directory decode; .json would be mis-parsed as YAML) and then
  re-registers DAL models on the live server (CLI import bypasses the
  running process — without this, record queries fail with "model does not
  exist").

## Harvesting

When a good structure is built live (via API or UI), capture it back into a
fixture instead of rebuilding next time:

```sh
cd server && ./build/gin-bin --env-file .env export compose-namespace <slug> 2>&1
```

Review the YAML, place it under `dev/fixtures/<slug>/`, commit.
