# Server gotchas

Facts about the Go backend that are easy to get wrong: store and DAL, RBAC and auth, the REST surface, envoy and provisioning, settings, codegen, and integration tests.

## Store and DAL

### store.Search ignores the context scope

`store.Search<X>(ctx, s, filter)` does not apply the project/tenant scope from the context; `scope.GetScopeFromContext` is a service-layer concern. Code calling the store directly sees every project's rows unless the filter narrows them. Each generated filter (`server/store/adapters/rdbms/filters.gen.go`, plus `filter.go` for roles) applies `ProjectID` and `NamespaceID` only when `> 0`, so a zero value means every project.

**How to apply:** direct store calls set `NamespaceID` or `ProjectID` on the filter explicitly; never rely on ctx scope.

### compose_record has two models

The generated `Record` model (`server/compose/model/models.gen.go`, from `record.cue`) creates the table and decides which columns exist and are NOT NULL. `ModuleToModel` / `moduleSystemFieldsToAttributes` (`server/compose/service/module.go`) builds the per-module model that writes rows. Nothing keeps them in sync: a column missing from the second is never written, and if it is NOT NULL every insert into a fresh table fails.

**How to apply:** after changing `record.cue` or a scoped-table migration, update `moduleSystemFieldsToAttributes` too. `TestModuleSystemFieldsCarryScopeColumns` guards the scope columns.

### A new attribute needs an upgrade fix

`createTablesFromModels` (`server/store/adapters/rdbms/upgrade.go`) creates a table only when it is absent and never diffs an existing one. A new attribute added via `.cue` + `make codegen` exists in Go and every query but not in an existing database; boot fails in provisioning with `pq: column "<x>" does not exist`.

**How to apply:** add an entry to `fixesPre`/`fixesPost` in `server/store/adapters/rdbms/upgrade_fixes.go` that calls `addColumn(ctx, s, <model>.Ident, <model>.Attributes.FindByIdent("<Attr>"))`. Prove it on a worktree DB (cloned, populated), not a fresh one.

### The DDL index surface

`TableLookup(...).Indexes` is always empty on every driver — the lookup fills Columns only. Postgres `IndexLookup` works. The upgrade only ever adds indexes the model declares and never drops ones it stopped declaring, so renaming an index in cue needs a hand-written fix in `upgrade_fixes.go`, or migrated databases keep enforcing the old one. `dropIndexes` (`upgrade.go`) attempts the drop and tolerates failure; postgres emits `DROP INDEX IF EXISTS`.

### Two cursor builders

`store/adapters/rdbms/cursor.go` builds paging cursors two ways: `CursorCondition` (raw SQL, most resources) and `CursorExpression` (any resource with a `json:` sortable — workflows, TAQs, agents — and the compose DAL). Both handle the `coalesce` and `isnull` sort modifiers.

**How to apply:** a new sort modifier touches `filter.SortExpr.SetModifier`, `generateSorting`, both cursor builders and the collect switch in `rdbms.go.tpl` (then `make codegen`). Test paging by comparing a cursor walk with one big page.

### Timestamps in cursors are normalised to UTC

RFC3339 writes a zone offset as hours and minutes only. Go's zero time in a zone with an LMT offset such as `+00:58:04` round-trips seconds off, which breaks the cursor's `column = value` branch for rows sharing a sort value: ascending lists come back short with no error, descending ones loop. `filter.PagingCursor.SetModifier` (`server/pkg/filter/pagination.go`) converts `time.Time` to UTC; both `Set` and `PagingCursorFrom` go through it.

**How to apply:** anywhere else a timestamp round-trips through RFC3339, convert it to UTC first.

### Virtual JSON sort attributes

Eleven resources keep their display name in the JSON `meta` column (`meta.name` or `meta.short`) and declare a virtual `name` attribute with `sortableJSON` and `store: false`. Codegen turns it into a `"json:meta.short"` sortable and, via `queryJSON` in `server.store.cue`, a `queryJSONExpr` query filter. `storeIdent` defaults to the attribute name, so a virtual attribute looks like a real column to codegen that ignores `sortableJSON`; a bare `"name"` in a `query:` list fails at SQL time.

**How to apply:** filter structs are hand-written (`system/types/*.go`, `automation/types/*.go`), so a new `query` needs the struct field, a `rest.yaml` param and the `Query: r.Query` mapping in the REST handler.

### Unique-violation errors name their constraint

`store.ErrNotUniqueOn(resource, fields...)` (`server/store/errors.go`) is raised by the generated `check<X>Constraints` and by the postgres `errorHandler`, which names the table and constraint; the mysql, mssql and sqlite handlers return the bare `store.ErrNotUnique`. It keeps `errors.KindDuplicateData`, and returns `*errors.Error` so callers keep `.Wrap()`.

**Why:** a service-level duplicate and a database violation from a stale index have opposite fixes; the message tells them apart.

### store.AfterCommit

`AfterCommit(ctx, fn)` (`server/store/tx.go`) queues fn until the surrounding `store.Tx` commits, or runs it at once outside one; it is retry- and nesting-safe. Generated service wrappers own the transaction, so an event dispatched from a hand-written `onUpdate`/`onDelete` hook is seen by handlers before the commit — and for `onUpdate`, before the write.

**How to apply:** a dispatch whose handler re-reads the store (cache refresh, re-query) goes in `AfterCommit`. `system/service/role.go` uses it.

### In-memory SQLite runs Tx without a transaction

`sqlite.Connect` (and so `ConnectInMemory`) sets `TxRetryLimit: -1`, and `(*Store).Tx` then runs the body directly. Service tests on that store cannot reproduce a bug about reading uncommitted state.

**How to apply:** test the ordering with a synchronous stub dispatcher (the real eventbus dispatches in a goroutine) and confirm visibility against the dev server.

### goqu GROUP BY over a value-bearing expression

goqu binds every Go value as a placeholder, so `Select(caseExpr).GroupBy(caseExpr)` gives SELECT and GROUP BY different placeholders and postgres rejects it ("must appear in the GROUP BY clause"). SQLite accepts it, so a green sqlite store test hides the failure.

**How to apply:** compute the expression in a subquery and group the outer select by its alias, as `server/store/adapters/rdbms/custom_stats.go` does. `goqu.C("s.created_at")` quotes the dotted name as one column; use `goqu.I(...)`.

### Seeing the SQL the DAL sends

The postgres log is not readable by the shell user, but the `postgres` DB user is superuser: on a worktree DB set `log_min_duration_statement = 0`, terminate its backends, note `pg_stat_file(<log>).size`, make the call, read with `pg_read_file(<log>, <size>, 400000)`, then reset. Every record list runs the page query plus a keyset probe for the next page. The DAL sorts `ASC NULLS FIRST, id ASC NULLS FIRST`, which a default (NULLS LAST) btree cannot serve — declare such indexes `NULLS FIRST`. A Record/User field sorts by the referenced ID.

### Related labels and label sort

Design for `?refs=` labels and sort-by-displayed-label on records is `server/devdocs/related-labels.md`. Measured: in-process label resolution is 8–10 ms for 50 rows with an ID-set filter; sorting by label via a join at 1M rows is 430 ms against 205 ms for sort-by-ID.

**How to apply:** read the doc's open rulings before building past them.

### Generated action errors are all KindInternal

`NamespaceErrNotFound()`, `ModuleErrNotAllowedToUpdate()` and the rest of the
`*_actions.gen.go` constructors are `errors.KindInternal`; what kind of failure
they are is in their `type` meta (`notFound`, `notAllowedToRead`, `invalidID`,
`staleData`, `generic`). `errors.IsNotFound(err)` is false for every one of them.

**How to apply:** classify a service error by `(*errors.Error).MetaValue("type")`
first and by kind second; `mcp.ClassifyError` in `server/system/agentic/mcp` is
the one place that does it for the MCP surface.

### The OpenAPI page is served under the API prefix

`server/docs/*.yaml` and its Swagger UI are mounted at `/api/docs/`
(`server/app/servers.go`, `options.CleanBase(baseUrl, apiBaseUrl, "docs")`), not
at `/docs/`, where the webapp's SPA fallback answers with HTML and HTTP 200. The
files are generated from `server/*/rest.yaml` by `make codegen` and embedded, so
a new document reaches the server only after a rebuild.

### The construct library is empty in a unit test

TAQ step refs and trigger events are registered by the automation handlers'
constructors (`h.reg.AddFunctions` in `compose/automation`, `system/automation`
and the rest) when services boot. `service.ConstructLibrary()` in a plain `go
test` answers nothing; `tests/automation` sees it filled only because its
`TestMain` boots the app against a database. A reference that needs the
catalogue is generated there (`docs/reference/taq/*.gen.md`), never from a unit
test.

## RBAC and auth

### Auth sessions cache the user

At login the auth server gob-serialises the whole `types.User` into `auth_sessions.data` (`server/auth/request/auth_user.go`). Nothing refreshes it, so DB changes to a user (flags, roles) do not apply to active sessions until re-login. Force it with `delete from auth_sessions where rel_user = <id>`. Trusted auth clients skip the authorize Allow/Deny screen.

### The FE user comes from the token response

`$Auth.user` (`lib/vue/src/plugins/auth.ts`) is built from the OAuth token response, never from `GET /system/users/{id}` — a plain user cannot read their own user record. A user update replaces meta wholesale, so a meta field not carried in the token is wiped by an FE save of `$Auth.user`.

**How to apply:** a new per-user meta field the FE needs goes into the token response in `server/auth/handlers/handle_oauth2.go` and is mapped in both places in `auth.ts` (info and exchange).

### A permission refusal is invisible to clients

A REST handler's error is HTTP 200 (`errors.ServeHTTP`, `server/pkg/errors/http.go`); only the token and scope middleware answer with real status codes. In production `meta` and `stack` are stripped, so `meta.type: notAllowedTo*` exists only in dev; only the translated, locale-dependent message survives.

**How to apply:** gate a control on a `can*` flag from the resource payload, never on a caught error. Where no flag answers the question, an empty result is ambiguous and should not be labelled a permission problem.

### Group-less users are denied

`rbac.service.Check` (`server/pkg/rbac/service.go`) merges the role-rule verdict with the user-group org tree. When rules give Inherit the tree decides, and `MemberBranch` errors for `userGroupID = 0` (`pkg/rbac/user_groups.go`), which becomes Deny. `POST /system/users/` leaves `userGroupID = 0` unless the caller passes one; `setDefaultUserGroupRefs` backfills only at provisioning.

### Revoking a baseline grant is a deny on authenticated

Provisioning never deletes rules, so a long-lived DB keeps grants `000_base` has dropped. A `deny:` on `authenticated` revokes them safely: `check()` walks role kinds in order `ContextRole, CommonRole, AuthenticatedRole, AnonymousRole` (`pkg/rbac/ruleset_checks.go`) and returns on the first Allow, so a role's allow wins before the authenticated tier is read. Pinned by `ruleset_checks_baseline_test.go`.

**How to apply:** grant something everybody needs to a role, not to `authenticated`; reach existing installs by widening the partial gate (see Envoy and provisioning).

### Resource strings need every segment

An RBAC resource needs one segment per ID in its `RbacResource()` (`server/compose/types/rbac.gen.go`), wildcards included: `corteza::compose:module/<ns>/*`, `module-field/<ns>/*/*`, `record/<ns>/*/*`, `chart/<ns>/*`. A short string is refused ("invalid resource path structure"), and `CPermissionsDialog` swallows that into an empty dialog. Top-level resources are genuinely one segment.

**How to apply:** check with `dev/agent/api.sh GET '/compose/permissions/<roleID>/rules?resource=<urlencoded>'`; `[]` is fine, an `Error:` string is a dead button.

### The RBAC role registry refreshes only on role events

`getSessionRoles` (`pkg/rbac/roles.go`) skips a session role missing from `rbac.service.roles`. That list is written only by `UpdateRbacRoles` (`system/service/role.go`) at boot and on `system:role` after-create/update/delete events; the hourly `rbac.Watch` reloads rules, not roles. Deleted and archived roles are excluded; their rules and memberships outlive them.

**How to apply:** if role changes look one event behind, suspect this registry.

### rbac.Deny is the zero value

`server/pkg/rbac/permissions.go`: `Allow = 1`, `Deny = 0`, `Inherit = -1`. A missed map lookup or unset field reads as Deny, so `require.Equal(t, rbac.Deny, got[k])` passes when the rule is absent.

**How to apply:** assert presence (`v, ok := got[k]`), or compare `Access.String()` against a lookup that reports a miss distinctly.

### Auditing what a role can do

`accessControl.List()` in each `server/<component>/service/access_control.gen.go` is every evaluable `(resource, operation)`; its `"any"` field is the wildcard string provisioning writes. `GET /<component>/permissions/trace?roleID[]=<id>` evaluates exactly those roles without logging in. A rule naming a nonexistent operation or with too few path segments imports cleanly and never matches; `TestProvisionedRulesAreEvaluable` (`server/system/envoy/`) checks base provisioning. A context role whose `meta.context.resourceTypes` is empty is inert for every resource (`pkg/rbac/roles.go`).

### Contextual roles need a specific resource

`getSessionRoles` skips every `ContextRole` when the resource has wildcards, and `Trace` returns `unknown-context` there. A trace on `type/*` never exercises a role's context expression.

**How to apply:** trace a specific resource and give the role a matching `meta.context.resourceTypes`; wildcard rules still match the specific resource.

### Can refuses wildcards; Trace reads them

`checkValidity` (`pkg/rbac/service.go`) returns false for any resource containing `*`, so `Can` and `/permissions/effective` report every op denied. `Trace` evaluates a wildcard: a rule at the same or broader scope matches, a narrower one does not. It skips the org-tree branch, so it can only under-report. `callerHolds` in `server/system/agentic/permission_handler.go` calls `rbac.CanPassOn` (`pkg/rbac/grant_ceiling.go`), which tries `Can` and falls back to `Trace`.

### Rule PATCH is per component

`PATCH /<service>/permissions/<roleID>/rules` validates every rule with that service's validator; mixing components fails with `unknown resource type`. Component-level resources need the trailing slash: `corteza::compose/`. Without `?resource=`, `GET .../rules` returns only rules on the component and type-wildcard resources (`accessControl.Resources()`).

**How to apply:** one PATCH per component, then a filtered GET per resource to confirm.

### Reaching a unify section as a restricted user

Opening `/workflow/:id/edit` needs `corteza::system:application/<ID>` `read` and `access`, `corteza::system/` `applications.search`, `corteza::automation/` `workflows.search` and `triggers.search`, and the workflow's `read` and `execute`. Without the application `read` the application list is empty, not an error, and the shell bounces to `/?denied=<section>`. Without `triggers.search` the editor renders no trigger nodes.

**How to apply:** on `/?denied=`, grant the application `read` before suspecting the router; for missing nodes check that kind's `*.search`.

### Undelete is authorised by CanDelete

Only compose Record, automation Workflow and NgAutomation have an `undelete` RBAC op and a `canUndelete*` flag. Every other resource authorises undelete with `CanDelete<Resource>`. A soft-deleted resource is readable with `deletedAt` set and its `canDelete*` flag.

**How to apply:** gate restore on `canUndelete*` for those three and `canDelete*` for everything else.

### ngAutomation execute is not enforced

`canExecuteNgAutomation` is declared and returned in the payload, but `ngAutomation.Exec`/`ExecAndWait` load the automation with `loadNgAutomation` and run no RBAC check at all — only tenant membership and the write capability (`checkScope`, `automation/service/ng_automation.gen.go`). The REST exec handler checks nothing either.

## REST API

### Update verbs differ per service

Compose updates are `POST` to the resource path (`/compose/.../record/{id}`, module, page, chart, namespace); `PUT` answers a bare HTTP 405 with no body. System updates are `PUT` (`/system/application/{id}`, `/system/roles/{id}`); `POST` there is the 405.

**How to apply:** a bare `HTTP 405` from `api.sh` means check the verb first.

### Undelete endpoints return OK, not the resource

Every generated `*Undelete` handler returns `api.OK()` → `{success:{message:"OK"}}`, so `restored || fallback` takes the truthy junk. Project archive/unarchive/update do return the resource.

**How to apply:** after an undelete, update the row the view already holds; check the `.gen.go` handler before assuming a return shape.

## Envoy and provisioning

### Envoy scope is compose-only

`getScopeNodes` in a non-compose `envoy/store_decode.gen.go` is a stub, so `ResourceFilter.Scope` is ignored and a scoped decode of a system or automation resource reads every row in the store. `matchup<X>` (create-vs-update on encode) is scope-aware for compose only; the system and automation ones match identifiers against every row in the store, so an encode updates a same-ID or same-handle row in any project. `schema.ProjectRefField` is a plain `ID`, not a ref, so envoy never repoints `ProjectID` on a clone. System branch copies are hand-written in `system/service/project_revision_clone.go`.

**How to apply:** narrow a non-compose decode with explicit `ResourceFilter.Identifiers` from a normal store query.

### envoyx.Identifiers drops zero

`Identifiers.Add` / `MakeIdentifiers` (`server/pkg/envoyx/node.go`) skip `""`, `0` and `"0"`. An "only these IDs" filter carried in Identifiers turns a bad ID into an empty set, which reads as no filter and selects everything.

**How to apply:** carry such a filter as `[]uint64`, as the record export does via `DecodeParams.Params["recordIDs"]` (`compose/envoy/store_decode.go`). Guarded by `TestRecordExportByID`.

### Provisioning never corrects an existing install

A full import runs only when the store holds zero RBAC rules (`canImportConfig`, `server/pkg/provision/config.go`). Otherwise `provisionPartialBase` (`partial.go`) re-imports `000_base` only when a named `baseMarkers` entry is missing, so a new base grant needs a marker. The `003_auth` partial directory is listed but `server/provision/003_auth/` does not exist; when its gate fires, the empty decode makes `collectUnimportedConfigs` return early and skip the later directories. Every import encodes with `envoyx.OnConflictSkip`, which keeps the stored row (roles included), and the application files also carry `skipIf: '!missing'`; only `contextRoleTypes` (`context_roles.go`) fills empty `meta.context.resourceTypes`.

**How to apply:** a change to an existing install's rules, rows or roles needs a migration or a widened gate.

## Settings

### Settings write key differs from read key

Writes (`PATCH /settings/`, `POST /settings/{key}`) address the kebab `kv:` name in `server/system/types/app_settings.go` (`ui.main-logo`); reads (`GET /settings/current`, `$Settings.get`) address the camel `json:` name (`ui.mainLogo`). A write under the read name stores a row nothing decodes, and the save reports success.

**How to apply:** take the write key from the `kv:` tag.

### Settings attachment upload

`POST /system/settings/{key}` with `-F upload=@file` needs curl, not `api.sh` (it forces JSON). A tiny synthetic PNG fails with `failedToProcessImage`, returned as plain text on HTTP 200, so the FE reports invalid JSON. Created attachment IDs are in `GET /system/actionlog/?resource=system:attachment&action=create`. Nothing tracks settings attachments: `DELETE /system/attachment/settings/{id}` leaves `server/var/store/system/{id}.png` and `{id}_preview.jpg` on disk.

## Codegen

### How cue codegen works

`make codegen` runs `cue eval server.*.cue | human-json-tpl-exec` over `codegen/assets/templates/` and regenerates the whole tree; check `git diff --stat` afterwards. Expression types and events come from `make codegen-legacy`, not `make codegen`.

- The rdbms aux `db:` tag uses the attribute key, not `storeIdent`; name the attribute after the column.
- `precision: 0` is falsy and omitted; codegen cannot emit `Precision: 0`.
- Service codegen emits the struct and method bodies; the constructor always lives in the hand-written companion file.
- Generated `Update` copies every field, zeroing ones the payload omits; guard with `service.omitUpdateFields`.
- `store/tests/all_test.go` calls a hand-written `test<Resources>` per resource; a missing one stops the whole package compiling.
- Expr types need hand-written `CastTo<Type>` and `Clone()`.
- Hand code never lives in `.gen.go`.
- Options regenerate alone from `codegen/server.options.cue`; never `cue fmt` `server/app/options/*.cue`. `defaultGoExpr: ""` is a no-op, checked by `pkg/options/envexample_check_test.go`.
- `make codegen` reorders `server/pkg/codegen/resource_schema.gen.json`; compare with `jq -S`.

## Tests

### Integration suites truncate whole tables

Suites under `server/tests/` built on `helpers.NewIntegrationTestApp` truncate tables between tests. The helper uses in-memory SQLite unless `INTEGRATION_DSN` names another database, and refuses one equal to `DB_DSN`.

**How to apply:** run with `LOCALE_DEVELOPMENT_MODE=true LOCALE_PATH=<repo>/locale`; against a real `INTEGRATION_DSN`, pass `go test -p 1`, since the packages share one database.

### go vet finds test-only compile breaks

`go build ./...` never compiles `_test.go`, so a test package referencing a renamed type stays broken unseen. `cd server && go vet ./...` reports it.

**How to apply:** run it after any type rename, and sweep for siblings when one broken test package is reported.
