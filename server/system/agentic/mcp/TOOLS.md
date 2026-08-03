# MCP tool coverage

Generated. Do not edit by hand — run:

```sh
cd server && go test ./tests/mcp/ -run TestToolsMatrix -update
```

Resource-to-tool naming is fixed by `RESOURCES.md`; the rules these
tools are held to are in `CONVENTIONS.md`.

## Registered tools (43)

| Tool | Group | Risk | Surface |
|---|---|---|---|
| `automation_taq_exec` | usage | write | both |
| `automation_taq_execution_trace` | usage | read | both |
| `automation_taq_executions` | usage | read | both |
| `automation_taq_lookup` | configuring | read | both |
| `automation_taq_undelete` | configuring | write | both |
| `automation_workflow_exec` | usage | write | both |
| `automation_workflow_lookup` | configuring | read | both |
| `automation_workflow_undelete` | configuring | write | both |
| `compose_chart_create` | configuring | write | both |
| `compose_chart_delete` | configuring | destructive | both |
| `compose_chart_lookup` | configuring | read | both |
| `compose_chart_undelete` | configuring | write | both |
| `compose_chart_update` | configuring | write | both |
| `compose_module_create` | configuring | write | both |
| `compose_module_delete` | configuring | destructive | both |
| `compose_module_lookup` | configuring | read | both |
| `compose_module_undelete` | configuring | write | both |
| `compose_module_update` | configuring | write | both |
| `compose_namespace_create` | configuring | write | both |
| `compose_namespace_delete` | configuring | destructive | both |
| `compose_namespace_lookup` | configuring | read | both |
| `compose_namespace_update` | configuring | write | both |
| `compose_page_block_schema` | configuring | read | both |
| `compose_page_create` | configuring | write | both |
| `compose_page_delete` | configuring | destructive | both |
| `compose_page_lookup` | configuring | read | both |
| `compose_page_remove_blocks` | configuring | write | both |
| `compose_page_reorder` | configuring | write | both |
| `compose_page_undelete` | configuring | write | both |
| `compose_page_update` | configuring | write | both |
| `compose_record_create` | usage | write | both |
| `compose_record_delete` | usage | destructive | both |
| `compose_record_lookup` | usage | read | both |
| `compose_record_undelete` | usage | write | both |
| `compose_record_update` | usage | write | both |
| `discovery_search` | usage | read | both |
| `system_reminder_create` | usage | write | both |
| `system_reminder_delete` | usage | destructive | both |
| `system_reminder_dismiss` | usage | write | both |
| `system_reminder_lookup` | usage | read | both |
| `system_reminder_snooze` | usage | write | both |
| `system_reminder_undismiss` | usage | write | both |
| `system_reminder_update` | usage | write | both |

## Totals

| Dimension | Value | Tools |
|---|---|---|
| group | configuring | 26 |
| group | usage | 17 |
| risk | destructive | 6 |
| risk | read | 12 |
| risk | write | 25 |
