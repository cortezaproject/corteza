---
title: How building works
description: How data, pages, automation, agents and permissions fit together when you build on Human.
---

<script setup>
import LayerStack from '../.vitepress/theme/components/diagrams/LayerStack.vue'
</script>

# How building works

Everything you build on Human is made from the same few layers. You rarely use
all of them at once, and you can add each one later. This page explains what
each layer is for and how they connect.

<LayerStack />

## Data: start with what you keep track of

An app begins as a **namespace** with one or more **modules**. A module is one
kind of thing, such as a contact, a case or an invoice, and its **fields**
describe it. Every entry you save is a **record**.

Fields carry most of the app's rules. A field can be required, hold several
values, link to a record in another module (a **Record** field) or to a person
(a **User** field). A field can also calculate its own value, clean up what
people type, and reject values that break a rule, all with
[expressions](/reference/expressions/).

## Interface: pages people work in

**Pages** are what people open. A page is laid out on a grid from **blocks**: a
record list, a single record, a chart, a metric, a calendar, a map, and more.
Every module can get a list page and a record page generated for it, which you
then rearrange in the page builder. **Charts** summarise module data and appear
on pages through a Chart block.

To put an app in everyone's menu on the home screen, an administrator adds it
as an **application**.

## Automation: what happens on its own

A **TAQ** (Trigger Action Query) runs when something happens: a record is
created, updated or deleted, a schedule comes round, an agent or chatbot calls
it, or someone runs it by hand. It then works through **steps**: finding and
changing records, sending notifications, branching on conditions, looping over
records, and prompting an agent. Each step can use the results of the steps
before it.

A TAQ runs as the person who triggered it, unless you set **Run as** to a
specific user.

## AI: agents that work with your data

An **agent** is a model with instructions and the tools it may use. Tools are
what connect it to the layers below: it can look up and change records, read
the data model, build pages, or start TAQs, as far as you allow. You can limit
an agent to a single namespace.

People talk to agents in the Assistant on the home screen or in an Agent Chat
block on a page. TAQs prompt them as a step, and chatbots put them in front of
visitors to your website.

## People and access: one model for everything

Every layer is covered by the same permission model. **Roles** are granted
operations on resources, such as reading a module or updating its records, and
anything not allowed is denied. People get permissions through the roles they
belong to.

Automations and agents do not get permissions of their own. They act as a
person, either the one who invoked them or the user you choose, so they can
never do more than that person can.

## A typical order

1. **Model the data**: a namespace, its modules and their fields.
2. **Generate the pages**, then shape them in the page builder and add charts.
3. **Set the permissions**: who can read, create and change what.
4. **Automate** the repetitive parts with TAQs.
5. **Add agents** where people would rather ask than click.

## Next

- [The platform](/platform/): each part of Human in more detail.
- [Core concepts](./concepts): the vocabulary in one place.
- [Expressions](/reference/expressions/): the formula language for fields and
  roles.
