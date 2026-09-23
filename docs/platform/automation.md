---
title: Automation (TAQ)
description: What a TAQ is, the triggers that start it, the steps, branches and loops it is built from, and how it runs.
---

# Automation (TAQ)

Automations in Human are **TAQs**, short for **Trigger Action Query**. A TAQ is
a flow drawn on a canvas: one or more **triggers** start it, **steps** do the
work, and **branches** and **loops** decide which steps run and how often. You
build TAQs in the **Automation (TAQ)** application.

```
Trigger ─► Step ─► Branch ─┬─► Step ─► End
                           └─► Loop ─► Step
```

## Anatomy of a TAQ

A TAQ has a **TAQ Name**, a **Description**, optional **Labels** for grouping,
and a **Run as** user (see [Who a TAQ runs as](#who-a-taq-runs-as)). The rest is
its flow.

The builder lays the flow out automatically: you never position nodes by hand,
the shape follows from how the steps connect. Every connection carries a **+**
to insert a step at that point, and any node can be replaced by another of the
same sort. Selecting a node opens its configuration beside the canvas.

## Triggers

A trigger says what starts the TAQ. A TAQ can have several; any one of them
starts it.

| Trigger                                                                | Starts the TAQ when                                                                               |
| ---------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- |
| Record events: Before/After Record Create, Update, Delete and Undelete | A record in a chosen module changes. _Before_ runs ahead of the change, _After_ once it is saved. |
| User events: Before/After User Create, Update, Delete and Suspend      | A user account changes.                                                                           |
| **Manual**                                                             | A person presses a button on a page, or someone runs it from the builder.                         |
| **Interval**                                                           | A schedule comes round. Pick a preset from **Every 30 minutes** to **Yearly**.                    |
| **Timestamp**                                                          | A chosen date and time is reached.                                                                |
| **Agent Invoked**                                                      | An agent that has been given this TAQ decides to run it.                                          |
| **Before Chatbot Step**, **After Chatbot Step**                        | A chatbot is about to run, or has finished, a step of its journey.                                |

A record trigger is narrowed to a namespace and a module, and can be narrowed
to a single record. A **Manual** trigger started from a page receives the
record, module, namespace and page it was started from and, from a record list,
the selected records and the list's current filter.

An **Agent Invoked** trigger can declare an input schema: named parameters, each
a string, number or boolean, that the agent must supply. They become values
later steps can use.

If a configured [connection](./administration#connections-and-data-sources)
defines webhooks, its events appear as triggers too.

## Steps

Steps are the actions. The step picker groups them:

| Group             | Steps                                                                                                          |
| ----------------- | -------------------------------------------------------------------------------------------------------------- |
| **Records**       | Find Record, Find Oldest Record, Find Newest Record, Create Record, Update Record, Clone Record, Delete Record |
| **Users**         | Create User, Update User, Suspend User, Delete User                                                            |
| **Notifications** | Send Notification, Send Record Notification (see [notifications](./home#notifications))                        |
| **Agents**        | Prompt Agent, Continue Conversation                                                                            |
| **Loops**         | Each Record, Count, While                                                                                      |
| **Branches**      | First Match, All Matches                                                                                       |
| **Workflows**     | Run Workflow                                                                                                   |
| **Corredor**      | Run Corredor script, available only when [`CORREDOR_ENABLED`](/reference/environment#CORREDOR_ENABLED) is on   |

A configured [connection](./administration#connections-and-data-sources) adds
its own operations as steps, grouped under the connection's name.

**Prompt Agent** sends a prompt to an agent and hands its answer, and the
conversation it started, to the next step. **Continue Conversation** sends a
follow-up into that same conversation. The agent does not need system
invocation for this: it runs as the TAQ does. See [Agents](./agents).

## Branches

A branch splits the flow into paths, each with a condition. Its **Branching
Mode** is one of two:

- **First Match** takes the first path whose condition is true.
- **All Matches** takes every path whose condition is true.

Paths are ordered and can be reordered. A condition compares a value with
**equals**, **not equals**, **less than**, **less or equal**, **greater than**,
**greater or equal**, **is empty** or **is not empty**, and conditions can be
combined into groups where all or any must match. A path with no condition is
the default.

## Loops

A loop repeats the steps in its body:

- **Each Record** runs once per record matching a filter in a module, with
  optional sort and limit. Each pass gets the record, its index and the total.
- **Count** runs from a first to a last number in steps.
- **While** runs as long as its condition, an
  [expression](/reference/expressions/), is true.

## Passing values between steps

Any step argument can be typed in, or bound to a value produced earlier in the
flow. Binding opens the **Step Results** panel, which lists what is available
at that point:

- **Invoker**, the user who triggered the TAQ, and **Runner**, the user it runs
  as.
- What the trigger carries. A record trigger, for example, provides the
  **Record**, the **OldRecord** (its values before the change), the **Module**
  and the **Namespace**.
- The outputs of each earlier step, such as the record a **Create Record** step
  made or the answer from **Prompt Agent**.

Records expand into their fields, so a step can use a single value from a
record found three steps earlier.

## Who a TAQ runs as

Every step runs with someone's permissions. By default a TAQ runs as whoever
triggered it, which the **Run as** setting shows as **Default (invoker)**. Set
**Run as** to a specific user, typically a service account, to run it with that
user's permissions instead. Either way, a TAQ can never do more than the user
it runs as is allowed to.

## Enabling and running

A TAQ runs only while it is **Enabled**. A disabled TAQ does not fire on its
triggers and refuses to be started by hand. A TAQ with no trigger, or with
unresolved problems in its configuration, does not fire on events either.

From the builder you can **Run** a TAQ on the spot. If there are unsaved
changes you are asked to save first. If the trigger expects input, such as a
record, you supply it before the run starts.

## The execution trace

A run started from the builder is traced on the canvas: it shows whether the
run completed or failed and how long it took. Nodes that ran are highlighted, along with the paths
taken, and selecting a node shows its **Input**, **Scope** and **Output** for
that run, including any error.

Human has no screen that lists past executions. The trace is available for runs
started from the builder.

## Permissions

TAQs have their own permissions: who may read, update, delete, run and grant
permissions on each one, or on all of them. The
[automation permissions reference](/reference/permissions/automation) lists
them.

## Workflows

Human also contains a separate designer, **Workflows**. It is switched off on
new instances, and these docs cover TAQs.

## Where to go next

- [Agents](./agents): let an agent start TAQs, or prompt one from a TAQ.
- [Chatbots](./chatbots): run a TAQ before or after a chatbot step.
- [Namespaces and data](./namespaces): the records most TAQs act on.
