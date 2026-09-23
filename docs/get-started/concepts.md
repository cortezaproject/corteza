---
title: Core concepts
description: The terms the rest of the Human docs use, and how they fit together.
---

# Core concepts

A short tour of the terms the rest of the docs use. Each one is a thing you can
open in the web app.

## Apps and data

**Namespace.** A container for one app. It holds the app's modules, pages and
charts, and its short name is the app's address: `/compose/namespace/<short name>`.

**Module.** One kind of thing your app keeps track of, such as a contact, a case
or an invoice. A module is defined by its fields.

**Field.** One piece of information on a module. The field kinds are Text,
Number, Checkbox (Y/N), Date and Time, Select / Dropdown, Email, URL, File,
User, Record and Location. A **Record** field links to a record in another module. That is how you
connect contacts to companies, for example.

**Record.** One entry in a module: one contact, one case.

**Page.** A screen in a namespace, laid out from blocks in the page builder. A
_record page_ shows a single record; every other page is free-form.

**Block.** One piece of a page. The kinds include Record list, Record, Chart,
Metric, Calendar, Map, Content, Tabs, Comments, Report, Agent Chat and more.

**Chart.** A visualisation of module data: bar, line, pie, doughnut and
scatter charts, plus funnel, gauge and radar. You place a chart on a page with a
Chart block.

```
Namespace
 ├─ Module ── Fields
 │    └─ Records
 ├─ Page ── Blocks (Record list, Chart, …)
 └─ Chart
```

## Automation

**TAQ (Trigger Action Query).** An automation built from triggers, function
steps, and the gateways and iterators that shape its flow. TAQs are drawn on a
canvas in the **Automation (TAQ)** app.

Human also includes a separate **Workflows** designer. It is switched off on new
instances, and these docs cover TAQs.

**Trigger.** What starts an automation: a record being created or changed, a
schedule, or an explicit call.

## AI

**LLM provider.** A connection to a model vendor such as Anthropic or Mistral,
or any service with an OpenAI-compatible API. Administrators add them once, and
agents pick one.

**Agent.** A model with instructions (the _system prompt_), guardrails, the
tools it may call and the topics it knows about. People chat with agents, and
automations and chatbots can run them.

**Chatbot.** A chat widget you embed on a website. It follows a journey of
steps, and a step can hand the conversation to an agent or to a person.

**Assistant.** The chat column on the home page. It offers the agents your
roles are allowed to see.

## People and access

**User**, **User group** and **Role.** People sign in as users. Permissions are
granted to roles, and users are members of roles.

**Permission.** One operation on one kind of resource, such as _read_ on a
module, allowed or denied for a role. Anything not allowed is denied. The
[permissions reference](/reference/permissions/) lists every operation.

Agents and automations act as a person: the one who invoked them, or a service
account you choose. They can never do more than that person can.

## Projects

**Project.** A governed bundle of related resources (a data model, pages,
automations, agents, chatbots and the people who work on them) that moves
through a lifecycle from _Draft_ to _Published_.

## Applications

**Application.** An entry in the menu on the home page. Human ships with
applications for its own tools, and administrators can add their own under
**Admin Area → System → Applications**, pointing at any URL, including one of your
namespaces.
