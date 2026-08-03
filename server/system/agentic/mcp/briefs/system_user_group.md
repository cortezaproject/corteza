# Brief: `system_user_group`

Tool file: `server/system/agentic/user_group_tools.go`
Handler:   `server/system/agentic/user_group_handler.go`
Group: `configuring`

Tool names use the snake_case resource segment `user_group` per `RESOURCES.md`
— `system_user_group_lookup`, not `system_userGroup_lookup`. Note this makes
them four-segment names; that is expected and the structural test allows it.

## 1. Service

`system/service/user_group.go` — `UserGroupService`,
`sysService.DefaultUserGroup`.

| Method | Signature | Becomes |
|---|---|---|
| `FindByID` | `(ctx, userGroupID uint64)` | single-fetch by ID |
| `FindByHandle` | `(ctx, handle string)` | single-fetch by handle |
| `Search` | `(ctx, types.UserGroupFilter)` | list mode |
| `Create` / `Update` | `(ctx, *types.UserGroup)` | `create` / `update` |
| `DeleteByID` | `(ctx, ID uint64) error` | `delete` |
| `UndeleteByID` | `(ctx, ID uint64) error` | `undelete` |
| `MemberList` | `(ctx, userGroupID uint64) (UserSet, error)` | domain op, read |
| `MemberAdd` | `(ctx, userGroupID, userID uint64) error` | domain op, write |

Risks: `lookup`, `member_list` read; `create`/`update`/`undelete`/`member_add`
write; `delete` destructive.

## 2. Authorization (§8.6) — PASSES

8 `ac.Can*` calls in `system/service/user_group.go`. Verify `MemberList` and
`MemberAdd` specifically before writing those two — membership operations are
the ones most likely to be checked at the controller instead.

## 3. Identifier strategy

**There is no `FindByAny`.** The six services that have one are `role`,
`namespace`, `template`, `user_group`… — *verify this*: the conventions §8.4
list includes `user_group`, but the interface above shows only `FindByID` and
`FindByHandle`. If `FindByAny` genuinely does not exist on this service, write
a small local resolver that tries `strconv.ParseUint` → `FindByID`, else
`FindByHandle`, exactly as `taq_handler.go` does. Do not invent a `FindByAny`
that is not there, and do not silently accept only IDs.

`undelete` takes a required `userGroupID`: `UserGroupFilter.Deleted` defaults
to `StateExcluded`.

## 4. Filter fields → params

`UserGroupFilter` (`system/types/user_group.go:10`):

| Field | Param? | Note |
|---|---|---|
| `Query` | yes | |
| `Name`, `Handle` | yes | |
| `MemberID` | yes, as `member` | a user ref, resolved via `DefaultUser.FindByAny` |
| `Deleted` | yes, `includeDeleted` boolean | |
| `Archived` | yes, `includeArchived` boolean | |
| `LabeledIDs`, `Labels`, `Check` | **no** | internal |

Plus `limit` / `pageCursor`.

## 5. List projection

`{userGroupID, name, handle, archivedAt, deletedAt}`. Do **not** include
members — a group's membership is the heavy field and the reason `member_list`
exists as its own tool.

## 6. Traps

**There is no `MemberRemove`.** The interface has `MemberAdd` and no removal
counterpart, unlike `role` which has both. Do not write
`system_user_group_member_remove` against a method that does not exist. If
removal matters, that is a service gap to file, not a tool to improvise.

**Groups feed RBAC.** `rbacUserGroupService` (`user_group.go:50`) syncs group
membership into the permission graph, so `member_add` changes what someone can
do. Same escalation caution as roles — say so in the description.

**Archived is a distinct state** from deleted, and both are excluded by
default. A caller looking for "the group I archived last week" finds nothing
unless `includeArchived` is set; the lookup description must say that.
