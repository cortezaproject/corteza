---
name: automation
description: Rules for executing TAQs and Workflows.
triggers:
  - automation_taq_lookup
  - automation_workflow_lookup
  - automation_workflow_exec
---

# TAQs

TAQs are pre-built automations. Allowed TAQs are exposed directly to you as individual tools prefixed with `automation_`. Run them — do not create or modify them.

1. Match the user's intent to the appropriate tool based on its name and description.
2. Check the input schema for the tool. If required inputs are missing, ask the user for them or generate them if instructed to do so.
3. Call the specific `automation_<id>` tool directly with the required arguments matching its JSON schema. Do not use generic execution or lookup verbs to execute TAQs.

Rules:

- Only execute TAQs you have been granted access to. If denied, stop and tell the user.
- Prefer running a TAQ over doing the same thing step by step.

# Workflows

Same rules as TAQs. If the workflow is in **AVAILABLE AUTOMATIONS**, use its `internal-id` directly with `automation_workflow_exec`. Only call `automation_workflow_lookup` if it's not listed there.
