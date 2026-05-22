// Translations contract for the shared CAgentChat component.
// Each caller (sidebar, page block, agent preview, ...) owns its own i18n keys
// and builds this object from its own prefix via makeAgentChatTranslations(t, prefix).

export interface AgentChatTranslations {
  title: string
  placeholder: string
  thinking: string
  empty: string
  newChat: string
  clearAllChats: string
  chatTab: (id: number) => string
  noAgents: string
  history: {
    button: string
    empty: string
    loading: string
    deleteTooltip: string
    untitled: string
  }
}

type TFn = (key: string, params?: Record<string, any>) => string

export function makeAgentChatTranslations(t: TFn, prefix: string): AgentChatTranslations {
  const k = (key: string, params?: Record<string, any>) => t(`${prefix}${key}`, params || {})
  return {
    title: k('title'),
    placeholder: k('placeholder'),
    thinking: k('thinking'),
    empty: k('empty'),
    newChat: k('newChat'),
    clearAllChats: k('clearAllChats'),
    chatTab: (id: number) => k('chatTab', { id }),
    noAgents: k('noAgents'),
    history: {
      button: k('history.button'),
      empty: k('history.empty'),
      loading: k('history.loading'),
      deleteTooltip: k('history.deleteTooltip'),
      untitled: k('history.untitled'),
    },
  }
}
