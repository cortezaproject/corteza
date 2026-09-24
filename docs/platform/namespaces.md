---
title: Namespaces and data
description: Namespaces, modules and fields, records, pages and blocks, charts, and how people use an app built in Human.
---

<script setup>
import DataModel from '../.vitepress/theme/components/diagrams/DataModel.vue'
</script>

# Namespaces and data

Apps in Human are built in **namespaces**. A namespace holds one app's data
model, its records and the pages people use to work with them. You find them
under **Namespaces** in the menu.

<DataModel layouts />

## Namespaces

A namespace has a **Full namespace name**, which people see, and a **Short
name**, which becomes its address: `/compose/namespace/<short name>`. Opening a
namespace shows its pages in the sidebar.

Each namespace has two faces:

- **Its pages**, which is what end users see and work in.
- **Its admin panel**, where builders manage the namespace's **Modules**,
  **Pages** and **Charts**. The admin panel opens only for people allowed to
  manage the namespace. Being allowed to edit a single page or module is not
  enough on its own.

## Modules and fields

A **module** is one kind of thing the app keeps track of: a contact, a case, an
invoice. A module is defined by its **fields**. In the module editor each field
has a **Name** (used in filters and automations), a **Title** (shown in forms
and lists), a type, and whether it is **Multi** (holds several values) or
**Required**.

<Screenshot name="module-editor" alt="The module editor, listing the fields of a Deal module" />

| Field type            | Holds                                              |
| --------------------- | -------------------------------------------------- |
| **Text**              | Text, on one line, several lines, or as rich text. |
| **Number**            | A number.                                          |
| **Checkbox (Y/N)**    | Yes or no, shown as a checkbox or a switch.        |
| **Date and Time**     | A date, a time, or both.                           |
| **Select / Dropdown** | One of a list of options you define.               |
| **Email**             | An email address.                                  |
| **URL**               | A web address.                                     |
| **File**              | Uploaded files.                                    |
| **User**              | A person who can sign in to Human.                 |
| **Record**            | A link to a record in another module.              |
| **Location**          | A point on a map.                                  |

**Record** and **User** fields are how modules connect. A _Deal_ module with a
Record field pointing at _Company_ links every deal to its company; a User
field makes someone the deal's owner, so filters and permissions can follow
that person.

A module also has settings that apply to all of its records, among them
**Record Revisions** (keep a history of every change), **Unique values** (treat a
record as a duplicate when chosen fields match an existing one) and **Data store** (where the
records are kept: the built-in database or a data source an administrator has
connected). Permissions can be set on the module, on its records and on
individual fields.

## Records

A **record** is one entry in a module: one contact, one case. Records are
created, edited and deleted from pages, by TAQs, by agents and through the API,
and every one of those paths checks the same permissions. Deleted records can
be restored.

## Pages, layouts and blocks

<Screenshot name="page-builder" alt="The page builder, showing a dashboard of metrics and charts" />

A **page** is one screen in a namespace. Pages form a tree, which becomes the
navigation in the namespace's sidebar. There are two kinds:

- A **record page** belongs to a module and shows one record of it. The same
  page is used to view, create and edit a record.
- Every other page is free-form: a dashboard, a list, a landing page.

The module editor offers **Create record page** and **Create record list page**,
which generate a form for one record and a page listing them all, already
linked together.

Pages are designed in the **Page builder**. A page is laid out on a grid from
**blocks**, arranged in one or more **layouts**. Which layout a person sees,
and whether a given block is shown to them, can depend on conditions and on
their roles, so one page can serve different people differently.

| Block                | Shows                                                                           |
| -------------------- | ------------------------------------------------------------------------------- |
| **Record list**      | Records of a module in a table, with search, filters, sorting and bulk actions. |
| **Record organizer** | Records as cards in columns; moving a card updates the record's fields.         |
| **Record**           | The fields of the current record. Record pages only.                            |
| **Record revisions** | The change history of the current record. Record pages only.                    |
| **Chart**            | One of the namespace's charts.                                                  |
| **Metric**           | Headline numbers computed from records.                                         |
| **Progress bar**     | A value against a minimum and maximum, fixed or computed from records.          |
| **Calendar**         | Records with dates, on a calendar.                                              |
| **Map**              | Records with a location, on a map.                                              |
| **Content**          | Rich text you write.                                                            |
| **File**             | Files as a list or gallery, to preview or download.                             |
| **IFrame**           | Another web page, embedded.                                                     |
| **Navigation**       | Links to other pages or addresses.                                              |
| **Tabs**             | Other blocks, grouped into tabs.                                                |
| **Comments**         | Comments stored as records, with replies, reactions and attachments.            |
| **Automation**       | Buttons that run a TAQ.                                                         |
| **Agent Chat**       | A chat with chosen agents. See [Agents](./agents).                              |
| **Chatbot**          | An operator inbox for chatbot sessions. See [Chatbots](./chatbots).             |

Block titles, descriptions and many block settings can include values from the
current record or the current user, so a block on a record page can say, for
example, which customer it is about.

## Charts

A **chart** turns module data into a picture. Charts belong to the namespace,
are built in its admin panel, and are placed on pages with a **Chart** block.

**New Chart** offers four kinds: a generic chart, which draws as **Bar**,
**Line**, **Pie**, **Doughnut** or **Scatter**, and funnel, gauge and radar
charts. A chart reads one module, groups records by a dimension (a field, or a
date bucketed by week, month, quarter or year) and computes metrics such as the number of records
or the sum of a field. Clicking a data point can open the records behind it.

## How people use an app

An end user never sees modules or the page builder. They open the app from the
menu or from **Namespaces**, move between its pages in the sidebar, and work
with records:

- In a **Record list** they search, filter and sort records, open one, and,
  where the page allows it, select several to edit or delete together, import
  records from a file, or export them.
- On a **record page** they view a record, switch to editing it, and save.
  Links to related records open those records.
- Buttons in **Automation** blocks and on record lists run TAQs on the records
  in front of them.

What each person can see and change, down to single fields, is decided by their
roles. The [compose permissions reference](/reference/permissions/compose) lists
every operation on namespaces, modules, records, pages and charts.

<Screenshot name="record-list" alt="A list page of deals, as people using the app see it" />

<Screenshot name="record-page" alt="A deal's record page, with a related list of the company's other deals" />

## Where to go next

- [Automation (TAQ)](./automation): act on records when they change.
- [Agents](./agents): let an agent read and update records for people.
- [Core concepts](/get-started/concepts): the vocabulary in one page.
