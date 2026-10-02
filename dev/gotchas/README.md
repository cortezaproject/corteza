# Gotchas

Non-obvious behaviour of this codebase, by area: what the code does that its
names, docs and tests do not make plain. Each entry states what is true and how
to work with it. Intent docs (`*.intent.md`) hold the contracts; these hold the
traps. Read the area file before working in it; add an entry when something
surprised you and cost a turn.

| file                             | covers                                                                                     |
| -------------------------------- | ------------------------------------------------------------------------------------------ |
| [server.md](server.md)           | store and DAL, RBAC and auth, REST verbs, envoy and provisioning, settings, codegen, tests |
| [compose.md](compose.md)         | namespaces, modules, records, pages, blocks, charts                                        |
| [taq.md](taq.md)                 | TAQ (Trigger Action Query) binding, expressions, scheduler, traces, workflow canvas        |
| [agentic.md](agentic.md)         | agents, tool access, caller context, chatbots, Human's own MCP server                      |
| [project.md](project.md)         | project lifecycle, revisions, AI systems and FRIA, deletion, isolation                     |
| [frontend.md](frontend.md)       | unify webapp, lib/vue, Vue 3, PrimeVue, Tailwind, CodeMirror, i18n                         |
| [dev-toolkit.md](dev-toolkit.md) | dev/agent scripts, dev MCP tools, worktrees and git, tests and e2e, browser checks         |
| [intent.md](intent.md)           | the intent tooling beyond what .intent/SPEC.md states                                      |
