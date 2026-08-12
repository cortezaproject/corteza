import { defineStore } from 'pinia'
import { inject, ref } from 'vue'

export const useAgentChatStore = defineStore('agentChat', () => {
  const $SystemAPI = inject<any>('$SystemAPI')
  const availableAgents = ref<any[]>([])
  const activeAgentID = ref<string | null>(null)
  const conversations = ref<Record<string, any[]>>({})
  const activeConversationIndex = ref<Record<string, number>>({})

  // History — server-persisted past conversations, lazy-loaded per agent.
  const history = ref<Record<string, any[]>>({})
  const historyLoading = ref<Record<string, boolean>>({})
  const historyLoadedAt = ref<Record<string, number>>({})

  function setAvailableAgents(agents: any[]) {
    availableAgents.value = [...agents].sort((a, b) => {
      const nameA = a.meta?.short || a.handle || ''
      const nameB = b.meta?.short || b.handle || ''
      return nameA.localeCompare(nameB)
    })
    if (
      availableAgents.value.length > 0 &&
      (!activeAgentID.value || !availableAgents.value.find(a => a.agentID === activeAgentID.value))
    ) {
      activeAgentID.value = availableAgents.value[0].agentID
    }
  }

  function setActiveAgentID(agentID: string) {
    activeAgentID.value = agentID
  }

  function initConversation(agentID: string) {
    if (!conversations.value[agentID] || conversations.value[agentID].length === 0) {
      conversations.value[agentID] = [
        {
          id: 1,
          messages: [],
          conversationID: null,
        },
      ]
      activeConversationIndex.value[agentID] = 0
    }
  }

  function getConversation(agentID: string) {
    if (!conversations.value[agentID]) {
      initConversation(agentID)
    }
    const idx = activeConversationIndex.value[agentID] || 0
    return conversations.value[agentID][idx]
  }

  function getAllConversations(agentID: string) {
    if (!conversations.value[agentID]) {
      initConversation(agentID)
    }
    return conversations.value[agentID]
  }

  function addMessage(agentID: string, message: any) {
    if (!conversations.value[agentID]) initConversation(agentID)
    const idx = activeConversationIndex.value[agentID] || 0
    conversations.value[agentID][idx].messages.push(message)
  }

  function setConversationID(agentID: string, conversationID: string) {
    if (!conversations.value[agentID]) initConversation(agentID)
    const idx = activeConversationIndex.value[agentID] || 0
    conversations.value[agentID][idx].conversationID = conversationID
  }

  function setActiveConversationIndex(agentID: string, index: number) {
    if (conversations.value[agentID] && conversations.value[agentID][index]) {
      activeConversationIndex.value[agentID] = index
    }
  }

  function startNewConversation(agentID: string) {
    if (!conversations.value[agentID]) {
      initConversation(agentID)
      return
    }

    // Determine highest ID to avoid duplicate keys
    const maxId = conversations.value[agentID].reduce((acc, curr) => Math.max(acc, curr.id), 0)

    conversations.value[agentID].push({
      id: maxId + 1,
      messages: [],
      conversationID: null,
    })

    activeConversationIndex.value[agentID] = conversations.value[agentID].length - 1
  }

  function closeConversation(agentID: string, index: number) {
    if (!conversations.value[agentID]) return

    conversations.value[agentID].splice(index, 1)

    if (conversations.value[agentID].length === 0) {
      initConversation(agentID)
    } else {
      let currentActive = activeConversationIndex.value[agentID]
      if (index === currentActive) {
        // If closing current active, step back 1 or stay at 0
        activeConversationIndex.value[agentID] = Math.max(0, currentActive - 1)
      } else if (index < currentActive) {
        // If closing a tab before the active one, shift active index down
        activeConversationIndex.value[agentID] = currentActive - 1
      }
    }
  }

  function clearConversation(agentID: string) {
    conversations.value[agentID] = [
      {
        id: 1,
        messages: [],
        conversationID: null,
      },
    ]
    activeConversationIndex.value[agentID] = 0
  }

  // --- History (server-persisted) ---

  async function loadHistory(agentID: string) {
    if (!agentID) return
    historyLoading.value = { ...historyLoading.value, [agentID]: true }
    try {
      const res = await $SystemAPI.aiConversationList({ agentID, deleted: 0, limit: 200 })
      const set = (res?.set || []) as any[]
      // Newest first; backend may already sort, but we don't trust it.
      set.sort((a, b) => {
        const ta = new Date(a.updatedAt || a.createdAt || 0).getTime()
        const tb = new Date(b.updatedAt || b.createdAt || 0).getTime()
        return tb - ta
      })
      history.value = { ...history.value, [agentID]: set }
      historyLoadedAt.value = { ...historyLoadedAt.value, [agentID]: Date.now() }
    } finally {
      historyLoading.value = { ...historyLoading.value, [agentID]: false }
    }
  }

  function getHistory(agentID: string): any[] {
    return history.value[agentID] || []
  }

  function isHistoryLoading(agentID: string): boolean {
    return !!historyLoading.value[agentID]
  }

  // Open a past conversation in a new tab (or focus existing tab if already open).
  function openConversationFromHistory(agentID: string, conv: any) {
    if (!agentID || !conv) return
    initConversation(agentID)

    const cid = String(conv.aiConversationID)
    const existingIdx = conversations.value[agentID].findIndex(
      (c: any) => c.conversationID && String(c.conversationID) === cid,
    )
    if (existingIdx >= 0) {
      activeConversationIndex.value[agentID] = existingIdx
      return
    }

    const messages = (conv.messages || []).map((m: any) => ({
      role: m.role === 'assistant' ? 'agent' : m.role,
      content: m.content || '',
    }))

    // If the only tab is an empty, fresh one, reuse it instead of stacking.
    const current = conversations.value[agentID]
    const last = current[current.length - 1]
    if (
      current.length === 1 &&
      last &&
      !last.conversationID &&
      (!last.messages || last.messages.length === 0)
    ) {
      last.conversationID = cid
      last.messages = messages
      activeConversationIndex.value[agentID] = current.length - 1
      return
    }

    const maxId = current.reduce((acc, curr) => Math.max(acc, curr.id), 0)
    current.push({
      id: maxId + 1,
      messages,
      conversationID: cid,
    })
    activeConversationIndex.value[agentID] = current.length - 1
  }

  // Soft-delete a conversation server-side, then prune locally.
  async function deleteConversation(agentID: string, conversationID: string) {
    if (!agentID || !conversationID) return
    await $SystemAPI.aiConversationDelete({ aiConversationID: conversationID })

    // Remove from history list
    if (history.value[agentID]) {
      history.value[agentID] = history.value[agentID].filter(
        (c: any) => String(c.aiConversationID) !== String(conversationID),
      )
    }

    // Remove from any open tabs
    if (conversations.value[agentID]) {
      for (let i = conversations.value[agentID].length - 1; i >= 0; i--) {
        const c = conversations.value[agentID][i]
        if (c.conversationID && String(c.conversationID) === String(conversationID)) {
          closeConversation(agentID, i)
        }
      }
    }
  }

  return {
    availableAgents,
    activeAgentID,
    conversations,
    activeConversationIndex,
    history,
    historyLoading,
    historyLoadedAt,
    setAvailableAgents,
    setActiveAgentID,
    initConversation,
    getConversation,
    getAllConversations,
    addMessage,
    setConversationID,
    setActiveConversationIndex,
    startNewConversation,
    closeConversation,
    clearConversation,
    loadHistory,
    getHistory,
    isHistoryLoading,
    openConversationFromHistory,
    deleteConversation,
  }
})
