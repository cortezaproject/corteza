// client/web/compose/src/stores/translator.ts
import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface ResourceTranslation {
  resource: string
  key: string
  lang: string
  message: string
}

// One row of the translation table. A key alone does not name a row: a set
// spans several resources and the same key (`title`) appears under each.
export interface TranslationTarget {
  resource: string
  key: string
}

export interface TranslatorConfig {
  resource: string
  titles: Record<string, string>
  highlight?: TranslationTarget
  fetcher: () => Promise<ResourceTranslation[]>
  updater: (changes: ResourceTranslation[]) => Promise<void>
  keyPrettifier?: (key: string) => string
}

export const useTranslatorStore = defineStore('translator', () => {
  const visible = ref(false)
  const config = ref<TranslatorConfig | null>(null)

  function open(c: TranslatorConfig): void {
    config.value = c
    visible.value = true
  }

  function close(): void {
    visible.value = false
    config.value = null
  }

  return { visible, config, open, close }
})
