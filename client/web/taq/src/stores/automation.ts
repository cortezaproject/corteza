import { automation } from '@cortezaproject/corteza-js-next'
import type { IconDef } from '@cortezaproject/corteza-js-next/src/automation/types/icon'
import { defineStore } from 'pinia'
import { ref } from 'vue'

const NgAutomation = automation.NgAutomation
type NgAutomationInstance = InstanceType<typeof NgAutomation>

// Segment element input definition
export interface SegmentInput {
  type?: string
  label?: string
  argument?: string
  placeholder?: string
  visual?: Record<string, unknown>
}

// Segment element
export interface SegmentElement {
  input?: SegmentInput
}

// Segment section
export interface SegmentSection {
  meta?: { short?: string }
  elements?: SegmentElement[]
}

// Segment definition
export interface Segment {
  meta?: { short?: string }
  sections?: SegmentSection[]
}

// Function type from automation catalog
export interface AutomationFunction {
  ref: string
  kind?: string
  groups?: string[]
  meta?: { short?: string; description?: string; icon?: IconDef | string }
  parameters?: Array<{
    argumentName: string
    types?: string[]
    required?: boolean
    aggregate?: boolean
  }>
  results?: Array<{
    argumentName: string
    types?: string[]
  }>
  segments?: Segment[]
}

// Constraint definition from construct library (allowed names + types)
export interface CatalogConstraint {
  name: string
  types: string[]
  meta?: Record<string, unknown>
}

// Trigger type from construct library catalog
export interface AutomationTrigger {
  resourceType: string
  eventType: string
  meta?: { short?: string; description?: string; icon?: IconDef | string }
  constraints?: CatalogConstraint[]
  properties?: Array<{
    name: string
    type: string
    meta?: { short?: string; description?: string }
  }>
  segments?: Segment[]
}

// API interface for type safety
interface AutomationAPI {
  ngAutomationList: (
    filter?: Record<string, unknown>,
  ) => Promise<{ set: Record<string, unknown>[] }>
  ngAutomationCreate: (data: Record<string, unknown>) => Promise<Record<string, unknown>>
  ngAutomationDelete: (params: { automationID: string }) => Promise<void>
  constructLibraryFunctions: () => Promise<{ set: AutomationFunction[] }>
  constructLibraryTriggers: () => Promise<{ set: AutomationTrigger[] }>
}

/**
 * Pinia store for shared TAQ state.
 *
 * This store manages:
 * - Automation list (for List view)
 * - Function/event type catalogs (shared across app)
 * - Loading/error states for list operations
 *
 * Note: Individual automation editing is handled by useFlowEditor composable.
 */
export const useAutomationStore = defineStore('automation', () => {
  // List of all automations (for List view)
  const list = ref<NgAutomationInstance[]>([])

  // Loading and error states
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Catalogs from backend (shared across app)
  const functions = ref<AutomationFunction[]>([])
  const triggers = ref<AutomationTrigger[]>([])
  const catalogReady = ref(false)

  /**
   * Fetch list of automations from API
   */
  async function fetchList(api: AutomationAPI, filter: Record<string, unknown> = { disabled: 1 }) {
    loading.value = true
    error.value = null
    try {
      const response = await api.ngAutomationList(filter)
      list.value = (response.set || []).map(
        (item: Record<string, unknown>) => new NgAutomation(item),
      )
      return list.value
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch automations'
      throw e
    } finally {
      loading.value = false
    }
  }

  /**
   * Create a new automation (used by List view)
   */
  async function create(api: AutomationAPI, data: Partial<NgAutomationInstance>) {
    loading.value = true
    error.value = null
    try {
      const response = await api.ngAutomationCreate(data as Record<string, unknown>)
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

  /**
   * Add or update an automation in the list
   */
  function updateInList(automation: NgAutomationInstance) {
    const idx = list.value.findIndex(a => a.automationID === automation.automationID)
    if (idx >= 0) {
      list.value[idx] = automation
    } else {
      list.value.push(automation)
    }
  }

  /**
   * Delete an automation
   */
  async function remove(api: AutomationAPI, automationID: string) {
    loading.value = true
    error.value = null
    try {
      await api.ngAutomationDelete({ automationID })
      list.value = list.value.filter(a => a.automationID !== automationID)
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to delete automation'
      throw e
    } finally {
      loading.value = false
    }
  }

  /**
   * Load functions catalog from construct library API
   */
  async function loadFunctions(api: AutomationAPI) {
    try {
      const response = await api.constructLibraryFunctions()
      functions.value = response.set || []
    } catch (e) {
      console.error('Failed to load functions:', e)
    }
  }

  /**
   * Load triggers catalog from construct library API
   */
  async function loadTriggers(api: AutomationAPI) {
    try {
      const response = await api.constructLibraryTriggers()
      triggers.value = response.set || []
    } catch (e) {
      console.error('Failed to load triggers:', e)
    }
  }

  /**
   * Load all catalogs
   */
  async function loadCatalog(api: AutomationAPI) {
    await Promise.all([loadFunctions(api), loadTriggers(api)])
    catalogReady.value = true
  }

  /**
   * Reset store state
   */
  function reset() {
    list.value = []
    loading.value = false
    error.value = null
  }

  return {
    // State
    list,
    loading,
    error,
    functions,
    triggers,
    catalogReady,

    // Actions
    fetchList,
    create,
    updateInList,
    remove,
    loadFunctions,
    loadTriggers,
    loadCatalog,
    reset,
  }
})
