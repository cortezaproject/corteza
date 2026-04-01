<template>
  <div v-if="issues && issues.length > 0" class="flex flex-col gap-4">
    <div class="flex items-center gap-2 mb-2">
      <h3 class="text-lg font-medium text-primary m-0">
        {{ $t('module.edit.issues.title', 'Configuration Issues') }}
      </h3>
    </div>

    <Message
      v-for="(issue, index) in issues"
      :key="index"
      :severity="issue.severity || 'warn'"
      :closable="false"
      class="w-full"
    >
      <div class="flex flex-col gap-1">
        <span class="font-bold">{{ issue.summary }}</span>
        <span v-if="issue.details" class="text-sm font-normal">{{ issue.details }}</span>
      </div>
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

// Compute simple structure warnings or fallback issues dynamically
const issues = computed(() => {
  const list = []
  
  if (!props.module) return list

  if (props.module.fields && props.module.fields.length === 0) {
    list.push({
      severity: 'warn',
      summary: 'No Fields Defined',
      details: 'This module does not have any fields configured. It will not be able to store any custom data.',
    })
  }

  const { recordDeDup, privacy } = props.module.config || {}
  
  if (recordDeDup?.enabled && (!recordDeDup.rules || recordDeDup.rules.length === 0)) {
    list.push({
      severity: 'info',
      summary: 'Duplicate Detection Enabled But Empty',
      details: 'Duplicate detection is enabled but no rules have been configured yet.',
    })
  }

  return list
})
</script>
