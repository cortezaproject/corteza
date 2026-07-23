---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/components/agent/CChatComposer.vue
touched-by: []
tests: []
---

# chatbot/

## Intention

Operator-facing chatbot session inbox: browse/watch sessions of selected
chatbots, read transcripts, reply. One component reused wherever an inbox is
needed (unify chatbot section, admin Sessions tab) — pure UX + SDK calls, no
opinion on what wraps it.

## Map

- `CChatbotInbox.vue` — session list + transcript + reply composer.
- `translations.ts` — `CHATBOT_INBOX_TRANSLATION_KEYS`, `makeChatbotInboxTranslations`.

## Contracts consumers rely on

- Props: `chatbotIDs` (empty → "pick chatbots" empty state), `statusFilter`
  (empty → all), `autoOpenFirst`, `showFilter`, `refreshRate` (seconds; 0
  disables polling), `translations`.
- Zero vue-i18n coupling: strings resolve from the `translations` prop only;
  missing keys render empty (visible, but never leak key names). Consumers
  build the object from their own `t()` via `makeChatbotInboxTranslations`.
- APIs via inject (`$SystemAPI`, `$Auth`); reply box reuses agent/CChatComposer.

## When changing this

New user-facing strings must be added to the translation keys/factory and to
every consumer's locale prefix — there are no English fallbacks in `l()`.
