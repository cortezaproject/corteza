---
title: Glossary
description: Every term Human uses, what it means, where it appears in the web app, and the name the API and MCP tools use for it.
---

# Glossary

One name per thing. The first column is the word the docs and the web app use;
the last is the prefix of the MCP tools and REST routes that act on it, so a
term read on screen can be found in the API without guessing.

## Apps and data

| Term            | Meaning                                                                                                                                                                                                           | In the web app                          | API / MCP prefix                                  |
| --------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------- | ------------------------------------------------- |
| **Namespace**   | The container one app is built in: its modules, records, pages and charts. Its slug is the app's address. The docs sometimes say "app" for what a namespace holds; the resource is always the namespace.          | Namespaces                              | `compose_namespace`                               |
| **Module**      | One kind of thing the app keeps, defined by its fields: a contact, a case, an invoice.                                                                                                                            | Namespace → Admin → Modules             | `compose_module`                                  |
| **Field**       | One piece of information on a module, of a kind: String, Number, Bool, DateTime, Select, Email, Url, File, User, Record, Geometry. The web app shows friendlier labels (Text, Checkbox, Date and Time, Location). | Module editor                           | part of the module                                |
| **Record**      | One entry in a module.                                                                                                                                                                                            | Record list and record pages            | `compose_record`                                  |
| **Page**        | A screen in a namespace, made of blocks. A **record page** shows one record of a module; every other page is free-form.                                                                                           | Namespace → Admin → Pages, page builder | `compose_page`                                    |
| **Page layout** | One arrangement of a page's blocks. A page has at least one; further layouts show when their visibility expression passes.                                                                                        | Page builder → Layouts                  | `compose_page_layout`                             |
| **Block**       | One piece of a page: Record, RecordList, Chart, Metric, Calendar, Content, Tabs, Comment, AgentChat and the rest. A block's options are what it shows; its position belongs to the layout.                        | Page builder                            | `compose_page_block_schema` describes the options |
| **Chart**       | A visualisation over module data, placed on a page with a Chart block.                                                                                                                                            | Namespace → Admin → Charts              | `compose_chart`                                   |
| **Application** | An entry in the home-screen menu, pointing at a URL. It is not a namespace: a namespace is where an app is built, an application is how it is listed.                                                             | Admin Area → System → Applications      | `system_application`                              |

## Automation

| Term         | Meaning                                                                                                                                                                                             | In the web app           | API / MCP prefix                                                                     |
| ------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------ | ------------------------------------------------------------------------------------ |
| **TAQ**      | Trigger Action Query: an automation built from triggers, function steps, gateways and iterators. The word "automation" alone means a TAQ in these docs. Internally the resource is `ng-automation`. | Automation (TAQ)         | `automation_taq`, `automation_taq_construct` for the step and trigger catalogue      |
| **Trigger**  | What starts a TAQ: a record event, a schedule, or an explicit call. Triggers are part of the TAQ definition and can also be managed on their own.                                                   | TAQ canvas               | `automation_trigger`, `automation_event_type` for the events a trigger can listen to |
| **Workflow** | The separate, older automation designer, switched off on new instances. Not a TAQ.                                                                                                                  | Workflows (when enabled) | `automation_workflow`                                                                |

## AI

| Term             | Meaning                                                                                                                                           | In the web app             | API / MCP prefix      |
| ---------------- | ------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------- | --------------------- |
| **LLM provider** | A connection to a model vendor, added once by an administrator and picked by agents.                                                              | Admin Area → LLM providers | `system_llm_provider` |
| **Agent**        | A model with instructions, guardrails, the tools it may call and the topics it knows. People chat with agents; automations and chatbots run them. | Agents, the Assistant      | `system_agent`        |
| **Skill**        | Working rules the MCP server hands an assistant before it uses a tool family.                                                                     | not shown                  | `system_skill`        |
| **Chatbot**      | A chat widget embedded on a website, following a journey of steps that can hand over to an agent or a person. Not an agent.                       | Chatbots                   | `system_chatbot`      |
| **Assistant**    | The chat column on the home page, offering the agents your roles may see.                                                                         | Home                       | uses `system_agent`   |
| **Discovery**    | Search across records and documents the caller may read.                                                                                          | Search                     | `discovery_search`    |

## People and access

| Term            | Meaning                                                                                              | In the web app            | API / MCP prefix     |
| --------------- | ---------------------------------------------------------------------------------------------------- | ------------------------- | -------------------- |
| **User**        | A person who signs in. System accounts (`kind: sys`) are built-in users.                             | Admin Area → Users        | `system_user`        |
| **User group**  | A node in the organisation tree that can carry roles. Membership is inherited down the tree.         | Admin Area → User groups  | `system_user_group`  |
| **Role**        | What permissions are granted to. Users and user groups are members of roles.                         | Admin Area → Roles        | `system_role`        |
| **Permission**  | One operation on one kind of resource, allowed or denied for a role. Anything not allowed is denied. | Permissions dialogs       | `system_permission`  |
| **Auth client** | An OAuth2 client an external application or script uses to obtain tokens.                            | Admin Area → Auth clients | `system_auth_client` |

## Everything else

| Term                | Meaning                                                                                                                               | In the web app                | API / MCP prefix   |
| ------------------- | ------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------- | ------------------ |
| **Project**         | A governed bundle of related resources that moves from Draft to Published. A **revision** is a new draft copy of a published project. | Projects                      | project tools      |
| **Reminder**        | A dated note to one user, shown in the web app and dismissable or snoozable.                                                          | Reminders                     | `system_reminder`  |
| **Theme**           | The light and dark colour and logo settings of the instance.                                                                          | Admin Area → Theming          | `system_theme`     |
| **Record revision** | The history of one record's values, when the module keeps it. Unrelated to a project revision.                                        | Record page → Revisions block | part of the record |

## Names that mean something else

- **Automation** on its own is a TAQ. "Automation (Workflows)" named the old designer in an earlier version; the docs now say **Workflows**.
- **App** is informal for what a namespace holds. **Application** is the menu entry.
- **Corteza** is the project Human grew from; a `corteza::` prefix in a resource name is that history, not a separate product.
