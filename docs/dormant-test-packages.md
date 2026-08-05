# Dormant test packages

Written 2026-08-05. Four Go packages had stopped compiling their test binaries.
`go test ./...` prints `[build failed]` for such a package and carries on, so
their tests had silently stopped running and nothing reported it.

`make test-compile` (wired into `make test`) now fails when this happens. It
runs `go test -run='^$' ./...`, which builds every test binary and runs none —
about 35 seconds for the tree.

## Resolved

| Package              | Cause                                                                                                          | Fix                                                                                                                                                                                                                 |
| -------------------- | -------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `federation/service` | `node`'s fields moved onto `nodeServices` in f64409476 (2026-07-13); `node_test.go` kept setting them directly | Construct through `nodeServices`. 9 of its 10 tests pass; the tenth panics — see task #13                                                                                                                           |
| `store/tests`        | `all_test.go` called `testProjectFriaScenarios`, which was never written                                       | Call removed with a comment. The FRIA types and the generated `ProjectFriaScenarios` store interface both exist, so writing `store/tests/project_fria_scenarios.go` mirroring `project_features.go` is all it needs |
| `tests/dal` (partly) | `Create` gained a metadata return; every call site still passed its single result to `NoError`                 | `errorOf` helper drops the metadata                                                                                                                                                                                 |

## Open — they target a type that no longer exists

`tests/dal` and `store/adapters/api/cred_registry` both reference
`types.ConnectionConfig` and `types.ConnectionConfigDAL`.

Commit **2853b5cdf** ("Update codegen types, actionlog, services", 2026-07-07)
replaced `ConnectionConfig` with `ConfiguredConnectionConfig`. **This was not a
rename.** The shapes have nothing in common:

```go
// before — what the tests still expect
ConnectionConfig{ DAL: &ConnectionConfigDAL{ ModelIdent: "compose_record" } }

// after — system/types/configured_connection.gen.go
ConfiguredConnectionConfig{ NamespaceID, DalConnectionID, CredentialID, Params, Discovery }
```

`ModelIdent` now lives in `system/types/dml_connection.gen.go` — the concept
moved into the Data Migration Layer.

So these tests assert against a pre-DML world. Repairing them means deciding
what they should assert against the DML model instead, which is a call for
whoever did that work — the engineering handover flags the area explicitly:
_"DML must use cache-only registration for external models"_, and _"Build-green
≠ behavior-preserved"_.

**Until then `make test` is red, deliberately.** These two packages have not run
since 7 July; a gate that passes while their coverage is off is the thing that
let this happen.

## Why this went unnoticed

Nothing failed. `go test ./...` exits non-zero for a build failure, but in a
tree this size the line scrolls past among thousands of `ok` lines, and no
tooling distinguished "this package's tests failed" from "this package's tests
did not run at all". Coverage can therefore drop to zero silently, which is
worse than a red test: a red test is information, a missing test is not.

`dev_test_run` in the developer MCP reports the same condition per-package —
it surfaces the compiler error rather than `[build failed]`, and flags an
aborted test binary — but that only helps someone already looking at that
package. The tree-wide gate is what catches it unprompted.
