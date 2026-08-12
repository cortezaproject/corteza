<!-- client/web/compose/src/components/Admin/Chart/ChartTranslator.vue -->
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

const props = defineProps<{
  chart: any
  namespace: any
  disabled?: boolean
}>()

const emit = defineEmits<{
  'update:chart': [chart: any]
}>()

const $ComposeAPI = inject('$ComposeAPI') as any
const { t } = useI18n()

const isEdit = computed(() => props.chart?.chartID && props.chart.chartID !== '0')

const resource = computed(
  () => `compose:chart/${props.namespace.namespaceID}/${props.chart.chartID}`,
)

const titles = computed(() => ({
  [resource.value]: t('translator.resources.chart.title', {
    handle: props.chart.name || props.chart.handle || props.chart.chartID,
  }),
}))

function fetcher() {
  return $ComposeAPI.chartListTranslations({
    namespaceID: props.namespace.namespaceID,
    chartID: props.chart.chartID,
  })
}

async function updater(changes: any[]) {
  await $ComposeAPI.chartUpdateTranslations({
    namespaceID: props.namespace.namespaceID,
    chartID: props.chart.chartID,
    translations: changes,
  })
  // Re-fetch the chart so backend-applied translations (nested config labels) are reflected
  const fresh = await $ComposeAPI.chartRead({
    namespaceID: props.namespace.namespaceID,
    chartID: props.chart.chartID,
  })
  emit('update:chart', fresh)
}
</script>
