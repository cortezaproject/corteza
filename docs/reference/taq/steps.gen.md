---
title: TAQ steps
description: Every step a TAQ can run, with the inputs it takes and the values it returns.
outline: [2, 3]
---

<!-- This file is auto-generated from the TAQ construct catalog (server/tests/automation/docs_catalog_test.go). -->

# TAQ steps

Steps do the work of a [TAQ](/platform/automation); branches and loops decide
which steps run and how often. These are the steps the builder offers on every
server, grouped as in the step picker. **Inputs** are what you set on the step;
**Outputs** are the values later steps can use. What starts a TAQ is in
[TAQ triggers](./triggers).

Steps marked **Needs Corredor** are listed on every server but can only be
added where Corredor is turned on. Operations that come from a connection
depend on the connections configured on your instance and are not listed here.

## Agents

### Prompt Agent {#agentPrompt}

<div class="env-meta">

`agentPrompt`

</div>

Send a prompt to an agent.

**Inputs**

| Input | Type | Required |
| --- | --- | --- |
| Agent | `ID` | Yes |
| Input | `String` | Yes |

**Outputs**

| Output | Type |
| --- | --- |
| `output` | `String` |
| `conversationID` | `ID` |

### Continue Conversation {#agentContinue}

<div class="env-meta">

`agentContinue`

</div>

Continue an existing agent conversation.

**Inputs**

| Input | Type | Required |
| --- | --- | --- |
| Conversation ID | `ID` | Yes |
| Input | `String` | Yes |

**Outputs**

| Output | Type |
| --- | --- |
| `output` | `String` |
| `conversationID` | `ID` |

## Branches

### First Match {#first-match}

Takes the first path whose condition is true.

Each path out of the branch carries its own condition, an [expression](/reference/expressions/) that decides whether the path is taken.

### All Matches {#all-matches}

Takes every path whose condition is true.

Each path out of the branch carries its own condition, an [expression](/reference/expressions/) that decides whether the path is taken.

## Corredor

### Run Corredor script {#corredorExec}

<div class="env-meta">

`corredorExec` · Needs Corredor (`CORREDOR_ENABLED`)

</div>

Run a manual Corredor server script and return what it returns.

**Inputs**

| Input | Type | Required | Description |
| --- | --- | --- | --- |
| Script | `String` | Yes | Name of the server script to run, e.g. /server-scripts/hello.js:default. |
| Arguments | `Vars` | | Values the script receives as its arguments. |

**Outputs**

| Output | Type |
| --- | --- |
| `results` | `Vars` |

## Loops

### Count {#loopSequence}

<div class="env-meta">

`loopSequence` · Loop

</div>

Iterate from first to last by step.

**Inputs**

| Input | Type | Required |
| --- | --- | --- |
| First | `Integer` | |
| Last | `Integer` | Yes |
| Step | `Integer` | |

**Outputs**

| Output | Type |
| --- | --- |
| `counter` | `Integer` |
| `isFirst` | `Boolean` |
| `isLast` | `Boolean` |

### Each Record {#composeRecordsEach}

<div class="env-meta">

`composeRecordsEach` · Loop

</div>

Iterate over records matching a query.

**Inputs**

| Input | Type | Required | Description |
| --- | --- | --- | --- |
| Namespace | `ID`, `Handle`, `ComposeNamespace` | Yes | |
| Module | `ID`, `Handle`, `ComposeModule` | Yes | Records are fetched from the specified module. |
| Query | `String` | | Optional filter string to match records. |
| Sort | `String` | | Sort expression (e.g. "createdAt DESC"). |
| Limit | `Integer` | | Maximum number of records to iterate. |

**Outputs**

| Output | Type |
| --- | --- |
| `record` | `ComposeRecord` |
| `index` | `Integer` |
| `total` | `Integer` |

### While {#loopDo}

<div class="env-meta">

`loopDo` · Loop

</div>

Iterate while condition is true.

**Inputs**

| Input | Type | Required | Description |
| --- | --- | --- | --- |
| While | `String` | Yes | Boolean expression evaluated each iteration |

## Notifications

### Send Notification {#notificationSend}

<div class="env-meta">

`notificationSend`

</div>

Sends a simple notification with title and description to a user.

**Inputs**

| Input | Type | Required | Description |
| --- | --- | --- | --- |
| Recipient | `ID`, `Handle`, `String` | Yes | User to send the notification to. Can be an ID, handle, or email address. |
| Title | `String` | Yes | |
| Description | `String` | | |

### Send Record Notification {#notificationSendRecord}

<div class="env-meta">

`notificationSendRecord`

</div>

Sends a notification that links to a specific record.

**Inputs**

| Input | Type | Required | Description |
| --- | --- | --- | --- |
| Recipient | `ID`, `Handle`, `String` | Yes | User to send the notification to. Can be an ID, handle, or email address. |
| Title | `String` | Yes | |
| Description | `String` | | |
| Namespace | `ID`, `Handle`, `ComposeNamespace` | Yes | |
| Module | `ID`, `Handle`, `ComposeModule` | Yes | |
| Record | `ID`, `ComposeRecord` | | |
| Open Mode | `String` | | Options: Same Tab, New Tab, Modal. |
| Edit | `Boolean` | | |

## Records

### Clone Record {#composeRecordsClone}

<div class="env-meta">

`composeRecordsClone`

</div>

Creates a copy of an existing record.

**Inputs**

| Input | Type | Required | Description |
| --- | --- | --- | --- |
| Namespace | `ID`, `Handle`, `ComposeNamespace` | Yes | |
| Module | `ID`, `Handle`, `ComposeModule` | Yes | Even with unique record ID across all modules, module needs to be known before doing any record operations. Mainly because records of different modules can be located in different stores. |
| Record | `ID`, `ComposeRecord` | Yes | |
| Override Values | `KV`, `KVV`, `Any` | | |

**Outputs**

| Output | Type |
| --- | --- |
| `record` | `ComposeRecord` |

### Create Record {#composeRecordsCreate}

<div class="env-meta">

`composeRecordsCreate`

</div>

Add new record to module.

**Inputs**

| Input | Type | Required | Description |
| --- | --- | --- | --- |
| Namespace | `ID`, `Handle`, `ComposeNamespace` | Yes | |
| Module | `ID`, `Handle`, `ComposeModule` | Yes | Even with unique record ID across all modules, module needs to be known before doing any record operations. Mainly because records of different modules can be located in different stores. |
| Values | `KV`, `KVV`, `Any` | Yes | |

**Outputs**

| Output | Type |
| --- | --- |
| `record` | `ComposeRecord` |

### Delete Record {#composeRecordsDelete}

<div class="env-meta">

`composeRecordsDelete`

</div>

Delete a record by ID.

**Inputs**

| Input | Type | Required | Description |
| --- | --- | --- | --- |
| Namespace | `ID`, `Handle`, `ComposeNamespace` | | |
| Module | `ID`, `Handle`, `ComposeModule` | | Even with unique record ID across all modules, module needs to be known before doing any record operations. Mainly because records of different modules can be located in different stores. |
| Record | `ID`, `ComposeRecord` | Yes | |

### Find Newest Record {#composeRecordsLast}

<div class="env-meta">

`composeRecordsLast`

</div>

Find the newest record in a module.

**Inputs**

| Input | Type | Required | Description |
| --- | --- | --- | --- |
| Namespace | `ID`, `Handle`, `ComposeNamespace` | Yes | |
| Module | `ID`, `Handle`, `ComposeModule` | Yes | Even with unique record ID across all modules, module needs to be known before doing any record operations. Mainly because records of different modules can be located in different stores. |

**Outputs**

| Output | Type |
| --- | --- |
| `record` | `ComposeRecord` |

### Find Oldest Record {#composeRecordsFirst}

<div class="env-meta">

`composeRecordsFirst`

</div>

Find the oldest record in a module.

**Inputs**

| Input | Type | Required | Description |
| --- | --- | --- | --- |
| Namespace | `ID`, `Handle`, `ComposeNamespace` | Yes | |
| Module | `ID`, `Handle`, `ComposeModule` | Yes | Even with unique record ID across all modules, module needs to be known before doing any record operations. Mainly because records of different modules can be located in different stores. |

**Outputs**

| Output | Type |
| --- | --- |
| `record` | `ComposeRecord` |

### Find Record {#composeRecordsLookup}

<div class="env-meta">

`composeRecordsLookup`

</div>

Lookup record by ID.

**Inputs**

| Input | Type | Required | Description |
| --- | --- | --- | --- |
| Namespace | `ID`, `Handle`, `ComposeNamespace` | Yes | |
| Module | `ID`, `Handle`, `ComposeModule` | Yes | Even with unique record ID across all modules, module needs to be known before doing any record operations. Mainly because records of different modules can be located in different stores. |
| Record | `ID`, `ComposeRecord` | Yes | |

**Outputs**

| Output | Type |
| --- | --- |
| `record` | `ComposeRecord` |

### Update Record {#composeRecordsUpdate}

<div class="env-meta">

`composeRecordsUpdate`

</div>

Update an existing record.

**Inputs**

| Input | Type | Required | Description |
| --- | --- | --- | --- |
| Namespace | `ID`, `Handle`, `ComposeNamespace` | Yes | |
| Module | `ID`, `Handle`, `ComposeModule` | Yes | Even with unique record ID across all modules, module needs to be known before doing any record operations. Mainly because records of different modules can be located in different stores. |
| Record | `ID`, `ComposeRecord` | Yes | |
| Values | `KV`, `KVV`, `Any` | Yes | |

**Outputs**

| Output | Type |
| --- | --- |
| `record` | `ComposeRecord` |

## Users

### Create User {#usersCreate}

<div class="env-meta">

`usersCreate`

</div>

Create a new user account.

**Inputs**

| Input | Type | Required |
| --- | --- | --- |
| Email | `String` | Yes |
| Name | `String` | |
| Handle | `String` | |
| Username | `String` | |

**Outputs**

| Output | Type |
| --- | --- |
| `user` | `User` |

### Delete User {#usersDelete}

<div class="env-meta">

`usersDelete`

</div>

Delete a user account.

**Inputs**

| Input | Type | Required |
| --- | --- | --- |
| User | `ID`, `Handle`, `String`, `User` | Yes |

### Suspend User {#usersSuspend}

<div class="env-meta">

`usersSuspend`

</div>

Suspend a user account.

**Inputs**

| Input | Type | Required |
| --- | --- | --- |
| User | `ID`, `Handle`, `String`, `User` | Yes |

### Update User {#usersUpdate}

<div class="env-meta">

`usersUpdate`

</div>

Update an existing user account.

**Inputs**

| Input | Type | Required |
| --- | --- | --- |
| User | `ID`, `Handle`, `String`, `User` | Yes |
| Email | `String` | |
| Name | `String` | |
| Handle | `String` | |
| Username | `String` | |

**Outputs**

| Output | Type |
| --- | --- |
| `user` | `User` |

## Workflows

### Run Workflow {#workflowExec}

<div class="env-meta">

`workflowExec`

</div>

Run a workflow and return its results.

**Inputs**

| Input | Type | Required | Description |
| --- | --- | --- | --- |
| Workflow | `ID`, `Handle` | Yes | Workflow to run, referenced by ID or handle. |
| Input | `Vars`, `Any` | | Scope passed to the workflow as its input. |

**Outputs**

| Output | Type |
| --- | --- |
| `results` | `Vars` |
