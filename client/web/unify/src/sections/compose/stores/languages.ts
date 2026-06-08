// client/web/compose/src/stores/languages.ts
import { defineStore } from 'pinia'
import { inject, ref } from 'vue'

export interface Language {
  tag: string
  name: string
  localizedName: string
}

export const useLanguagesStore = defineStore('languages', () => {
  const $SystemAPI = inject('$SystemAPI') as any

  const set = ref<Language[]>([])
  const loaded = ref(false)
  const loading = ref(false)

  async function load(): Promise<void> {
    if (loaded.value || loading.value) return
    loading.value = true
    try {
      const { set: langs = [] } = await $SystemAPI.localeList()
      set.value = langs
      loaded.value = true
    } finally {
      loading.value = false
    }
  }

  return { set, loaded, loading, load }
})
