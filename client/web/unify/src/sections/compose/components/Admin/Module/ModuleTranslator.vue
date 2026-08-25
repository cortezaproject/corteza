<!-- client/web/compose/src/components/Admin/Module/ModuleTranslator.vue -->
<template>
  <CTranslatorButton
    v-if="isEdit"
    :resource="resource"
    :titles="titles"
    :fetcher="fetcher"
    :updater="updater"
    :key-prettifier="keyPrettifier"
    :disabled="disabled"
  />
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import { useI18n } from 'vue-i18n'
import { compose } from '@planetcrust/human-js'
import CTranslatorButton from '@/sections/compose/components/Translator/CTranslatorButton.vue'
import {
  applyModuleTranslations,
  moduleFieldKeyLabel,
} from '@/sections/compose/lib/resource-translations'
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

const isEdit = computed(() => props.module?.moduleID && props.module.moduleID !== '0')

const resource = computed(
  () => `compose:module/${props.namespace.namespaceID}/${props.module.moduleID}`,
)

const titles = computed(() => {
  const tt: Record<string, string> = {}
  tt[resource.value] = t('translator.resources.module.title', {
    handle: props.module.name || props.module.handle || props.module.moduleID,
  })
  for (const field of props.module.fields || []) {
    const fRes = `compose:module-field/${props.namespace.namespaceID}/${props.module.moduleID}/${field.fieldID}`
    tt[fRes] = t('translator.resources.module.field.title', { name: field.label || field.name })
  }
  return tt
})

// The set spans the module and every one of its fields; a field key reads as
// its own name rather than the raw dotted path.
function keyPrettifier(key: string): string {
  return moduleFieldKeyLabel(key, t)
}

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
  // Reconstructed, not handed back as the plain clone it was applied to: the
  // editor holds a compose.Module and calls its methods.
  const updated = JSON.parse(JSON.stringify(props.module))
  applyModuleTranslations(updated, fresh, currentLanguage.value)
  emit('update:module', new compose.Module(updated))
}
</script>
