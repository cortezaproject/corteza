---
name: access_control
description: Roles, user groups and users — which of the two grants access, how a permission rule is written, and the ceiling on what an agent may grant.
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
  - system_permission_schema
  - system_permission_lookup
  - system_permission_grant
  - system_permission_revoke
---

# Access control

Two hierarchies, and they are not interchangeable.

- A **role** is a permission container. Permission rules are granted to roles.
  A user or a user group added to a role receives everything that role allows.
- A **user group** is a tree that permissions flow **down**: a member of a
  group gets what the group and all its ancestors are allowed to do.

Both are privilege changes and both take effect immediately.

## A new role grants nothing until you write a rule

Creating a role and adding members gives nobody anything. The rule is the
grant, and there are four tools for it:

- `system_permission_schema` — every resource type, the shape of its resource
  string, and the operations valid on it. Read it first; nothing here is
  guessable and a wrong string is refused rather than stored.
- `system_permission_lookup` — the rules a role holds, and `yourAccess`: what
  you personally hold on a resource.
- `system_permission_grant` — write allow/deny rules.
- `system_permission_revoke` — remove rules, back to inherit.

`system_role_clone_rules` is still there and still **replaces** the target's
whole rule set. Use grant for anything narrower.

## You can only grant what you hold

Every write checks the operation against your own access on that exact
resource, and refuses what you do not have — on revoke as well as on grant.
It is not a policy, it is the same RBAC evaluation the server runs when you
act yourself, so there is no argument that gets around it. You also need the
`grant` permission on the component; holding an operation is not the same as
being allowed to delegate it.

When a grant is refused, read `yourAccess` from `system_permission_lookup` for
that resource: that list is the ceiling, and the answer for the user is
usually "ask someone who has it", not another attempt.

## Writing a rule

A rule is a resource string, an operation, and allow or deny.

- The resource is the whole string: `corteza::compose:module/511/512`, not
  `module/*`. Every path segment counts — a module has two, a record has three.
  A component-wide rule ends in a slash: `corteza::compose/`.
- `*` in a segment covers every resource at that level, and you can only write
  one if you hold the operation at that same breadth.
- Operations live on the resource that owns them, which is rarely the obvious
  one: `record.create` is on the **module**, `module.create` is on the
  **namespace**. Letting a role into an application needs `access` and `read`
  together, on the application.
- Nothing here covers federation.

A rule takes effect immediately for new sessions. Someone already signed in
keeps the access their session was built with until they sign in again.

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

## Auth client secrets are shown once

`system_auth_client_create` returns the client secret the server generates, and
that is the only time anything shows it. There is no tool that reads a secret
back or regenerates one, on purpose: a credential already in use is not
something to pull into a transcript. A client whose secret has been lost is
replaced, not recovered — or a person reads it in the admin UI.

So record the secret where it is going before moving on, and do not repeat it
anywhere it does not need to be.

A `client_credentials` client acts as one nominated user and carries that
user's permissions, so `impersonateUser` decides what the credential can do.
Pick the narrowest account that can do the job.
