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

**`Reorder` re-weights every application, and that is correct here.**
Investigated: `onReorder` (`application.go:199`) checks
`CanUpdateApplication` per ID, then `store.ReorderApplications`
(`store/adapters/rdbms/custom_applications.go:33`) searches with an empty
filter. That is *not* the page-reorder bug — applications are a flat global
list with no parent dimension, so "everything" is the right scope. The page bug
was a filter that failed to constrain to root pages; there is nothing analogous
to constrain to.

**But leftover ordering is non-deterministic, in both.** Applications not named
in `order` are re-weighted by ranging over a `map[uint64]bool`
(`custom_applications.go:70`), and Go randomises map iteration. So reordering a
subset silently shuffles everything else, differently each call.
`custom_compose_pages.go` has the identical pattern and it survived the scope
fix. The `reorder` tool's description must tell the caller to pass the complete
list, and this is worth fixing in the store for both resources.

**Flags are per-user *or* global, and the caller picks by an argument.**
Investigated: `checkFlag` (`application.go:212`) branches on `ownedBy` — zero
means a global flag gated by `CanGlobalApplicationFlag`, non-zero must equal
the caller's own identity and is gated by `CanSelfApplicationFlag`. So a flag
tool has two distinct modes with different permissions and different blast
radius, and a caller cannot flag on someone else's behalf.

If you write flag tools, the mode must be explicit in the schema rather than
inferred from an ID, and the description must say which one writes shared
state. A single tool that silently picks a mode is the outcome to avoid.

**No interface to extend.** `DefaultApplication` is a concrete pointer. If a
needed method is unexported, raise it rather than adding one, and do not reach
into the store directly from a handler.
