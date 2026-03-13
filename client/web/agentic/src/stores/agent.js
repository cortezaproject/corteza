import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAgentStore = defineStore('agent', () => {
  const list = ref([])
  const loading = ref(false)

  async function fetchList(api) {
    loading.value = true
    try {
      const response = await api.agentList({ limit: 0 })
      list.value = response.set || []
      return list.value
    } catch (e) {
      console.error('Failed to fetch agents:', e)
      list.value = []
    } finally {
      loading.value = false
    }
  }

  function removeFromList(agentID) {
    list.value = list.value.filter(a => a.agentID !== agentID)
  }

  function updateInList(agent) {
    const idx = list.value.findIndex(a => a.agentID === agent.agentID)
    if (idx >= 0) {
      list.value[idx] = agent
    } else {
      list.value.push(agent)
    }
  }

  return {
    list,
    loading,
    fetchList,
    removeFromList,
    updateInList,
  }
})
