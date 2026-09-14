---
name: agent_building
description: Creating an AI agent and the chatbot that fronts it — the sections that replace wholesale, and the three settings that decide whether either answers anything.
triggers:
  - system_agent_create
  - system_agent_update
  - system_agent_lookup
  - system_chatbot_create
  - system_chatbot_update
  - system_chatbot_lookup
  - system_llm_provider_lookup
  - system_agent_exec
---

# Agents and chatbots

An **agent** is a model, a system prompt and a set of tools. A **chatbot** is
the embeddable widget in front of one. A chatbot with no agent answers nothing,
so build the agent first.

## Every section replaces

An agent's configuration is five objects — `meta`, `behavior`, `execution`,
`access`, `invocation` — and a chatbot's is three: `handoff`, `styling`,
`scenarios`. Sending one **replaces that whole section**. To change one field,
read the resource first with the matching `_lookup` and send an edited copy of
the section, never a fragment.

`scenarios` is a collection and replaces the same way: send every scenario the
chatbot should have, not just the new one.

## An agent with no tools reaches nothing

`access` is deny-by-default. An agent whose `tools` list is empty still answers
from its prompt and its knowledge bases, and cannot look anything up or change
anything. Grant it something in the same call:

```
{"allow":[{"namespaceID":"<id>"}],"tools":[{"name":"compose_record_lookup","permission":"always","description":"Read candidates"}]}
```

Granting a tool never widens what the agent may see — it runs as the person who
invoked it, so RBAC still applies on top.

The `namespace` argument is a shortcut: it sets `access.allow` to that
namespace. It is **ignored when you also send `access`**, and it grants no tool
on its own.

The platform context — what namespaces, modules and records are, and the rules
for working with them — is not a setting. The server adds it to the prompt of
every agent that holds at least one tool in `access.tools`, and leaves it out of
an agent that holds none.

## The model is settled when the agent is written

`system_llm_provider_lookup` lists the providers this instance has and, with
`models` set, the model names a provider advertises. Both go in
`execution.model`.

Creation resolves both and stores what it resolved, so an agent always carries
the provider and model it will actually run on.

- Leaving `llmProviderID` out is fine where one provider is active — that one is
  written onto the agent. Where several are, the write is refused naming them.
- A model name is required. If the agent names none and the provider has no
  `config.model` of its own, the write is refused and the provider's models are
  listed in the error.

`system_agent_update` settles it the same way, so an edit cannot leave an agent
in a state creation would have refused.

## Run it before you call it done

`system_agent_exec` sends the agent a message and returns what it said. Storing
an agent proves nothing: a prompt that ignores its tools, a tool that was never
granted and a model that cannot be resolved all look identical until a run.

Read `toolCalls`, not just `output`. An agent that answered plausibly having
called nothing made the answer up from its prompt — an empty `toolCalls` on a
question that needed data is the tell.

A `status` of `awaiting_approval` is not a failure: a tool granted with
permission `ask` stops the run every time. `pendingApproval` says which tool and
with what arguments; send the call again with that name in `approvedTools` to
carry on.

An agent cannot call `system_agent_exec` to start another agent. Chaining is a
TAQ step (`agentPrompt`) or a workflow step (`agentRun`).

## Three things decide whether the chatbot answers

1. A `conversation` scenario carrying an `agentID`. Without one the **whole
   save is rejected** — "conversation scenario is missing an agent" — not just
   that scenario.
2. The agent's `invocation.system` must be enabled with a `serviceAccount`. A
   chatbot runs its agent under a service account, and without it the visitor's
   first message fails with "agent not configured". `system_chatbot_create`
   checks this and says so in its returned note — read it.
3. `enabled` defaults to **false** on a new chatbot. It keeps its whole
   configuration and serves nobody until you turn it on.

`allowedOrigins` is the access control on the widget, not decoration: an empty
list restricts nothing, so a chatbot with none is usable from any site that has
its widget key. The key is generated on create and returned; you cannot choose
it, and `system_chatbot_regenerate_key` invalidates the previous one the moment
it returns.
