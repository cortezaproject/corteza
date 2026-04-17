import { defineStore } from 'pinia'
import { reactive } from 'vue'

export const useLabelsStore = defineStore('labels', () => {
  const namespaces = reactive({})
  const modules = reactive({})

  function getNamespace(namespaceID) {
    return namespaces[namespaceID]
  }

  function getModule(moduleID) {
    return modules[moduleID]
  }

  async function resolveNamespace({ namespaceID, api }) {
    // Return cached value if available
    if (namespaces[namespaceID] !== undefined) {
      return namespaces[namespaceID]
    }

    // Mark as loading to prevent duplicate requests
    namespaces[namespaceID] = null

    try {
      const namespace = await api.namespaceRead({ namespaceID })
      const name = namespace.name || namespace.slug || namespaceID
      namespaces[namespaceID] = name
      return name
    } catch {
      namespaces[namespaceID] = namespaceID
      return namespaceID
    }
  }

  async function resolveModule({ moduleID, namespaceID, api }) {
    // Return cached value if available
    if (modules[moduleID] !== undefined) {
      return modules[moduleID]
    }

    // Mark as loading to prevent duplicate requests
    modules[moduleID] = null

    try {
      const module = await api.moduleRead({ namespaceID, moduleID })
      const name = module.name || module.handle || moduleID
      modules[moduleID] = name
      return name
    } catch {
      modules[moduleID] = moduleID
      return moduleID
    }
  }

  async function resolveMultipleNamespaces({ namespaceIDs, api }) {
    return Promise.all(
      namespaceIDs.map(namespaceID =>
        resolveNamespace({ namespaceID, api }),
      ),
    )
  }

  async function resolveMultipleModules({ modules: modList, api }) {
    return Promise.all(
      modList.map(({ moduleID, namespaceID }) =>
        resolveModule({ moduleID, namespaceID, api }),
      ),
    )
  }

  return {
    namespaces,
    modules,
    getNamespace,
    getModule,
    resolveNamespace,
    resolveModule,
    resolveMultipleNamespaces,
    resolveMultipleModules,
  }
})
