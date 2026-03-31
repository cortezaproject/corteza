<template>
  <div class="flex flex-col h-full">
    <!-- Header -->
    <div class="flex items-center justify-between px-3 py-2 border-b border-surface">
      <div class="flex items-center gap-2">
        <Tag
          :severity="hasError ? 'danger' : 'success'"
          :icon="hasError ? 'pi pi-times-circle' : 'pi pi-check-circle'"
          :value="hasError ? $t('builder.trace.failed') : $t('builder.trace.completed')"
          class="text-xs"
        />
      </div>
      <Button icon="pi pi-times" text rounded size="small" @click="handleClose" />
    </div>

    <!-- Duration -->
    <div v-if="duration" class="flex items-center gap-1 px-3 py-2 border-b border-surface text-xs text-muted-color">
      <i class="pi pi-clock text-xs" />
      <span>{{ duration }}</span>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-auto px-3 py-2">
      <!-- Error -->
      <div v-if="errorMessage" class="mb-3">
        <div class="text-xs font-semibold text-red-500 mb-1">{{ $t('builder.trace.error') }}</div>
        <pre class="trace-json bg-emphasis rounded-border p-2 text-xs text-red-500 overflow-auto max-h-24">{{ errorMessage }}</pre>
      </div>

      <!-- Input -->
      <div class="mb-3">
        <div class="text-xs font-semibold text-muted-color mb-1">{{ $t('builder.trace.input') }}</div>
        <pre class="trace-json bg-emphasis rounded-border p-2 text-xs overflow-auto max-h-48">{{ formatJson(frame?.input) }}</pre>
      </div>

      <!-- Output -->
      <div>
        <div class="text-xs font-semibold text-muted-color mb-1">{{ $t('builder.trace.output') }}</div>
        <pre class="trace-json bg-emphasis rounded-border p-2 text-xs overflow-auto max-h-48">{{ formatJson(frame?.output) }}</pre>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  frame: { type: Object, default: null },
  executionError: { type: String, default: '' },
})

const emit = defineEmits(['close'])

const { t } = useI18n()

const hasError = computed(() => !!props.frame?.error || !!props.executionError)

const errorMessage = computed(() => props.frame?.error || props.executionError || '')

const duration = computed(() => {
  if (!props.frame?.startedAt || !props.frame?.endedAt) return null
  const ms = new Date(props.frame.endedAt).getTime() - new Date(props.frame.startedAt).getTime()
  if (ms < 1000) return `${ms}ms`
  return `${(ms / 1000).toFixed(1)}s`
})

function handleClose() {
  emit('close')
}

function formatJson(data) {
  if (!data || (typeof data === 'object' && Object.keys(data).length === 0)) {
    return '—'
  }
  try {
    return JSON.stringify(data, null, 2)
  } catch {
    return String(data)
  }
}
</script>

<style scoped>
.trace-json {
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
  white-space: pre-wrap;
  word-break: break-all;
  line-height: 1.4;
}
</style>
