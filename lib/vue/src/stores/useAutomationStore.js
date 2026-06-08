import { automation } from '@planetcrust/human-js'
import { defineStore } from 'pinia'
import { inject, ref } from 'vue'

const { NgAutomation } = automation

export const useAutomationStore = defineStore('automation', () => {
  const $AutomationAPI = inject('$AutomationAPI')

  const list = ref([])
  const loading = ref(false)
  const error = ref(null)

  const functions = ref([])
  const triggers = ref([])
  const catalogReady = ref(false)

  async function fetchList(filter = { disabled: 1 }) {
    loading.value = true
    error.value = null
    try {
      const response = await $AutomationAPI.ngAutomationList(filter)
      list.value = (response.set || []).map(item => new NgAutomation(item))
      return list.value
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch automations'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function create(data) {
    loading.value = true
    error.value = null
    try {
      const response = await $AutomationAPI.ngAutomationCreate(data)
      const created = new NgAutomation(response)
      list.value.push(created)
      return created
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to create automation'
      throw e
    } finally {
      loading.value = false
    }
  }

  function updateInList(automation) {
    const idx = list.value.findIndex(a => a.automationID === automation.automationID)
    if (idx >= 0) {
      list.value[idx] = automation
    } else {
      list.value.push(automation)
    }
  }

  function removeFromList(automationID) {
    list.value = list.value.filter(a => a.automationID !== automationID)
  }

  async function remove(automationID) {
    loading.value = true
    error.value = null
    try {
      await $AutomationAPI.ngAutomationDelete({ automationID })
      list.value = list.value.filter(a => a.automationID !== automationID)
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to delete automation'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function loadFunctions() {
    try {
      const response = await $AutomationAPI.constructLibraryFunctions()
      functions.value = response.set || []
    } catch (e) {
      console.error('Failed to load functions:', e)
    }
  }

  async function loadTriggers() {
    try {
      const response = await $AutomationAPI.constructLibraryTriggers()
      triggers.value = response.set || []
    } catch (e) {
      console.error('Failed to load triggers:', e)
    }
  }

  async function loadCatalog() {
    await Promise.all([loadFunctions(), loadTriggers()])
    catalogReady.value = true
  }

  function reset() {
    list.value = []
    loading.value = false
    error.value = null
  }

  return {
    list,
    loading,
    error,
    functions,
    triggers,
    catalogReady,

    fetchList,
    create,
    updateInList,
    removeFromList,
    remove,
    loadFunctions,
    loadTriggers,
    loadCatalog,
    reset,
  }
})
