<template>
  <PageBlock :block="block">
    <div v-if="loading" class="flex items-center justify-center h-full">
      <ProgressSpinner style="width: 28px; height: 28px" />
    </div>

    <div v-else class="w-full h-full">
      <!-- Map placeholder — requires Leaflet integration -->
      <div class="flex flex-col items-center justify-center h-full p-5 text-muted-color">
        <i class="pi pi-map text-4xl mb-3" />
        <p class="text-center mb-2 font-medium text-color">
          {{ $t('block.geometry.label') }}
        </p>
        <p class="text-center text-sm">
          {{ feedSummary }}
        </p>
      </div>
    </div>
  </PageBlock>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import PageBlock from './PageBlock.vue'

const { t } = useI18n()

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
  record: { type: Object, default: undefined },
})

const loading = ref(false)

const options = computed(() => props.block.options || {})
const feeds = computed(() => options.value.feeds || [])

const feedSummary = computed(() => {
  const count = feeds.value.length
  if (!count) return t('block.noConfiguration')
  return `${count} source${count > 1 ? 's' : ''} configured`
})
</script>
