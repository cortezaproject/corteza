# Brief: `system_user`

Tool file: `server/system/agentic/user_tools.go`
Handler:   `server/system/agentic/user_handler.go`
Group: `configuring` · Package: `system/agentic` (reuse the existing
`toolRegistrar` in `discovery_handler.go`; do not redeclare it)

## 1. Service

`system/service/user.go` — `UserService` interface, `sysService.DefaultUser`.

| Method | Signature | Becomes |
|---|---|---|
| `FindByAny` | `(ctx, identifier any) (*types.User, error)` | single-fetch in `lookup` |
| `Search` | `(ctx, types.UserFilter) (UserSet, UserFilter, error)` | list mode in `lookup` |
| `Create` | `(ctx, *types.User) (*types.User, error)` | `create` |
| `Update` | `(ctx, *types.User) (*types.User, error)` | `update` |
| `DeleteByID` | `(ctx, id uint64) error` | `delete` |
| `UndeleteByID` | `(ctx, id uint64) error` | `undelete` |
| `Suspend` / `Unsuspend` | `(ctx, id uint64) error` | domain ops |
| `SetPassword` | `(ctx, userID uint64, password string) error` | **do not expose** — see traps |
| `ToggleEmailConfirmation` | `(ctx, userID uint64, confirm bool) error` | domain op, optional |
| `CreateWithAvatar` / `UpdateWithAvatar` | take an `io.Reader` | **cannot** be tools — no way to pass a stream |

Risks: `lookup` read; `create`/`update`/`undelete`/`suspend`/`unsuspend` write;
`delete` destructive.

## 2. Authorization (§8.6) — PASSES

16 `ac.Can*` calls in `system/service/user.go`. Reads and writes are both
gated. No service-user elevation. Safe to build on.

## 3. Identifier strategy

`FindByAny` exists (`user.go:121`) and accepts ID, handle, email or username.
So `lookup`, `update`, `delete` take a `user` ref — name-or-handle-or-ID — not
a bare `userID`.

**Except `undelete`.** Check first whether `Search`/`FindByAny` can see a
deleted user: `UserFilter.Deleted` defaults to `filter.StateExcluded`. If a
deleted user is unreachable through `FindByAny`, `undelete` takes a required
`userID` string, and its description says where to get the ID — exactly as the
compose undeletes do. Verify, do not assume.

## 4. Filter fields → params

`UserFilter` (`system/types/user.go:77`):

| Field | Param? | Note |
|---|---|---|
| `Query` | yes | free text |
| `Email`, `Handle`, `Username` | yes | exact-match narrowing |
| `RoleID []string` | yes, as `role` | accept a role ref, resolve via `DefaultRole.FindByAny` |
| `UserGroupID` | yes, as `userGroup` | same treatment |
| `Kind` | yes | `UserKind`; enumerate the valid values in the description |
| `AllKinds` | yes, boolean | `mcp.WithBoolean` |
| `Deleted` | yes, as `includeDeleted` boolean | maps to `filter.StateInclusive` |
| `Suspended` | yes, as `includeSuspended` if present on the filter | verify the field exists |
| `TenantID`, `ProjectID` | **no** | tenancy is out of scope (§2) |
| `LabeledIDs`, `Labels`, `Check` | **no** | internal |

Plus `limit` and `pageCursor` (§8.1), wired through
`filter.NewPaging(page.Limit, page.Cursor)`.

## 5. List projection

Users carry meta and settings that are irrelevant when picking one from a list.
Project to `{userID, email, handle, name, kind, suspendedAt, deletedAt}`.
Single-fetch returns the raw service type.

## 6. Traps

**Never expose `SetPassword`.** It is a credential-setting operation with no
confirmation step, reachable by any agent in a write-capable session. It is not
on the §8.6 deny-list because the service authorizes it correctly — the
objection is that an LLM should not be setting passwords at all. If you think
it belongs, raise it rather than writing it.

**Email is an identity, not a field.** `Update` on a user whose email changes
may trigger confirmation flows. Read `beforeUpdate` before writing the update
tool, and say in the description what changing an email actually does.

**`Delete` vs `DeleteByID`** both exist. Use `DeleteByID`; `Delete` is the
older shape.

**Avatar methods are unreachable** from MCP — they take `io.Reader`. Do not
write a tool that pretends to accept an avatar path or URL.

**Suspend is not delete.** A suspended user still exists and still appears in
listings. The two descriptions must say which one the caller wants: suspend to
stop someone signing in, delete to remove them.
