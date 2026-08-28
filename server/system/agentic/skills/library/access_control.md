---
name: access_control
description: Roles, user groups and users — which of the two grants access, what a new role does not do, and the permission surface these tools do not have.
triggers:
  - system_role_create
  - system_role_update
  - system_role_member_add
  - system_role_member_remove
  - system_role_clone_rules
  - system_user_group_create
  - system_user_group_update
  - system_user_group_member_add
  - system_user_create
---

# Access control

Two hierarchies, and they are not interchangeable.

- A **role** is a permission container. Permission rules are granted to roles.
  A user or a user group added to a role receives everything that role allows.
- A **user group** is a tree that permissions flow **down**: a member of a
  group gets what the group and all its ancestors are allowed to do.

Both are privilege changes and both take effect immediately.

## There is no permission tool

Nothing on this surface reads, writes or lists a permission **rule**. You can
create a role and put people in it; you cannot grant that role a single
permission through these tools. Two consequences:

- A new role grants nothing. It is an empty container until somebody assigns it
  rules in the admin UI, or until you copy another role's whole rule set onto
  it with `system_role_clone_rules`.
- `system_role_clone_rules` is the only way to give a role permissions from
  here, and it **replaces** the target's rules rather than adding to them.

Say so plainly when a user asks for "permission to do X" — offer the role and
the members, and tell them the rules themselves are a UI step.

## A user belongs to exactly one group

`system_user_group_member_add` is a move, not an addition: whatever group the
user was in, they are no longer in, and they lose what came with it.

There is no `system_user_group_member_remove` and there is no removal operation
anywhere in the product to build one on. To take a user out of a group, add
them to a different one. Leaving a user with no group is worse than it sounds —
the group tree turns Inherit into Deny for them, so a group-less user gets only
what a role rule explicitly allows.

Roles are the opposite: `system_role_member_remove` exists, and a user may hold
any number of roles.

## Order of operations

1. `system_user_create` if the person does not exist. Only the email is
   required; a handle is derived when you omit it.
2. `system_role_create` — name and handle must each be unique.
3. `system_role_member_add` with the user's handle, email or ID.

`system_role_archive` takes a role out of use without deleting it, and its
permissions stop applying to its members. `system_role_delete` is soft and
`system_role_undelete` restores the role _and_ everything it granted.

## Auth clients are deliberately absent

There is no `system_auth_client_create`, and there is deliberately no tool that
reveals or regenerates a client secret. The service can do all three; the
omission is not an oversight. Creating a client mints a working credential and
returns it, and a credential that reaches a model's context is in the
transcript and the logs for good. Direct the user to the admin UI rather than
looking for another route.
