---
title: Home and apps
description: The home screen, the app menu, the application registry behind it, and notifications.
---

# Home and apps

The home screen at `/` is the first thing a signed-in person sees. It puts
three things side by side: the applications they can open, the Assistant, and
their notifications. The two outer columns can be resized by dragging, and
Human remembers the widths in that browser.

## The three columns

| Column            | What it holds                                                                                                      |
| ----------------- | ------------------------------------------------------------------------------------------------------------------ |
| **Menu**          | The applications available to you, with a search box to filter them by name.                                       |
| **Assistant**     | A chat with the agents your roles are allowed to see. Conversations are kept, and you can start new ones.          |
| **Notifications** | Messages addressed to you, split into **Unread** and **All**, with actions to mark them read, delete or mute them. |

Everywhere else in Human the same three things are reachable from the top bar.
The home screen shows them inline instead.

## The app menu

On a fresh install the menu lists **Projects**, **Agentic**, **Chatbot**,
**Automation (TAQ)**, **Namespaces** and **Admin Area**. Each entry opens one
part of the platform:

| Entry                | Opens                                                            |
| -------------------- | ---------------------------------------------------------------- |
| **Projects**         | [Projects](./projects)                                           |
| **Agentic**          | The agent list and editor. See [Agents](./agents).               |
| **Chatbot**          | Chatbots and their sessions. See [Chatbots](./chatbots).         |
| **Automation (TAQ)** | The TAQ list and builder. See [Automation](./automation).        |
| **Namespaces**       | Your apps. See [Namespaces and data](./namespaces).              |
| **Admin Area**       | Instance administration. See [Administration](./administration). |

The first account to sign up on a new instance gets every administrative role,
so it sees all of these.

People who may change applications can drag entries to reorder the menu. The
order is saved on the server, so it changes for everyone.

## Applications

Each menu entry is an **application**: a record in the application registry,
kept under **Admin Area → System → Applications**. Human ships with
applications for its own parts, and administrators can add their own that point
at any URL, including one of your namespaces.

An application has:

- A **Name** and **Description** for administrators.
- **Enabled**. When an application is disabled, users cannot open it. It still
  appears in the menu, greyed out.
- In the **App selector** panel: the **Name** and **Logo** people see in the
  menu, the **URL** it opens, and **Listed**. Only listed applications appear
  in the menu.

Two permissions decide what a person gets:

- **read** on the application decides whether it appears in their menu.
- **access** decides whether they can open it. An entry they can see but not
  open is shown greyed out, like a disabled one.

If someone follows a link into a part of Human they have no access to, they
land on the home screen with a message naming the address they asked for and
telling them to ask an administrator for access. This gate is a convenience, not the security
boundary: every request behind it checks its own permissions.

Human also ships two applications that are switched off and unlisted on new
instances. One of them is **Automation (Workflows)**, a separate automation
designer; these docs cover TAQs.

## Notifications

A notification is a short message sent to one person: a title, a description
and, for a record notification, a link that opens the record. TAQs send them
with the **Send Notification** and **Send Record Notification** steps, so an
automation can tell someone that a record needs their attention.

New notifications arrive live, without reloading. In the Notifications column
or panel you can:

- switch between **Unread** and **All**,
- **Mark as read**, **Mark as unread** or **Delete** one notification, or
  **Mark all as read**,
- **Mute notifications**,
- **Enable desktop notifications**, so your browser shows one when a
  notification arrives while you are not looking at Human.

## Where to go next

- [Namespaces and data](./namespaces): what the Namespaces application holds.
- [Agents](./agents): which agents appear in the Assistant, and why.
- [Administration](./administration): the application registry and the rest of
  the Admin Area.
