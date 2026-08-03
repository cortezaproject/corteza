# Brief: `system_application`

Tool file: `server/system/agentic/application_tools.go`
Handler:   `server/system/agentic/application_handler.go`
Group: `configuring`

## 1. Service

`system/service/application.go` — concrete `*application`,
`sysService.DefaultApplication`.

Read the file before writing: several methods are unexported (`onSearch`,
`onFlag`, `onUnflag`, `onReorder`, `beforeCreate`, `beforeUpdate`), so the
exported surface has to be confirmed from the `Default*` var's method set
rather than an interface declaration — there is no `ApplicationService`
interface.

Confirmed exported: `Delete(ctx, ID)`, `Undelete(ctx, ID)`. Confirm
`Search`, `FindByID`, `Create`, `Update`, `Reorder`, `Flag`, `Unflag` before
using them, and use `DeleteByID`/`UndeleteByID` if those exist alongside the
non-`ByID` forms — the rest of the codebase prefers the `ByID` shape.

Risks: `lookup` read; `create`/`update`/`undelete`/`reorder`/`flag`/`unflag`
write; `delete` destructive.

## 2. Authorization (§8.6) — verify before writing

6 `ac.Can*` calls in `system/service/application.go`. That is enough for CRUD,
but check the flag operations specifically: `checkFlag` (`application.go:212`)
suggests flags have their own access rule, and `onFlag`/`onUnflag` take an
`ownedBy` argument, which implies a per-user dimension. Read `checkFlag` and
report what it enforces. If flags turn out to be per-user state rather than
application state, they may not belong on the configuring surface at all.

## 3. Identifier strategy

**No `FindByAny` and no `FindByHandle`** — applications have a `Name`, and
`ApplicationFilter` has both `Name` and `Query`.

Write a local resolver: numeric → `FindByID`; otherwise `Search` with `Name`
set, erroring clearly on none or several. Do not fall back to `Query`, which is
a substring match and would silently pick the wrong application.

`undelete` takes a required `applicationID`: `ApplicationFilter.Deleted`
defaults to `StateExcluded`.

## 4. Filter fields → params

`ApplicationFilter` (`system/types/applications.go:9`):

| Field | Param? | Note |
|---|---|---|
| `Query` | yes | substring |
| `Name` | yes | exact |
| `Deleted` | yes, `includeDeleted` boolean | |
| `Flags []string` | yes | only if the flag ops are in scope after §2 |
| `ApplicationID []string` | **no** | the ref param covers single fetch |
| `LabeledIDs`, `Labels`, `FlaggedIDs`, `IncFlags`, `Check` | **no** | internal |

Plus `limit` / `pageCursor`.

## 5. List projection

`{applicationID, name, enabled, weight, deletedAt}`. Applications carry a
`Unify` config block describing how they appear in the launcher — heavy and
irrelevant when picking one from a list, so project it away and return it from
single-fetch only.

## 6. Traps

**`Reorder` takes an order of IDs** and, judging by `onReorder`, re-weights
what it is given. Read it against the page-reorder bug that was just fixed in
`store/adapters/rdbms/custom_compose_pages.go`, where reordering a subset
silently re-weighted everything: check whether `onReorder` has the same shape.
If it does, that is a bug to report, not a tool to write around.

**Flags may be per-user.** See §2. Until `checkFlag` is understood, do not
write `flag`/`unflag` tools — a tool that quietly writes per-user state while
appearing to configure a shared resource is worse than a missing tool.

**No interface to extend.** `DefaultApplication` is a concrete pointer. If a
needed method is unexported, raise it rather than adding one, and do not reach
into the store directly from a handler.
