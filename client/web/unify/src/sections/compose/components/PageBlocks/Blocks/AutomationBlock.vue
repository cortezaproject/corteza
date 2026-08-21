<template>
  <PageBlock :block="block" :record="record">
    <AutomationButtons
      v-if="buttons.length"
      :buttons="buttons"
      :namespace="namespace"
      :page="page"
      :module="pageModule"
      :record="record"
      container-class="flex flex-wrap gap-2 p-3"
      button-class="flex-auto min-w-[150px] whitespace-normal"
    />
    <div v-else class="flex items-center justify-center h-full p-3 text-muted-color italic">
      {{ $t('block.automation.noScripts') }}
    </div>
  </PageBlock>
</template>

<script setup>
import { computed } from 'vue'
import { useModuleStore } from '@planetcrust/human-vue'
import PageBlock from './PageBlock.vue'
import AutomationButtons from '../Shared/AutomationButtons.vue'

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
  record: { type: Object, default: undefined },
})

const moduleStore = useModuleStore()

const buttons = computed(() => props.block.options?.buttons || [])

// Blocks are handed namespace, page and record but never a module, so the one a
// manual automation runs against is read off whichever of the two carries it.
const pageModule = computed(() => {
  const moduleID = props.record?.moduleID || props.page?.moduleID
  if (!moduleID || moduleID === '0') return null
  return moduleStore.getByID(moduleID) || null
})
</script>
