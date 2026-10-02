# Claude Code in this repo

@AGENTS.md

The shared rules above apply in full. What follows is only what Claude Code adds.

- **Skills** carry the procedures: `/dev-task` (any piece of work that is not a
  one-liner), `/dev-change` (one change), `/intent-task` and `/intent-audit`
  (the intent system's opt-in), `/dev-api`, `/dev-seed`, `/sys-design`
  (building whole systems in Human), `/corredor-script`, `/orchestrate`
  (several independent issues in parallel lanes — invoking it is what
  authorises spawning subagents; nothing else does).
- **Reporting**: `/dev-task` and `/dev-change` reply in the standard block —
  sections, confidence grade, open things asked as an interview. The format is
  `.claude/reporting.md`. Questions go through `AskUserQuestion`, priced.
- **MCP servers**: `human-dev` (`.mcp.json`) is the repo toolkit —
  `dev_test_run`, `dev_format_run`, `dev_commit_create` (enforces the commit
  convention), `dev_intent_*`, `dev_server_status`, `dev_ui_verify`.
  `human-local`, this checkout's own `/api/mcp` with its ~125 configurator
  tools, attaches only with `make claude MCP=1`; `mcp__human-local__*` is the
  one Human MCP family that reaches the dev server.
- **Dev MCP edits** do not reach the running session; restart the client, or
  verify through a Go test.
