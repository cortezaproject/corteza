---
title: TAQ triggers
description: Every trigger a TAQ can start from, with the inputs it takes and the values it provides.
outline: [2, 3]
---

<!-- This file is auto-generated from the TAQ construct catalog (server/tests/automation/docs_catalog_test.go). -->

# TAQ triggers

A trigger decides when a [TAQ](/platform/automation) runs. These are the
triggers the builder offers on every server, grouped as in the trigger picker.
**Inputs** are what you set on the trigger; **Provides** lists the values it
hands to the steps that follow. The steps are in [TAQ steps](./steps).

Triggers that come from a connection, such as its webhooks, depend on the
connections configured on your instance and are not listed here.

## Agents

### Agent Invoked {#agent-invoked}

<div class="env-meta">

`automation:trigger:agentic` · `onAgentic`

</div>

Triggered by an Agent providing a variable payload.

The values the agent passes are declared per TAQ, in the trigger's input schema.

## Chatbot

### After Chatbot Step {#after-chatbot-step}

<div class="env-meta">

`corteza::system:chatbot` · `onAfterStep`

</div>

Triggered after a chatbot scenario step completes. steps contains state from all steps including the current one.

**Provides**

| Value | Type | Description |
| --- | --- | --- |
| `chatbotID` | `ID` | Chatbot ID |
| `sessionID` | `ID` | Session ID |
| `stepID` | `ID` | Step ID |
| `agentID` | `ID` | Agent ID |
| `steps` | Any | State of completed scenario steps keyed by scenarioID. Each entry: steps.&lt;scenarioID>.form.{fields,submitted} or steps.&lt;scenarioID>.conversation.{conversationID,handoff} |

### Before Chatbot Step {#before-chatbot-step}

<div class="env-meta">

`corteza::system:chatbot` · `onBeforeStep`

</div>

Triggered before a chatbot scenario step executes. steps contains state from all previously completed steps.

**Provides**

| Value | Type | Description |
| --- | --- | --- |
| `chatbotID` | `ID` | Chatbot ID |
| `sessionID` | `ID` | Session ID |
| `stepID` | `ID` | Step ID |
| `agentID` | `ID` | Agent ID |
| `steps` | Any | State of completed scenario steps keyed by scenarioID. Each entry: steps.&lt;scenarioID>.form.{fields,submitted} or steps.&lt;scenarioID>.conversation.{conversationID,handoff} |

## Records

### After Record Create {#after-record-create}

<div class="env-meta">

`compose:record` · `afterCreate`

</div>

Triggered after record is created.

**Inputs**

| Input | Type | Required |
| --- | --- | --- |
| Namespace | `ID`, `Handle`, `ComposeNamespace` | Yes |
| Module | `ID`, `Handle`, `ComposeModule` | Yes |

**Provides**

| Value | Type |
| --- | --- |
| `record` | `ComposeRecord` |
| `oldRecord` | `ComposeRecord` |
| `module` | `ComposeModule` |
| `namespace` | `ComposeNamespace` |
| `recordValueErrors` | `ComposeRecordValueErrorSet` |
| `selected` | Any |

### After Record Delete {#after-record-delete}

<div class="env-meta">

`compose:record` · `afterDelete`

</div>

Triggered after record is deleted.

**Inputs**

| Input | Type | Required |
| --- | --- | --- |
| Namespace | `ID`, `Handle`, `ComposeNamespace` | Yes |
| Module | `ID`, `Handle`, `ComposeModule` | Yes |

**Provides**

| Value | Type |
| --- | --- |
| `record` | `ComposeRecord` |
| `oldRecord` | `ComposeRecord` |
| `module` | `ComposeModule` |
| `namespace` | `ComposeNamespace` |
| `recordValueErrors` | `ComposeRecordValueErrorSet` |
| `selected` | Any |

### After Record Undelete {#after-record-undelete}

<div class="env-meta">

`compose:record` · `afterUndelete`

</div>

Triggered after record is undeleted.

**Inputs**

| Input | Type | Required |
| --- | --- | --- |
| Namespace | `ID`, `Handle`, `ComposeNamespace` | Yes |
| Module | `ID`, `Handle`, `ComposeModule` | Yes |

**Provides**

| Value | Type |
| --- | --- |
| `record` | `ComposeRecord` |
| `oldRecord` | `ComposeRecord` |
| `module` | `ComposeModule` |
| `namespace` | `ComposeNamespace` |
| `recordValueErrors` | `ComposeRecordValueErrorSet` |
| `selected` | Any |

### After Record Update {#after-record-update}

<div class="env-meta">

`compose:record` · `afterUpdate`

</div>

Triggered after record is updated.

**Inputs**

| Input | Type | Required |
| --- | --- | --- |
| Namespace | `ID`, `Handle`, `ComposeNamespace` | Yes |
| Module | `ID`, `Handle`, `ComposeModule` | Yes |

**Provides**

| Value | Type |
| --- | --- |
| `record` | `ComposeRecord` |
| `oldRecord` | `ComposeRecord` |
| `module` | `ComposeModule` |
| `namespace` | `ComposeNamespace` |
| `recordValueErrors` | `ComposeRecordValueErrorSet` |
| `selected` | Any |

### Before Record Create {#before-record-create}

<div class="env-meta">

`compose:record` · `beforeCreate`

</div>

Triggered before record is created.

**Inputs**

| Input | Type | Required |
| --- | --- | --- |
| Namespace | `ID`, `Handle`, `ComposeNamespace` | Yes |
| Module | `ID`, `Handle`, `ComposeModule` | Yes |

**Provides**

| Value | Type |
| --- | --- |
| `record` | `ComposeRecord` |
| `oldRecord` | `ComposeRecord` |
| `module` | `ComposeModule` |
| `namespace` | `ComposeNamespace` |
| `recordValueErrors` | `ComposeRecordValueErrorSet` |
| `selected` | Any |

### Before Record Delete {#before-record-delete}

<div class="env-meta">

`compose:record` · `beforeDelete`

</div>

Triggered before record is deleted.

**Inputs**

| Input | Type | Required |
| --- | --- | --- |
| Namespace | `ID`, `Handle`, `ComposeNamespace` | Yes |
| Module | `ID`, `Handle`, `ComposeModule` | Yes |

**Provides**

| Value | Type |
| --- | --- |
| `record` | `ComposeRecord` |
| `oldRecord` | `ComposeRecord` |
| `module` | `ComposeModule` |
| `namespace` | `ComposeNamespace` |
| `recordValueErrors` | `ComposeRecordValueErrorSet` |
| `selected` | Any |

### Before Record Undelete {#before-record-undelete}

<div class="env-meta">

`compose:record` · `beforeUndelete`

</div>

Triggered before record is undeleted.

**Inputs**

| Input | Type | Required |
| --- | --- | --- |
| Namespace | `ID`, `Handle`, `ComposeNamespace` | Yes |
| Module | `ID`, `Handle`, `ComposeModule` | Yes |

**Provides**

| Value | Type |
| --- | --- |
| `record` | `ComposeRecord` |
| `oldRecord` | `ComposeRecord` |
| `module` | `ComposeModule` |
| `namespace` | `ComposeNamespace` |
| `recordValueErrors` | `ComposeRecordValueErrorSet` |
| `selected` | Any |

### Before Record Update {#before-record-update}

<div class="env-meta">

`compose:record` · `beforeUpdate`

</div>

Triggered before record is updated.

**Inputs**

| Input | Type | Required |
| --- | --- | --- |
| Namespace | `ID`, `Handle`, `ComposeNamespace` | Yes |
| Module | `ID`, `Handle`, `ComposeModule` | Yes |

**Provides**

| Value | Type |
| --- | --- |
| `record` | `ComposeRecord` |
| `oldRecord` | `ComposeRecord` |
| `module` | `ComposeModule` |
| `namespace` | `ComposeNamespace` |
| `recordValueErrors` | `ComposeRecordValueErrorSet` |
| `selected` | Any |

## System

### Interval {#interval}

<div class="env-meta">

`system` · `onInterval`

</div>

Triggered on interval.

**Inputs**

| Input | Type | Description |
| --- | --- | --- |
| Interval | `String` | Select or type cron expression. Options: Every 30 minutes, Every hour, Every 2 hours, Every 6 hours, Every 12 hours, Every 24 hours, Every 2 days, Weekly, Monthly, Yearly. |

### Manual {#manual}

<div class="env-meta">

`compose:record` · `onManual`

</div>

Triggered manually.

**Inputs**

| Input | Type |
| --- | --- |
| Namespace | `ID`, `Handle`, `ComposeNamespace` |
| Module | `ID`, `Handle`, `ComposeModule` |

**Provides**

| Value | Type | Description |
| --- | --- | --- |
| `record` | `ComposeRecord` | |
| `module` | `ComposeModule` | |
| `namespace` | `ComposeNamespace` | |
| `page` | `ComposePage` | The compose page the automation was started from |
| `selected` | `Array` | The records ticked in the record list the button was pressed in |
| `filter` | `String` | The query the record list was showing when the button was pressed |

### Timestamp {#timestamp}

<div class="env-meta">

`system` · `onTimestamp`

</div>

Triggered at timestamp.

**Inputs**

| Input | Type |
| --- | --- |
| Timestamp | `String` |

## Users

### After User Create {#after-user-create}

<div class="env-meta">

`system:user` · `afterCreate`

</div>

Triggered after user is created.

**Provides**

| Value | Type |
| --- | --- |
| `user` | `User` |
| `oldUser` | `User` |

### After User Delete {#after-user-delete}

<div class="env-meta">

`system:user` · `afterDelete`

</div>

Triggered after user is deleted.

**Provides**

| Value | Type |
| --- | --- |
| `user` | `User` |
| `oldUser` | `User` |

### After User Suspend {#after-user-suspend}

<div class="env-meta">

`system:user` · `afterSuspend`

</div>

Triggered after user is suspended.

**Provides**

| Value | Type |
| --- | --- |
| `user` | `User` |
| `oldUser` | `User` |

### After User Update {#after-user-update}

<div class="env-meta">

`system:user` · `afterUpdate`

</div>

Triggered after user is updated.

**Provides**

| Value | Type |
| --- | --- |
| `user` | `User` |
| `oldUser` | `User` |

### Before User Create {#before-user-create}

<div class="env-meta">

`system:user` · `beforeCreate`

</div>

Triggered before user is created.

**Provides**

| Value | Type |
| --- | --- |
| `user` | `User` |
| `oldUser` | `User` |

### Before User Delete {#before-user-delete}

<div class="env-meta">

`system:user` · `beforeDelete`

</div>

Triggered before user is deleted.

**Provides**

| Value | Type |
| --- | --- |
| `user` | `User` |
| `oldUser` | `User` |

### Before User Suspend {#before-user-suspend}

<div class="env-meta">

`system:user` · `beforeSuspend`

</div>

Triggered before user is suspended.

**Provides**

| Value | Type |
| --- | --- |
| `user` | `User` |
| `oldUser` | `User` |

### Before User Update {#before-user-update}

<div class="env-meta">

`system:user` · `beforeUpdate`

</div>

Triggered before user is updated.

**Provides**

| Value | Type |
| --- | --- |
| `user` | `User` |
| `oldUser` | `User` |
