import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useWorkflowStore = defineStore('workflow', () => {
  const list = ref([])
  const loading = ref(false)

  async function fetchList(api, params = {}) {
    loading.value = true
    try {
      const response = await api.workflowList(params)
      list.value = response.set || []
      return response
    } catch (e) {
      console.error('Failed to fetch workflows:', e)
      list.value = []
      return { set: [], filter: {} }
    } finally {
      loading.value = false
    }
  }

  function removeFromList(workflowID) {
    list.value = list.value.filter(w => w.workflowID !== workflowID)
  }

  function updateInList(workflow) {
    const idx = list.value.findIndex(w => w.workflowID === workflow.workflowID)
    if (idx >= 0) {
      list.value[idx] = workflow
    } else {
      list.value.push(workflow)
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
