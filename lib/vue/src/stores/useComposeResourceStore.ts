import { defineStore } from 'pinia'
import { inject, ref } from 'vue'

export const useComposeResourceStore = defineStore(
  'composeResource',
  () => {
    const $ComposeAPI = inject('$ComposeAPI') as any

    const namespaces = ref(new Map<string, any>())
    const modules = ref(new Map<string, any>())
    const charts = ref(new Map<string, any>())

    function cacheNamespace(ns: any) {
      if (ns?.namespaceID) {
        namespaces.value.set(ns.namespaceID, ns)
      }
    }

    function cacheModule(mod: any) {
      if (mod?.moduleID && mod?.namespaceID) {
        modules.value.set(`${mod.namespaceID}:${mod.moduleID}`, mod)
      }
    }

    function cacheChart(chart: any) {
      if (chart?.chartID && chart?.namespaceID) {
        charts.value.set(`${chart.namespaceID}:${chart.chartID}`, chart)
      }
    }

    function getNamespace(id: string) {
      return namespaces.value.get(id)
    }

    function getModule(nsID: string, modID: string) {
      return modules.value.get(`${nsID}:${modID}`)
    }

    function getChart(nsID: string, chartID: string) {
      return charts.value.get(`${nsID}:${chartID}`)
    }

    async function resolveNamespace(id: string) {
      if (!id || !$ComposeAPI) return undefined

      const cached = namespaces.value.get(id)
      if (cached) return cached

      const ns = await $ComposeAPI.namespaceRead({ namespaceID: id })
      cacheNamespace(ns)
      return ns
    }

    async function resolveModule(nsID: string, modID: string) {
      if (!nsID || !modID || !$ComposeAPI) return undefined

      const key = `${nsID}:${modID}`
      const cached = modules.value.get(key)
      if (cached) return cached

      const mod = await $ComposeAPI.moduleRead({
        namespaceID: nsID,
        moduleID: modID,
      })
      cacheModule(mod)
      return mod
    }

    async function resolveChart(nsID: string, chartID: string) {
      if (!nsID || !chartID || !$ComposeAPI) return undefined

      const key = `${nsID}:${chartID}`
      const cached = charts.value.get(key)
      if (cached) return cached

      const chart = await $ComposeAPI.chartRead({
        namespaceID: nsID,
        chartID,
      })
      cacheChart(chart)
      return chart
    }

    function searchNamespaces(params: Record<string, any> = {}) {
      if (!$ComposeAPI) {
        return {
          response: () => Promise.resolve({ set: [] }),
          cancel: () => {},
        }
      }

      const { response, cancel } = $ComposeAPI.namespaceListCancellable(params)

      return {
        response: async () => {
          const result = await response()
          const set = result.set || []
          set.forEach(cacheNamespace)
          return result
        },
        cancel,
      }
    }

    function searchModules(
      nsID: string,
      params: Record<string, any> = {},
    ) {
      if (!nsID || !$ComposeAPI) {
        return {
          response: () => Promise.resolve({ set: [] }),
          cancel: () => {},
        }
      }

      const { response, cancel } = $ComposeAPI.moduleListCancellable({
        namespaceID: nsID,
        ...params,
      })

      return {
        response: async () => {
          const result = await response()
          const set = result.set || []
          set.forEach(cacheModule)
          return result
        },
        cancel,
      }
    }

    function searchCharts(
      nsID: string,
      params: Record<string, any> = {},
    ) {
      if (!nsID || !$ComposeAPI) {
        return {
          response: () => Promise.resolve({ set: [] }),
          cancel: () => {},
        }
      }

      const { response, cancel } = $ComposeAPI.chartListCancellable({
        namespaceID: nsID,
        ...params,
      })

      return {
        response: async () => {
          const result = await response()
          const set = result.set || []
          set.forEach(cacheChart)
          return result
        },
        cancel,
      }
    }

    return {
      namespaces,
      modules,
      charts,
      getNamespace,
      getModule,
      getChart,
      resolveNamespace,
      resolveModule,
      resolveChart,
      searchNamespaces,
      searchModules,
      searchCharts,
    }
  },
)
