<!-- client/web/compose/src/components/Admin/Module/ModuleTranslator.vue -->
<template>
  <CTranslatorButton
    v-if="isEdit"
    :resource="resource"
    :titles="titles"
    :fetcher="fetcher"
    :updater="updater"
    :disabled="disabled"
  />
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import { useI18n } from 'vue-i18n'
import CTranslatorButton from '@/sections/compose/components/Translator/CTranslatorButton.vue'
import { applyModuleTranslations } from '@/sections/compose/lib/resource-translations'
import { useResourceTranslations } from '@/sections/compose/composables/useResourceTranslations'

const props = defineProps<{
  module: any
  namespace: any
  disabled?: boolean
}>()

const emit = defineEmits<{ 'update:module': [module: any] }>()

const $ComposeAPI = inject('$ComposeAPI') as any
const { t } = useI18n()
const { currentLanguage } = useResourceTranslations()

const isEdit = computed(() =>
  props.module?.moduleID && props.module.moduleID !== '0',
)

const resource = computed(() =>
  `compose:module/${props.namespace.namespaceID}/${props.module.moduleID}`,
)

const titles = computed(() => {
  const tt: Record<string, string> = {}
  tt[resource.value] = t('translator.resources.module.title', { handle: props.module.name || props.module.handle || props.module.moduleID })
  for (const field of props.module.fields || []) {
    const fRes = `compose:module-field/${props.namespace.namespaceID}/${props.module.moduleID}/${field.fieldID}`
    tt[fRes] = t('translator.resources.module.field.title', { name: field.label || field.name })
  }
  return tt
})

function fetcher() {
  return $ComposeAPI.moduleListTranslations({
    namespaceID: props.namespace.namespaceID,
    moduleID: props.module.moduleID,
  })
}

async function updater(changes: any[]) {
  await $ComposeAPI.moduleUpdateTranslations({
    namespaceID: props.namespace.namespaceID,
    moduleID: props.module.moduleID,
    translations: changes,
  })
  // Apply to module object so field labels update immediately for current language
  const fresh = await fetcher()
  const updated = JSON.parse(JSON.stringify(props.module))
  applyModuleTranslations(updated, fresh, currentLanguage.value)
  emit('update:module', updated)
}
</script>
