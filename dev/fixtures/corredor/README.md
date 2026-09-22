# Corredor fixture extension

A Corteza-shaped Corredor extension the worktree harness mounts for every
slot that runs Corredor (`worktree.sh up`). Its scripts target the
`agent-sandbox` namespace from `dev/fixtures/agent-sandbox` (seed it with
`dev/agent/seed.sh agent-sandbox`), so they stay inert on any other data.

Layout follows the Corredor extension contract:

- `server-scripts/` — run inside Corredor, triggered by server events or by
  a page button (`on('manual')`).
- `client-scripts/<bundle>/` — bundled by Corredor and executed in the
  browser; `compose` and `admin` are the bundles the webapp loads per section,
  `unify` is loaded in every section.

Each script is a plain Corteza automation script (`export default { label,
triggers, exec }`) and doubles as the parity check: if one stops working, a
Corteza extension would break the same way.
