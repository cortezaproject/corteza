import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useChatbotStore = defineStore('chatbot', () => {
  const list = ref([])
  const loading = ref(false)

  async function fetchList(api) {
    loading.value = true
    try {
      const response = await api.chatbotList({ limit: 0 })
      list.value = response.set || []
      return list.value
    } catch (e) {
      console.error('Failed to fetch chatbots:', e)
      list.value = []
    } finally {
      loading.value = false
    }
  }

  function removeFromList(chatbotID) {
    list.value = list.value.filter(c => c.chatbotID !== chatbotID)
  }

  function updateInList(cb) {
    const idx = list.value.findIndex(c => c.chatbotID === cb.chatbotID)
    if (idx >= 0) {
      list.value[idx] = cb
    } else {
      list.value.push(cb)
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
