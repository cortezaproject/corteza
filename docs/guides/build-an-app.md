---
title: Build an app
description: Model data in modules, link them, and build list, record and dashboard pages with charts.
---

# Build an app

This guide builds a small sales tracker: companies, and the deals you have with
them, plus a dashboard. Along the way it covers every part of an app. If you
have not run Human yet, start with the [Quickstart](/get-started/quickstart).

## Create the namespace

A namespace holds one app. Open **Namespaces** from the menu and click **New
Namespace**.

- **Full namespace name** is what people see, e.g. `Sales`.
- **Short name** becomes part of the address, e.g. `sales` gives
  `/compose/namespace/sales`. Keep it short and lowercase.

## Model the data

Open the namespace's **Modules** and click **New Module** for each of these.

**Company**

| Field    | Kind              |
| -------- | ----------------- |
| Name     | Text              |
| Website  | URL               |
| Industry | Select / Dropdown |

**Deal**

| Field      | Kind              |
| ---------- | ----------------- |
| Title      | Text              |
| Company    | Record            |
| Value      | Number            |
| Stage      | Select / Dropdown |
| Close date | Date and Time     |
| Owner      | User              |

Two field kinds do more than hold a value:

- A **Record** field links to another module. Point the deal's _Company_ field
  at the _Company_ module, and each deal can be opened from its company.
- A **User** field links to a person who can sign in. Use it for owners and
  assignees, so permissions and filters can follow them.

Give **Stage** the options `Lead`, `Proposal`, `Won` and `Lost`.

## Create the pages

In each module, click **Create record page**, then **Create record list page**.
That gives every module a list of its records and a form for one record, linked
together.

To change a page, open **Pages** in the namespace and choose **Page builder** on
it. Pages are laid out from blocks on a grid. Use **Add block** to add one, and
drag its corner to resize it.

## Add a dashboard

First make the charts. Open the namespace's **Charts** and click **New Chart**.

1. Pick **Generic chart**, then choose the _Deal_ module.
2. Under **Dimensions**, choose _Stage_. Under **Metrics**, pick **Number of records**,
   or sum _Value_.
3. Pick a chart type (Bar, Line, Pie, Doughnut or Scatter) and save.

Then click **New Page** in **Pages**, call it `Dashboard`, and open it in the
page builder. Add a **Chart** block for each chart, and a **Metric** block or
two for headline numbers such as open deal value.

::: tip Funnel, gauge and radar charts
**New Chart** also offers funnel, gauge and radar charts. A funnel over the
_Stage_ field is a quick way to show a pipeline.
:::

## Share it

People open an app at `/compose/namespace/<short name>`. To put it in everyone's
menu on the home page, an administrator adds it as an application:

1. **Admin Area → System → Applications → New Application**.
2. Give it a name and a logo, and set its URL to the namespace's address.
3. Turn on **Listed** so it appears in the menu.

Who can see and change what is set with roles. See the
[permissions reference](/reference/permissions/compose) for every operation on
namespaces, modules, records and pages.

## Next steps

- [Agents](./agents): give an agent the _Records_ tools and let it answer
  questions about your deals.
- [Core concepts](/get-started/concepts): how automations fit in.
