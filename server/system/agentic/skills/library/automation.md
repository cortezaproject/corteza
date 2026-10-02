---
name: automation
description: Rules for running an existing TAQ or workflow — over MCP with automation_taq_exec, or as an agent through the per-TAQ tools it was granted. Writing one is taq_authoring.
triggers:
  - automation_taq_lookup
  - automation_taq_exec
  - automation_taq_executions
  - automation_taq_execution_trace
  - automation_workflow_lookup
  - automation_workflow_exec
---

# Running automations

TAQs (Trigger Action Queries) are pre-built automations. This skill is about
running one; writing one is `taq_authoring`. Unless the request is to build an
automation, run what exists rather than changing it, and prefer running a TAQ
over doing the same thing step by step.

## Over MCP

1. `automation_taq_lookup` finds the TAQ by handle or name and shows its
   triggers and whether it is `runnable`.
2. `automation_taq_exec` runs it: `taq` is the handle or ID, `input` must
   satisfy the chosen trigger's input schema (the lookup shows it), and the
   first trigger is used unless `entryPoint` names another.
3. `exec` reports status only. Read `automation_taq_execution_trace` for what
   each step received and produced; `automation_taq_executions` lists past
   runs.

"manager: executable not found" means the TAQ has validation `issues` and is
not runnable — not that the ID is wrong. A disabled TAQ also refuses to run.

## As an agent

An agent holds the TAQs it may run as individual tools named `automation_<id>`,
listed under **AVAILABLE AUTOMATIONS** in its context. Match the request to a
tool by its name and description, check its input schema, ask for or generate
the required inputs, and call that tool directly rather than a generic exec.
Only run what you were granted; if a call is denied, stop and tell the user.

## Workflows

The same rules apply to the older Workflows designer: use the `internal-id`
listed under AVAILABLE AUTOMATIONS with `automation_workflow_exec`, and call
`automation_workflow_lookup` only when it is not listed.
