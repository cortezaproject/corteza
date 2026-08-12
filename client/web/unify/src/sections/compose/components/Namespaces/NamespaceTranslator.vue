<!-- client/web/compose/src/components/Namespaces/NamespaceTranslator.vue -->
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
import { applyNamespaceTranslations } from '@/sections/compose/lib/resource-translations'
import { useResourceTranslations } from '@/sections/compose/composables/useResourceTranslations'

const props = defineProps<{
  namespace: any
  disabled?: boolean
}>()

const emit = defineEmits<{
  'update:namespace': [namespace: any]
}>()

const $ComposeAPI = inject('$ComposeAPI') as any
const { t } = useI18n()
const { currentLanguage } = useResourceTranslations()

const isEdit = computed(() => props.namespace?.namespaceID && props.namespace.namespaceID !== '0')

const resource = computed(() => `compose:namespace/${props.namespace.namespaceID}`)

const titles = computed(() => ({
  [resource.value]: t('translator.resources.namespace.title', {
    name: props.namespace.name || props.namespace.slug || props.namespace.namespaceID,
  }),
}))

function fetcher() {
  return $ComposeAPI.namespaceListTranslations({
    namespaceID: props.namespace.namespaceID,
  })
}

async function updater(changes: any[]) {
  await $ComposeAPI.namespaceUpdateTranslations({
    namespaceID: props.namespace.namespaceID,
    translations: changes,
  })
  const fresh = await $ComposeAPI.namespaceListTranslations({
    namespaceID: props.namespace.namespaceID,
  })
  const updated = JSON.parse(JSON.stringify(props.namespace))
  applyNamespaceTranslations(updated, fresh, currentLanguage.value)
  emit('update:namespace', updated)
}
</script>
