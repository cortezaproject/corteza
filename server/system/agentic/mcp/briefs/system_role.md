# Brief: `system_role`

Tool file: `server/system/agentic/role_tools.go`
Handler:   `server/system/agentic/role_handler.go`
Group: `configuring`

## 1. Service

`system/service/role.go` — `RoleService`, `sysService.DefaultRole`.

| Method | Signature | Becomes |
|---|---|---|
| `FindByAny` | `(ctx, identifier any) (*types.Role, error)` | single-fetch in `lookup` |
| `Search` | `(ctx, types.RoleFilter) (RoleSet, RoleFilter, error)` | list mode |
| `Create` / `Update` | `(ctx, *types.Role) (*types.Role, error)` | `create` / `update` |
| `DeleteByID` | `(ctx, ID uint64) error` | `delete` |
| `UndeleteByID` | `(ctx, ID uint64) error` | `undelete` |
| `Archive` / `Unarchive` | `(ctx, ID uint64) error` | domain ops |
| `MemberList` | `(ctx, roleID uint64) (RoleMemberSet, error)` | domain op, read |
| `MemberAdd` / `MemberRemove` | `(ctx, roleID, userID uint64) error` | domain ops, write |
| `MemberAddGroup` / `MemberRemoveGroup` | `(ctx, roleID, userGroupID uint64) error` | domain ops, write |
| `Membership` | `(ctx, userID uint64) (RoleMemberSet, error)` | **see traps** |
| `CloneRules` | `(ctx, ID uint64, cloneTo ...uint64) error` | domain op, write |
| `IsSystem` / `IsClosed` | predicates | not tools; use them in descriptions and guards |

Risks: `lookup`, `member_list` read; `create`/`update`/`undelete`/`archive`/
`unarchive`/`member_*`/`clone_rules` write; `delete` destructive.

## 2. Authorization (§8.6) — PASSES, with one exception

11 `ac.Can*` calls in `system/service/role.go`. But **`role.Membership` is on
the §8.6 deny-list** — verify it before writing a tool for it. If it has no
check of its own, do not write `system_role_membership`; file it as a gap.
`MemberList` is the checked equivalent and is what a caller usually wants.

## 3. Identifier strategy

`FindByAny` exists (`role.go:188`) — accepts ID, name or handle. Use a `role`
ref for lookup/update/delete/archive/member ops.

`undelete` takes a required `roleID`: `RoleFilter.Deleted` defaults to
`StateExcluded`, so a deleted role is not findable by name. Verify and say so
in the description.

Member ops take a `role` ref plus a `user` ref (resolve via
`DefaultUser.FindByAny`) or a `userGroup` ref — not raw IDs. A caller naming
roles and users by handle is the whole point.

## 4. Filter fields → params

`RoleFilter` (`system/types/role.go:10`):

| Field | Param? | Note |
|---|---|---|
| `Query` | yes | |
| `Name`, `Handle` | yes | |
| `MemberID` | yes, as `member` | a user ref, resolved to an ID |
| `UserGroupID` | yes, as `userGroup` | a group ref |
| `Deleted` | yes, `includeDeleted` boolean | |
| `Archived` | yes, `includeArchived` boolean | |
| `Resource` | **no** | json:"-", internal to RBAC evaluation |
| `TenantID`, `ProjectID` | **no** | out of scope |

Plus `limit` / `pageCursor`.

## 5. List projection

`{roleID, name, handle, archivedAt, deletedAt}` — plus `isSystem` from
`IsSystem`, because whether a role is a system role is the single most
important thing a caller needs before touching it.

## 6. Traps

**System and closed roles.** `IsSystem` and `IsClosed` exist because some roles
cannot be edited or deleted. Every write tool's description must say that
system roles are refused, and the projection must expose `isSystem` so the
model can tell before it tries. Skipping this produces an agent that repeatedly
attempts an impossible edit.

**This is the privilege-escalation surface.** Roles carry permissions, so
`create`, `update`, `member_add` and `clone_rules` are how an agent could widen
its own access. This was ruled in scope deliberately (full CRUD,
destructive-tagged), so build it — but the descriptions must be plain about
what each op grants, and none of them should read as routine.

**`CloneRules` copies permissions wholesale** from one role to another. It is
the fastest way to escalate in the whole surface. Risk `write` by the letter of
§4, but say clearly in the description what it does.

**Archive is not delete and not suspend.** Three near-synonyms live in this
resource. The descriptions have to distinguish them explicitly.
