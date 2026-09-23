---
title: Agents
description: LLM providers, what an agent is made of, how tool access works, and everywhere an agent can run.
---

# Agents

An **agent** is a language model given a job: instructions, knowledge about
your domain and data, and a set of tools it may call to look things up or make
changes in Human. People chat with agents, and TAQs and chatbots can run them.
You build agents in the **Agentic** application.

Whoever runs an agent, **it only ever does what that person is already allowed
to do**. Giving an agent a tool never widens anyone's access.

## LLM providers

Agents need a model to talk to. An **LLM provider** is a connection to a model
vendor, added once by an administrator under **Admin Area → System → LLM
Providers**. A provider has a name, a **Provider** type, a **Prompt URL** (the
API endpoint) and an **API Key**.

| Provider               | Connects to                               |
| ---------------------- | ----------------------------------------- |
| **Anthropic**          | Anthropic's API                           |
| **Mistral**            | Mistral's API                             |
| **Other (Custom URL)** | Any service with an OpenAI-compatible API |

The API key is write-only: once saved it is never shown again, only replaced
with **Update API Key**. Agents pick a provider and one of the models it offers.

## The agent editor

The editor has three tabs, **Configuration**, **Inspect** and **History**, with
a live chat beside them so you can test the agent as you change it.
**Configuration** is split into panels.

### General

The agent's **Name**, **Handle**, **Description** (shown to the people who use
it), **Labels** and **Status**. Only active agents can be invoked.

### Execution

The **LLM Provider** and **Model**, and the limits of a run: **Temperature**,
**Max Iterations** (how many think-and-act cycles it may take), **Token Limit**
for the whole conversation, **Agent Response Limit** for one reply, and a
**Timeout**.

### Behavior

- **System Prompt**: the instructions sent to the model at the start of every
  conversation. Who the agent is for, what it should do and what it must not.
- **Guardrails**: short rules the agent must follow, such as "Stay on topic".

### Knowledge Base

**Topics** give an agent domain knowledge. A topic has a **Title**, a free-text
**Description** (policies, context, instructions) and a **Compose Context**:
namespaces and modules whose structure the agent should understand. Topics are
shared, so several agents can use the same one.

### Access

What the agent can reach and do.

- **Works in** confines every tool to one namespace. Left empty, a tool reaches
  whatever the person invoking the agent can reach.
- **Tools** are the operations the agent may call. See
  [Tools and their modes](#tools-and-their-modes).
- **TAQs** lists the TAQs the agent may start, each with an optional
  description that tells the model when to use it. A listed TAQ is started
  through its **Agent Invoked** trigger, and the parameters that trigger
  declares are what the agent passes in. See [Automation](./automation).
- **Workflows** does the same for workflows.

### Invocation

Who and what may start the agent.

- **Enable User Chat Invocation** lets people chat with it.
- **Agent Sidebar Roles** decides which roles see it in the Assistant. This is
  about visibility only; it grants no access.
- **Enable System Invocation** lets the system run it without a person in the
  conversation, as a chatbot does. It then needs a **Service Account**: the
  user it acts as.

## Tools and their modes

The tools are the same operations Human offers over MCP: finding and changing
records, reading and changing the data model, pages and charts, running TAQs,
managing users and roles, and more. They are grouped by subject:

| Subject                 | Covers                                                              |
| ----------------------- | ------------------------------------------------------------------- |
| **Records**             | The rows a module holds, and the aggregates over them.              |
| **Data model**          | The namespaces and modules that decide what a record is.            |
| **Pages and charts**    | The pages, layouts and charts a namespace is read through.          |
| **Automation**          | TAQs and workflows, the triggers that start them, and their runs.   |
| **People and access**   | Users, groups and roles, and the auth clients that connect as them. |
| **Agents and chatbots** | Other agents, and the chatbots that put them in front of visitors.  |
| **Reminders**           | Reminders, which are only ever the invoking person's own.           |
| **Workspace**           | The applications in the sidebar, and the theme.                     |

Each tool also has a risk level: it **Reads**, **Writes** or **Deletes**. For a
whole subject, or for a single tool, you choose a mode:

| Mode               | Means                                                                       |
| ------------------ | --------------------------------------------------------------------------- |
| **Always allow**   | The agent calls the tool without asking.                                    |
| **Needs approval** | The chat asks the person before the tool runs, every time or once per chat. |
| **Blocked**        | The agent cannot use the tool.                                              |

Access is deny by default: a tool the agent has not been given is blocked. A
record or module tool can be narrowed further with **Limit to modules**, and any tool can carry a **Note for the model** saying when to use
it. An agent with no tools answers only from what it knows.

## Testing and inspecting

Once saved, an agent can be tested in the chat beside the editor.

- **Inspect** shows each run step by step: every tool call with its arguments
  and result, what was sent to the model, and the tokens used.
- **History** lists past conversations; opening one continues it.

## Where agents run

| Where                   | How                                                                                                                                                           |
| ----------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **The Assistant**       | The middle column of the [home screen](./home), and the agent panel in the top bar elsewhere. Needs user chat invocation and a sidebar role the person holds. |
| **An Agent Chat block** | A chat with chosen agents on any page. Needs user chat invocation. See [Namespaces and data](./namespaces#pages-layouts-and-blocks).                          |
| **A TAQ**               | The **Prompt Agent** step runs an agent and passes its answer on. It runs as the TAQ does and does not need system invocation.                                |
| **A chatbot**           | A **Conversation** step in a chatbot's journey. Needs system invocation and a service account. See [Chatbots](./chatbots).                                    |
| **An MCP client**       | A client connected to [Human's MCP server](#human-as-an-mcp-server) can run an agent as the signed-in user. Needs user chat invocation.                       |

In each case the agent acts with the permissions of whoever invoked it: the
person chatting, the user the TAQ runs as, or the chatbot agent's service
account. An agent cannot start another agent directly; to chain agents, use
**Prompt Agent** steps in a TAQ.

## Human as an MCP server

Human is also an MCP (Model Context Protocol) server. An AI client that speaks
MCP can connect to Human, sign in as a Human user and use the same tools agents
use, under that user's permissions. Nothing an MCP client does goes around
Human's permission checks.

## Permissions

Agents, LLM providers and chatbots have their own permissions: who may read,
change, delete and use each one. The
[system permissions reference](/reference/permissions/system) lists them.
