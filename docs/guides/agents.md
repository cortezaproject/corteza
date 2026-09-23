---
title: Agents
description: Create an AI agent, choose its model and tools, and run it from chat, automations and chatbots.
---

# Agents

An agent is a model with instructions, the tools it may call and the topics it
knows about. People chat with agents in the Assistant, and automations and
chatbots can run them too. Whoever runs an agent, **it only ever does what that
person is already allowed to do.** Granting an agent a tool never widens
anyone's access.

## Before you start: an LLM provider

Agents need a model. An administrator adds one under **Admin Area → System → LLM
Providers → New LLM Provider**:

| Setting                                      | What to enter                                                                                       |
| -------------------------------------------- | --------------------------------------------------------------------------------------------------- |
| **Provider**                                 | **Anthropic**, **Mistral**, or **Other (Custom URL)** for any service with an OpenAI-compatible API |
| **API Key**                                  | The vendor's key. Replace it later with **Update API Key**.                                         |
| **Default Model**                            | The model new agents start with                                                                     |
| **Temperature**, **Timeout**, **Max Tokens** | Defaults an agent can override                                                                      |

## Create an agent

Open **Agentic** from the menu and click **Create agent**. The editor has these
panels.

### General and Behavior

- **Name**, and a description so people know what the agent is for.
- **System Prompt**: the instructions sent to the model at the start of every
  conversation. Say who the agent is for, what it should do, and what it must
  not do.
- **Guardrails**: short rules the agent must follow, such as _"Stay on topic"_.

### Execution

Pick the **LLM Provider** and **Model**, and set limits: the most steps the
agent may take, a timeout, and how many tokens it may use.

### Access

This is where you decide what the agent can do.

- **Works in** confines every tool to one namespace. Leave it empty and a tool
  reaches whatever the invoking person can reach.
- **Tools** are grouped by section: _Records_, _Data model_, _Pages and charts_,
  _Automation_, _People and access_ and more. Set each section, or each tool, to
  **Always allow**, **Needs approval** or **Blocked**. With _Needs approval_,
  the chat asks the person before the tool runs.
- A tool can be narrowed further with **Limit to modules**, and given a **Note
  for the model** that says when to use it.
- **TAQs** lists the automations the agent may start.

With no tools at all, an agent answers only from what it knows. It cannot look
anything up or change anything.

### Knowledge Base

Attach **Topics** to give the agent domain knowledge and an understanding of
your data.

### Invocation

- **Enable User Chat Invocation** lets people chat with the agent.
- **Agent Sidebar Roles** decides which roles see it in the Assistant on the
  home page.
- **Enable System Invocation** lets a chatbot run it. It then needs a **Service
  Account**: the user it acts as when no person is present.

## Test and inspect

Save the agent, then chat with it in the editor.

- **Inspect** shows each run step by step: every tool call with its arguments
  and result, the context sent to the model, and the tokens used.
- **History** lists past conversations.

## Run an agent from elsewhere

- **The Assistant.** The middle column of the home page. People see the agents
  their roles allow.
- **A page.** Add an **Agent Chat** block to any page in the page builder.
- **Automations.** A TAQ's **Prompt Agent** step runs an agent and hands its
  answer to the next step.
- **Chatbots.** A step in a chatbot's journey can hand the visitor to an agent.
  The agent needs system invocation and a service account.
- **Your own AI client.** Human is also an MCP server, so Claude and other MCP
  clients can work with your Human data directly, under your permissions.

## Next steps

- [Build an app](./build-an-app): give your agent some data to work with.
- [Permissions reference](/reference/permissions/system): the operations on
  agents, chatbots and LLM providers.
