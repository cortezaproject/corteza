<template>
  <div class="flex flex-col h-full">
    <!-- Header -->
    <div class="shrink-0 flex items-center justify-between px-3 py-2 border-b border-surface">
      <div class="flex items-center gap-2">
        <h4 class="text-sm font-semibold text-color truncate max-w-[200px]" :title="stepName">
          {{ stepName || 'Trace Panel' }}
        </h4>
      </div>
    </div>

    <!-- Meta/Duration -->
    <div class="shrink-0 flex items-center gap-2 px-3 py-1.5 border-b border-surface">
      <Tag
        :severity="hasError ? 'danger' : 'success'"
        :icon="hasError ? 'pi pi-times-circle' : 'pi pi-check-circle'"
        :value="hasError ? $t('builder.trace.failed') : $t('builder.trace.completed')"
        class="text-xs"
      />
      <div v-if="duration" class="flex items-center gap-1 text-xs text-muted-color ml-1">
        <i class="pi pi-clock text-[10px]" />
        <span>{{ duration }}</span>
      </div>
    </div>

    <!-- Error -->
    <div v-if="errorMessage" class="shrink-0 px-3 pt-2">
      <div class="text-xs font-semibold text-red-500 mb-1">{{ $t('builder.trace.error') }}</div>
      <div class="bg-emphasis rounded-border p-2 text-xs text-red-500 whitespace-pre-wrap break-words">{{ errorMessage }}</div>
    </div>

    <!-- Tabs: Input / Output / Scope -->
    <div class="flex-1 min-h-0 flex flex-col px-3 pt-2 pb-2 overflow-hidden">
      <Tabs v-model:value="activeTab" class="flex-1 min-h-0 flex flex-col">
        <TabList class="shrink-0">
          <Tab value="input">{{ $t('builder.trace.input') }}</Tab>
          <Tab value="output">{{ $t('builder.trace.output') }}</Tab>
          <Tab v-if="hasScope" value="scope">{{ $t('builder.trace.scope') }}</Tab>
        </TabList>
        <TabPanels class="flex-1 min-h-0 overflow-hidden !p-0 flex flex-col">
          <TabPanel value="input" class="flex flex-col gap-2 pt-2 h-full">
            <InputText
              v-model="searchInput"
              :placeholder="$t('builder.trace.searchPlaceholder')"
              size="small"
              class="w-full shrink-0"
            />
            <div class="flex-1 min-h-0 overflow-auto">
              <TraceValueTree :value="frame?.args" :search="searchInput" />
            </div>
          </TabPanel>
          <TabPanel value="output" class="flex flex-col gap-2 pt-2 h-full">
            <InputText
              v-model="searchOutput"
              :placeholder="$t('builder.trace.searchPlaceholder')"
              size="small"
              class="w-full shrink-0"
            />
            <div class="flex-1 min-h-0 overflow-auto">
              <TraceValueTree :value="frame?.output" :search="searchOutput" />
            </div>
          </TabPanel>
          <TabPanel v-if="hasScope" value="scope" class="flex flex-col gap-2 pt-2 h-full">
            <InputText
              v-model="searchScope"
              :placeholder="$t('builder.trace.searchPlaceholder')"
              size="small"
              class="w-full shrink-0"
            />
            <div class="flex-1 min-h-0 overflow-auto">
              <TraceValueTree :value="frame?.input" :search="searchScope" />
            </div>
          </TabPanel>
        </TabPanels>
      </Tabs>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import TraceValueTree from './TraceValueTree.vue'

const props = defineProps({
  frame: { type: Object, default: null },
  executionError: { type: String, default: '' },
  stepName: { type: String, default: '' },
})

const activeTab = ref('input')
const searchInput = ref('')
const searchOutput = ref('')
const searchScope = ref('')

const hasError = computed(() => !!props.frame?.error || !!props.executionError)

const errorMessage = computed(() => {
  const raw = props.frame?.error || props.executionError || ''
  if (!raw) return ''
  if (typeof raw === 'object') return raw.message || String(raw)
  if (typeof raw === 'string') {
    const trimmed = raw.trim()
    if (trimmed.startsWith('{')) {
      try {
        const parsed = JSON.parse(trimmed)
        if (parsed && typeof parsed === 'object' && parsed.message) return parsed.message
      } catch { /* fall through */ }
    }
    return raw
  }
  return String(raw)
})

const hasScope = computed(() => {
  const s = props.frame?.input
  return !!s && typeof s === 'object' && Object.keys(s).length > 0
})

const duration = computed(() => {
  if (!props.frame?.startedAt || !props.frame?.endedAt) return null
  const ms = new Date(props.frame.endedAt).getTime() - new Date(props.frame.startedAt).getTime()
  if (ms < 1000) return `${ms}ms`
  return `${(ms / 1000).toFixed(1)}s`
})

// Reset searches & active tab when switching frames.
watch(() => props.frame?.id, () => {
  searchInput.value = ''
  searchOutput.value = ''
  searchScope.value = ''
  activeTab.value = 'input'
})

// If Scope tab becomes unavailable while active, fall back.
watch(hasScope, (v) => {
  if (!v && activeTab.value === 'scope') activeTab.value = 'input'
})
</script>

<style scoped>
.trace-json {
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
  white-space: pre-wrap;
  word-break: break-all;
  line-height: 1.4;
}
</style>
