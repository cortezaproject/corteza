<!-- client/web/compose/src/components/Translator/CTranslatorButton.vue -->
<template>
  <Button
    v-if="showTranslatorButton"
    icon="pi pi-language"
    size="small"
    :disabled="disabled"
    v-tooltip.bottom="$t('translator.button.label')"
    @click="handleClick"
  />
</template>

<script setup lang="ts">
import {
  useTranslatorStore,
  type TranslationTarget,
  type TranslatorConfig,
} from '@/sections/compose/stores/translator'
import { useResourceTranslations } from '@/sections/compose/composables/useResourceTranslations'

const props = defineProps<{
  resource: string
  titles: Record<string, string>
  fetcher: () => Promise<any[]>
  updater: (changes: any[]) => Promise<void>
  highlight?: TranslationTarget
  keyPrettifier?: (key: string) => string
  disabled?: boolean
}>()

const translatorStore = useTranslatorStore()
const { showTranslatorButton } = useResourceTranslations()

function handleClick(): void {
  translatorStore.open({
    resource: props.resource,
    titles: props.titles,
    highlight: props.highlight,
    fetcher: props.fetcher,
    updater: props.updater,
    keyPrettifier: props.keyPrettifier,
  } satisfies TranslatorConfig)
}
</script>
