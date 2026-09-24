---
title: Administration
description: The Admin Area — people, roles and the permission model, the application registry, LLM providers, connections, email, templates and the action log.
---

<script setup>
import PermissionOrder from '../.vitepress/theme/components/diagrams/PermissionOrder.vue'
</script>

# Administration

The **Admin Area** is where an instance is configured as a whole: who can sign
in, what they may do, which applications they see, which models and external
services Human can reach, and what has happened. You find it under **Admin
Area** in the menu.

The Admin Area's sidebar is grouped into **System**, **Compose**,
**Automation**, **Federation** (only when federation is enabled) and **User
interface**. Each person sees only the entries their permissions reach, so an
administrator with a narrow role sees a short menu.

## People

### Users

A **user** is an account that can sign in. Administrators manage its identity
(name, email, handle), its credentials, the roles it is a member of, and its
lifecycle: a user can be suspended, which blocks sign-in without deleting
anything, or deleted and later restored.

The first account to sign up on a new instance gets every administrative role.

### Roles

A **role** is what permissions are granted to. Users are members of roles, and
a user's permissions are the combination of the rules on all of their roles.
Roles can be archived or deleted; the rules of an inactive role stop applying
and come back if it is restored.

Some roles are not assigned by hand:

- **Context roles** apply to a user only in relation to a particular resource.
  Human ships with **Owner**, **Creator**, **Updater** and **Deleter**, which
  apply to the user who owns, created, last updated or deleted a record,
  and **Agent Creator**, which applies to the agent that created a resource.
  They let you say, for example, "people may edit the records they own".
- **Automatic roles** are given to every signed-in user, or to every anonymous
  visitor. Which roles these are is set by
  [`RBAC_AUTHENTICATED_ROLES`](/reference/environment#RBAC_AUTHENTICATED_ROLES)
  and [`RBAC_ANONYMOUS_ROLES`](/reference/environment#RBAC_ANONYMOUS_ROLES).
- **Bypass roles** skip permission checks entirely: their members may do
  anything. They are set by
  [`RBAC_BYPASS_ROLES`](/reference/environment#RBAC_BYPASS_ROLES).

### User groups

**User groups** arrange users in a hierarchy, with every user in exactly one
group. Groups are used to organise people and to scope things to them, and a
group can carry roles, which then apply to its members.

## The permission model

<Screenshot name="permissions" alt="The permissions grid, with roles as columns and operations as rows" />

A **permission** is one operation on one kind of resource, such as _read_ on a
module or _execute_ on a TAQ. Each role can have a rule for each operation:

| Rule        | Means                                                   |
| ----------- | ------------------------------------------------------- |
| **Allow**   | Members of the role may do it.                          |
| **Deny**    | Members of the role may not do it.                      |
| **Inherit** | No rule; the decision is left to other rules and roles. |

A rule can apply to every resource of a kind (every module) or to one resource
(this module). When Human checks an operation, it decides like this:

<PermissionOrder />

Permissions are edited as a grid of roles against operations, where clicking a
cell cycles it through allow, deny and inherit. There is one grid per
component: **Admin Area → System → Permissions**, and the **Permissions**
entries under **Compose** and **Automation**. Most resources also have their
own permissions button, for rules that apply to that resource alone. The grid
can also evaluate permissions for a user or a combination of roles, to see
what they would actually be allowed to do.

Everything in Human is checked this way, including what TAQs and agents do:
they act as a user and get that user's permissions. The
[permissions reference](/reference/permissions/) lists every operation.

## Applications

**Applications** is the registry behind the menu on the home screen: each entry
has a name, a logo, a URL, and settings for whether it is **Enabled** and
**Listed**. See [Home and apps](./home#applications).

## AI

**LLM Providers** are the connections to model vendors that agents use: an
Anthropic, Mistral or OpenAI-compatible endpoint and its API key. See
[Agents](./agents#llm-providers).

## Connections and data sources

- **Connections** link Human to external services. A configured connection can
  add its operations as TAQ steps and its webhook events as TAQ triggers, so
  automations can work with that service. See [Automation](./automation).
- **Data Sources** are external databases that module records can be stored
  in, alongside the built-in one. A module's **Data store** setting picks where
  its records live. See [Namespaces and data](./namespaces#modules-and-fields).

## Sign-in and integration

- **Auth Clients** are OAuth2 clients: other applications allowed to
  authenticate against Human.
- **Auth Settings** covers how people sign in: internal accounts, password
  rules, multi-factor authentication and external identity providers.
- **Integration Gateway** defines custom API endpoints on Human, each a chain
  of filters that process the request, with a profiler to inspect traffic.
- **Messaging Queues** configures the message queues used by Human's event
  bus and automation.

## Email and templates

- **Email settings** holds the outgoing mail (SMTP) server, with **Test SMTP
  Server** to check it before relying on it for invitations, password resets
  and notifications. The mail server can also be set with
  [environment variables](/reference/environment#email-sending).
- **Templates** are reusable documents, such as emails and PDFs, built from
  content and partial templates.
- **Code Snippets** are pieces of code, such as scripts or markup, injected
  into the web app.

## Labels

**Projects (Labels)** shows every resource that carries a given label, across
namespaces, agents and TAQs, so resources that belong together can be found and
managed in one place.

## Action log

The **Action log** is a read-only audit trail: who did what, to which
resource, and when. It cannot be edited from the web app. Whether it is kept, and in
which database, is set with the
[action log environment variables](/reference/environment#actionlog).

## The other groups

- **Compose** holds settings and the permission grid for namespaces and
  everything in them.
- **Automation** lists TAQs, workflows and their sessions, and Corredor scripts,
  and holds the automation permission grid.
- **User interface** configures the look of the web app for everyone: theming,
  the navigation chrome, and the map provider used by location fields and map
  blocks.
