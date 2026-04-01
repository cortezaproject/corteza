import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAgentSidebarStore = defineStore('agentSidebar', () => {
  const visible = ref(false)
  const availableAgents = ref<any[]>([])
  const activeAgentID = ref<string | null>(null)
  const conversations = ref<Record<string, any[]>>({})
  const activeConversationIndex = ref<Record<string, number>>({})

  function toggleVisibility() {
    visible.value = !visible.value
  }

  function setVisible(value: boolean) {
    visible.value = value
  }

  function setAvailableAgents(agents: any[]) {
    availableAgents.value = agents
    if (agents.length > 0 && (!activeAgentID.value || !agents.find(a => a.agentID === activeAgentID.value))) {
      activeAgentID.value = agents[0].agentID
    }
  }

  function setActiveAgentID(agentID: string) {
    activeAgentID.value = agentID
  }

  function initConversation(agentID: string) {
    if (!conversations.value[agentID]) {
      conversations.value[agentID] = [{
        id: 1,
        messages: [],
        conversationID: null,
      }]
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
    conversations.value[agentID] = [{
      id: 1,
      messages: [],
      conversationID: null,
    }]
    activeConversationIndex.value[agentID] = 0
  }

  return {
    visible,
    availableAgents,
    activeAgentID,
    conversations,
    activeConversationIndex,
    toggleVisibility,
    setVisible,
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
    clearConversation
  }
})
