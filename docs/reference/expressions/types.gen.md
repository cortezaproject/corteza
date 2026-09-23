---
title: Expression types
description: The value types of the expression language and the fields of each structured type.
outline: [2, 3]
---

<!-- This file is auto-generated from server/automation/automation/expr_types.yaml, server/compose/automation/expr_types.yaml, server/pkg/expr/expr_types.yaml, server/system/automation/expr_types.yaml. -->

# Expression types

Every value in an expression, workflow variable or function argument has one
of these types. Structured types have fields, read with a dot:
`record.values.name`, `user.email`.

## Value types

`Action`, `Any`, `Array`, `Boolean`, `Bytes`, `ComposeRecordValueErrorSet`, `ComposeRecordValues`, `DateTime`, `DocumentType`, `Duration`, `EmailMessage`, `Float`, `Handle`, `ID`, `Integer`, `KV`, `KVV`, `LabelValue`, `Meta`, `RbacResource`, `Reader`, `RenderOptions`, `String`, `UnsignedInteger`, `Url`, `Vars`

## Structured types

### `Attachment` {#attachment}

| Field        | Type       | Notes     |
| ------------ | ---------- | --------- |
| `ID`         | `ID`       | read-only |
| `kind`       | `String`   |           |
| `url`        | `Handle`   |           |
| `previewUrl` | `Handle`   |           |
| `name`       | `Handle`   |           |
| `createdAt`  | `DateTime` | read-only |
| `updatedAt`  | `DateTime` | read-only |
| `deletedAt`  | `DateTime` | read-only |

### `ComposeModule` {#composemodule}

| Field         | Type         | Notes                      |
| ------------- | ------------ | -------------------------- |
| `ID`          | `ID`         | read-only; also `moduleID` |
| `namespaceID` | `ID`         | read-only                  |
| `name`        | `String`     |                            |
| `handle`      | `Handle`     |                            |
| `labels`      | `LabelValue` |                            |
| `createdAt`   | `DateTime`   | read-only                  |
| `updatedAt`   | `DateTime`   | read-only                  |
| `deletedAt`   | `DateTime`   | read-only                  |

### `ComposeNamespace` {#composenamespace}

| Field       | Type         | Notes                         |
| ----------- | ------------ | ----------------------------- |
| `ID`        | `ID`         | read-only; also `namespaceID` |
| `name`      | `String`     |                               |
| `slug`      | `Handle`     | also `handle`                 |
| `labels`    | `LabelValue` |                               |
| `createdAt` | `DateTime`   | read-only                     |
| `updatedAt` | `DateTime`   | read-only                     |
| `deletedAt` | `DateTime`   | read-only                     |

### `ComposePage` {#composepage}

| Field         | Type         | Notes                    |
| ------------- | ------------ | ------------------------ |
| `ID`          | `ID`         | read-only; also `pageID` |
| `moduleID`    | `ID`         | read-only                |
| `namespaceID` | `ID`         | read-only                |
| `title`       | `String`     |                          |
| `handle`      | `Handle`     |                          |
| `description` | `String`     |                          |
| `labels`      | `LabelValue` |                          |
| `createdAt`   | `DateTime`   | read-only                |
| `updatedAt`   | `DateTime`   | read-only                |
| `deletedAt`   | `DateTime`   | read-only                |

### `ComposeRecord` {#composerecord}

| Field         | Type                  | Notes                      |
| ------------- | --------------------- | -------------------------- |
| `ID`          | `ID`                  | read-only; also `recordID` |
| `moduleID`    | `ID`                  | read-only                  |
| `namespaceID` | `ID`                  | read-only                  |
| `values`      | `ComposeRecordValues` |                            |
| `meta`        | `Meta`                |                            |
| `ownedBy`     | `ID`                  |                            |
| `createdAt`   | `DateTime`            | read-only                  |
| `createdBy`   | `ID`                  | read-only                  |
| `updatedAt`   | `DateTime`            | read-only                  |
| `updatedBy`   | `ID`                  | read-only                  |
| `deletedAt`   | `DateTime`            | read-only                  |
| `deletedBy`   | `ID`                  | read-only                  |

### `HttpRequest` {#httprequest}

| Field      | Type     | Notes |
| ---------- | -------- | ----- |
| `Method`   | `String` |       |
| `URL`      | `Url`    |       |
| `Header`   | `KVV`    |       |
| `Body`     | `Reader` |       |
| `Form`     | `KVV`    |       |
| `PostForm` | `KVV`    |       |

### `QueueMessage` {#queuemessage}

| Field     | Type     | Notes |
| --------- | -------- | ----- |
| `Queue`   | `String` |       |
| `Payload` | `Bytes`  |       |

### `Reminder` {#reminder}

| Field         | Type              | Notes                        |
| ------------- | ----------------- | ---------------------------- |
| `ID`          | `ID`              | read-only; also `reminderID` |
| `resource`    | `String`          |                              |
| `snoozeCount` | `UnsignedInteger` | read-only                    |
| `assignedTo`  | `ID`              |                              |
| `assignedBy`  | `ID`              | read-only                    |
| `assignedAt`  | `DateTime`        | read-only                    |
| `dismissedBy` | `ID`              | read-only                    |
| `dismissedAt` | `DateTime`        | read-only                    |
| `remindAt`    | `DateTime`        |                              |
| `createdAt`   | `DateTime`        | read-only                    |
| `updatedAt`   | `DateTime`        | read-only                    |
| `deletedAt`   | `DateTime`        | read-only                    |
| `payload`     | `Bytes`           |                              |

### `RenderedDocument` {#rendereddocument}

| Field      | Type     | Notes |
| ---------- | -------- | ----- |
| `document` | `Reader` |       |
| `name`     | `string` |       |
| `type`     | `string` |       |

### `Role` {#role}

| Field        | Type         | Notes                    |
| ------------ | ------------ | ------------------------ |
| `ID`         | `ID`         | read-only; also `roleID` |
| `name`       | `String`     |                          |
| `handle`     | `Handle`     |                          |
| `labels`     | `LabelValue` |                          |
| `createdAt`  | `DateTime`   | read-only                |
| `updatedAt`  | `DateTime`   | read-only                |
| `archivedAt` | `DateTime`   | read-only                |
| `deletedAt`  | `DateTime`   | read-only                |

### `Template` {#template}

| Field        | Type                            | Notes                        |
| ------------ | ------------------------------- | ---------------------------- |
| `ID`         | `ID`                            | read-only; also `templateID` |
| `handle`     | `Handle`                        |                              |
| `language`   | `String`                        |                              |
| `type`       | `DocumentType`                  |                              |
| `partial`    | `Boolean`                       |                              |
| `meta`       | [`TemplateMeta`](#templatemeta) |                              |
| `template`   | `String`                        |                              |
| `labels`     | `LabelValue`                    |                              |
| `ownerID`    | `ID`                            | read-only                    |
| `createdAt`  | `DateTime`                      | read-only                    |
| `updatedAt`  | `DateTime`                      | read-only                    |
| `deletedAt`  | `DateTime`                      | read-only                    |
| `lastUsedAt` | `DateTime`                      | read-only                    |

### `TemplateMeta` {#templatemeta}

| Field              | Type     | Notes |
| ------------------ | -------- | ----- |
| `short`            | `String` |       |
| `description`      | `String` |       |
| `headerTemplateID` | `ID`     |       |
| `footerTemplateID` | `ID`     |       |

### `User` {#user}

| Field            | Type         | Notes                    |
| ---------------- | ------------ | ------------------------ |
| `ID`             | `ID`         | read-only; also `userID` |
| `username`       | `String`     |                          |
| `email`          | `String`     |                          |
| `name`           | `String`     |                          |
| `handle`         | `Handle`     |                          |
| `emailConfirmed` | `Boolean`    |                          |
| `labels`         | `LabelValue` |                          |
| `createdAt`      | `DateTime`   | read-only                |
| `updatedAt`      | `DateTime`   | read-only                |
| `suspendedAt`    | `DateTime`   | read-only                |
| `deletedAt`      | `DateTime`   | read-only                |
