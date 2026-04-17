import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

const resourcePrefix = 'corteza::'

interface Rule {
  resource: string
  operation: string
  allow: boolean
}

interface APIClient {
  permissionsEffective: (_a?: Record<string, unknown>) => Promise<Rule[]>
}

export const useRBACStore = defineStore('rbac', () => {
  const loaded = ref(false)
  const rules = ref<Rule[]>([])

  function can(resource: string, operation: string): boolean {
    const rule = rules.value.find(
      (r) => r.resource === resourcePrefix + resource && r.operation === operation,
    )
    return rule?.allow ?? false
  }

  async function load(apis: APIClient[]) {
    loaded.value = false
    const results = await Promise.all(
      apis.map((api) => api.permissionsEffective({}).catch(() => [] as Rule[])),
    )
    rules.value = results
      .flat()
      .filter(({ resource }) => resource.startsWith(resourcePrefix))
    loaded.value = true
  }

  function clear() {
    rules.value = []
    loaded.value = false
  }

  return {
    loaded: computed(() => loaded.value),
    rules: computed(() => rules.value),
    can,
    load,
    clear,
  }
})
