<template>
  <PageBlock :block="block">
    <div v-if="block.options?.url || block.options?.srcField" class="h-full">
      <iframe
        :src="resolvedUrl"
        class="w-full h-full border-0"
        :title="block.title || $t('block.iframe.label')"
        sandbox="allow-scripts allow-same-origin allow-popups allow-forms"
      />
    </div>
    <div v-else class="flex items-center justify-center h-full p-3 text-muted-color italic">
      {{ $t('block.iframe.noInput') }}
    </div>
  </PageBlock>
</template>

<script setup>
import { computed, inject } from 'vue'
import PageBlock from './PageBlock.vue'

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

// If srcField is set and we have a record context, use that
const recordViewContext = inject('recordViewContext', null)

const resolvedUrl = computed(() => {
  const { url, srcField } = props.block.options || {}

  if (srcField && recordViewContext?.record?.value) {
    const fieldValue = recordViewContext.record.value.values?.[srcField]
    if (fieldValue) return fieldValue
  }

  return url || ''
})
</script>
