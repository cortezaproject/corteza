---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/stores/useAgentChatStore.ts
touched-by: []
tests: []
---

# agent/

## Intention

The AI-agent chat UX, reused verbatim by every surface that talks to agents
(topbar sidebar, compose page block, agent preview). All conversation state
lives in `useAgentChatStore`; these components are the shared view over it.

## Map

- `CAgentChat.vue` — the full chat: agent picker, tabs, messages, composer, history.
- `CAgentSidebar.vue` / `CAgentSidebarButton.vue` — right-sidebar host + topbar toggle.
- `CChatMessages.vue` / `CChatComposer.vue` — message list and input; composer is
  also reused by chatbot/ (inbox reply box).
- `CConversationTabs.vue` — multi-conversation tab strip.
- `translations.ts` — `AgentChatTranslations` contract + `makeAgentChatTranslations(t, prefix)`.

## Contracts consumers rely on

- Zero i18n coupling: every caller passes a `translations` object built from
  its own key prefix via `makeAgentChatTranslations`. New strings extend the
  interface and every call-site prefix.
- CAgentChat props: `allowedAgentIDs` (allowlist — when set, agents resolve
  via API and bypass sidebar roles), `defaultAgentID`, `autoResume`, and
  `contextProvider` (callable whose result is sent as `context` on agentExec —
  how the page block attaches record/page IDs without UI changes).
- APIs arrive via inject (`$SystemAPI`, `$Auth`), never imports.

## When changing this

Multiple surfaces share one store — verify sidebar, page block, and preview
after changes; a store-shape change breaks all three at once.
