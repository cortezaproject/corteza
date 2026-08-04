# MCP tool coverage

Generated. Do not edit by hand — run:

```sh
cd server && go test ./tests/mcp/ -run TestToolsMatrix -update
```

Resource-to-tool naming is fixed by `RESOURCES.md`; the rules these
tools are held to are in `CONVENTIONS.md`.

## Registered tools (93)

| Tool | Group | Risk | Surface |
|---|---|---|---|
| `automation_event_type_lookup` | configuring | read | both |
| `automation_taq_create` | configuring | write | both |
| `automation_taq_delete` | configuring | destructive | both |
| `automation_taq_exec` | usage | write | both |
| `automation_taq_execution_trace` | usage | read | both |
| `automation_taq_executions` | usage | read | both |
| `automation_taq_lookup` | configuring | read | both |
| `automation_taq_undelete` | configuring | write | both |
| `automation_taq_update` | configuring | write | both |
| `automation_trigger_create` | configuring | write | both |
| `automation_trigger_delete` | configuring | destructive | both |
| `automation_trigger_lookup` | configuring | read | both |
| `automation_trigger_undelete` | configuring | write | both |
| `automation_trigger_update` | configuring | write | both |
| `automation_workflow_create` | configuring | write | both |
| `automation_workflow_delete` | configuring | destructive | both |
| `automation_workflow_exec` | usage | write | both |
| `automation_workflow_lookup` | configuring | read | both |
| `automation_workflow_undelete` | configuring | write | both |
| `automation_workflow_update` | configuring | write | both |
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
| `system_application_create` | configuring | write | both |
| `system_application_delete` | configuring | destructive | both |
| `system_application_flag` | configuring | write | both |
| `system_application_lookup` | configuring | read | both |
| `system_application_reorder` | configuring | write | both |
| `system_application_undelete` | configuring | write | both |
| `system_application_unflag` | configuring | write | both |
| `system_application_update` | configuring | write | both |
| `system_auth_client_delete` | configuring | destructive | both |
| `system_auth_client_lookup` | configuring | read | both |
| `system_auth_client_undelete` | configuring | write | both |
| `system_auth_client_update` | configuring | write | both |
| `system_reminder_create` | usage | write | both |
| `system_reminder_delete` | usage | destructive | both |
| `system_reminder_dismiss` | usage | write | both |
| `system_reminder_lookup` | usage | read | both |
| `system_reminder_snooze` | usage | write | both |
| `system_reminder_undismiss` | usage | write | both |
| `system_reminder_update` | usage | write | both |
| `system_role_archive` | configuring | write | both |
| `system_role_clone_rules` | configuring | destructive | both |
| `system_role_create` | configuring | write | both |
| `system_role_delete` | configuring | destructive | both |
| `system_role_lookup` | configuring | read | both |
| `system_role_member_add` | configuring | write | both |
| `system_role_member_list` | configuring | read | both |
| `system_role_member_remove` | configuring | write | both |
| `system_role_unarchive` | configuring | write | both |
| `system_role_undelete` | configuring | write | both |
| `system_role_update` | configuring | write | both |
| `system_user_create` | configuring | write | both |
| `system_user_delete` | configuring | destructive | both |
| `system_user_group_create` | configuring | write | both |
| `system_user_group_delete` | configuring | destructive | both |
| `system_user_group_lookup` | configuring | read | both |
| `system_user_group_member_add` | configuring | write | both |
| `system_user_group_member_list` | configuring | read | both |
| `system_user_group_undelete` | configuring | write | both |
| `system_user_group_update` | configuring | write | both |
| `system_user_lookup` | configuring | read | both |
| `system_user_set_email_confirmed` | configuring | write | both |
| `system_user_suspend` | configuring | write | both |
| `system_user_undelete` | configuring | write | both |
| `system_user_unsuspend` | configuring | write | both |
| `system_user_update` | configuring | write | both |

## Totals

| Dimension | Value | Tools |
|---|---|---|
| group | configuring | 76 |
| group | usage | 17 |
| risk | destructive | 15 |
| risk | read | 21 |
| risk | write | 57 |
