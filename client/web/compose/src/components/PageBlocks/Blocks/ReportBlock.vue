<template>
  <PageBlock :block="block">
    <div v-if="loading" class="flex items-center justify-center h-full">
      <ProgressSpinner style="width: 28px; height: 28px" />
    </div>

    <div v-else class="flex items-center justify-center h-full p-3 text-muted-color">
      <div class="text-center">
        <i class="pi pi-chart-bar text-4xl mb-3" />
        <p class="font-medium text-color mb-2">
          {{ $t('block.report.label') }}
        </p>
        <p v-if="options.reportID && options.reportID !== '0'" class="text-sm">
          Report ID: {{ options.reportID }}
          <span v-if="options.scenarioID"> — {{ options.scenarioID }}</span>
        </p>
        <p v-else class="text-sm italic">
          {{ $t('block.noConfiguration') }}
        </p>
      </div>
    </div>
  </PageBlock>
</template>

<script setup>
import { ref, computed } from 'vue'
import PageBlock from './PageBlock.vue'

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
  record: { type: Object, default: undefined },
})

const loading = ref(false)
const options = computed(() => props.block.options || {})
</script>
