---
title: The platform
description: A map of Human's parts — home, namespaces, automation, agents, chatbots, projects, administration and federation — and how they relate.
---

# The platform

Human is one web app with several parts. Each part is an application in the
menu on the home screen, and each is gated by the same role-based permissions.
This page is a map. Every section links to a page that explains that part in
more depth.

If a term is new to you, [Core concepts](/get-started/concepts) defines it in a
sentence.

## Home and apps

The screen at `/` is where everyone starts: the menu of applications on the
left, the Assistant in the middle and notifications on the right. What appears
in the menu is decided by the application registry and by each person's roles.
[Home and apps →](./home)

## Namespaces and data

A namespace is one app built in Human: its modules (the kinds of records it
keeps), the records themselves, and the pages, blocks and charts people work
in. Builders design namespaces in their admin panel; end users only see the
pages. [Namespaces and data →](./namespaces)

## Automation (TAQ)

A TAQ (Trigger Action Query) is an automation drawn on a canvas. Triggers start
it, steps do the work, and branches and loops shape the flow between them. TAQs
react to record and user events, run on a schedule, or are started by a person,
an agent or a chatbot. [Automation (TAQ) →](./automation)

## Agents

An agent is a model from an LLM provider, given instructions, knowledge and a
set of tools. People talk to agents in the Assistant and on pages, and TAQs and
chatbots can run them. An agent never does more than the person it acts for.
[Agents →](./agents)

## Chatbots

A chatbot is a chat widget for visitors outside Human. It walks each visitor
through a journey of steps, one of which can be a conversation with an agent,
and can hand the conversation to a human operator. [Chatbots →](./chatbots)

## Projects

A project bundles related resources into one governed unit.
[Projects →](./projects)

## Administration

The Admin Area holds everything that is configured once for the whole
instance: users, roles and permissions, the application registry, LLM
providers, connections, email and the audit log.
[Administration →](./administration)

## Federation

Federation shares records between two Human servers: one exposes a module, the
other maps it onto a module of its own and receives its records as they are
created, changed and deleted. [Federation →](./federation)

## How the parts relate

- **Namespaces hold the data.** TAQs, agents and chatbots act on records in
  namespaces; they do not keep business data of their own.
- **TAQs connect things.** A record change can start a TAQ, a TAQ can prompt an
  agent, an agent can start a TAQ, and a chatbot step can run a TAQ before or
  after it.
- **Agents are reused everywhere.** One agent can answer in the Assistant, sit
  on a page in an Agent Chat block, run as a TAQ step and hold the conversation
  in a chatbot.
- **One permission model covers all of it.** Every read and write, whoever or
  whatever makes it, is checked against the roles of the person it runs as. The
  [permissions reference](/reference/permissions/) lists every operation.
