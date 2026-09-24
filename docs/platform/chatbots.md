---
title: Chatbots
description: What a chatbot is, the journey of steps it takes visitors through, handoff to human operators, styling, and the widget key and allowed origins.
---

# Chatbots

A **chatbot** is a chat widget for people outside Human, such as visitors to
your website. Where an [agent](./agents) is open-ended, a chatbot follows a
**journey**: a fixed sequence of steps that can greet the visitor, collect
details, ask for consent and then hand them to an agent for a live
conversation. At any point in a conversation the visitor can ask for a person,
and an operator in Human can take over.

You build chatbots in the **Chatbot** application.

## The chatbot editor

The editor has a **Config** tab and, once the chatbot is saved, a **Sessions**
tab. Beside it runs a live preview: the real widget, working against the real
chatbot, so every change can be tried as you make it.

<Screenshot name="chatbot-editor" alt="The chatbot editor, showing a journey step and the live preview" />

**Config** has three panels.

### General

The chatbot's **Name** and **Handle**, **Enable chatbot** to switch the widget
on or off, the **Widget key** and **Allowed origins** (see
[below](#widget-key-and-allowed-origins)), and the **Session TTL**: how long a
visitor's session lasts, such as `2h` or `24h`.

### Journey

A journey is an ordered list of steps the visitor moves through. Each step has
a name, a handle and a type:

| Step type          | What the visitor gets                                                                                      |
| ------------------ | ---------------------------------------------------------------------------------------------------------- |
| **Static message** | A message from the chatbot, optionally moving on by itself after a delay.                                  |
| **Form**           | A form to fill in. Fields are **Text**, **Multi-line text**, **Email** or **Number**, and can be required. |
| **Consent**        | A text to accept or reject, which can link to external documents. Rejecting ends the session.              |
| **Conversation**   | A live conversation with an **Agent**, optionally opened with an initial prompt.                           |

The agent in a **Conversation** step must have **Enable System Invocation** on
and a **Service Account** set, because no Human user is present: the agent acts
as that service account. See [Agents](./agents#invocation).

Every step can run a TAQ around it. A **Before automation** runs before the
step and stops it if the TAQ fails; an **After automation** runs once the step
completes, and its errors are only logged. These TAQs start from the **Before
Chatbot Step** and **After Chatbot Step** triggers, which receive the
chatbot, the session, and what the visitor entered in earlier steps. See
[Automation](./automation).

### Styling

How the widget looks: the chat panel's **Title** and **Logo**, the launcher
button's **Button label**, **Icon**, **Button position**, **Button shape** and
**Button size**, whether it should **Start open**, its colors, and the text and
heading sizes. Images can be uploaded once the chatbot has been saved.

## Handoff to a person

With **Enable handoff** on, the widget shows a "Talk to human" button during
conversation steps. Pressing it asks for an operator. Two TAQs can be attached
to the handoff: **When handoff requested** and **When operator accepts**, for
example to notify the team.

## Sessions and operators

Every visitor's conversation is a **session**. Operators work in the
**Sessions** inbox, which is available in three places: the **Sessions** tab of
one chatbot, the Chatbot application's own **Sessions** view across every
chatbot you can see, and a **Chatbot** block on any page, so an inbox can sit
inside one of your apps.

Sessions are sorted by state:

| State                 | Means                                                          |
| --------------------- | -------------------------------------------------------------- |
| **Awaiting operator** | The visitor asked for a person and nobody has accepted yet.    |
| **In progress**       | An operator has accepted and is replying.                      |
| **Active**            | The visitor is moving through the journey without an operator. |
| **Closed**            | The session has ended.                                         |

An operator **Accept**s a waiting session to take it over; from then on the
visitor sees the operator's replies in real time, under the operator's
**Display name**. **Resolve** ends the handoff and lets the agent carry on at
the same step. People allowed to manage a session can also **Force advance
step** or **Force close session**. Sessions started from the editor's preview are marked as previews
and can be filtered out.

## Widget key and allowed origins

Each chatbot has a **Widget key**, generated when it is first saved, which
identifies the chatbot to its widget. **Regenerate widget key** replaces it and
ends every live session of that chatbot.

**Allowed origins** lists the websites allowed to load the widget, each as an
origin such as `https://www.example.com`. The origin a request comes from must
match an entry exactly. An entry of `*` allows any origin; avoid it in
production. A disabled chatbot refuses every request, whatever its origins.

## Permissions

Chatbots and their sessions have their own permissions, listed in the
[system permissions reference](/reference/permissions/system).
