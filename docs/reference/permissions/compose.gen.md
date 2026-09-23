---
title: Compose permissions
outline: [2, 2]
---

<!-- This file is auto-generated from server/app and the component definitions. -->

# Compose permissions

Every operation is denied unless a role is granted it.

## Component

| Operation | Description |
| --- | --- |
| `grant` | Manage compose permissions |
| `settings.read` | Read settings |
| `settings.manage` | Manage settings |
| `namespace.create` | Create namespace |
| `namespaces.search` | List, search or filter namespaces |
| `resource-translations.manage` | List, search, create, or update resource translations |
| `email-notifications.send` | Send email notifications |

## Chart

| Operation | Description |
| --- | --- |
| `read` | read |
| `update` | update |
| `delete` | delete |

## Module

| Operation | Description |
| --- | --- |
| `read` | read |
| `update` | update |
| `delete` | delete |
| `record.create` | Create record |
| `owned-record.create` | Create record with custom owner |
| `records.search` | List, search or filter records |

## Module Field

| Operation | Description |
| --- | --- |
| `record.value.read` | Read field value on records |
| `record.value.update` | Update field value on records |

## Namespace

| Operation | Description |
| --- | --- |
| `read` | read |
| `update` | update |
| `delete` | delete |
| `export` | Access to export the entire namespace |
| `manage` | Access to namespace admin panel |
| `module.create` | Create module on namespace |
| `modules.search` | List, search or filter module on namespace |
| `modules.export` | Export modules on namespace |
| `chart.create` | Create chart on namespace |
| `charts.search` | List, search or filter chart on namespace |
| `charts.export` | Export charts on namespace |
| `page.create` | Create page on namespace |
| `pages.search` | List, search or filter pages on namespace |

## Page

| Operation | Description |
| --- | --- |
| `read` | read |
| `update` | update |
| `delete` | delete |
| `page-layout.create` | Create page layout on namespace |
| `page-layouts.search` | List, search or filter page layouts on namespace |

## Page Layout

| Operation | Description |
| --- | --- |
| `read` | read |
| `update` | update |
| `delete` | delete |

## Record

| Operation | Description |
| --- | --- |
| `read` | read |
| `update` | update |
| `delete` | delete |
| `undelete` | undelete |
| `owner.manage` | owner.manage |
| `revisions.search` | revisions.search |
