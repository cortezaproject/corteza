<template>
  <div v-if="issues && issues.length > 0" class="flex flex-col gap-4">
    <Message
      v-for="(issue, index) in issues"
      :key="index"
      severity="error"
      :closable="false"
      class="w-full"
    >
      {{ issue.issue || issue }}
    </Message>
  </div>
  <div v-else class="flex flex-col items-center justify-center p-8 text-muted-color">
    <i class="pi pi-check-circle text-4xl mb-4 text-green-500"></i>
    <p class="m-0 text-lg">{{ $t('module.edit.issues.noIssues', 'No configuration issues detected.') }}</p>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  module: {
    type: Object,
    required: true,
  },
})

// Use actual server-side issues from module.issues (matching Corteza behavior)
const issues = computed(() => {
  return (props.module?.issues || [])
})
</script>
