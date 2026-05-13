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
import CTranslatorButton from '@/components/Translator/CTranslatorButton.vue'

const props = defineProps<{
  namespace: any
  disabled?: boolean
}>()

const $ComposeAPI = inject('$ComposeAPI') as any
const { t } = useI18n()

const isEdit = computed(() =>
  props.namespace?.namespaceID && props.namespace.namespaceID !== '0',
)

const resource = computed(() =>
  `compose:namespace/${props.namespace.namespaceID}`,
)

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
}
</script>
