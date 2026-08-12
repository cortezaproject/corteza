import { defineStore } from 'pinia'
import { inject, reactive } from 'vue'

export const useLabelsStore = defineStore('labels', () => {
  const $ComposeAPI = inject('$ComposeAPI')

  const namespaces = reactive({})
  const modules = reactive({})

  function getNamespace(namespaceID) {
    return namespaces[namespaceID]
  }

  function getModule(moduleID) {
    return modules[moduleID]
  }

  async function resolveNamespace({ namespaceID }) {
    if (namespaces[namespaceID] !== undefined) {
      return namespaces[namespaceID]
    }

    namespaces[namespaceID] = null

    try {
      const namespace = await $ComposeAPI.namespaceRead({ namespaceID })
      const name = namespace.name || namespace.slug || namespaceID
      namespaces[namespaceID] = name
      return name
    } catch {
      namespaces[namespaceID] = namespaceID
      return namespaceID
    }
  }

  async function resolveModule({ moduleID, namespaceID }) {
    if (modules[moduleID] !== undefined) {
      return modules[moduleID]
    }

    modules[moduleID] = null

    try {
      const module = await $ComposeAPI.moduleRead({ namespaceID, moduleID })
      const name = module.name || module.handle || moduleID
      modules[moduleID] = name
      return name
    } catch {
      modules[moduleID] = moduleID
      return moduleID
    }
  }

  async function resolveMultipleNamespaces({ namespaceIDs }) {
    return Promise.all(namespaceIDs.map(namespaceID => resolveNamespace({ namespaceID })))
  }

  async function resolveMultipleModules({ modules: modList }) {
    return Promise.all(
      modList.map(({ moduleID, namespaceID }) => resolveModule({ moduleID, namespaceID })),
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
